package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	OncaPosX = 45.0
)

type Particle struct {
	X, Y       float64
	VX, VY     float64
	Life, Max  int
	Size       float64
	Color      color.RGBA
}

type Onca struct {
	Y           float64
	VelocityY   float64
	IsJumping   bool
	IsCrouching bool
	JumpCount   int
	JumpHolding bool
	CoyoteTimer int
	Particles   []*Particle
}

func NewOnca() *Onca {
	return &Onca{
		Y:           0,
		VelocityY:   0,
		IsJumping:   false,
		IsCrouching: false,
		JumpCount:   0,
		JumpHolding: false,
		CoyoteTimer: 0,
		Particles:   make([]*Particle, 0),
	}
}

func (o *Onca) Jump() (jumped bool, isDouble bool) {
	if !o.IsJumping || o.CoyoteTimer > 0 {
		o.IsJumping = true
		o.IsCrouching = false
		o.JumpCount = 1
		o.VelocityY = -6.8
		o.JumpHolding = true
		o.CoyoteTimer = 0
		o.spawnDust(OncaPosX+10, 0, 6)
		return true, false
	} else if o.JumpCount == 1 {
		o.JumpCount = 2
		o.VelocityY = -6.2
		o.JumpHolding = true
		o.spawnJumpBurst(OncaPosX+16, o.Y+18)
		return true, true
	}
	return false, false
}

func (o *Onca) ReleaseJump() {
	o.JumpHolding = false
	if o.IsJumping && o.VelocityY < -2.5 {
		o.VelocityY = -2.5
	}
}

func (o *Onca) SetCrouch(crouch bool) {
	if !o.IsCrouching && crouch && !o.IsJumping {
		o.spawnDust(OncaPosX+30, 0, 4)
	}
	o.IsCrouching = crouch
}

func (o *Onca) FastDrop() {
	if o.IsJumping {
		o.VelocityY += 0.9
	}
}

func (o *Onca) Update() {
	if o.IsJumping {
		gravity := 0.36
		if !o.JumpHolding && o.VelocityY < 0 {
			gravity = 0.55
		}
		o.Y += o.VelocityY
		o.VelocityY += gravity

		if o.Y >= 0 {
			o.Y = 0
			o.IsJumping = false
			o.JumpCount = 0
			o.VelocityY = 0
			o.spawnDust(OncaPosX+12, 0, 8)
		}
	} else {
		o.CoyoteTimer = 6
	}

	if o.CoyoteTimer > 0 {
		o.CoyoteTimer--
	}

	alive := o.Particles[:0]
	for _, p := range o.Particles {
		p.X += p.VX
		p.Y += p.VY
		p.Life++
		if p.Life < p.Max {
			alive = append(alive, p)
		}
	}
	o.Particles = alive
}

func (o *Onca) spawnDust(x, yOffset float64, count int) {
	for i := 0; i < count; i++ {
		p := &Particle{
			X:     x + float64(i*3),
			Y:     155.0 + yOffset,
			VX:    -1.5 - float64(i)*0.4,
			VY:    -0.4 - float64(i%3)*0.3,
			Life:  0,
			Max:   14 + i*2,
			Size:  2.0 + float64(i%2),
			Color: color.RGBA{R: 190, G: 180, B: 170, A: 200},
		}
		o.Particles = append(o.Particles, p)
	}
}

func (o *Onca) spawnJumpBurst(x, y float64) {
	for i := 0; i < 6; i++ {
		angle := float64(i) * math.Pi / 3.0
		p := &Particle{
			X:     x,
			Y:     155.0 - 24.0 + y,
			VX:    math.Cos(angle) * 1.8,
			VY:    math.Sin(angle) * 1.2,
			Life:  0,
			Max:   12,
			Size:  2.5,
			Color: color.RGBA{R: 255, G: 235, B: 170, A: 230},
		}
		o.Particles = append(o.Particles, p)
	}
}

func (o *Onca) Reset() {
	o.Y = 0
	o.VelocityY = 0
	o.IsJumping = false
	o.IsCrouching = false
	o.JumpCount = 0
	o.JumpHolding = false
	o.CoyoteTimer = 0
	o.Particles = o.Particles[:0]
}

func (o *Onca) GetBounds(groundY float64) (x, y, w, h float64) {
	w = 34.0
	h = 22.0
	if o.IsCrouching {
		w = 42.0
		h = 12.0
	}
	x = OncaPosX
	y = groundY - h + o.Y
	return x, y, w, h
}

