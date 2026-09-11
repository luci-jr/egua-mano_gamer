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
	X           float64
	Y           float64
	VelocityY   float64
	IsJumping   bool
	IsCrouching bool
	IsRunning   bool
	FacingRight bool
	AimUp       bool
	IsAttacking bool
	AttackTimer int

	JumpCount   int
	JumpHolding bool
	CoyoteTimer int
	RunTicks    int

	IsSwinging   bool
	SwingingVine *Vine

	Particles []*Particle
}

func NewOnca() *Onca {
	return &Onca{
		X:            45.0,
		Y:            0,
		VelocityY:    0,
		IsJumping:    false,
		IsCrouching:  false,
		IsRunning:    false,
		FacingRight:  true,
		AimUp:        false,
		IsAttacking:  false,
		AttackTimer:  0,
		JumpCount:    0,
		JumpHolding:  false,
		CoyoteTimer:  0,
		RunTicks:     0,
		IsSwinging:   false,
		SwingingVine: nil,
		Particles:    make([]*Particle, 0),
	}
}

func (o *Onca) MoveForward(speed float64) {
	o.X += speed
	o.FacingRight = true
	o.IsRunning = true
	if o.X > 260.0 {
		o.X = 260.0
	}
}

func (o *Onca) MoveBackward(speed float64) {
	o.X -= speed
	o.FacingRight = false
	o.IsRunning = true
	if o.X < 15.0 {
		o.X = 15.0
	}
}

func (o *Onca) StopRunning() {
	o.IsRunning = false
}

func (o *Onca) SetAimUp(aim bool) {
	o.AimUp = aim
}

func (o *Onca) TriggerAttack() {
	o.IsAttacking = true
	o.AttackTimer = 14
}

func (o *Onca) GetShootOrigin(groundY float64) (x, y, vx, vy float64) {
	speed := 8.0
	if o.AimUp {
		x = o.X + 22.0
		y = groundY - 24.0 + o.Y
		vx = 0.0
		vy = -speed
		return x, y, vx, vy
	}

	dir := 1.0
	x = o.X + 36.0
	if !o.FacingRight {
		dir = -1.0
		x = o.X - 4.0
	}

	if o.IsCrouching {
		y = groundY - 8.0 + o.Y
	} else {
		y = groundY - 14.0 + o.Y
	}

	vx = dir * speed
	vy = 0.0
	return x, y, vx, vy
}

func (o *Onca) Jump() (jumped bool, isDouble bool) {
	if o.IsSwinging {
		o.ReleaseVine()
		return true, false
	}
	if !o.IsJumping || o.CoyoteTimer > 0 {
		o.IsJumping = true
		o.IsCrouching = false
		o.JumpCount = 1
		o.VelocityY = -7.0 // Salto ágil felino
		o.JumpHolding = true
		o.CoyoteTimer = 0
		o.spawnDust(o.X+10, 0, 6)
		return true, false
	} else if o.JumpCount == 1 {
		o.JumpCount = 2
		o.VelocityY = -6.4
		o.JumpHolding = true
		o.spawnJumpBurst(o.X+16, o.Y+18)
		return true, true
	}
	return false, false
}

func (o *Onca) GrabVine(v *Vine) {
	o.IsSwinging = true
	o.SwingingVine = v
	o.IsJumping = false
	o.IsCrouching = false
	o.JumpCount = 0
	o.VelocityY = 0
	o.CoyoteTimer = 0
}

func (o *Onca) ReleaseVine() bool {
	if !o.IsSwinging || o.SwingingVine == nil {
		return false
	}
	o.VelocityY = -7.0 - math.Sin(o.SwingingVine.Angle)*2.5
	o.X += 16.0
	o.IsSwinging = false
	o.IsJumping = true
	o.JumpCount = 1
	o.JumpHolding = true
	o.SwingingVine.Grabbed = false
	o.SwingingVine = nil
	o.spawnJumpBurst(o.X+16, o.Y+12)
	return true
}

func (o *Onca) ReleaseJump() {
	o.JumpHolding = false
	if o.IsJumping && o.VelocityY < -2.5 {
		o.VelocityY = -2.5
	}
}

func (o *Onca) SetCrouch(crouch bool) {
	if !o.IsCrouching && crouch && !o.IsJumping && !o.IsSwinging {
		o.spawnDust(o.X+30, 0, 4)
	}
	o.IsCrouching = crouch
	if crouch {
		o.IsRunning = false
	}
}

