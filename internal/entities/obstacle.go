package entities

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ObstacleType int

const (
	TypeGround ObstacleType = 0
	TypeAir    ObstacleType = 1
)

type Obstacle struct {
	X           float64
	Type        ObstacleType
	screenWidth float64
	groundY     float64
}

func NewObstacle(screenWidth, groundY float64) *Obstacle {
	return &Obstacle{
		X:           screenWidth,
		Type:        TypeGround,
		screenWidth: screenWidth,
		groundY:     groundY,
	}
}

func (obs *Obstacle) Update(speed float64) bool {
	obs.X -= speed
	if obs.X < -35 {
		obs.X = obs.screenWidth + float64(rand.Intn(45))
		obs.Type = ObstacleType(rand.Intn(2))
		return true
	}
	return false
}

func (obs *Obstacle) Reset() {
	obs.X = obs.screenWidth
	obs.Type = TypeGround
}

func (obs *Obstacle) GetBounds() (x, y, w, h float64) {
	if obs.Type == TypeGround {
		w = 22.0
		h = 24.0
		x = obs.X
		y = obs.groundY - h
	} else {
		w = 28.0
		h = 16.0
		x = obs.X
		y = obs.groundY - 34.0
	}
	return x, y, w, h
}

func (obs *Obstacle) CheckCollision(oncaX, oncaY, oncaW, oncaH float64) bool {
	obsX, obsY, obsW, obsH := obs.GetBounds()
	overlapX := obsX < oncaX+oncaW && obsX+obsW > oncaX
	overlapY := obsY < oncaY+oncaH && obsY+obsH > oncaY
	return overlapX && overlapY
}

