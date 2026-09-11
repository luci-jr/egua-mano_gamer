package entities

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ObstacleType int

const (
	TypeGround ObstacleType = 0 // Paneiro de açaí / Cesto artesanal
	TypeAir    ObstacleType = 1 // Urubu / Gaivota / Arara (aéreo - desvia agachando)
	TypeTall   ObstacleType = 2 // Paneiros empilhados (alto - incentiva pulo duplo)
	TypePuddle ObstacleType = 3 // Poça da chuva das 4h (rasteiro e rápido)
	TypeJacare ObstacleType = 4 // Jacaré-Açu Amazônico com bocarra, dentes e escamas (substitui barris)
	TypeSnake  ObstacleType = 5 // Cobra-Coral / Sucuri ondulando com língua bífida (rasteira)
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
	case TypeJacare:
		w = 32.0
		h = 17.0
		x = obs.X
		y = obs.groundY - h
	case TypeSnake:
		w = 28.0
		h = 12.0
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
	case TypeJacare:
		obs.drawJacare(screen, ticks, stage)
	case TypeSnake:
		obs.drawSnake(screen, ticks, stage)
	}
}

func (obs *Obstacle) drawGround(screen *ebiten.Image, ticks int, stage int) {
	obsRealY := obs.groundY - 24.0

	switch stage {
	case 2:
		// Na Estação das Docas, ao invés de barris genéricos, um Jacaré descansando no cais!
		obs.drawJacare(screen, ticks, stage)

	case 3:
		// Cesto artesanal com castanhas e cupuaçus (Theatro da Paz)
		ebitenutil.DrawRect(screen, obs.X+2, obsRealY+10, 18, 14, color.RGBA{R: 120, G: 85, B: 50, A: 255})
		ebitenutil.DrawRect(screen, obs.X+4, obsRealY+5, 7, 7, color.RGBA{R: 235, G: 155, B: 30, A: 255})
		ebitenutil.DrawRect(screen, obs.X+6, obsRealY+4, 5, 4, color.RGBA{R: 215, G: 60, B: 35, A: 255})
		ebitenutil.DrawRect(screen, obs.X+11, obsRealY+6, 7, 7, color.RGBA{R: 240, G: 175, B: 35, A: 255})

	default:
		// Paneiro tradicional de açaí (Ver-o-Peso)
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
		// Gaivota do cais da Baía do Guajará
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
		// Arara / Maritaca Verde da Praça da República
		ebitenutil.DrawRect(screen, obs.X+6, birdY+4, 16, 7, color.RGBA{R: 35, G: 165, B: 60, A: 255})
		ebitenutil.DrawRect(screen, obs.X+1, birdY+3, 6, 6, color.RGBA{R: 50, G: 190, B: 75, A: 255})
		ebitenutil.DrawRect(screen, obs.X-2, birdY+5, 4, 3, color.RGBA{R: 245, G: 210, B: 40, A: 255})

		if (ticks/7)%2 == 0 {
			ebitenutil.DrawRect(screen, obs.X+9, birdY-5, 14, 9, color.RGBA{R: 25, G: 130, B: 45, A: 255})
		} else {
			ebitenutil.DrawRect(screen, obs.X+9, birdY+9, 14, 8, color.RGBA{R: 25, G: 130, B: 45, A: 255})
		}

	default:
		// Urubu clássico do Ver-o-Peso
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
		// Jacaré grande com bocarra aberta empinada na beira do cais
		obs.drawJacare(screen, ticks, stage)

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

// drawJacare desenha um autêntico Jacaré-Açu da Amazônia com escamas, cristas, dentes afiados e bocarra
func (obs *Obstacle) drawJacare(screen *ebiten.Image, ticks int, stage int) {
	obsRealY := obs.groundY - 17.0

	darkGreen := color.RGBA{R: 32, G: 68, B: 30, A: 255}
	midGreen := color.RGBA{R: 52, G: 104, B: 48, A: 255}
	scuteGreen := color.RGBA{R: 18, G: 46, B: 18, A: 255}
	bellyYellow := color.RGBA{R: 135, G: 155, B: 75, A: 255}
	mouthRed := color.RGBA{R: 195, G: 45, B: 55, A: 255}
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	eyeYellow := color.RGBA{R: 250, G: 220, B: 35, A: 255}
	black := color.RGBA{R: 12, G: 18, B: 12, A: 255}

	// 1. Cauda longa escamosa com cristas pontudas (à direita)
	ebitenutil.DrawRect(screen, obs.X+20, obsRealY+8, 11, 6, darkGreen)
	ebitenutil.DrawRect(screen, obs.X+26, obsRealY+10, 5, 3, midGreen)
	ebitenutil.DrawRect(screen, obs.X+22, obsRealY+6, 2, 3, scuteGreen)
	ebitenutil.DrawRect(screen, obs.X+25, obsRealY+6, 2, 3, scuteGreen)
	ebitenutil.DrawRect(screen, obs.X+28, obsRealY+8, 2, 2, scuteGreen)

	// 2. Corpo do réptil
	ebitenutil.DrawRect(screen, obs.X+9, obsRealY+6, 13, 8, darkGreen)
	ebitenutil.DrawRect(screen, obs.X+10, obsRealY+5, 11, 2, midGreen)
	ebitenutil.DrawRect(screen, obs.X+11, obsRealY+13, 10, 2, bellyYellow)

	// Cristas dorsais serrilhadas
	ebitenutil.DrawRect(screen, obs.X+11, obsRealY+3, 2, 3, scuteGreen)
	ebitenutil.DrawRect(screen, obs.X+14, obsRealY+3, 2, 3, scuteGreen)
	ebitenutil.DrawRect(screen, obs.X+17, obsRealY+3, 2, 3, scuteGreen)

	// Patas apoiadas no chão
	ebitenutil.DrawRect(screen, obs.X+10, obsRealY+14, 3, 3, scuteGreen)
	ebitenutil.DrawRect(screen, obs.X+19, obsRealY+14, 3, 3, scuteGreen)

	// 3. Olho amarelo reptiliano vigilante
	ebitenutil.DrawRect(screen, obs.X+7, obsRealY+2, 3, 3, eyeYellow)
	if (ticks/45)%6 != 0 {
		ebitenutil.DrawRect(screen, obs.X+8, obsRealY+3, 1, 2, black) // Pupila vertical
	}

	// 4. Cabeça e Bocarra (virada para a esquerda ameaçando a Onça)
	// Mandíbula superior
	ebitenutil.DrawRect(screen, obs.X, obsRealY+5, 9, 3, darkGreen)
	ebitenutil.DrawRect(screen, obs.X+1, obsRealY+4, 7, 1, midGreen)
	ebitenutil.DrawRect(screen, obs.X, obsRealY+5, 1, 1, black) // Narina

	// Abertura dinâmica da bocarra (abre e fecha ligeiramente)
	mouthGap := 2.0
	if (ticks/16)%2 == 0 {
		mouthGap = 3.0
	}
	ebitenutil.DrawRect(screen, obs.X+1, obsRealY+8, 7, mouthGap, mouthRed)

	// Dentes pontiagudos brancos afiados
	ebitenutil.DrawRect(screen, obs.X+1, obsRealY+7, 1, 2, white)
	ebitenutil.DrawRect(screen, obs.X+4, obsRealY+7, 1, 2, white)
	ebitenutil.DrawRect(screen, obs.X+7, obsRealY+7, 1, 2, white)

	// Mandíbula inferior
	lowerY := obsRealY + 8.0 + mouthGap
	ebitenutil.DrawRect(screen, obs.X+1, lowerY, 8, 3, darkGreen)
	ebitenutil.DrawRect(screen, obs.X+2, lowerY-1, 1, 2, white)
	ebitenutil.DrawRect(screen, obs.X+5, lowerY-1, 1, 2, white)
}

// drawSnake desenha a temível Cobra-Coral Amazônica rastejando com corpo ondulante e língua bífida
func (obs *Obstacle) drawSnake(screen *ebiten.Image, ticks int, stage int) {
	obsRealY := obs.groundY - 12.0

	redCoral := color.RGBA{R: 235, G: 38, B: 38, A: 255}
	black := color.RGBA{R: 20, G: 20, B: 24, A: 255}
	yellow := color.RGBA{R: 252, G: 218, B: 42, A: 255}
	white := color.RGBA{R: 245, G: 245, B: 250, A: 255}
	eyeYellow := color.RGBA{R: 255, G: 240, B: 60, A: 255}
	tongueRed := color.RGBA{R: 255, G: 35, B: 55, A: 255}

	// Animação senoidal de rastejo nos segmentos
	wave := int((ticks / 5) % 4)
	wY := func(segment int) float64 {
		offset := (wave + segment) % 4
		if offset == 1 || offset == 2 {
			return 1.0
		}
		return 0.0
	}

	// 1. Cabeça triangular da serpente (à esquerda, encarando a onça)
	headY := obsRealY + wY(0) + 1.0
	ebitenutil.DrawRect(screen, obs.X+1, headY+1, 5, 5, redCoral)
	ebitenutil.DrawRect(screen, obs.X, headY+2, 2, 3, redCoral)
	ebitenutil.DrawRect(screen, obs.X+2, headY, 3, 2, black)

	// Olho brilhante
	ebitenutil.DrawRect(screen, obs.X+2, headY+1, 2, 2, eyeYellow)
	ebitenutil.DrawRect(screen, obs.X+2, headY+2, 1, 1, black)

	// Língua bífida vermelha saindo e tremulando veloz
	if (ticks/6)%2 == 0 {
		ebitenutil.DrawRect(screen, obs.X-3, headY+3, 3, 1, tongueRed)
		ebitenutil.DrawRect(screen, obs.X-4, headY+2, 1, 1, tongueRed)
		ebitenutil.DrawRect(screen, obs.X-4, headY+4, 1, 1, tongueRed)
	}

	// 2. Anéis característicos da Cobra-Coral Amazônica (Vermelho - Preto - Amarelo - Preto)
	rings := []struct {
		dx  float64
		col color.RGBA
	}{
		{5, black},
		{7, yellow},
		{9, black},
		{11, redCoral},
		{14, black},
		{16, yellow},
		{18, black},
		{20, redCoral},
		{23, black},
		{25, yellow},
	}

	for i, r := range rings {
		segY := obsRealY + wY(i+1) + 2.0
		ebitenutil.DrawRect(screen, obs.X+r.dx, segY, 3, 5, r.col)
		ebitenutil.DrawRect(screen, obs.X+r.dx, segY+4, 3, 1, white) // Ventre esbranquiçado
	}

	// Cauda afilada
	tailY := obsRealY + wY(7) + 3.0
	ebitenutil.DrawRect(screen, obs.X+27, tailY, 2, 3, black)
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
			{X: screenWidth + 30, Type: TypeJacare, screenWidth: screenWidth, groundY: groundY},
			{X: screenWidth + 185, Type: TypeAir, screenWidth: screenWidth, groundY: groundY},
			{X: screenWidth + 340, Type: TypeSnake, screenWidth: screenWidth, groundY: groundY},
		},
	}
	return m
}

func (m *ObstacleManager) Update(speed float64) int {
	passedCount := 0
	for _, obs := range m.Obstacles {
		obs.X -= speed
		if obs.X < -40 {
			furthestX := m.screenWidth
			for _, other := range m.Obstacles {
				if other != obs && other.X > furthestX {
					furthestX = other.X
				}
			}
			obs.X = furthestX + 135.0 + float64(rand.Intn(55))
			// Sorteia entre os 6 tipos proceduralmente
			obs.Type = ObstacleType(rand.Intn(6))
			obs.collided = false
			passedCount++
		}
	}
	return passedCount
}

func (m *ObstacleManager) CheckCollision(oncaX, oncaY, oncaW, oncaH float64) (bool, ObstacleType) {
	for _, obs := range m.Obstacles {
		if !obs.collided && obs.CheckCollision(oncaX, oncaY, oncaW, oncaH) {
			obs.collided = true
			return true, obs.Type
		}
	}
	return false, TypeGround
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
	types := []ObstacleType{TypeJacare, TypeAir, TypeSnake}
	for i, obs := range m.Obstacles {
		obs.X = m.screenWidth + 25.0 + float64(i)*spacing
		obs.Type = types[i%len(types)]
		obs.collided = false
	}
}
