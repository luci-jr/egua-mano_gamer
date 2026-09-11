package entities

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ObstacleType int

const (
	TypeGround ObstacleType = 0 // Paneiro de açaí / Caixote / Cesto padrão
	TypeAir    ObstacleType = 1 // Urubu / Gaivota / Arara (aéreo - desvia agachando)
	TypeTall   ObstacleType = 2 // Paneiros empilhados / Caixas duplas (alto - incentiva pulo duplo)
	TypePuddle ObstacleType = 3 // Poça da chuva das 4h / Casca escorregadia (rasteiro e rápido)
)

type Obstacle struct {
	X           float64
	Type        ObstacleType
	screenWidth float64
	groundY     float64
	collided    bool
}

func NewObstacle(screenWidth, groundY float64) *Obstacle {
	return &Obstacle{
		X:           screenWidth,
		Type:        TypeGround,
		screenWidth: screenWidth,
		groundY:     groundY,
		collided:    false,
	}
}

func (obs *Obstacle) GetBounds() (x, y, w, h float64) {
	switch obs.Type {
	case TypeGround:
		w = 22.0
		h = 24.0
		x = obs.X
		y = obs.groundY - h
	case TypeAir:
		w = 28.0
		h = 16.0
		x = obs.X
		y = obs.groundY - 34.0
	case TypeTall:
		w = 22.0
		h = 36.0
		x = obs.X
		y = obs.groundY - h
	case TypePuddle:
		w = 26.0
		h = 8.0
		x = obs.X
		y = obs.groundY - h
	default:
		w = 22.0
		h = 24.0
		x = obs.X
		y = obs.groundY - h
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
	switch obs.Type {
	case TypeGround:
		obs.drawGround(screen, ticks, stage)
	case TypeAir:
		obs.drawAir(screen, ticks, stage)
	case TypeTall:
		obs.drawTall(screen, ticks, stage)
	case TypePuddle:
		obs.drawPuddle(screen, ticks, stage)
	}
}

func (obs *Obstacle) drawGround(screen *ebiten.Image, ticks int, stage int) {
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
}

func (obs *Obstacle) drawAir(screen *ebiten.Image, ticks int, stage int) {
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

func (obs *Obstacle) drawTall(screen *ebiten.Image, ticks int, stage int) {
	obsRealY := obs.groundY - 36.0

	switch stage {
	case 2:
		// Caixotes portuários empilhados (Docas)
		woodDark := color.RGBA{R: 100, G: 55, B: 28, A: 255}
		woodLight := color.RGBA{R: 140, G: 85, B: 45, A: 255}
		metalGrey := color.RGBA{R: 180, G: 185, B: 190, A: 255}

		// Caixa inferior
		ebitenutil.DrawRect(screen, obs.X+2, obsRealY+18, 18, 18, woodDark)
		ebitenutil.DrawRect(screen, obs.X+4, obsRealY+20, 14, 14, woodLight)
		ebitenutil.DrawRect(screen, obs.X+1, obsRealY+26, 20, 2, metalGrey)

		// Caixa superior
		ebitenutil.DrawRect(screen, obs.X+3, obsRealY+2, 16, 16, woodDark)
		ebitenutil.DrawRect(screen, obs.X+5, obsRealY+4, 12, 12, woodLight)
		ebitenutil.DrawRect(screen, obs.X+2, obsRealY+9, 18, 2, metalGrey)

	case 3:
		// Cestos nobres empilhados (Theatro da Paz)
		gold := color.RGBA{R: 215, G: 160, B: 40, A: 255}
		brown := color.RGBA{R: 120, G: 75, B: 40, A: 255}
		redNut := color.RGBA{R: 190, G: 50, B: 30, A: 255}

		ebitenutil.DrawRect(screen, obs.X+2, obsRealY+17, 18, 19, brown)
		ebitenutil.DrawRect(screen, obs.X+4, obsRealY+19, 14, 15, gold)
		ebitenutil.DrawRect(screen, obs.X+3, obsRealY+2, 16, 15, brown)
		ebitenutil.DrawRect(screen, obs.X+5, obsRealY+4, 12, 11, gold)
		ebitenutil.DrawRect(screen, obs.X+6, obsRealY, 8, 4, redNut)

	default:
		// Pilha dupla de paneiros de açaí (Ver-o-Peso)
		basketStraw := color.RGBA{R: 185, G: 135, B: 75, A: 255}
		basketDark := color.RGBA{R: 135, G: 90, B: 45, A: 255}
		acaiPurple := color.RGBA{R: 45, G: 15, B: 48, A: 255}
		acaiLight := color.RGBA{R: 85, G: 30, B: 85, A: 255}
		palmGreen := color.RGBA{R: 40, G: 145, B: 50, A: 255}

		// Paneiro inferior
		ebitenutil.DrawRect(screen, obs.X+2, obsRealY+18, 18, 18, basketStraw)
		ebitenutil.DrawRect(screen, obs.X+4, obsRealY+20, 14, 14, basketDark)
		ebitenutil.DrawRect(screen, obs.X+1, obsRealY+26, 20, 2, basketDark)

		// Paneiro superior
		ebitenutil.DrawRect(screen, obs.X+3, obsRealY+3, 16, 15, basketStraw)
		ebitenutil.DrawRect(screen, obs.X+5, obsRealY+5, 12, 11, basketDark)
		ebitenutil.DrawRect(screen, obs.X+2, obsRealY+10, 18, 2, basketDark)

		// Açaí derramando no topo
		ebitenutil.DrawRect(screen, obs.X+4, obsRealY+1, 14, 5, acaiPurple)
		ebitenutil.DrawRect(screen, obs.X+6, obsRealY, 10, 3, acaiLight)

		// Folha de palmeira saindo do topo
		ebitenutil.DrawRect(screen, obs.X+16, obsRealY-4, 3, 7, palmGreen)
		ebitenutil.DrawRect(screen, obs.X+18, obsRealY-6, 4, 3, palmGreen)
	}
}

func (obs *Obstacle) drawPuddle(screen *ebiten.Image, ticks int, stage int) {
	obsRealY := obs.groundY - 8.0

	switch stage {
	case 2:
		// Poça salina do cais das Docas
		waterDeep := color.RGBA{R: 20, G: 70, B: 90, A: 230}
		waterMid := color.RGBA{R: 45, G: 125, B: 145, A: 230}
		waterFoam := color.RGBA{R: 210, G: 240, B: 245, A: 255}

		ebitenutil.DrawRect(screen, obs.X+2, obsRealY+2, 22, 5, waterDeep)
		ebitenutil.DrawRect(screen, obs.X+5, obsRealY+3, 16, 3, waterMid)
		if (ticks/12)%2 == 0 {
			ebitenutil.DrawRect(screen, obs.X+8, obsRealY+3, 5, 1, waterFoam)
		} else {
			ebitenutil.DrawRect(screen, obs.X+14, obsRealY+4, 5, 1, waterFoam)
		}

	case 3:
		// Poça cristalina sobre pedras portuguesas (Theatro da Paz)
		waterDeep := color.RGBA{R: 40, G: 80, B: 130, A: 230}
		waterLight := color.RGBA{R: 90, G: 155, B: 220, A: 230}
		waterGlint := color.RGBA{R: 240, G: 245, B: 255, A: 255}

		ebitenutil.DrawRect(screen, obs.X+3, obsRealY+2, 20, 5, waterDeep)
		ebitenutil.DrawRect(screen, obs.X+6, obsRealY+3, 14, 3, waterLight)
		ebitenutil.DrawRect(screen, obs.X+9, obsRealY+3, 4, 1, waterGlint)

	default:
		// Poça da chuva das 4h da tarde de Belém (Ver-o-Peso)
		rainDeep := color.RGBA{R: 35, G: 75, B: 145, A: 240}
		rainMid := color.RGBA{R: 70, G: 140, B: 215, A: 240}
		rainReflect := color.RGBA{R: 200, G: 230, B: 255, A: 255}

		ebitenutil.DrawRect(screen, obs.X+2, obsRealY+2, 22, 5, rainDeep)
		ebitenutil.DrawRect(screen, obs.X+5, obsRealY+3, 16, 3, rainMid)
		ebitenutil.DrawRect(screen, obs.X+1, obsRealY+4, 24, 2, rainDeep)

		if (ticks/10)%2 == 0 {
			ebitenutil.DrawRect(screen, obs.X+7, obsRealY+3, 6, 1, rainReflect)
			ebitenutil.DrawRect(screen, obs.X+16, obsRealY+4, 3, 1, rainReflect)
		} else {
			ebitenutil.DrawRect(screen, obs.X+10, obsRealY+4, 7, 1, rainReflect)
			ebitenutil.DrawRect(screen, obs.X+4, obsRealY+3, 3, 1, rainReflect)
		}
	}
}

// ObstacleManager gerencia múltiplos obstáculos simultâneos na tela com espaçamento justo e seguro
type ObstacleManager struct {
	Obstacles   []*Obstacle
	screenWidth float64
	groundY     float64
}

func NewObstacleManager(screenWidth, groundY float64) *ObstacleManager {
	m := &ObstacleManager{
		screenWidth: screenWidth,
		groundY:     groundY,
		Obstacles: []*Obstacle{
			{X: screenWidth + 30, Type: TypeGround, screenWidth: screenWidth, groundY: groundY},
			{X: screenWidth + 185, Type: TypeAir, screenWidth: screenWidth, groundY: groundY},
			{X: screenWidth + 340, Type: TypePuddle, screenWidth: screenWidth, groundY: groundY},
		},
	}
	return m
}

func (m *ObstacleManager) Update(speed float64) int {
	passedCount := 0
	for _, obs := range m.Obstacles {
		obs.X -= speed
		if obs.X < -40 {
			// Localiza a maior coordenada X entre os outros obstáculos para nunca sobrepor
			furthestX := m.screenWidth
			for _, other := range m.Obstacles {
				if other != obs && other.X > furthestX {
					furthestX = other.X
				}
			}
			// Garante distância mínima de 135px e variação de até 55px
			obs.X = furthestX + 135.0 + float64(rand.Intn(55))
			obs.Type = ObstacleType(rand.Intn(4))
			obs.collided = false
			passedCount++
		}
	}
	return passedCount
}

func (m *ObstacleManager) CheckCollision(oncaX, oncaY, oncaW, oncaH float64) bool {
	for _, obs := range m.Obstacles {
		if !obs.collided && obs.CheckCollision(oncaX, oncaY, oncaW, oncaH) {
			obs.collided = true
			return true
		}
	}
	return false
}

func (m *ObstacleManager) Draw(screen *ebiten.Image, ticks int, stage int) {
	for _, obs := range m.Obstacles {
		if obs.X > -40 && obs.X < m.screenWidth+40 {
			obs.Draw(screen, ticks, stage)
		}
	}
}

func (m *ObstacleManager) Reset() {
	spacing := 160.0
	for i, obs := range m.Obstacles {
		obs.X = m.screenWidth + 25.0 + float64(i)*spacing
		obs.Type = ObstacleType(i % 4)
		obs.collided = false
	}
}
