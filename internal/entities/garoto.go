package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GarotoState int

const (
	StateIdle GarotoState = iota
	StateRunning
	StateJumping
	StateFalling
	StateCrouching
)

type Garoto struct {
	X            float64
	Y            float64 // Offset vertical do pulo (Y <= 0 quando no ar ou em plataforma)
	GroundOffset float64 // Offset de chão sólido (0 para piso da rua; negativo para plataformas elevadas)
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
}

func NewGaroto() *Garoto {
	return &Garoto{
		X:            45.0,
		Y:            0,
		GroundOffset: 0,
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

func (g *Garoto) SetPositionX(x float64) {
	g.X = x
}

func (g *Garoto) GetPosition() (x, y float64) {
	return g.X, g.Y
}

func (g *Garoto) SetGroundOffset(offset float64) {
	g.GroundOffset = offset
	g.Y = offset
	g.VelocityY = 0
	g.IsJumping = false
	g.JumpCount = 0
}

func (g *Garoto) GetGroundOffset() float64 {
	return g.GroundOffset
}

func (g *Garoto) GetVelocityY() float64 {
	return g.VelocityY
}

func (g *Garoto) FallFromPlatform() {
	if g.GroundOffset != 0 {
		g.GroundOffset = 0
		if g.Y < 0 {
			g.IsJumping = true
			if g.VelocityY < 0 {
				g.VelocityY = 0
			}
		}
	}
}

func (g *Garoto) MoveForward(speed float64) {
	g.X += speed
	g.FacingRight = true
	g.IsRunning = true
	if g.X > 260.0 {
		g.X = 260.0
	}
}

func (g *Garoto) MoveBackward(speed float64) {
	g.X -= speed
	g.FacingRight = false
	g.IsRunning = true
	if g.X < 15.0 {
		g.X = 15.0
	}
}

func (g *Garoto) StopRunning() {
	g.IsRunning = false
}

func (g *Garoto) SetAimUp(aim bool) {
	g.AimUp = aim
}

func (g *Garoto) TriggerAttack() {
	g.IsAttacking = true
	g.AttackTimer = 12
}

func (g *Garoto) GetShootOrigin(groundY float64) (x, y, vx, vy float64) {
	speed := 7.5
	if g.AimUp {
		// Tiro vertical para o alto (anti-aéreo contra urubus e pássaros)
		x = g.X + 8.0
		y = groundY - 32.0 + g.Y
		vx = 0.0
		vy = -speed
		return x, y, vx, vy
	}

	dir := 1.0
	x = g.X + 18.0
	if !g.FacingRight {
		dir = -1.0
		x = g.X - 4.0
	}

	if g.IsCrouching {
		y = groundY - 8.0 + g.Y
	} else {
		y = groundY - 17.0 + g.Y
	}

	vx = dir * speed
	vy = 0.0
	return x, y, vx, vy
}

func (g *Garoto) Jump() (jumped bool, isDouble bool) {
	if g.IsSwinging {
		g.ReleaseVine()
		return true, false
	}
	if !g.IsJumping || g.CoyoteTimer > 0 {
		g.GroundOffset = 0 // Pulo liberta da plataforma atual
		g.IsJumping = true
		g.IsCrouching = false
		g.JumpCount = 1
		g.VelocityY = -6.8
		g.JumpHolding = true
		g.CoyoteTimer = 0
		g.spawnDust(g.X+6, g.Y, 6)
		return true, false
	} else if g.JumpCount == 1 {
		g.JumpCount = 2
		g.VelocityY = -6.2
		g.JumpHolding = true
		g.spawnJumpBurst(g.X+7, g.Y+20)
		return true, true
	}
	return false, false
}

func (g *Garoto) GrabVine(v *Vine) {
	g.IsSwinging = true
	g.SwingingVine = v
	g.IsJumping = false
	g.IsCrouching = false
	g.JumpCount = 0
	g.VelocityY = 0
	g.CoyoteTimer = 0
}

func (g *Garoto) ReleaseVine() bool {
	if !g.IsSwinging || g.SwingingVine == nil {
		return false
	}
	g.VelocityY = -6.8 - math.Sin(g.SwingingVine.Angle)*2.5
	g.X += 16.0
	g.IsSwinging = false
	g.IsJumping = true
	g.JumpCount = 1
	g.JumpHolding = true
	g.SwingingVine.Grabbed = false
	g.SwingingVine = nil
	g.spawnJumpBurst(g.X+8, g.Y+12)
	return true
}

func (g *Garoto) ReleaseJump() {
	g.JumpHolding = false
	if g.IsJumping && g.VelocityY < -2.5 {
		g.VelocityY = -2.5
	}
}

func (g *Garoto) SetCrouch(crouch bool) {
	if !g.IsCrouching && crouch && !g.IsJumping && !g.IsSwinging {
		g.spawnDust(g.X+10, g.Y, 4)
	}
	g.IsCrouching = crouch
	if crouch {
		g.IsRunning = false
	}
}

func (g *Garoto) FastDrop() {
	if g.IsJumping && !g.IsSwinging {
		g.VelocityY += 0.9
	}
}

func (g *Garoto) Update() {
	if g.IsRunning && !g.IsSwinging {
		g.RunTicks++
		if g.RunTicks%10 == 0 && !g.IsJumping {
			g.spawnDust(g.X+4, g.Y, 2)
		}
	} else {
		g.RunTicks = 0
	}

	if g.AttackTimer > 0 {
		g.AttackTimer--
		if g.AttackTimer == 0 {
			g.IsAttacking = false
		}
	}

	if g.IsSwinging {
		if g.SwingingVine != nil && g.SwingingVine.Active {
			tipX, tipY := g.SwingingVine.GetTipPosition()
			g.X = tipX - 8.0
			g.Y = tipY - 155.0 + 8.0
			g.VelocityY = 0
			g.IsJumping = false
			g.FacingRight = g.SwingingVine.AngleVelocity >= 0
		} else {
			g.IsSwinging = false
			g.IsJumping = true
			g.VelocityY = 1.0
		}
	} else if g.IsJumping {
		gravity := 0.36
		if !g.JumpHolding && g.VelocityY < 0 {
			gravity = 0.55
		}
		g.Y += g.VelocityY
		g.VelocityY += gravity

		targetY := g.GroundOffset
		if g.Y >= targetY {
			g.Y = targetY
			g.IsJumping = false
			g.JumpCount = 0
			g.VelocityY = 0
			g.spawnDust(g.X+6, targetY, 8)
		}
	} else {
		if g.Y < g.GroundOffset {
			g.IsJumping = true
		}
		g.CoyoteTimer = 6
	}

	if g.CoyoteTimer > 0 {
		g.CoyoteTimer--
	}

	alive := g.Particles[:0]
	for _, p := range g.Particles {
		p.X += p.VX
		p.Y += p.VY
		p.Life++
		if p.Life < p.Max {
			alive = append(alive, p)
		}
	}
	g.Particles = alive
}

func (g *Garoto) spawnDust(x, yOffset float64, count int) {
	dir := -1.0
	if !g.FacingRight {
		dir = 1.0
	}
	for i := 0; i < count; i++ {
		p := &Particle{
			X:     x + float64(i*2),
			Y:     155.0 + yOffset,
			VX:    dir * (1.2 + float64(i)*0.3),
			VY:    -0.3 - float64(i%2)*0.2,
			Life:  0,
			Max:   12 + i*2,
			Size:  2.0,
			Color: color.RGBA{R: 205, G: 195, B: 180, A: 200},
		}
		g.Particles = append(g.Particles, p)
	}
}

func (g *Garoto) spawnJumpBurst(x, y float64) {
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
			Color: color.RGBA{R: 255, G: 220, B: 120, A: 230},
		}
		g.Particles = append(g.Particles, p)
	}
}

