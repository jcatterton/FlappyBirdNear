package Game

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	Network "github.com/jcatterton/GoNeat/GoNeat"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
)

var score = 0

func Start() {
	pixelgl.Run(run)
}

func run() {
	linesVisible := false
	brainVisible := true

	cfg := pixelgl.WindowConfig{
		Title:  "Flappy Bird",
		Bounds: pixel.R(0, 0, 500, 750),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	birdImage, err := loadPicture("./bird.png")
	if err != nil {
		panic(err)
	}

	pipeImage, err := loadPicture("./pipe.png")
	if err != nil {
		panic(err)
	}

	backgroundImage, err := loadPicture("./background.png")
	if err != nil {
		panic(err)
	}
	background := pixel.NewSprite(backgroundImage, backgroundImage.Bounds())

	/*bird := Bird{
		375,
		0,
		*pixel.NewSprite(birdImage, birdImage.Bounds()),
	}*/

	pipes := make([]*Pipe, 6)
	for i := range pipes {
		if i%2 == 0 {
			pipes[i] = &Pipe{float64(rand.Intn(350) + 100), 600 + float64(200*i), true, *pixel.NewSprite(pipeImage, pipeImage.Bounds())}
		} else {
			pipes[i] = pipes[i-1].CreateSisterPipe()
		}
	}

	imd := imdraw.New(nil)

	pop := Network.InitPopulation(4, 1, 3, 5)

	birds := make([]Bird, len(pop.GetAllGenomes()))
	for i := range birds {
		birds[i].height = 375
		birds[i].yVel = 0
		birds[i].sprite = *pixel.NewSprite(birdImage, birdImage.Bounds())
		birds[i].dead = false
	}

	for !win.Closed() {
		background.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

		score++

		win.SetTitle("Flappy Bird - " + strconv.Itoa(score))

		if win.JustPressed(pixelgl.Key1) {
			linesVisible = !linesVisible
		}
		if win.JustPressed(pixelgl.Key2) {
			brainVisible = !brainVisible
		}

		for x := range pop.GetAllGenomes() {
			if !birds[x].dead {
				if err := pop.GetAllGenomes()[x].TakeInput(
					[]float64{
						birds[x].GetYVel(),
						birds[x].GetInformationOnNextPipes(pipes)[0],
						birds[x].GetInformationOnNextPipes(pipes)[1],
						birds[x].GetInformationOnNextPipes(pipes)[2],
					}); err != nil {
					panic(err)
				}

				pop.GetAllGenomes()[x].FeedForward()
				output := pop.GetAllGenomes()[x].GetOutputs()[0]

				if output > 0.5 {
					birds[x].Jump()
				}
				birds[x].Fall()
				birds[x].Draw(win, pixel.IM.Moved(pixel.V(win.Bounds().W()/2, birds[x].height)).
					Rotated(pixel.V(win.Bounds().W()/2, birds[x].height), birds[x].yVel/35))

				if linesVisible {
					imd.Clear()
					imd.Color = colornames.White

					imd.Push(birds[x].Bounds().Center)
					imd.Push(pixel.V(birds[x].Bounds().Center.X+birds[x].GetInformationOnNextPipes(pipes)[0], birds[x].GetHeight()))
					imd.Line(2)

					imd.Push(birds[x].Bounds().Center)
					imd.Push(pixel.V(250, birds[x].Bounds().Center.Y+birds[x].Bounds().Radius+birds[x].GetInformationOnNextPipes(pipes)[1]))
					imd.Line(2)

					imd.Push(birds[x].Bounds().Center)
					imd.Push(pixel.V(250, birds[x].Bounds().Center.Y-birds[x].Bounds().Radius+birds[x].GetInformationOnNextPipes(pipes)[2]))
					imd.Line(2)

					imd.Draw(win)
				}
			}

			if brainVisible {
				for i := range pop.GetAllGenomes() {
					if !birds[i].dead {
						drawGenome(pop.GetAllGenomes()[i], win)
						break
					}
				}
			}

			if checkForCollisions(birds[x], pipes) && !birds[x].dead {
				pop.GetAllGenomes()[x].SetFitness(float64(score))
				birds[x].dead = true
			}
		}

		for i := range pipes {
			pipes[i].Draw(win, pixel.IM)
			pipes[i].MoveLeft()
			if pipes[i].xPos <= -200 && i%2 == 0 {
				pipes[i] = &Pipe{float64(rand.Intn(350) + 100), float64(1000), true, *pixel.NewSprite(pipeImage, pipeImage.Bounds())}
				pipes[i+1] = pipes[i].CreateSisterPipe()
			}
		}

		if allBirdsDead(birds) {
			for x := range pop.GetSpecies() {
				pop.GetSpecies()[x].SetChampion()
				pop.GetSpecies()[x].CullTheWeak()
				pop.GetSpecies()[x].OrderByFitness()
			}
			pop.SetGrandChampion()
			pop.ExtinctionEvent()
			pop.Mutate()
			log.Println(pop.GetGeneration(), " - ", pop.GetGrandChampion().GetFitness())
			score = 0

			for i := range birds {
				birds[i].height = 375
				birds[i].yVel = 0
				birds[i].sprite = *pixel.NewSprite(birdImage, birdImage.Bounds())
				birds[i].dead = false
			}
			for i := range pipes {
				if i%2 == 0 {
					pipes[i] = &Pipe{float64(rand.Intn(350) + 100), 600 + float64(200*i), true, *pixel.NewSprite(pipeImage, pipeImage.Bounds())}
				} else {
					pipes[i] = pipes[i-1].CreateSisterPipe()
				}
			}
		}

		win.Update()
		win.Clear(colornames.Black)
	}
}

func allBirdsDead(birds []Bird) bool {
	for i := range birds {
		if birds[i].dead == false {
			return false
		}
	}
	return true
}

func drawGenome(g *Network.Genome, win *pixelgl.Window) {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(0, 0), basicAtlas)

	imd := imdraw.New(nil)

	w := win.Bounds().W() / 2
	h := win.Bounds().H() / 3

	for i := 0; i < g.GetLayers(); i++ {
		for j := range g.GetNodesWithLayer(i + 1) {
			if g.GetNodesWithLayer(i + 1)[j].IsActivated() {
				imd.Color = pixel.RGB(0, 1, 0)
			} else {
				imd.Color = pixel.RGB(1, 0, 0)
			}
			imd.Push(pixel.V(
				(float64(i)+0.5)*(w/float64(g.GetLayers())),
				(float64(j)+0.5)*(h/float64(len(g.GetNodesWithLayer(i+1))))))
			imd.Circle(5, 20)

			for k := range g.GetNodesWithLayer(i + 1)[j].GetOutwardConnections() {
				imd.Color = pixel.RGB(1, 1, 1)
				imd.Push(
					pixel.V(
						(float64(i)+0.5)*(w/float64(g.GetLayers()))+10,
						(float64(j)+0.5)*(h/float64(len(g.GetNodesWithLayer(i+1))))),
					pixel.V(
						(float64(g.GetNodesWithLayer(i + 1)[j].GetOutwardConnections()[k].GetNodeB().GetLayer())-0.5)*(w/float64(g.GetLayers()))-10,
						(float64(Network.NodeIndex(g.GetNodesWithLayer(g.GetNodesWithLayer(i + 1)[j].GetOutwardConnections()[k].GetNodeB().GetLayer()),
							g.GetNodesWithLayer(i + 1)[j].GetOutwardConnections()[k].GetNodeB()))+0.5)*(h/float64(len(g.GetNodesWithLayer(g.GetNodesWithLayer(i + 1)[j].GetOutwardConnections()[k].GetNodeB().GetLayer()))))))
				imd.Line(2)
			}

			basicTxt.Color = colornames.White
			_, err := fmt.Fprintf(basicTxt, strconv.Itoa(g.GetNodesWithLayer(i + 1)[j].GetInnovationNumber()))
			if err != nil {
				panic(err)
			}
			basicTxt.Draw(win, pixel.IM.Moved(pixel.V(
				(float64(i)+0.5)*(w/float64(g.GetLayers()))-1,
				(float64(j)+0.5)*(h/float64(len(g.GetNodesWithLayer(i+1))))+20)))
			basicTxt.Clear()
		}
	}
	imd.Draw(win)
}

func checkForCollisions(bird Bird, pipes []*Pipe) bool {
	if bird.height <= 0 {
		return true
	}
	for i := range pipes {
		if bird.Bounds().IntersectRect(pipes[i].Bounds()).X != 0 || bird.Bounds().IntersectRect(pipes[i].Bounds()).Y != 0 {
			return true
		}
	}
	return false
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