func (obs *Obstacle) Draw(screen *ebiten.Image, ticks int, stage int) {
	if obs.Type == TypeGround {
		obsRealY := obs.groundY - 24.0

		switch stage {
		case 2:
			ebitenutil.DrawRect(screen, obs.X+2, obsRealY+5, 18, 19, color.RGBA{R: 110, G: 65, B: 35, A: 255})
			ebitenutil.DrawRect(screen, obs.X+4, obsRealY+7, 14, 15, color.RGBA{R: 135, G: 80, B: 45, A: 255})
			ebitenutil.DrawRect(screen, obs.X+1, obsRealY+8, 20, 2, color.RGBA{R: 215, G: 175, B: 55, A: 255})
			ebitenutil.DrawRect(screen, obs.X+1, obsRealY+17, 20, 2, color.RGBA{R: 215, G: 175, B: 55, A: 255})

		case 3:
			ebitenutil.DrawRect(screen, obs.X+2, obsRealY+10, 18, 14, color.RGBA{R: 120, G: 85, B: 50, A: 255})
			ebitenutil.DrawRect(screen, obs.X+4, obsRealY+5, 7, 7, color.RGBA{R: 235, G: 155, B: 30, A: 255})
			ebitenutil.DrawRect(screen, obs.X+6, obsRealY+4, 5, 4, color.RGBA{R: 215, G: 60, B: 35, A: 255})
			ebitenutil.DrawRect(screen, obs.X+11, obsRealY+6, 7, 7, color.RGBA{R: 240, G: 175, B: 35, A: 255})

		default:
			basketStraw := color.RGBA{R: 185, G: 135, B: 75, A: 255}
			basketDark := color.RGBA{R: 135, G: 90, B: 45, A: 255}
			acaiPurple := color.RGBA{R: 45, G: 15, B: 48, A: 255}
			acaiLight := color.RGBA{R: 85, G: 30, B: 85, A: 255}
			palmGreen := color.RGBA{R: 40, G: 145, B: 50, A: 255}

			ebitenutil.DrawRect(screen, obs.X+2, obsRealY+7, 18, 17, basketStraw)
			ebitenutil.DrawRect(screen, obs.X+4, obsRealY+9, 14, 13, basketDark)

			for weave := 0; weave < 3; weave++ {
				wy := obsRealY + 9 + float64(weave*5)
				ebitenutil.DrawRect(screen, obs.X+1, wy, 20, 2, basketDark)
			}

			ebitenutil.DrawRect(screen, obs.X+1, obsRealY+2, 20, 8, acaiPurple)
			ebitenutil.DrawRect(screen, obs.X+3, obsRealY+1, 16, 6, acaiLight)
			ebitenutil.DrawRect(screen, obs.X+5, obsRealY, 12, 3, acaiPurple)

			ebitenutil.DrawRect(screen, obs.X+17, obsRealY-3, 3, 7, palmGreen)
			ebitenutil.DrawRect(screen, obs.X+19, obsRealY-5, 4, 3, palmGreen)
			ebitenutil.DrawRect(screen, obs.X+15, obsRealY-1, 3, 2, palmGreen)

			ebitenutil.DrawRect(screen, obs.X-3, obsRealY+21, 3, 3, acaiPurple)
			ebitenutil.DrawRect(screen, obs.X-1, obsRealY+22, 2, 2, acaiLight)
		}
	} else {
		birdY := obs.groundY - 34.0

		switch stage {
		case 2:
			ebitenutil.DrawRect(screen, obs.X+6, birdY+4, 16, 7, color.RGBA{R: 240, G: 240, B: 245, A: 255})
			ebitenutil.DrawRect(screen, obs.X+1, birdY+3, 6, 6, color.RGBA{R: 245, G: 245, B: 250, A: 255})
			ebitenutil.DrawRect(screen, obs.X-2, birdY+5, 4, 2, color.RGBA{R: 240, G: 180, B: 30, A: 255})
			ebitenutil.DrawRect(screen, obs.X+2, birdY+4, 2, 2, color.RGBA{R: 30, G: 30, B: 30, A: 255})

			if (ticks/8)%2 == 0 {
				ebitenutil.DrawRect(screen, obs.X+9, birdY-5, 14, 9, color.RGBA{R: 200, G: 205, B: 215, A: 255})
			} else {
				ebitenutil.DrawRect(screen, obs.X+9, birdY+9, 14, 8, color.RGBA{R: 200, G: 205, B: 215, A: 255})
			}

		case 3:
			ebitenutil.DrawRect(screen, obs.X+6, birdY+4, 16, 7, color.RGBA{R: 35, G: 165, B: 60, A: 255})
			ebitenutil.DrawRect(screen, obs.X+1, birdY+3, 6, 6, color.RGBA{R: 50, G: 190, B: 75, A: 255})
			ebitenutil.DrawRect(screen, obs.X-2, birdY+5, 4, 3, color.RGBA{R: 245, G: 210, B: 40, A: 255})

			if (ticks/7)%2 == 0 {
				ebitenutil.DrawRect(screen, obs.X+9, birdY-5, 14, 9, color.RGBA{R: 25, G: 130, B: 45, A: 255})
			} else {
				ebitenutil.DrawRect(screen, obs.X+9, birdY+9, 14, 8, color.RGBA{R: 25, G: 130, B: 45, A: 255})
			}

		default:
			urubuBlack := color.RGBA{R: 25, G: 25, B: 28, A: 255}
			urubuHead := color.RGBA{R: 155, G: 110, B: 115, A: 255}
			urubuBeak := color.RGBA{R: 85, G: 80, B: 80, A: 255}
			featherGrey := color.RGBA{R: 65, G: 65, B: 70, A: 255}

			ebitenutil.DrawRect(screen, obs.X+6, birdY+3, 16, 8, urubuBlack)

			ebitenutil.DrawRect(screen, obs.X+1, birdY+2, 6, 7, urubuHead)
			ebitenutil.DrawRect(screen, obs.X-3, birdY+4, 5, 4, urubuBeak)
			ebitenutil.DrawRect(screen, obs.X+2, birdY+3, 2, 2, color.RGBA{R: 245, G: 245, B: 245, A: 255})

			if (ticks/9)%2 == 0 {
				ebitenutil.DrawRect(screen, obs.X+8, birdY-6, 15, 10, urubuBlack)
				ebitenutil.DrawRect(screen, obs.X+16, birdY-8, 7, 7, featherGrey)
			} else {
				ebitenutil.DrawRect(screen, obs.X+8, birdY+9, 15, 9, urubuBlack)
				ebitenutil.DrawRect(screen, obs.X+16, birdY+13, 7, 6, featherGrey)
			}

			ebitenutil.DrawRect(screen, obs.X+21, birdY+5, 7, 6, urubuBlack)
		}
	}
}