func (o *Onca) FastDrop() {
	if o.IsJumping && !o.IsSwinging {
		o.VelocityY += 0.9
	}
}

func (o *Onca) Update() {
	if o.IsRunning && !o.IsSwinging {
		o.RunTicks++
		if o.RunTicks%8 == 0 && !o.IsJumping {
			o.spawnDust(o.X+4, 0, 2)
		}
	} else {
		o.RunTicks = 0
	}

	if o.AttackTimer > 0 {
		o.AttackTimer--
		if o.AttackTimer == 0 {
			o.IsAttacking = false
		}
	}

	if o.IsSwinging {
		if o.SwingingVine != nil && o.SwingingVine.Active {
			tipX, tipY := o.SwingingVine.GetTipPosition()
			o.X = tipX - 16.0
			o.Y = tipY - 155.0 + 8.0
			o.VelocityY = 0
			o.IsJumping = false
			o.FacingRight = o.SwingingVine.AngleVelocity >= 0
		} else {
			o.IsSwinging = false
			o.IsJumping = true
			o.VelocityY = 1.0
		}
	} else if o.IsJumping {
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
			o.spawnDust(o.X+12, 0, 8)
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
	o.X = 45.0
	o.Y = 0
	o.VelocityY = 0
	o.IsJumping = false
	o.IsCrouching = false
	o.IsRunning = false
	o.FacingRight = true
	o.AimUp = false
	o.IsAttacking = false
	o.AttackTimer = 0
	o.JumpCount = 0
	o.JumpHolding = false
	o.CoyoteTimer = 0
	o.RunTicks = 0
	o.IsSwinging = false
	o.SwingingVine = nil
	o.Particles = o.Particles[:0]
}

func (o *Onca) GetBounds(groundY float64) (x, y, w, h float64) {
	if o.IsSwinging {
		w = 26.0
		h = 28.0
		x = o.X
		y = groundY - h + o.Y
		return x, y, w, h
	}
	w = 34.0
	h = 22.0
	if o.IsCrouching {
		w = 42.0
		h = 12.0
	}
	x = o.X
	y = groundY - h + o.Y
	return x, y, w, h
}

func (o *Onca) IsPlayerJumping() bool   { return o.IsJumping }
func (o *Onca) IsPlayerCrouching() bool { return o.IsCrouching }
func (o *Onca) IsPlayerSwinging() bool  { return o.IsSwinging }
func (o *Onca) GetJumpHolding() bool    { return o.JumpHolding }
func (o *Onca) GetJumpCount() int       { return o.JumpCount }
func (o *Onca) GetHeroKind() int        { return HeroOnca }
func (o *Onca) GetName() string         { return "ONÇA-PINTADA" }

func (o *Onca) Draw(screen *ebiten.Image, groundY float64, ticks int, invincibleTicks int) {
	for _, p := range o.Particles {
		alpha := uint8(float64(p.Color.A) * (1.0 - float64(p.Life)/float64(p.Max)))
		c := color.RGBA{R: p.Color.R, G: p.Color.G, B: p.Color.B, A: alpha}
		ebitenutil.DrawRect(screen, p.X, p.Y, p.Size, p.Size, c)
	}

	if invincibleTicks > 0 && (invincibleTicks/4)%2 != 0 {
		return
	}

	posX := o.X

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

	// 0. ESTADO: BALANÇANDO NO CIPÓ (Onça pendurada estilo Pitfall)
	if o.IsSwinging {
		oncaY := groundY - 26.0 + o.Y
		// Cipó descendo
		cVine := color.RGBA{R: 85, G: 130, B: 45, A: 255}
		ebitenutil.DrawRect(screen, posX+14, oncaY-8, 3, 14, cVine)

		// Patas dianteiras agarradas ao cipó
		ebitenutil.DrawRect(screen, posX+11, oncaY-2, 4, 6, cGold)
		ebitenutil.DrawRect(screen, posX+15, oncaY-2, 4, 6, cGold)
		ebitenutil.DrawRect(screen, posX+11, oncaY-4, 8, 3, cCream)

		// Cabeça erguida
		headX := posX + 16
		headY := oncaY + 2
		ebitenutil.DrawRect(screen, headX, headY, 10, 9, cGold)
		ebitenutil.DrawRect(screen, headX+1, headY+1, 8, 6, cGoldLight)
		ebitenutil.DrawRect(screen, headX+5, headY+5, 5, 4, cCream)
		ebitenutil.DrawRect(screen, headX+8, headY+4, 2, 2, cNose)
		ebitenutil.DrawRect(screen, headX+5, headY+2, 2, 2, cEye)
		ebitenutil.DrawRect(screen, headX+2, headY-2, 3, 3, cEarPink)

		// Corpo inclinado pendurado verticalmente
		ebitenutil.DrawRect(screen, posX+8, oncaY+6, 14, 16, cGold)
		ebitenutil.DrawRect(screen, posX+9, oncaY+7, 12, 14, cGoldLight)
		ebitenutil.DrawRect(screen, posX+10, oncaY+10, 8, 10, cCream)
		drawRosette(posX+10, oncaY+8)
		drawRosette(posX+12, oncaY+14)

		// Patas traseiras abraçando o cipó
		ebitenutil.DrawRect(screen, posX+10, oncaY+20, 5, 5, cGold)
		ebitenutil.DrawRect(screen, posX+15, oncaY+20, 5, 5, cGold)

		// Rabo curvado no ar
		tailWave := math.Sin(float64(ticks)*0.2) * 2.0
		ebitenutil.DrawRect(screen, posX+4+tailWave, oncaY+18, 5, 3, cGold)
		ebitenutil.DrawRect(screen, posX+2+tailWave, oncaY+20, 4, 4, cSpotBlack)
		return
	}

	if o.IsCrouching {
		oncaY := groundY - 12.0 + o.Y

		ebitenutil.DrawRect(screen, posX, oncaY+2, 38, 9, cGold)
		ebitenutil.DrawRect(screen, posX+4, oncaY+3, 30, 6, cGoldLight)
		ebitenutil.DrawRect(screen, posX+6, oncaY+8, 26, 3, cCream)

		drawRosette(posX+8, oncaY+3)
		drawRosette(posX+18, oncaY+3)
		drawRosette(posX+28, oncaY+4)

		headX := posX + 34
		headY := oncaY + 1
		ebitenutil.DrawRect(screen, headX, headY, 9, 8, cGold)
		ebitenutil.DrawRect(screen, headX+3, headY+5, 6, 4, cCream)
		ebitenutil.DrawRect(screen, headX+9, headY+4, 2, 3, cNose)
		ebitenutil.DrawRect(screen, headX+5, headY+2, 3, 3, cEye)
		ebitenutil.DrawRect(screen, headX+6, headY+2, 1, 3, color.RGBA{R: 10, G: 10, B: 10, A: 255})
		ebitenutil.DrawRect(screen, headX-2, headY-1, 4, 3, cSpotBlack)
		ebitenutil.DrawRect(screen, headX-1, headY, 2, 2, cEarPink)

		step := float64((ticks / 4) % 2)
		ebitenutil.DrawRect(screen, posX+4+step*2, oncaY+9, 7, 3, cGold)
		ebitenutil.DrawRect(screen, posX+3+step*2, oncaY+11, 8, 2, cCream)
		ebitenutil.DrawRect(screen, posX+30-step*2, oncaY+9, 7, 3, cGold)
		ebitenutil.DrawRect(screen, posX+31-step*2, oncaY+11, 8, 2, cCream)

		ebitenutil.DrawRect(screen, posX-12, oncaY+5, 13, 3, cGold)
		ebitenutil.DrawRect(screen, posX-15, oncaY+4, 4, 3, cSpotBlack)

		if o.IsAttacking {
			ebitenutil.DrawRect(screen, headX+10, headY+4, 2, 4, color.RGBA{R: 200, G: 50, B: 50, A: 255})
			ebitenutil.DrawRect(screen, headX+12, headY+3, 2, 6, color.RGBA{R: 255, G: 200, B: 50, A: 200})
		}
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

	ebitenutil.DrawRect(screen, posX+3, oncaY+4+bodyYOffset, 27, 12, cGold)
	ebitenutil.DrawRect(screen, posX+5, oncaY+5+bodyYOffset, 23, 8, cGoldLight)
	ebitenutil.DrawRect(screen, posX+6, oncaY+13+bodyYOffset, 20, 4, cCream)

	drawRosette(posX+6, oncaY+6+bodyYOffset)
	drawRosette(posX+14, oncaY+7+bodyYOffset)
	drawRosette(posX+21, oncaY+6+bodyYOffset)
	ebitenutil.DrawRect(screen, posX+10, oncaY+12+bodyYOffset, 3, 2, cSpotBlack)
	ebitenutil.DrawRect(screen, posX+18, oncaY+12+bodyYOffset, 3, 2, cSpotBlack)

	headX := posX + 26
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

	if o.IsAttacking {
		// Boca rugindo e onda sônica
		ebitenutil.DrawRect(screen, headX+10, headY+5, 3, 5, color.RGBA{R: 180, G: 30, B: 30, A: 255})
		ebitenutil.DrawRect(screen, headX+11, headY+6, 1, 3, color.RGBA{R: 250, G: 240, B: 240, A: 255}) // Caninos
		ebitenutil.DrawRect(screen, headX+13, headY+4, 2, 7, color.RGBA{R: 255, G: 215, B: 50, A: 220})
		ebitenutil.DrawRect(screen, headX+16, headY+2, 2, 11, color.RGBA{R: 255, G: 160, B: 30, A: 170})
	}

	tailWave := math.Sin(float64(ticks)*0.25) * 1.8
	if o.IsJumping {
		tailWave = -2.5
	}
	ebitenutil.DrawRect(screen, posX-1, oncaY+9+bodyYOffset, 5, 4, cGold)
	ebitenutil.DrawRect(screen, posX-5, oncaY+5+bodyYOffset, 5, 5, cGold)
	ebitenutil.DrawRect(screen, posX-7, oncaY+tailWave+bodyYOffset, 4, 6, cGold)
	ebitenutil.DrawRect(screen, posX-5, oncaY-3+tailWave+bodyYOffset, 4, 5, cSpotBlack)
	ebitenutil.DrawRect(screen, posX-3, oncaY-4+tailWave+bodyYOffset, 3, 3, cSpotBlack)

	switch gallopFrame {
	case 0:
		ebitenutil.DrawRect(screen, posX-2, oncaY+14, 6, 6, cGold)
		ebitenutil.DrawRect(screen, posX-5, oncaY+18, 5, 5, cGold)
		ebitenutil.DrawRect(screen, posX-7, oncaY+21, 5, 2, cCream)

		ebitenutil.DrawRect(screen, posX+24, oncaY+13, 6, 6, cGold)
		ebitenutil.DrawRect(screen, posX+28, oncaY+17, 5, 5, cGold)
		ebitenutil.DrawRect(screen, posX+30, oncaY+21, 5, 2, cCream)

	case 1:
		ebitenutil.DrawRect(screen, posX+2, oncaY+15, 6, 5, cGold)
		ebitenutil.DrawRect(screen, posX+1, oncaY+18, 5, 4, cGold)
		ebitenutil.DrawRect(screen, posX, oncaY+21, 5, 2, cCream)

		ebitenutil.DrawRect(screen, posX+22, oncaY+14, 6, 5, cGold)
		ebitenutil.DrawRect(screen, posX+25, oncaY+17, 5, 5, cGold)
		ebitenutil.DrawRect(screen, posX+26, oncaY+21, 5, 2, cCream)

	case 2:
		ebitenutil.DrawRect(screen, posX+8, oncaY+14, 6, 5, cGold)
		ebitenutil.DrawRect(screen, posX+10, oncaY+17, 5, 5, cGold)
		ebitenutil.DrawRect(screen, posX+11, oncaY+21, 5, 2, cCream)

		ebitenutil.DrawRect(screen, posX+17, oncaY+14, 6, 5, cGold)
		ebitenutil.DrawRect(screen, posX+18, oncaY+17, 5, 5, cGold)
		ebitenutil.DrawRect(screen, posX+19, oncaY+21, 5, 2, cCream)

	case 3:
		ebitenutil.DrawRect(screen, posX+5, oncaY+14, 6, 5, cGold)
		ebitenutil.DrawRect(screen, posX+3, oncaY+18, 5, 4, cGold)
		ebitenutil.DrawRect(screen, posX+2, oncaY+21, 5, 2, cCream)

		ebitenutil.DrawRect(screen, posX+20, oncaY+13, 6, 6, cGold)
		ebitenutil.DrawRect(screen, posX+23, oncaY+17, 5, 5, cGold)
		ebitenutil.DrawRect(screen, posX+24, oncaY+21, 5, 2, cCream)
	}
}
