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
	X            float64
	Y            float64
	GroundOffset float64
	VelocityY    float64
	IsJumping    bool
	IsCrouching  bool
	IsRunning    bool
	FacingRight  bool
	AimUp        bool
	IsAttacking  bool
	AttackTimer  int

	JumpCount   int
	JumpHolding bool
	CoyoteTimer int
	RunTicks    int

	IsSwinging   bool
	SwingingVine *Vine

	Particles []*Particle

	StarPowerTimer int
}

func NewOnca() *Onca {
	return &Onca{
		X:              45.0,
		Y:              0,
		GroundOffset:   0,
		VelocityY:      0,
		IsJumping:      false,
		IsCrouching:    false,
		IsRunning:      false,
		FacingRight:    true,
		AimUp:          false,
		IsAttacking:    false,
		AttackTimer:    0,
		JumpCount:      0,
		JumpHolding:    false,
		CoyoteTimer:    0,
		RunTicks:       0,
		IsSwinging:     false,
		SwingingVine:   nil,
		Particles:      make([]*Particle, 0),
		StarPowerTimer: 0,
	}
}

func (o *Onca) SetPositionX(x float64) {
	o.X = x
}

func (o *Onca) GetPosition() (x, y float64) {
	return o.X, o.Y
}

func (o *Onca) SetGroundOffset(offset float64) {
	o.GroundOffset = offset
	o.Y = offset
	o.VelocityY = 0
	o.IsJumping = false
	o.JumpCount = 0
}

func (o *Onca) GetGroundOffset() float64 {
	return o.GroundOffset
}

func (o *Onca) GetVelocityY() float64 {
	return o.VelocityY
}