func (o *Onca) Draw(screen *ebiten.Image, groundY float64, ticks int, invincibleTicks int) {
	for _, p := range o.Particles {
		alpha := uint8(float64(p.Color.A) * (1.0 - float64(p.Life)/float64(p.Max)))
		c := color.RGBA{R: p.Color.R, G: p.Color.G, B: p.Color.B, A: alpha}
		ebitenutil.DrawRect(screen, p.X, p.Y, p.Size, p.Size, c)
	}

	if invincibleTicks > 0 && (invincibleTicks/4)%2 != 0 {
		return
	}

	cGold := color.RGBA{R: 235, G: 160, B: 35, A: 255}
	cGoldLight := color.RGBA{R: 248, G: 185, B: 65, A: 255}
	cCream := color.RGBA{R: 245, G: 230, B: 195, A: 255}
	cSpotBlack := color.RGBA{R: 35, G: 18, B: 10, A: 255}
	cSpotCenter := color.RGBA{R: 190, G: 95, B: 20, A: 255}
	cEye := color.RGBA{R: 85, G: 205, B: 110, A: 255}
	cEarPink := color.RGBA{R: 215, G: 120, B: 120, A: 255}
	cNose := color.RGBA{R: 70, G: 30, B: 30, A: 255}

	drawRosette := func(rx, ry float64) {
		ebitenutil.DrawRect(screen, rx, ry, 5, 4, cSpotBlack)
		ebitenutil.DrawRect(screen, rx+1, ry+1, 3, 2, cSpotCenter)
	}

	if o.IsCrouching {
		oncaY := groundY - 12.0 + o.Y

		ebitenutil.DrawRect(screen, OncaPosX, oncaY+2, 38, 9, cGold)
		ebitenutil.DrawRect(screen, OncaPosX+4, oncaY+3, 30, 6, cGoldLight)
		ebitenutil.DrawRect(screen, OncaPosX+6, oncaY+8, 26, 3, cCream)

		drawRosette(OncaPosX+8, oncaY+3)
		drawRosette(OncaPosX+18, oncaY+3)
		drawRosette(OncaPosX+28, oncaY+4)

		headX := OncaPosX + 34
		headY := oncaY + 1
		ebitenutil.DrawRect(screen, headX, headY, 9, 8, cGold)
		ebitenutil.DrawRect(screen, headX+3, headY+5, 6, 4, cCream)
		ebitenutil.DrawRect(screen, headX+9, headY+4, 2, 3, cNose)
		ebitenutil.DrawRect(screen, headX+5, headY+2, 3, 3, cEye)
		ebitenutil.DrawRect(screen, headX+6, headY+2, 1, 3, color.RGBA{R: 10, G: 10, B: 10, A: 255})
		ebitenutil.DrawRect(screen, headX-2, headY-1, 4, 3, cSpotBlack)
		ebitenutil.DrawRect(screen, headX-1, headY, 2, 2, cEarPink)

		step := float64((ticks / 4) % 2)
		ebitenutil.DrawRect(screen, OncaPosX+4+step*2, oncaY+9, 7, 3, cGold)
		ebitenutil.DrawRect(screen, OncaPosX+3+step*2, oncaY+11, 8, 2, cCream)
		ebitenutil.DrawRect(screen, OncaPosX+30-step*2, oncaY+9, 7, 3, cGold)
		ebitenutil.DrawRect(screen, OncaPosX+31-step*2, oncaY+11, 8, 2, cCream)

		ebitenutil.DrawRect(screen, OncaPosX-12, oncaY+5, 13, 3, cGold)
		ebitenutil.DrawRect(screen, OncaPosX-15, oncaY+4, 4, 3, cSpotBlack)
		return
	}

	oncaY := groundY - 22.0 + o.Y
	gallopFrame := (ticks / 5) % 4
	if o.IsJumping {
		gallopFrame = 0
	}

	bodyYOffset := 0.0
	if gallopFrame == 1 || gallopFrame == 3 {
		bodyYOffset = 1.0
	}

	ebitenutil.DrawRect(screen, OncaPosX+3, oncaY+4+bodyYOffset, 27, 12, cGold)
	ebitenutil.DrawRect(screen, OncaPosX+5, oncaY+5+bodyYOffset, 23, 8, cGoldLight)
	ebitenutil.DrawRect(screen, OncaPosX+6, oncaY+13+bodyYOffset, 20, 4, cCream)

	drawRosette(OncaPosX+6, oncaY+6+bodyYOffset)
	drawRosette(OncaPosX+14, oncaY+7+bodyYOffset)
	drawRosette(OncaPosX+21, oncaY+6+bodyYOffset)
	ebitenutil.DrawRect(screen, OncaPosX+10, oncaY+12+bodyYOffset, 3, 2, cSpotBlack)
	ebitenutil.DrawRect(screen, OncaPosX+18, oncaY+12+bodyYOffset, 3, 2, cSpotBlack)

	headX := OncaPosX + 26
	headY := oncaY + 2 + bodyYOffset
	ebitenutil.DrawRect(screen, headX, headY, 10, 10, cGold)
	ebitenutil.DrawRect(screen, headX+1, headY+1, 8, 7, cGoldLight)
	ebitenutil.DrawRect(screen, headX+4, headY+6, 6, 5, cCream)
	ebitenutil.DrawRect(screen, headX+9, headY+5, 2, 3, cNose)

	ebitenutil.DrawRect(screen, headX+5, headY+3, 3, 3, cEye)
	ebitenutil.DrawRect(screen, headX+6, headY+3, 1, 3, color.RGBA{R: 15, G: 15, B: 15, A: 255})
	ebitenutil.DrawRect(screen, headX+5, headY+2, 1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255})

	ebitenutil.DrawRect(screen, headX+1, headY-3, 4, 4, cSpotBlack)
	ebitenutil.DrawRect(screen, headX+2, headY-2, 2, 3, cEarPink)

	tailWave := math.Sin(float64(ticks)*0.25) * 1.8
	if o.IsJumping {
		tailWave = -2.5
	}
	ebitenutil.DrawRect(screen, OncaPosX-1, oncaY+9+bodyYOffset, 5, 4, cGold)
	ebitenutil.DrawRect(screen, OncaPosX-5, oncaY+5+bodyYOffset, 5, 5, cGold)
	ebitenutil.DrawRect(screen, OncaPosX-7, oncaY+tailWave+bodyYOffset, 4, 6, cGold)
	ebitenutil.DrawRect(screen, OncaPosX-5, oncaY-3+tailWave+bodyYOffset, 4, 5, cSpotBlack)
	ebitenutil.DrawRect(screen, OncaPosX-3, oncaY-4+tailWave+bodyYOffset, 3, 3, cSpotBlack)

	switch gallopFrame {
	case 0:
		ebitenutil.DrawRect(screen, OncaPosX-2, oncaY+14, 6, 6, cGold)
		ebitenutil.DrawRect(screen, OncaPosX-5, oncaY+18, 5, 5, cGold)
		ebitenutil.DrawRect(screen, OncaPosX-7, oncaY+21, 5, 2, cCream)

		ebitenutil.DrawRect(screen, OncaPosX+24, oncaY+13, 6, 6, cGold)
		ebitenutil.DrawRect(screen, OncaPosX+28, oncaY+17, 5, 5, cGold)
		ebitenutil.DrawRect(screen, OncaPosX+30, oncaY+21, 5, 2, cCream)

	case 1:
		ebitenutil.DrawRect(screen, OncaPosX+2, oncaY+15, 6, 5, cGold)
		ebitenutil.DrawRect(screen, OncaPosX+1, oncaY+18, 5, 4, cGold)
		ebitenutil.DrawRect(screen, OncaPosX, oncaY+21, 5, 2, cCream)

		ebitenutil.DrawRect(screen, OncaPosX+22, oncaY+14, 6, 5, cGold)
		ebitenutil.DrawRect(screen, OncaPosX+25, oncaY+17, 5, 5, cGold)
		ebitenutil.DrawRect(screen, OncaPosX+26, oncaY+21, 5, 2, cCream)

	case 2:
		ebitenutil.DrawRect(screen, OncaPosX+8, oncaY+14, 6, 5, cGold)
		ebitenutil.DrawRect(screen, OncaPosX+10, oncaY+17, 5, 5, cGold)
		ebitenutil.DrawRect(screen, OncaPosX+11, oncaY+21, 5, 2, cCream)

		ebitenutil.DrawRect(screen, OncaPosX+17, oncaY+14, 6, 5, cGold)
		ebitenutil.DrawRect(screen, OncaPosX+18, oncaY+17, 5, 5, cGold)
		ebitenutil.DrawRect(screen, OncaPosX+19, oncaY+21, 5, 2, cCream)

	case 3:
		ebitenutil.DrawRect(screen, OncaPosX+5, oncaY+14, 6, 5, cGold)
		ebitenutil.DrawRect(screen, OncaPosX+3, oncaY+18, 5, 4, cGold)
		ebitenutil.DrawRect(screen, OncaPosX+2, oncaY+21, 5, 2, cCream)

		ebitenutil.DrawRect(screen, OncaPosX+20, oncaY+13, 6, 6, cGold)
		ebitenutil.DrawRect(screen, OncaPosX+23, oncaY+17, 5, 5, cGold)
		ebitenutil.DrawRect(screen, OncaPosX+24, oncaY+21, 5, 2, cCream)
	}
}