func (g *Garoto) Reset() {
	g.X = 45.0
	g.Y = 0
	g.GroundOffset = 0
	g.VelocityY = 0
	g.IsJumping = false
	g.IsCrouching = false
	g.IsRunning = false
	g.FacingRight = true
	g.AimUp = false
	g.IsAttacking = false
	g.AttackTimer = 0
	g.JumpCount = 0
	g.JumpHolding = false
	g.CoyoteTimer = 0
	g.RunTicks = 0
	g.IsSwinging = false
	g.SwingingVine = nil
	g.Particles = g.Particles[:0]
}

func (g *Garoto) GetBounds(groundY float64) (x, y, w, h float64) {
	if g.IsSwinging {
		w = 18.0
		h = 24.0
		x = g.X
		y = groundY - h + g.Y
		return x, y, w, h
	}
	if g.IsCrouching {
		w = 20.0
		h = 16.0
		x = g.X - 2.0
		y = groundY - h + g.Y
		return x, y, w, h
	}
	w = 16.0
	h = 28.0
	x = g.X
	y = groundY - h + g.Y
	return x, y, w, h
}

func (g *Garoto) IsPlayerJumping() bool   { return g.IsJumping }
func (g *Garoto) IsPlayerCrouching() bool { return g.IsCrouching }
func (g *Garoto) IsPlayerSwinging() bool  { return g.IsSwinging }
func (g *Garoto) GetJumpHolding() bool    { return g.JumpHolding }
func (g *Garoto) GetJumpCount() int       { return g.JumpCount }
func (g *Garoto) GetHeroKind() int        { return HeroGaroto }
func (g *Garoto) GetName() string         { return "GAROTO CURUMIM" }


