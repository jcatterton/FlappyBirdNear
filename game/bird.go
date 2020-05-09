package Game

import "github.com/faiface/pixel"

type Bird struct {
	height float64
	yVel   float64
	sprite pixel.Sprite
}

func (b *Bird) Draw(t pixel.Target, m pixel.Matrix) {
	b.sprite.Draw(t, m)
}

func (b *Bird) Bounds() pixel.Circle {
	return pixel.C(
		pixel.V(250, b.height),
		b.sprite.Picture().Bounds().H()/2,
	)
}

func (b *Bird) Fall() {
	if b.yVel > -50.0 {
		b.yVel--
	} else {
		b.yVel = -50.0
	}

	b.height = b.height - (b.yVel * -0.2)

	if b.height >= 720 {
		b.height = 720
	}
}

func (b *Bird) Jump() {
	b.yVel = 30.0
}

func (b *Bird) GetHeight() float64 {
	return b.height
}

func (b *Bird) GetYVel() float64 {
	return b.yVel
}

func (b *Bird) GetInformationOnNextPipes(pipes []*Pipe) []float64 {
	for i := range pipes {
		if pipes[i].xPos > 150 && pipes[i].xPos < 550 && i%2 == 0 {
			return []float64{pipes[i].xPos - 300, pipes[i+1].height - (b.height + b.Bounds().Radius), pipes[i].height - (b.height - b.Bounds().Radius)}
		}
	}
	return []float64{500.0, 100, -100}
}
