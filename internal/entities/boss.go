package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type BossState int

const (
	BossStateInactive BossState = iota
	BossStateEntering
	BossStateFighting
	BossStateHurt
	BossStateDefeated
)

type Boss struct {
	Stage       int
	Name        string
	HP          int
	MaxHP       int
	X           float64
	Y           float64
	BaseY       float64
	Width       float64
	Height      float64
	VelX        float64
	VelY        float64
	State       BossState
	StateTicks  int
	AttackTimer int
	HurtTicks   int
	ScreenWidth float64
	GroundY     float64
	Defeated    bool
}

func NewBoss(stage int, screenWidth, groundY float64) *Boss {
	name := "URUBU-REI DO VER-O-PESO"
	maxHP := 8
	switch stage {
	case 2:
		name = "COBRA-GRANDE DAS DOCAS"
		maxHP = 10
	case 3:
		name = "CURUPIRA PROTETOR"
		maxHP = 12
	}

	return &Boss{
		Stage:       stage,
		Name:        name,
		HP:          maxHP,
		MaxHP:       maxHP,
		X:           screenWidth + 60,
		Y:           groundY - 70,
		BaseY:       groundY - 70,
		Width:       44.0,
		Height:      36.0,
		VelX:        -1.5,
		VelY:        0,
		State:       BossStateEntering,
		StateTicks:  0,
		AttackTimer: 0,
		HurtTicks:   0,
		ScreenWidth: screenWidth,
		GroundY:     groundY,
		Defeated:    false,
	}
}

func (b *Boss) IsActive() bool {
	return b.State != BossStateInactive
}

func (b *Boss) GetBounds() (x, y, w, h float64) {
	if b.Defeated || b.State == BossStateInactive {
		return 0, 0, 0, 0
	}
	return b.X, b.Y, b.Width, b.Height
}

func (b *Boss) TakeDamage(amount int) bool {
	if b.State == BossStateDefeated || b.State == BossStateInactive || b.HurtTicks > 0 {
		return false
	}
	b.HP -= amount
	b.HurtTicks = 20
	if b.HP <= 0 {
		b.HP = 0
		b.State = BossStateDefeated
		b.StateTicks = 0
		b.Defeated = true
		return true
	}
	return false
}

func (b *Boss) Update(playerX, playerY float64) (isDone bool) {
	if b.State == BossStateInactive {
		return false
	}

	b.StateTicks++
	if b.HurtTicks > 0 {
		b.HurtTicks--
	}

	switch b.State {
	case BossStateEntering:
		b.X += b.VelX
		targetX := b.ScreenWidth - 70
		if b.X <= targetX {
			b.X = targetX
			b.State = BossStateFighting
			b.StateTicks = 0
			b.AttackTimer = 60
		}

	case BossStateFighting:
		b.AttackTimer++
		// Movimento oscilatório suave (voo ou ondulação amazônica)
		wave := math.Sin(float64(b.StateTicks) * 0.08)
		b.Y = b.BaseY + wave*18.0

		// Padrão de investida e recuo a cada ciclo
		cycle := (b.AttackTimer / 60) % 4
		switch cycle {
		case 1: // Investida em direção ao jogador
			if b.X > playerX+40 {
				b.X -= 1.2
			}
		case 3: // Recuo estratégico
			targetX := b.ScreenWidth - 60
			if b.X < targetX {
				b.X += 1.0
			}
		}

	case BossStateDefeated:
		// Queda com efeito dramático
		b.Y += 1.5
		b.X += 0.5
		if b.StateTicks >= 90 {
			b.State = BossStateInactive
			return true // Chefe derrotado, avança fase
		}
	}

	return false
}

func (b *Boss) Draw(screen *ebiten.Image, ticks int) {
	if b.State == BossStateInactive {
		return
	}

	// Pisca se estiver em hurt ou derrotado
	if b.HurtTicks > 0 && (b.HurtTicks/3)%2 != 0 {
		return
	}
	if b.State == BossStateDefeated && (b.StateTicks/4)%2 != 0 {
		return
	}

	bx := b.X
	by := b.Y

	switch b.Stage {
	case 1:
		b.drawUrubuRei(screen, bx, by, ticks)
	case 2:
		b.drawCobraGrande(screen, bx, by, ticks)
	case 3:
		b.drawCurupira(screen, bx, by, ticks)
	}

	// Barra de Vida do Chefe no Topo da Tela
	b.drawHUD(screen)
}