func (o *Onca) FallFromPlatform() {
	if o.GroundOffset != 0 {
		o.GroundOffset = 0
		if o.Y < 0 {
			o.IsJumping = true
			if o.VelocityY < 0 {
				o.VelocityY = 0
			}
		}
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
		o.GroundOffset = 0 // Pulo liberta da plataforma atual
		o.IsJumping = true
		o.IsCrouching = false
		o.JumpCount = 1
		o.VelocityY = -7.0 // Salto ágil felino
		o.JumpHolding = true
		o.CoyoteTimer = 0
		o.spawnDust(o.X+10, o.Y, 6)
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
		o.spawnDust(o.X+30, o.Y, 4)
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
			o.spawnDust(o.X+4, o.Y, 2)
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
			gFacing := o.SwingingVine.AngleVelocity >= 0
			o.FacingRight = gFacing
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

		targetY := o.GroundOffset
		if o.Y >= targetY {
			o.Y = targetY
			o.IsJumping = false
			o.JumpCount = 0
			o.VelocityY = 0
			o.spawnDust(o.X+12, targetY, 8)
		}
	} else {
		if o.Y < o.GroundOffset {
			o.IsJumping = true
		}
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
	o.GroundOffset = 0
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
	o.StarPowerTimer = 0
}

func (o *Onca) Bounce(strength float64) {
	o.GroundOffset = 0
	o.IsJumping = true
	o.IsCrouching = false
	o.JumpCount = 1
	if strength == 0 {
		strength = -6.4
	}
	o.VelocityY = strength
	o.JumpHolding = true
	o.spawnJumpBurst(o.X+16, o.Y+18)
}

func (o *Onca) SetStarPower(timer int) {
	o.StarPowerTimer = timer
}

func (o *Onca) GetStarPower() int {
	return o.StarPowerTimer
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

	// Paleta cromática da Onça-Pintada
	cGold := color.RGBA{R: 235, G: 160, B: 35, A: 255}
	cGoldLight := color.RGBA{R: 248, G: 185, B: 65, A: 255}
	cCream := color.RGBA{R: 245, G: 230, B: 195, A: 255}
	cSpotBlack := color.RGBA{R: 35, G: 18, B: 10, A: 255}
	cSpotCenter := color.RGBA{R: 190, G: 95, B: 20, A: 255}
	cEye := color.RGBA{R: 85, G: 205, B: 110, A: 255}
	cEarPink := color.RGBA{R: 215, G: 120, B: 120, A: 255}
	cNose := color.RGBA{R: 70, G: 30, B: 30, A: 255}
	cMouthDark := color.RGBA{R: 140, G: 25, B: 25, A: 255}
	cTongue := color.RGBA{R: 220, G: 50, B: 60, A: 255}
	cTooth := color.RGBA{R: 255, G: 255, B: 250, A: 255}
	cRoarWave1 := color.RGBA{R: 255, G: 255, B: 255, A: 220} // Onda de pressão de ar sônica
	cRoarWave2 := color.RGBA{R: 235, G: 110, B: 40, A: 180}  // Laranja terracota
	cRoarWave3 := color.RGBA{R: 200, G: 60, B: 30, A: 140}   // Carmesim rugido

	// Efeito Star Power (Guaraná da Amazônia): Onça mística reluzente em arco-íris estilo Super Mario
	if o.StarPowerTimer > 0 {
		rainbowColors := []color.RGBA{
			{R: 255, G: 235, B: 55, A: 255},  // Ouro radiante
			{R: 55, G: 245, B: 240, A: 255},  // Ciano místico
			{R: 255, G: 85, B: 225, A: 255},  // Rosa cósmico
			{R: 255, G: 255, B: 255, A: 255}, // Branco celestial
			{R: 80, G: 255, B: 110, A: 255},  // Verde esmeralda
		}
		starCol := rainbowColors[(ticks/4)%len(rainbowColors)]
		starColLight := rainbowColors[(ticks/4+1)%len(rainbowColors)]
		cGold = starCol
		cGoldLight = starColLight
		cSpotCenter = rainbowColors[(ticks/4+2)%len(rainbowColors)]

		// Aura luminosa ao redor da Onça
		auraAlpha := uint8(110 + (ticks%5)*25)
		auraCol := color.RGBA{R: starCol.R, G: starCol.G, B: starCol.B, A: auraAlpha}
		baseY := groundY - 24.0 + o.Y
		ebitenutil.DrawRect(screen, posX-2, baseY-2, 38, 28, auraCol)
	}

	// Helper com espelhamento horizontal relativo ao centro da onça
	centerX := posX + 19.0
	drawBox := func(relX, relY, w, h float64, oncaBaseY float64, col color.Color) {
		var finalX float64
		if o.FacingRight {
			finalX = centerX - 19.0 + relX
		} else {
			finalX = centerX + 19.0 - relX - w
		}
		ebitenutil.DrawRect(screen, finalX, oncaBaseY+relY, w, h, col)
	}

	drawRosette := func(rx, ry float64, oncaBaseY float64) {
		drawBox(rx, ry, 5, 4, oncaBaseY, cSpotBlack)
		drawBox(rx+1, ry+1, 3, 2, oncaBaseY, cSpotCenter)
	}

	// 1. ESTADO: AGACHADA (Crouching)
	if o.IsCrouching {
		oncaY := groundY - 12.0 + o.Y

		drawBox(0, 2, 38, 9, oncaY, cGold)
		drawBox(4, 3, 30, 6, oncaY, cGoldLight)
		drawBox(6, 8, 26, 3, oncaY, cCream)

		drawRosette(8, 3, oncaY)
		drawRosette(18, 3, oncaY)
		drawRosette(28, 4, oncaY)

		headRelX := 32.0
		headRelY := 1.0
		drawBox(headRelX, headRelY, 9, 8, oncaY, cGold)
		drawBox(headRelX+3, headRelY+5, 6, 4, oncaY, cCream)
		drawBox(headRelX+5, headRelY+2, 3, 3, oncaY, cEye)
		drawBox(headRelX+6, headRelY+2, 1, 3, oncaY, color.RGBA{R: 10, G: 10, B: 10, A: 255})
		drawBox(headRelX-2, headRelY-1, 4, 3, oncaY, cSpotBlack)
		drawBox(headRelX-1, headRelY, 2, 2, oncaY, cEarPink)

		// Ataque agachada
		if o.IsAttacking {
			drawBox(headRelX+8, headRelY+4, 4, 4, oncaY, cMouthDark)
			drawBox(headRelX+9, headRelY+6, 2, 2, oncaY, cTongue)
			drawBox(headRelX+10, headRelY+3, 1, 3, oncaY, cTooth)
			drawBox(headRelX+12, headRelY+2, 2, 6, oncaY, cRoarWave1)
			drawBox(headRelX+15, headRelY+1, 2, 8, oncaY, cRoarWave2)
		} else {
			drawBox(headRelX+8, headRelY+4, 2, 3, oncaY, cNose)
		}

		// Patas agachadas
		step := float64((ticks / 4) % 2)
		drawBox(4+step*2, 9, 7, 3, oncaY, cGold)
		drawBox(3+step*2, 11, 8, 2, oncaY, cCream)
		drawBox(28-step*2, 9, 7, 3, oncaY, cGold)
		drawBox(29-step*2, 11, 8, 2, oncaY, cCream)

		// Rabo abaixado
		drawBox(-10, 5, 11, 3, oncaY, cGold)
		drawBox(-13, 4, 4, 3, oncaY, cSpotBlack)
		return
	}

	// 2. ESTADO: EM PÉ / CORRENDO / SALTANDO
	oncaY := groundY - 22.0 + o.Y
	gallopFrame := (ticks / 5) % 4
	if o.IsJumping {
		gallopFrame = 0
	}

	bodyYOffset := 0.0
	if gallopFrame == 1 || gallopFrame == 3 {
		bodyYOffset = 1.0
	}

	// Tronco da Onça
	drawBox(3, 4+bodyYOffset, 27, 12, oncaY, cGold)
	drawBox(5, 5+bodyYOffset, 23, 8, oncaY, cGoldLight)
	drawBox(6, 13+bodyYOffset, 20, 4, oncaY, cCream)

	drawRosette(6, 6+bodyYOffset, oncaY)
	drawRosette(14, 7+bodyYOffset, oncaY)
	drawRosette(21, 6+bodyYOffset, oncaY)
	drawBox(10, 12+bodyYOffset, 3, 2, oncaY, cSpotBlack)
	drawBox(18, 12+bodyYOffset, 3, 2, oncaY, cSpotBlack)

	// Cabeça da Onça
	headRelX := 26.0
	headRelY := 2.0 + bodyYOffset
	if o.AimUp {
		headRelY = -3.0 + bodyYOffset
	}

	drawBox(headRelX, headRelY, 10, 10, oncaY, cGold)
	drawBox(headRelX+1, headRelY+1, 8, 7, oncaY, cGoldLight)
	drawBox(headRelX+4, headRelY+6, 6, 5, oncaY, cCream)

	// Orelha
	drawBox(headRelX+1, headRelY-3, 4, 4, oncaY, cSpotBlack)
	drawBox(headRelX+2, headRelY-2, 2, 3, oncaY, cEarPink)

	// Olhos (semicerrados em fúria se atacando)
	if o.IsAttacking {
		drawBox(headRelX+5, headRelY+3, 3, 2, oncaY, cEye)
		drawBox(headRelX+6, headRelY+3, 1, 2, oncaY, color.RGBA{R: 10, G: 10, B: 10, A: 255})
	} else {
		drawBox(headRelX+5, headRelY+3, 3, 3, oncaY, cEye)
		drawBox(headRelX+6, headRelY+3, 1, 3, oncaY, color.RGBA{R: 15, G: 15, B: 15, A: 255})
		drawBox(headRelX+5, headRelY+2, 1, 1, oncaY, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	}

	// Animação de Focinho e Boca no Ataque
	if o.IsAttacking {
		if o.AimUp {
			// Rugido para cima (Anti-aéreo)
			drawBox(headRelX+3, headRelY-4, 5, 4, oncaY, cMouthDark)
			drawBox(headRelX+4, headRelY-3, 3, 2, oncaY, cTongue)
			drawBox(headRelX+3, headRelY-4, 1, 2, oncaY, cTooth)
			drawBox(headRelX+7, headRelY-4, 1, 2, oncaY, cTooth)
			// Ondas sônicas verticais subindo
			drawBox(headRelX+1, headRelY-7, 8, 2, oncaY, cRoarWave1)
			drawBox(headRelX-1, headRelY-11, 12, 2, oncaY, cRoarWave2)
			drawBox(headRelX-3, headRelY-15, 16, 2, oncaY, cRoarWave3)
		} else {
			// Rugido horizontal com mandíbula escancarada e bote
			drawBox(headRelX+8, headRelY+3, 4, 3, oncaY, cGold)
			drawBox(headRelX+10, headRelY+2, 2, 2, oncaY, cNose)
			// Interior da boca
			drawBox(headRelX+8, headRelY+5, 5, 6, oncaY, cMouthDark)
			drawBox(headRelX+9, headRelY+8, 3, 2, oncaY, cTongue)
			// Caninos afiados
			drawBox(headRelX+11, headRelY+5, 1, 3, oncaY, cTooth)
			drawBox(headRelX+11, headRelY+8, 1, 2, oncaY, cTooth)
			// Mandíbula inferior abaixada
			drawBox(headRelX+8, headRelY+10, 4, 2, oncaY, cGold)

			// Bote da pata dianteira com garras afiadas estendidas
			drawBox(headRelX+2, headRelY+10, 7, 4, oncaY, cGold)
			drawBox(headRelX+8, headRelY+10, 2, 1, oncaY, cTooth) // Garra 1
			drawBox(headRelX+8, headRelY+12, 2, 1, oncaY, cTooth) // Garra 2

			// Ondas sônicas concêntricas do Rugido
			drawBox(headRelX+14, headRelY+3, 2, 8, oncaY, cRoarWave1)
			drawBox(headRelX+17, headRelY+1, 2, 12, oncaY, cRoarWave2)
			drawBox(headRelX+21, headRelY-1, 2, 16, oncaY, cRoarWave3)
		}
	} else {
		// Focinho relaxado
		drawBox(headRelX+9, headRelY+5, 2, 3, oncaY, cNose)
	}

	// Cauda (Rabo balançando)
	tailWave := math.Sin(float64(ticks)*0.25) * 1.8
	if o.IsJumping {
		tailWave = -2.5
	}
	drawBox(-1, 9+bodyYOffset, 5, 4, oncaY, cGold)
	drawBox(-5, 5+bodyYOffset, 5, 5, oncaY, cGold)
	drawBox(-7, tailWave+bodyYOffset, 4, 6, oncaY, cGold)
	drawBox(-5, -3+tailWave+bodyYOffset, 4, 5, oncaY, cSpotBlack)
	drawBox(-3, -4+tailWave+bodyYOffset, 3, 3, oncaY, cSpotBlack)

	// Patas animadas com ciclo de galope 4-frames
	switch gallopFrame {
	case 0:
		drawBox(-2, 14, 6, 6, oncaY, cGold)
		drawBox(-5, 18, 5, 5, oncaY, cGold)
		drawBox(-7, 21, 5, 2, oncaY, cCream)

		drawBox(24, 13, 6, 6, oncaY, cGold)
		drawBox(28, 17, 5, 5, oncaY, cGold)
		drawBox(30, 21, 5, 2, oncaY, cCream)

	case 1:
		drawBox(2, 15, 6, 5, oncaY, cGold)
		drawBox(1, 18, 5, 4, oncaY, cGold)
		drawBox(0, 21, 5, 2, oncaY, cCream)

		drawBox(22, 14, 6, 5, oncaY, cGold)
		drawBox(25, 17, 5, 5, oncaY, cGold)
		drawBox(26, 21, 5, 2, oncaY, cCream)

	case 2:
		drawBox(8, 14, 6, 5, oncaY, cGold)
		drawBox(10, 17, 5, 5, oncaY, cGold)
		drawBox(11, 21, 5, 2, oncaY, cCream)

		drawBox(17, 14, 6, 5, oncaY, cGold)
		drawBox(18, 17, 5, 5, oncaY, cGold)
		drawBox(19, 21, 5, 2, oncaY, cCream)

	case 3:
		drawBox(5, 14, 6, 5, oncaY, cGold)
		drawBox(3, 18, 5, 4, oncaY, cGold)
		drawBox(2, 21, 5, 2, oncaY, cCream)

		drawBox(20, 13, 6, 6, oncaY, cGold)
		drawBox(23, 17, 5, 5, oncaY, cGold)
		drawBox(24, 21, 5, 2, oncaY, cCream)
	}
}
