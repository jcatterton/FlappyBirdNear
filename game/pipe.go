package Game

import (
	"github.com/faiface/pixel"
	"math"
)

type Pipe struct {
	height	float64
	xPos	float64
	bottom	bool
	sprite	pixel.Sprite
}

func (p *Pipe) CreateSisterPipe() *Pipe {
	return &Pipe{height: p.height + 200, xPos: p.xPos, bottom: false, sprite: p.sprite}
}

func (p *Pipe) MoveLeft() {
	p.xPos = p.xPos - 2
}

func (p *Pipe) Draw(t pixel.Target, m pixel.Matrix) {
	if p.bottom {
		p.sprite.Draw(t, m.Moved(pixel.V(p.xPos, p.height - p.sprite.Picture().Bounds().H()/2)))
	} else {
		p.sprite.Draw(t, m.Moved(pixel.V(p.xPos, p.height + p.sprite.Picture().Bounds().H()/2)).Rotated(pixel.V(p.xPos, p.height + p.sprite.Picture().Bounds().H()/2), math.Pi))
	}
}

func (p *Pipe) Bounds() pixel.Rect {
	if p.bottom {
		return pixel.R(p.xPos - p.sprite.Picture().Bounds().W()/2, 0, p.xPos + p.sprite.Picture().Bounds().W()/2, p.height)
	}
	return pixel.R(p.xPos - p.sprite.Picture().Bounds().W()/2, p.height, p.xPos + p.sprite.Picture().Bounds().W()/2, 1000)
}
