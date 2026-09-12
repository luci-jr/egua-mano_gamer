package entities

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ProjKindAcai = 0
	ProjKindRoar = 1
)

type Projectile struct {
	X         float64
	Y         float64
	VX        float64
	VY        float64
	Kind      int
	IsSpecial bool
	HitsLeft  int
	Active    bool
	Life      int
	MaxLife   int
}

func NewProjectile(x, y, vx, vy float64, kind int, isSpecial bool) *Projectile {
	hits := 1
	maxLife := 85
	if isSpecial {
		maxLife = 95
		if kind == ProjKindRoar {
			hits = 3 // Mega Rugido Alfa perfura até 3 obstáculos
		} else {
			hits = 2 // Super Caroço de Açaí perfura até 2 obstáculos
		}
	}
	return &Projectile{
		X:         x,
		Y:         y,
		VX:        vx,
		VY:        vy,
		Kind:      kind,
		IsSpecial: isSpecial,
		HitsLeft:  hits,
		Active:    true,
		Life:      0,
		MaxLife:   maxLife,
	}
}

func (p *Projectile) GetBounds() (x, y, w, h float64) {
	if p.IsSpecial {
		if p.Kind == ProjKindRoar {
			return p.X - 8.0, p.Y - 8.0, 18.0, 18.0
		}
		return p.X - 6.0, p.Y - 6.0, 13.0, 13.0
	}
	if p.Kind == ProjKindRoar {
		return p.X - 4.0, p.Y - 5.0, 10.0, 10.0
	}
	return p.X - 3.0, p.Y - 3.0, 7.0, 7.0
}

type ScorePopup struct {
	X, Y    float64
	Text    string
	Life    int
	MaxLife int
	Color   color.RGBA
}

type ProjectileManager struct {
	Projectiles []*Projectile
	Popups      []*ScorePopup
	Particles   []*Particle
	screenWidth float64
}

func NewProjectileManager(screenWidth float64) *ProjectileManager {
	return &ProjectileManager{
		Projectiles: make([]*Projectile, 0, 16),
		Popups:      make([]*ScorePopup, 0, 16),
		Particles:   make([]*Particle, 0, 64),
		screenWidth: screenWidth,
	}
}

func (pm *ProjectileManager) ShootAcai(x, y, vx, vy float64) {
	pm.Shoot(x, y, vx, vy, ProjKindAcai, false)
}

func (pm *ProjectileManager) ShootRoar(x, y, vx, vy float64) {
	pm.Shoot(x, y, vx, vy, ProjKindRoar, false)
}

func (pm *ProjectileManager) ShootSpecialAcai(x, y, vx, vy float64) {
	pm.Shoot(x, y, vx*1.2, vy*1.2, ProjKindAcai, true)
}

func (pm *ProjectileManager) ShootSpecialRoar(x, y, vx, vy float64) {
	pm.Shoot(x, y, vx*1.2, vy*1.2, ProjKindRoar, true)
}

func (pm *ProjectileManager) Shoot(x, y, vx, vy float64, kind int, isSpecial bool) {
	// Limite máximo de disparos simultâneos na tela
	activeCount := 0
	for _, p := range pm.Projectiles {
		if p.Active {
			activeCount++
		}
	}
	if activeCount >= 6 {
		return
	}

	pm.Projectiles = append(pm.Projectiles, NewProjectile(x, y, vx, vy, kind, isSpecial))

	// Partículas de disparo
	sparkColor := color.RGBA{R: 245, G: 215, B: 85, A: 240}
	if isSpecial {
		sparkColor = color.RGBA{R: 255, G: 240, B: 120, A: 255}
	} else if kind == ProjKindRoar {
		sparkColor = color.RGBA{R: 255, G: 165, B: 40, A: 255}
	}

	count := 5
	if isSpecial {
		count = 10
	}
	for i := 0; i < count; i++ {
		angle := float64(i)*0.5 - 1.0
		if vx < 0 {
			angle = math.Pi - angle
		} else if vy < 0 {
			angle = -math.Pi/2.0 + (float64(i)-2.0)*0.35
		}
		speed := 2.0 + float64(i)*0.4
		if isSpecial {
			speed *= 1.4
		}
		pm.Particles = append(pm.Particles, &Particle{
			X:     x,
			Y:     y,
			VX:    math.Cos(angle) * speed,
			VY:    math.Sin(angle) * speed,
			Life:  0,
			Max:   14,
			Size:  2.8,
			Color: sparkColor,
		})
	}
}

func (pm *ProjectileManager) SpawnHitBurst(x, y float64, col color.RGBA, count int) {
	for i := 0; i < count; i++ {
		angle := float64(i) * (2.0 * math.Pi / float64(count))
		speed := 1.5 + float64(i%3)*0.8
		pm.Particles = append(pm.Particles, &Particle{
			X:     x,
			Y:     y,
			VX:    math.Cos(angle) * speed,
			VY:    math.Sin(angle) * speed,
			Life:  0,
			Max:   14 + (i%3)*4,
			Size:  2.5,
			Color: col,
		})
	}
}