func (b *Boss) drawUrubuRei(screen *ebiten.Image, bx, by float64, ticks int) {
	cDarkFeather := color.RGBA{R: 25, G: 25, B: 30, A: 255}
	cLightFeather := color.RGBA{R: 220, G: 220, B: 230, A: 255}
	cCrownGold := color.RGBA{R: 255, G: 215, B: 0, A: 255}
	cRedBeak := color.RGBA{R: 225, G: 45, B: 35, A: 255}
	cOrangeEye := color.RGBA{R: 255, G: 140, B: 0, A: 255}

	// 1. Corpo Imponente do Urubu-Rei
	ebitenutil.DrawRect(screen, bx+8, by+8, 28, 20, cDarkFeather)
	ebitenutil.DrawRect(screen, bx+12, by+10, 20, 14, cLightFeather)

	// 2. Asas Grandes com Batimento
	wingFlap := int((float64(ticks) / 5.0)) % 2
	if wingFlap == 0 {
		ebitenutil.DrawRect(screen, bx+2, by+2, 14, 10, cDarkFeather)
		ebitenutil.DrawRect(screen, bx+26, by+2, 14, 10, cDarkFeather)
	} else {
		ebitenutil.DrawRect(screen, bx+2, by+12, 14, 10, cDarkFeather)
		ebitenutil.DrawRect(screen, bx+26, by+12, 14, 10, cDarkFeather)
	}

	// 3. Cabeça e Colar Vermelho / Laranja
	ebitenutil.DrawRect(screen, bx, by+10, 10, 10, cRedBeak)
	ebitenutil.DrawRect(screen, bx-4, by+13, 6, 5, cOrangeEye)

	// 4. Coroa Dourada Real do Urubu-Rei
	ebitenutil.DrawRect(screen, bx, by+4, 10, 4, cCrownGold)
	ebitenutil.DrawRect(screen, bx+1, by+1, 2, 4, cCrownGold)
	ebitenutil.DrawRect(screen, bx+4, by, 2, 5, cCrownGold)
	ebitenutil.DrawRect(screen, bx+7, by+1, 2, 4, cCrownGold)
}

func (b *Boss) drawCobraGrande(screen *ebiten.Image, bx, by float64, ticks int) {
	cDarkGreen := color.RGBA{R: 15, G: 80, B: 35, A: 255}
	cLightGreen := color.RGBA{R: 40, G: 160, B: 70, A: 255}
	cGlowingYellow := color.RGBA{R: 255, G: 230, B: 50, A: 255}

	ebitenutil.DrawRect(screen, bx, by+4, 40, 24, cDarkGreen)
	ebitenutil.DrawRect(screen, bx+4, by+8, 32, 16, cLightGreen)
	ebitenutil.DrawRect(screen, bx-2, by+10, 8, 8, cGlowingYellow)
}

func (b *Boss) drawCurupira(screen *ebiten.Image, bx, by float64, ticks int) {
	cSkin := color.RGBA{R: 180, G: 110, B: 70, A: 255}
	cHairFlame := color.RGBA{R: 235, G: 60, B: 20, A: 255}

	ebitenutil.DrawRect(screen, bx+8, by+10, 24, 22, cSkin)
	ebitenutil.DrawRect(screen, bx+6, by, 28, 12, cHairFlame)
}

func (b *Boss) drawHUD(screen *ebiten.Image) {
	barW := 140.0
	barH := 8.0
	barX := (b.ScreenWidth - barW) / 2.0
	barY := 26.0

	// Fundo da barra
	ebitenutil.DrawRect(screen, barX-2, barY-2, barW+4, barH+4, color.RGBA{R: 10, G: 10, B: 15, A: 220})
	ebitenutil.DrawRect(screen, barX, barY, barW, barH, color.RGBA{R: 60, G: 20, B: 20, A: 255})

	// Preenchimento de vida proporcional
	hpPercent := float64(b.HP) / float64(b.MaxHP)
	if hpPercent < 0 {
		hpPercent = 0
	}
	fillW := barW * hpPercent
	ebitenutil.DrawRect(screen, barX, barY, fillW, barH, color.RGBA{R: 235, G: 45, B: 45, A: 255})

	// Nome do Chefe
	ebitenutil.DebugPrintAt(screen, b.Name, int(barX+4), int(barY-12))
}