// Draw renderiza o Garoto aventureiro de Belém em autêntico estilo 16-bit SNES
func (g *Garoto) Draw(screen *ebiten.Image, groundY float64, ticks int, invincibleTicks int) {
	// 1. Partículas
	for _, p := range g.Particles {
		alpha := uint8(float64(p.Color.A) * (1.0 - float64(p.Life)/float64(p.Max)))
		c := color.RGBA{R: p.Color.R, G: p.Color.G, B: p.Color.B, A: alpha}
		ebitenutil.DrawRect(screen, p.X, p.Y, p.Size, p.Size, c)
	}

	// Efeito de invulnerabilidade (piscar transparente)
	if invincibleTicks > 0 && (invincibleTicks/4)%2 != 0 {
		return
	}

	// Paleta de Cores do Garoto
	cSkin := color.RGBA{R: 228, G: 168, B: 125, A: 255}       // Tom bronzeado amazônico
	cSkinDark := color.RGBA{R: 195, G: 135, B: 95, A: 255}    // Sombra da pele
	cHair := color.RGBA{R: 30, G: 25, B: 28, A: 255}          // Cabelo preto
	cBandana := color.RGBA{R: 235, G: 50, B: 45, A: 255}      // Faixa vermelha retrô na cabeça
	cBandanaShade := color.RGBA{R: 180, G: 30, B: 30, A: 255} // Sombra da faixa
	cShirt := color.RGBA{R: 35, G: 130, B: 220, A: 255}       // Camisa azul cobalto
	cShirtLight := color.RGBA{R: 70, G: 165, B: 245, A: 255}  // Brilho da camisa
	cBelt := color.RGBA{R: 110, G: 65, B: 35, A: 255}         // Cinto de couro
	cBuckle := color.RGBA{R: 245, G: 200, B: 55, A: 255}      // Fivela dourada
	cPants := color.RGBA{R: 190, G: 160, B: 110, A: 255}      // Bermuda cáqui
	cPantsDark := color.RGBA{R: 135, G: 105, B: 65, A: 255}   // Sombra da bermuda
	cSandalDark := color.RGBA{R: 60, G: 32, B: 18, A: 255}     // Couro escuro da sandália
	cSandalLight := color.RGBA{R: 175, G: 110, B: 55, A: 255}  // Tira de couro da sandália
	cSandalSole := color.RGBA{R: 245, G: 240, B: 225, A: 255}  // Sola clara de alto contraste com o chão
	cWoodDark := color.RGBA{R: 105, G: 55, B: 25, A: 255}       // Casca rústica de madeira amazônica (sombra do galho)
	cWood := color.RGBA{R: 155, G: 95, B: 45, A: 255}           // Corpo da forquilha de madeira talhada
	cWoodLight := color.RGBA{R: 205, G: 145, B: 85, A: 255}     // Brilho do galho polido da baladeira
	cRubber := color.RGBA{R: 245, G: 125, B: 25, A: 255}        // Liga de borracha elástica clássica (âmbar/alaranjada)
	cLeatherPouch := color.RGBA{R: 75, G: 38, B: 15, A: 255}    // Malha / sapata de couro onde assenta o caroço
	cAcaiBerry := color.RGBA{R: 60, G: 15, B: 68, A: 255}       // Caroço de açaí maduro
	cAcaiGleam := color.RGBA{R: 135, G: 45, B: 145, A: 255}     // Brilho do açaí
	cVine := color.RGBA{R: 85, G: 130, B: 45, A: 255}         // Rolo de cipó nas costas
	cEye := color.RGBA{R: 20, G: 20, B: 20, A: 255}

	px := g.X
	py := groundY - 28.0 + g.Y

	// Helper para desenhar com espelhamento horizontal relativo ao centro do garoto
	centerX := px + 8.0
	drawBox := func(relX, relY, w, h float64, col color.Color) {
		var finalX float64
		if g.FacingRight {
			finalX = centerX - 8.0 + relX
		} else {
			finalX = centerX + 8.0 - relX - w
		}
		ebitenutil.DrawRect(screen, finalX, py+relY, w, h, col)
	}

	// ==========================================
	// 0. ESTADO: BALANÇANDO NO CIPÓ (Swinging - Pitfall)
	// ==========================================
	if g.IsSwinging {
		// Cipó descendo pelas mãos
		drawBox(7, -8, 3, 10, cVine)

		// Braços erguidos segurando o cipó acima da cabeça
		drawBox(5, -6, 3, 9, cSkin)
		drawBox(9, -6, 3, 9, cSkin)

		// Cabeça olhando na direção do balanço
		drawBox(4, 1, 9, 8, cSkin)
		drawBox(2, -1, 13, 3, cHair)
		drawBox(2, 0, 3, 7, cHair)
		drawBox(3, 2, 12, 2, cBandana)
		drawBox(0, 3, 3, 2, cBandanaShade)
		drawBox(10, 3, 2, 2, cEye)

		// Tronco
		drawBox(4, 9, 9, 8, cShirt)
		drawBox(5, 10, 4, 6, cShirtLight)
		drawBox(4, 16, 9, 2, cBelt)
		drawBox(8, 16, 2, 2, cBuckle)

		// Baladeira Regional em Y presa no cinto
		drawBox(1, 11, 2, 4, cWood)        // Cabo da forquilha preso no cós
		drawBox(0, 14, 2, 3, cWoodDark)    // Haste esquerda do Y
		drawBox(3, 14, 2, 3, cWoodLight)   // Haste direita do Y
		drawBox(1, 15, 2, 2, cRubber)      // Ligas de borracha pendendo

		// Pernas balançando no ar conforme o ângulo do cipó
		legOffset := 0.0
		if g.SwingingVine != nil {
			legOffset = g.SwingingVine.Angle * 7.0
		}
		drawBox(3-legOffset, 17, 5, 4, cPants)
		drawBox(2-legOffset*1.2, 21, 4, 4, cSkin)
		drawBox(1-legOffset*1.3, 24, 5, 3, cSandalDark)
		drawBox(1-legOffset*1.3, 26, 5, 2, cSandalSole)

		drawBox(8+legOffset, 17, 5, 4, cPantsDark)
		drawBox(9+legOffset*1.2, 21, 4, 4, cSkinDark)
		drawBox(10+legOffset*1.3, 24, 5, 3, cSandalDark)
		drawBox(10+legOffset*1.3, 26, 5, 2, cSandalSole)
		return
	}

	// ==========================================
	// 1. ESTADO: AGACHADO (Crouching)
	// ==========================================
	if g.IsCrouching {
		crouchOffsetY := 12.0

		// Cabelo e Cabeça abaixada
		drawBox(1, crouchOffsetY+2, 11, 7, cSkin)
		drawBox(0, crouchOffsetY+0, 13, 3, cHair)
		drawBox(0, crouchOffsetY+1, 3, 5, cHair)
		drawBox(0, crouchOffsetY+3, 13, 2, cBandana)
		drawBox(7, crouchOffsetY+4, 2, 2, cEye)

		// Tronco encolhido
		drawBox(2, crouchOffsetY+7, 10, 5, cShirt)
		drawBox(1, crouchOffsetY+10, 11, 2, cBelt)

		// Pernas dobradas em cócoras com pés e solas claras visíveis
		drawBox(1, crouchOffsetY+9, 12, 4, cPants)
		drawBox(0, crouchOffsetY+11, 6, 4, cSkin)
		drawBox(7, crouchOffsetY+11, 6, 4, cSkinDark)
		drawBox(-1, crouchOffsetY+13, 6, 3, cSandalDark)
		drawBox(-1, crouchOffsetY+15, 6, 1, cSandalSole)
		drawBox(7, crouchOffsetY+13, 6, 3, cSandalDark)
		drawBox(7, crouchOffsetY+15, 6, 1, cSandalSole)

		// Braço e Baladeira Regional armada horizontalmente na frente do joelho
		drawBox(8, crouchOffsetY+8, 4, 3, cSkin)
		drawBox(12, crouchOffsetY+7, 2, 3, cSkinDark) // Mão segurando a empunhadura
		// Forquilha em Y de madeira regional
		drawBox(13, crouchOffsetY+7, 2, 4, cWood)      // Cabo da baladeira
		drawBox(14, crouchOffsetY+5, 2, 2, cWoodDark)  // Bifurcação central
		drawBox(13, crouchOffsetY+2, 2, 3, cWood)      // Haste superior do Y
		drawBox(16, crouchOffsetY+2, 2, 3, cWoodLight) // Haste inferior/frontal do Y
		// Ligas elásticas e semente de açaí
		drawBox(12, crouchOffsetY+3, 1, 2, cRubber)
		drawBox(15, crouchOffsetY+3, 1, 2, cRubber)
		drawBox(11, crouchOffsetY+5, 2, 2, cLeatherPouch)
		drawBox(11, crouchOffsetY+5, 1, 1, cAcaiBerry)
		return
	}

	// ==========================================
	// 2. ESTADO: NO AR (Jumping / Falling)
	// ==========================================
	if g.IsJumping {
		hairOffset := -1.0
		if g.VelocityY > 0 {
			hairOffset = 1.0
		}

		// Cipó enrolado nas costas
		drawBox(0, 10, 4, 8, cVine)

		// Cabeça & Faixa ao vento
		drawBox(2, 2, 10, 8, cSkin)
		drawBox(0, 0+hairOffset, 13, 3, cHair)
		drawBox(0, 1, 3, 7, cHair)
		drawBox(1, 3, 12, 2, cBandana)
		drawBox(-2, 3, 3, 2, cBandanaShade) // Ponta da faixa esvoaçando atrás
		drawBox(8, 4, 2, 2, cEye)

		// Tronco
		drawBox(2, 10, 9, 8, cShirt)
		drawBox(3, 11, 4, 6, cShirtLight)
		drawBox(2, 17, 9, 2, cBelt)
		drawBox(6, 17, 2, 2, cBuckle)

		// Braços e Baladeira Regional no Ar (Salto Acrobático)
		if g.AimUp {
			// Tiro vertical anti-aéreo para o alto
			drawBox(7, 3, 3, 7, cSkin)
			drawBox(7, 0, 3, 3, cSkinDark)
			// Forquilha em Y apontando para o céu
			drawBox(7, -2, 3, 3, cWood)
			drawBox(7, -4, 3, 2, cWoodDark)
			drawBox(5, -7, 2, 3, cWood)
			drawBox(10, -7, 2, 3, cWoodLight)
			drawBox(5, -8, 2, 1, cWoodDark)
			drawBox(10, -8, 2, 1, cWoodDark)
			// Ligas elásticas e semente de açaí
			drawBox(6, -6, 1, 3, cRubber)
			drawBox(9, -6, 1, 3, cRubber)
			drawBox(7, -4, 3, 2, cLeatherPouch)
			drawBox(7, -4, 2, 2, cAcaiBerry)
		} else if g.IsAttacking {
			// Disparo no ar: pose ágil de tiro com baladeira esticada
			drawBox(8, 10, 8, 3, cSkin)
			drawBox(15, 9, 2, 4, cSkinDark)
			// Forquilha em Y
			drawBox(16, 9, 2, 4, cWood)
			drawBox(16, 7, 2, 2, cWoodDark)
			drawBox(15, 4, 2, 3, cWood)
			drawBox(18, 4, 2, 3, cWoodLight)
			drawBox(15, 3, 2, 1, cWoodDark)
			drawBox(18, 3, 2, 1, cWoodDark)
			// Borracha esticada até a mão traseira
			drawBox(5, 9, 3, 3, cSkinDark)
			drawBox(7, 6, 8, 1, cRubber)
			drawBox(7, 8, 8, 1, cRubber)
			drawBox(5, 7, 2, 2, cLeatherPouch)
			drawBox(5, 7, 1, 1, cAcaiBerry)
		} else {
			// Empunhando a baladeira em Y durante o salto
			drawBox(9, 11, 5, 3, cSkin)
			drawBox(13, 10, 2, 3, cSkinDark)
			drawBox(14, 9, 2, 4, cWood)
			drawBox(14, 7, 2, 2, cWoodDark)
			drawBox(13, 4, 2, 3, cWood)
			drawBox(16, 4, 2, 3, cWoodLight)
			drawBox(13, 3, 2, 1, cWoodDark)
			drawBox(16, 3, 2, 1, cWoodDark)
			drawBox(12, 5, 2, 2, cRubber)
			drawBox(15, 5, 2, 2, cRubber)
			drawBox(11, 7, 2, 2, cLeatherPouch)
			drawBox(11, 7, 1, 1, cAcaiBerry)
		}

		// Pernas abertas no ar em salto acrobático 16-bit com sandálias delineadas
		// Perna dianteira dobrada com joelho erguido
		drawBox(6, 17, 6, 4, cPants)
		drawBox(9, 19, 5, 4, cSkin)
		drawBox(10, 22, 5, 3, cSandalDark)
		drawBox(10, 24, 5, 2, cSandalSole)

		// Perna traseira esticada para trás no salto
		drawBox(-1, 17, 5, 4, cPantsDark)
		drawBox(-4, 20, 5, 4, cSkinDark)
		drawBox(-6, 23, 5, 3, cSandalDark)
		drawBox(-6, 25, 5, 2, cSandalSole)
		return
	}

	// ==========================================
	// 3. ESTADO: CHÃO (Idle ou Running)
	// ==========================================
	bobY := 0.0
	legFrame := 0

	if g.IsRunning {
		// Ciclo de corrida de 4 frames
		legFrame = (g.RunTicks / 5) % 4
		if legFrame == 1 || legFrame == 3 {
			bobY = -1.0
		}
	} else {
		// Respiração idle
		if (ticks/24)%2 == 0 {
			bobY = 1.0
		}
	}

	// Cipó enrolado a tiracolo nas costas
	drawBox(0, 10+bobY, 4, 8, cVine)

	// Cabeça
	drawBox(3, 2+bobY, 10, 8, cSkin)
	drawBox(1, 0+bobY, 13, 3, cHair)
	drawBox(1, 1+bobY, 3, 7, cHair)
	drawBox(2, 3+bobY, 12, 2, cBandana)
	drawBox(-1, 4+bobY, 3, 2, cBandanaShade) // Ponta da faixa amarrada atrás
	drawBox(9, 4+bobY, 2, 2, cEye)
	drawBox(10, 7+bobY, 2, 1, cSkinDark)

	// Tronco (Regata Azul e Cinto com fivela)
	drawBox(3, 10+bobY, 9, 8, cShirt)
	drawBox(4, 11+bobY, 4, 6, cShirtLight)
	drawBox(3, 17+bobY, 9, 2, cBelt)
	drawBox(7, 17+bobY, 2, 2, cBuckle)

	// Braço dianteiro segurando a autêntica Baladeira Regional Paraense (Estilingue de Forquilha em Y)
	if g.AimUp {
		// Tiro Anti-Aéreo: Apontando a baladeira para o céu contra urubus e pássaros
		drawBox(8, 3+bobY, 3, 7, cSkin)
		drawBox(8, 0+bobY, 3, 3, cSkinDark)
		// Forquilha em Y voltada para o alto
		drawBox(8, -2+bobY, 3, 3, cWood)
		drawBox(8, -4+bobY, 3, 2, cWoodDark)
		drawBox(6, -7+bobY, 2, 3, cWood)
		drawBox(11, -7+bobY, 2, 3, cWoodLight)
		drawBox(6, -8+bobY, 2, 1, cWoodDark)
		drawBox(11, -8+bobY, 2, 1, cWoodDark)
		// Ligas elásticas e caroço de açaí
		drawBox(7, -6+bobY, 1, 3, cRubber)
		drawBox(10, -6+bobY, 1, 3, cRubber)
		drawBox(8, -3+bobY, 3, 2, cLeatherPouch)
		drawBox(8, -3+bobY, 2, 2, cAcaiBerry)
		drawBox(9, -3+bobY, 1, 1, cAcaiGleam)
	} else if g.IsAttacking {
		// ============================================================
		// POSTURA OFICIAL DE ATIRADOR DE BALADEIRA (AÇAÍ RETRO ARCADE)
		// ============================================================
		// Braço dianteiro esticado à frente segurando a forquilha
		drawBox(9, 10+bobY, 8, 3, cSkin)
		drawBox(16, 9+bobY, 2, 4, cSkinDark)

		// Forquilha em Y bem nítida e imponente
		drawBox(17, 9+bobY, 2, 5, cWood)      // Cabo do galho
		drawBox(17, 7+bobY, 2, 2, cWoodDark)  // Bifurcação
		drawBox(16, 4+bobY, 2, 3, cWood)      // Haste superior
		drawBox(19, 4+bobY, 2, 3, cWoodLight) // Haste frontal
		drawBox(16, 3+bobY, 2, 1, cWoodDark)  // Amarração superior
		drawBox(19, 3+bobY, 2, 1, cWoodDark)  // Amarração inferior

		// Braço traseiro puxado até a bochecha esticando a borracha
		drawBox(5, 9+bobY, 4, 3, cSkinDark)
		drawBox(4, 8+bobY, 2, 3, cSkin) // Mão pinçando a sapata

		// Tiras elásticas esticadas em linha reta tensa de disparo
		drawBox(6, 6+bobY, 5, 1, cRubber)
		drawBox(11, 5+bobY, 5, 1, cRubber)
		drawBox(6, 8+bobY, 5, 1, cRubber)
		drawBox(11, 7+bobY, 5, 1, cRubber)

		// Malha de couro e caroço de açaí pronto para voar
		drawBox(4, 7+bobY, 3, 2, cLeatherPouch)
		drawBox(4, 7+bobY, 2, 2, cAcaiBerry)
		drawBox(4, 7+bobY, 1, 1, cAcaiGleam)
	} else {
		// Empunhadura de repouso / corrida: forquilha em Y sempre à mostra
		drawBox(10, 11+bobY, 5, 3, cSkin)
		drawBox(14, 10+bobY, 2, 4, cSkinDark)

		// Forquilha em Y de madeira regional
		drawBox(15, 10+bobY, 2, 4, cWood)       // Cabo do galho
		drawBox(15, 8+bobY, 2, 2, cWoodDark)    // Bifurcação
		drawBox(14, 5+bobY, 2, 3, cWood)        // Haste esquerda do Y
		drawBox(17, 5+bobY, 2, 3, cWoodLight)   // Haste direita do Y
		drawBox(14, 4+bobY, 2, 1, cWoodDark)    // Amarração esquerda
		drawBox(17, 4+bobY, 2, 1, cWoodDark)    // Amarração direita

		// Ligas de borracha descendo até a malha de couro
		drawBox(13, 6+bobY, 2, 1, cRubber)
		drawBox(12, 7+bobY, 2, 1, cRubber)
		drawBox(16, 6+bobY, 2, 1, cRubber)
		drawBox(15, 7+bobY, 2, 1, cRubber)

		// Sapata de couro segurando o caroço de açaí
		drawBox(11, 8+bobY, 3, 2, cLeatherPouch)
		drawBox(11, 8+bobY, 2, 2, cAcaiBerry)
		drawBox(12, 8+bobY, 1, 1, cAcaiGleam)
	}

	// Pernas & Animação de Passadas Expressivas (estilo galope da onça com passada ampla)
	if !g.IsRunning {
		// Postura em pé firme (Idle respirando) com base sólida e pés no chão
		drawBox(1, 17+bobY, 5, 4, cPantsDark)
		drawBox(7, 17+bobY, 5, 4, cPants)
		drawBox(2, 20+bobY, 4, 5, cSkinDark)
		drawBox(8, 20+bobY, 4, 5, cSkin)
		drawBox(1, 24+bobY, 5, 3, cSandalDark)
		drawBox(1, 26+bobY, 6, 2, cSandalSole)
		drawBox(7, 24+bobY, 5, 3, cSandalDark)
		drawBox(7, 26+bobY, 6, 2, cSandalSole)
		return
	}

	// Ciclo de corrida de 4 quadros amplo e dinâmico
	switch legFrame {
	case 0:
		// Passo 1: Passada aberta máxima - Perna dianteira esticada à frente, traseira empurrando atrás
		// Perna Traseira
		drawBox(-2, 17+bobY, 5, 4, cPantsDark)
		drawBox(-4, 20+bobY, 5, 4, cSkinDark)
		drawBox(-6, 23+bobY, 5, 3, cSandalDark)
		drawBox(-5, 23+bobY, 2, 2, cSandalLight)
		drawBox(-6, 25+bobY, 5, 2, cSandalSole)

		// Perna Dianteira
		drawBox(6, 17+bobY, 6, 4, cPants)
		drawBox(10, 20+bobY, 5, 5, cSkin)
		drawBox(11, 24+bobY, 6, 3, cSandalDark)
		drawBox(12, 24+bobY, 3, 2, cSandalLight)
		drawBox(11, 26+bobY, 7, 2, cSandalSole)

	case 1:
		// Passo 2: Passagem - Perna de apoio sob o corpo, perna traseira avançando dobrada no ar
		// Perna de Apoio
		drawBox(4, 17+bobY, 6, 4, cPants)
		drawBox(5, 20+bobY, 5, 5, cSkin)
		drawBox(4, 24+bobY, 6, 3, cSandalDark)
		drawBox(5, 24+bobY, 3, 1, cSandalLight)
		drawBox(4, 26+bobY, 7, 2, cSandalSole)

		// Perna Traseira Avançando
		drawBox(0, 17+bobY, 5, 4, cPantsDark)
		drawBox(-1, 20+bobY, 5, 4, cSkinDark)
		drawBox(-3, 23+bobY, 4, 3, cSandalDark)
		drawBox(-3, 25+bobY, 4, 2, cSandalSole)

	case 2:
		// Passo 3: Passada aberta máxima invertida - Perna esquerda à frente, direita atrás
		// Perna Traseira
		drawBox(-2, 17+bobY, 5, 4, cPants)
		drawBox(-4, 20+bobY, 5, 4, cSkin)
		drawBox(-6, 23+bobY, 5, 3, cSandalDark)
		drawBox(-5, 23+bobY, 2, 2, cSandalLight)
		drawBox(-6, 25+bobY, 5, 2, cSandalSole)

		// Perna Dianteira
		drawBox(6, 17+bobY, 6, 4, cPantsDark)
		drawBox(10, 20+bobY, 5, 5, cSkinDark)
		drawBox(11, 24+bobY, 6, 3, cSandalDark)
		drawBox(12, 24+bobY, 3, 2, cSandalLight)
		drawBox(11, 26+bobY, 7, 2, cSandalSole)

	case 3:
		// Passo 4: Passagem invertida
		// Perna de Apoio
		drawBox(4, 17+bobY, 6, 4, cPantsDark)
		drawBox(5, 20+bobY, 5, 5, cSkinDark)
		drawBox(4, 24+bobY, 6, 3, cSandalDark)
		drawBox(4, 26+bobY, 7, 2, cSandalSole)

		// Perna Traseira Avançando
		drawBox(0, 17+bobY, 5, 4, cPants)
		drawBox(-1, 20+bobY, 5, 4, cSkin)
		drawBox(-3, 23+bobY, 4, 3, cSandalDark)
		drawBox(-3, 25+bobY, 4, 2, cSandalSole)
	}
}