func (pm *ProjectileManager) AddScorePopup(x, y float64, score int) {
	pm.Popups = append(pm.Popups, &ScorePopup{
		X:       x,
		Y:       y,
		Text:    fmt.Sprintf("+%d", score),
		Life:    0,
		MaxLife: 32,
		Color:   color.RGBA{R: 255, G: 225, B: 60, A: 255},
	})
}

func (pm *ProjectileManager) AddTextPopup(x, y float64, text string, col color.RGBA) {
	pm.Popups = append(pm.Popups, &ScorePopup{
		X:       x,
		Y:       y,
		Text:    text,
		Life:    0,
		MaxLife: 38,
		Color:   col,
	})
}

func (pm *ProjectileManager) Update(screenWidth float64) {
	// 1. Atualizar Projéteis
	active := pm.Projectiles[:0]
	for _, p := range pm.Projectiles {
		if !p.Active {
			continue
		}
		p.X += p.VX
		p.Y += p.VY
		p.Life++

		// Partícula de rastro
		if p.Life%2 == 0 {
			trailCol := color.RGBA{R: 160, G: 60, B: 180, A: 180} // Roxo açaí
			if p.Kind == ProjKindRoar {
				trailCol = color.RGBA{R: 255, G: 190, B: 50, A: 190} // Ouro rugido
			}
			pm.Particles = append(pm.Particles, &Particle{
				X:     p.X - p.VX*0.35,
				Y:     p.Y - p.VY*0.35,
				VX:    -p.VX * 0.08,
				VY:    -p.VY * 0.08,
				Life:  0,
				Max:   8,
				Size:  2.0,
				Color: trailCol,
			})
		}

		// Checa limites de tela
		if p.X < -15 || p.X > screenWidth+15 || p.Y < -20 || p.Y > 210 || p.Life >= p.MaxLife {
			p.Active = false
			continue
		}
		active = append(active, p)
	}
	pm.Projectiles = active

	// 2. Atualizar Popups de Pontos Flutuantes
	activePopups := pm.Popups[:0]
	for _, pop := range pm.Popups {
		pop.Y -= 0.6 // Flutua suavemente para cima
		pop.Life++
		if pop.Life < pop.MaxLife {
			activePopups = append(activePopups, pop)
		}
	}
	pm.Popups = activePopups

	// 3. Atualizar Partículas de Impacto
	aliveParticles := pm.Particles[:0]
	for _, pt := range pm.Particles {
		pt.X += pt.VX
		pt.Y += pt.VY
		pt.Life++
		if pt.Life < pt.Max {
			aliveParticles = append(aliveParticles, pt)
		}
	}
	pm.Particles = aliveParticles
}

func (pm *ProjectileManager) Draw(screen *ebiten.Image, ticks int) {
	// 1. Partículas
	for _, pt := range pm.Particles {
		alpha := uint8(float64(pt.Color.A) * (1.0 - float64(pt.Life)/float64(pt.Max)))
		c := color.RGBA{R: pt.Color.R, G: pt.Color.G, B: pt.Color.B, A: alpha}
		ebitenutil.DrawRect(screen, pt.X, pt.Y, pt.Size, pt.Size, c)
	}

	// 2. Projéteis
	cAcaiOuter := color.RGBA{R: 50, G: 15, B: 55, A: 255}       // Casca do açaí roxo profundo
	cAcaiInner := color.RGBA{R: 130, G: 45, B: 150, A: 255}     // Polpa de açaí vibrante
	cAcaiCore := color.RGBA{R: 255, G: 230, B: 140, A: 255}     // Brilho do caroço em alta velocidade

	cRoarOuter := color.RGBA{R: 220, G: 120, B: 25, A: 220}     // Onda de choque âmbar
	cRoarInner := color.RGBA{R: 255, G: 215, B: 50, A: 250}     // Arco sonoro dourado
	cRoarCore := color.RGBA{R: 255, G: 250, B: 210, A: 255}     // Núcleo branco energia

	for _, p := range pm.Projectiles {
		if !p.Active {
			continue
		}

		if p.Kind == ProjKindRoar {
			// Onda de Choque Sônica (Rugido da Onça) - Arcos concêntricos dourados
			px := p.X
			py := p.Y
			isUp := p.VY < 0

			if p.IsSpecial {
				// Mega Rugido Alfa: Onda tripla gigante com brilho radiante
				cSpecialGold := color.RGBA{R: 255, G: 235, B: 80, A: 255}
				cSpecialCrimson := color.RGBA{R: 255, G: 85, B: 30, A: 240}
				if isUp {
					ebitenutil.DrawRect(screen, px-12, py+4, 24, 3, cSpecialCrimson)
					ebitenutil.DrawRect(screen, px-10, py+1, 20, 3, cSpecialGold)
					ebitenutil.DrawRect(screen, px-7, py-2, 14, 3, cRoarCore)
					ebitenutil.DrawRect(screen, px-3, py-5, 6, 3, cRoarCore)
				} else {
					dir := 1.0
					if p.VX < 0 {
						dir = -1.0
					}
					ebitenutil.DrawRect(screen, px, py-9, 3, 18, cSpecialCrimson)
					ebitenutil.DrawRect(screen, px+dir*3, py-7, 3, 14, cSpecialGold)
					ebitenutil.DrawRect(screen, px+dir*6, py-4, 3, 8, cRoarCore)
					ebitenutil.DrawRect(screen, px-dir*3, py-5, 2, 10, cSpecialCrimson)
				}
				continue
			}

			if isUp {
				// Arco virado para cima: ^
				ebitenutil.DrawRect(screen, px-6, py+2, 12, 2, cRoarOuter)
				ebitenutil.DrawRect(screen, px-5, py, 10, 2, cRoarInner)
				ebitenutil.DrawRect(screen, px-3, py-2, 6, 2, cRoarCore)
				ebitenutil.DrawRect(screen, px-1, py-3, 2, 2, cRoarCore)
			} else {
				// Arco virado para a direita: )
				dir := 1.0
				if p.VX < 0 {
					dir = -1.0
				}
				// Onda 1 (principal)
				ebitenutil.DrawRect(screen, px, py-5, 2, 10, cRoarOuter)
				ebitenutil.DrawRect(screen, px+dir*2, py-4, 2, 8, cRoarInner)
				ebitenutil.DrawRect(screen, px+dir*4, py-2, 2, 4, cRoarCore)
				// Onda 2 (traseira menor)
				ebitenutil.DrawRect(screen, px-dir*3, py-3, 2, 6, cRoarOuter)
				ebitenutil.DrawRect(screen, px-dir*2, py-2, 2, 4, cRoarInner)
			}
			continue
		}

		if p.IsSpecial {
			// Super Semente de Açaí Dourada Energizada (10x10 brilhante)
			px := p.X - 5.0
			py := p.Y - 5.0
			cAura := color.RGBA{R: 255, G: 200, B: 40, A: 160}
			cGoldAcai := color.RGBA{R: 245, G: 165, B: 25, A: 255}
			ebitenutil.DrawRect(screen, px-1, py-1, 12, 12, cAura)
			ebitenutil.DrawRect(screen, px+1, py, 8, 10, cAcaiOuter)
			ebitenutil.DrawRect(screen, px, py+1, 10, 8, cAcaiOuter)
			ebitenutil.DrawRect(screen, px+2, py+2, 6, 6, cGoldAcai)
			ebitenutil.DrawRect(screen, px+3, py+3, 4, 4, cAcaiCore)
			continue
		}

		// Sementes de Açaí Tradicionais do Garoto (com rotação visual)
		rot := (ticks + int(p.X)) % 4
		px := p.X - 3.0
		py := p.Y - 3.0

		// Esfera 6x6 pixel art com cantos recortados
		ebitenutil.DrawRect(screen, px+1, py, 4, 6, cAcaiOuter)
		ebitenutil.DrawRect(screen, px, py+1, 6, 4, cAcaiOuter)
		ebitenutil.DrawRect(screen, px+1, py+1, 4, 4, cAcaiInner)

		// Brilho pulsante
		switch rot {
		case 0:
			ebitenutil.DrawRect(screen, px+2, py+1, 2, 2, cAcaiCore)
		case 1:
			ebitenutil.DrawRect(screen, px+3, py+2, 2, 2, cAcaiCore)
		case 2:
			ebitenutil.DrawRect(screen, px+2, py+3, 2, 2, cAcaiCore)
		case 3:
			ebitenutil.DrawRect(screen, px+1, py+2, 2, 2, cAcaiCore)
		}
	}

	// 3. Popups de Pontos (+50, +100) flutuando no ar
	for _, pop := range pm.Popups {
		alphaRatio := 1.0 - float64(pop.Life)/float64(pop.MaxLife)
		if alphaRatio < 0 {
			alphaRatio = 0
		}
		// Sombra preta para legibilidade
		ebitenutil.DebugPrintAt(screen, pop.Text, int(pop.X)+1, int(pop.Y)+1)
		ebitenutil.DebugPrintAt(screen, pop.Text, int(pop.X), int(pop.Y))
	}
}

func (pm *ProjectileManager) Reset() {
	pm.Projectiles = pm.Projectiles[:0]
	pm.Popups = pm.Popups[:0]
	pm.Particles = pm.Particles[:0]
}
