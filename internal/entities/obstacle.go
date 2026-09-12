package entities

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ObstacleType int

const (
	TypeGround ObstacleType = 0 // Paneiro de açaí tradicional / Cesto artesanal
	TypeAir    ObstacleType = 1 // Urubu / Gaivota / Arara (aéreo - desvia agachando)
	TypeJacare ObstacleType = 2 // Jacaré-Açu Amazônico com bocarra, dentes e escamas
	TypeSnake  ObstacleType = 3 // Cobra-Coral Amazônica ondulando com língua bífida
	TypeMudPit ObstacleType = 4 // Poço de Lama Movediça / Areia Movediça (Estilo Pitfall)
	TypeBench  ObstacleType = 5 // Banco de Praça Colonial de Belém (Madeira de lei & ferro) - Plataforma segura!
)

type Obstacle struct {
	X           float64
	Type        ObstacleType
	screenWidth float64
	groundY     float64
	Collided    bool
	Defeated    bool
	DefeatTicks int
}

func NewObstacle(screenWidth, groundY float64) *Obstacle {
	return &Obstacle{
		X:           screenWidth,
		Type:        TypeGround,
		screenWidth: screenWidth,
		groundY:     groundY,
		Collided:    false,
		Defeated:    false,
		DefeatTicks: 0,
	}
}

func (obs *Obstacle) GetBounds() (x, y, w, h float64) {
	if obs.Defeated {
		return 0, 0, 0, 0
	}
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
	case TypeMudPit:
		w = 40.0
		h = 10.0
		x = obs.X
		y = obs.groundY - 2.0
	case TypeBench:
		w = 38.0
		h = 16.0
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

func (obs *Obstacle) CheckCollision(playerX, playerY, playerW, playerH float64) bool {
	obsX, obsY, obsW, obsH := obs.GetBounds()
	overlapX := obsX < playerX+playerW && obsX+obsW > playerX
	overlapY := obsY < playerY+playerH && obsY+obsH > playerY
	return overlapX && overlapY
}

func (obs *Obstacle) IsPlayerOnTop(playerX, playerW float64) bool {
	if obs.Defeated {
		return false
	}
	obsX, _, obsW, _ := obs.GetBounds()
	playerCenterX := playerX + playerW/2.0
	return playerCenterX >= obsX-4.0 && playerCenterX <= obsX+obsW+4.0
}

func (obs *Obstacle) Draw(screen *ebiten.Image, ticks int, stage int) {
	if obs.Defeated {
		if (obs.DefeatTicks/3)%2 != 0 {
			return
		}
	}
	switch obs.Type {
	case TypeGround:
		obs.drawGround(screen, ticks, stage)
	case TypeAir:
		obs.drawAir(screen, ticks, stage)
	case TypeJacare:
		obs.drawJacare(screen, ticks, stage)
	case TypeSnake:
		obs.drawSnake(screen, ticks, stage)
	case TypeMudPit:
		obs.drawMudPit(screen, ticks)
	case TypeBench:
		obs.drawBench(screen, ticks, stage)
	}
}

func (obs *Obstacle) drawGround(screen *ebiten.Image, ticks int, stage int) {
	obsRealY := obs.groundY - 24.0

	// Paleta de Cores: Palha trançada, Açaí roxo fresco e folhas da Amazônia
	basketStraw := color.RGBA{R: 185, G: 135, B: 75, A: 255}
	basketDark := color.RGBA{R: 135, G: 90, B: 45, A: 255}
	acaiPurple := color.RGBA{R: 45, G: 15, B: 48, A: 255}
	acaiLight := color.RGBA{R: 95, G: 30, B: 100, A: 255}
	palmGreen := color.RGBA{R: 40, G: 155, B: 50, A: 255}
	energyGold := color.RGBA{R: 255, G: 225, B: 70, A: 255}

	// 1. Cesto de Palha do Paneiro de Açaí
	ebitenutil.DrawRect(screen, obs.X+2, obsRealY+7, 18, 17, basketStraw)
	ebitenutil.DrawRect(screen, obs.X+4, obsRealY+9, 14, 13, basketDark)

	for weave := 0; weave < 3; weave++ {
		wy := obsRealY + 9 + float64(weave*5)
		ebitenutil.DrawRect(screen, obs.X+1, wy, 20, 2, basketDark)
	}

	// 2. Frutos frescos de Açaí empilhados no cesto
	ebitenutil.DrawRect(screen, obs.X+1, obsRealY+2, 20, 8, acaiPurple)
	ebitenutil.DrawRect(screen, obs.X+3, obsRealY+1, 16, 6, acaiLight)
	ebitenutil.DrawRect(screen, obs.X+5, obsRealY, 12, 3, acaiPurple)

	// Folhas verdes frescas de palmeira de açaí
	ebitenutil.DrawRect(screen, obs.X+17, obsRealY-3, 3, 7, palmGreen)
	ebitenutil.DrawRect(screen, obs.X+19, obsRealY-5, 4, 3, palmGreen)
	ebitenutil.DrawRect(screen, obs.X+15, obsRealY-1, 3, 2, palmGreen)

	// Bagos de açaí caídos na base
	ebitenutil.DrawRect(screen, obs.X-3, obsRealY+21, 3, 3, acaiPurple)
	ebitenutil.DrawRect(screen, obs.X-1, obsRealY+22, 2, 2, acaiLight)

	// 3. Brilho sutil de energia vital no açaí (sem poluição de corações flutuantes)
	if !obs.Collided {
		if (ticks/10)%2 == 0 {
			ebitenutil.DrawRect(screen, obs.X+4, obsRealY-2, 2, 2, energyGold)
			ebitenutil.DrawRect(screen, obs.X+16, obsRealY-1, 2, 2, palmGreen)
		}
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

	// Abertura dinâmica da bocarra
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
			{X: screenWidth + 30, Type: TypeGround, screenWidth: screenWidth, groundY: groundY},
			{X: screenWidth + 185, Type: TypeAir, screenWidth: screenWidth, groundY: groundY},
			{X: screenWidth + 340, Type: TypeJacare, screenWidth: screenWidth, groundY: groundY},
		},
	}
	return m
}

func (obs *Obstacle) drawMudPit(screen *ebiten.Image, ticks int) {
	cMudDark := color.RGBA{R: 35, G: 22, B: 14, A: 255}
	cMudMid := color.RGBA{R: 62, G: 42, B: 26, A: 255}
	cMudLight := color.RGBA{R: 105, G: 72, B: 42, A: 255}
	cMoss := color.RGBA{R: 42, G: 125, B: 45, A: 255}

	pitW := 40.0
	pitH := 10.0
	py := obs.groundY - 2.0

	// Poço escavado no solo (lamaçal movediço estilo Pitfall)
	ebitenutil.DrawRect(screen, obs.X, py, pitW, pitH, cMudDark)
	ebitenutil.DrawRect(screen, obs.X+2, py+2, pitW-4, pitH-2, cMudMid)

	// Bordas com vegetação pantanosa
	ebitenutil.DrawRect(screen, obs.X-2, py-2, 4, 4, cMoss)
	ebitenutil.DrawRect(screen, obs.X+pitW-2, py-2, 4, 4, cMoss)

	// Bolhas de lodo borbulhando no pântano
	bubblePhase := (ticks / 9) % 3
	switch bubblePhase {
	case 0:
		ebitenutil.DrawRect(screen, obs.X+8, py+2, 4, 3, cMudLight)
		ebitenutil.DrawRect(screen, obs.X+9, py+1, 2, 1, cMudLight)
	case 1:
		ebitenutil.DrawRect(screen, obs.X+20, py+3, 5, 4, cMudLight)
		ebitenutil.DrawRect(screen, obs.X+21, py+2, 3, 1, cMudLight)
	case 2:
		ebitenutil.DrawRect(screen, obs.X+30, py+2, 4, 3, cMudLight)
	}
}

func (obs *Obstacle) drawBench(screen *ebiten.Image, ticks int, stage int) {
	benchY := obs.groundY - 16.0
	bx := obs.X

	// Paleta de Cores: Ferro fundido colonial de Belém + Madeira de Lei
	cIron := color.RGBA{R: 28, G: 55, B: 38, A: 255}       // Ferro ornamental verde escuro colonial
	cIronDark := color.RGBA{R: 16, G: 32, B: 22, A: 255}   // Sombra do ferro
	cWood := color.RGBA{R: 165, G: 100, B: 45, A: 255}     // Madeira de lei castanho dourado
	cWoodDark := color.RGBA{R: 115, G: 65, B: 25, A: 255}  // Sombra das ripas
	cWoodLight := color.RGBA{R: 205, G: 135, B: 70, A: 255} // Brilho superior da madeira
	cShadow := color.RGBA{R: 10, G: 15, B: 20, A: 120}     // Sombra de contato no solo

	// 1. Sombra suave sob o banco no piso
	ebitenutil.DrawRect(screen, bx+2, obs.groundY-2, 34, 2, cShadow)

	// 2. Pés e braços de ferro trabalhado (estilo Praça da República / Theatro da Paz)
	ebitenutil.DrawRect(screen, bx+3, benchY+7, 3, 9, cIron)
	ebitenutil.DrawRect(screen, bx+32, benchY+7, 3, 9, cIron)
	ebitenutil.DrawRect(screen, bx+1, benchY+14, 6, 2, cIronDark)
	ebitenutil.DrawRect(screen, bx+30, benchY+14, 6, 2, cIronDark)

	// Apoio de braço curvo nas laterais
	ebitenutil.DrawRect(screen, bx+2, benchY+3, 3, 5, cIron)
	ebitenutil.DrawRect(screen, bx+1, benchY+2, 5, 2, cIron)
	ebitenutil.DrawRect(screen, bx+33, benchY+3, 3, 5, cIron)
	ebitenutil.DrawRect(screen, bx+32, benchY+2, 5, 2, cIron)

	// Haste de suporte vertical do encosto
	ebitenutil.DrawRect(screen, bx+4, benchY-1, 2, 7, cIronDark)
	ebitenutil.DrawRect(screen, bx+32, benchY-1, 2, 7, cIronDark)

	// 3. Ripas de Madeira do Assento (plataforma sólida de aterrissagem)
	ebitenutil.DrawRect(screen, bx+2, benchY+6, 34, 2, cWoodLight)
	ebitenutil.DrawRect(screen, bx+2, benchY+8, 34, 2, cWood)
	ebitenutil.DrawRect(screen, bx+2, benchY+10, 34, 2, cWoodDark)

	// 4. Ripas de Madeira do Encosto (estilo praça colonial)
	ebitenutil.DrawRect(screen, bx+4, benchY-1, 30, 2, cWoodLight)
	ebitenutil.DrawRect(screen, bx+4, benchY+2, 30, 2, cWood)
}

func (m *ObstacleManager) Update(speed float64) int {
	passedCount := 0
	validTypes := []ObstacleType{TypeGround, TypeJacare, TypeBench, TypeSnake, TypeAir}

	for _, obs := range m.Obstacles {
		if obs.Defeated {
			obs.DefeatTicks--
			obs.X -= speed * 0.4
			if obs.DefeatTicks <= 0 || (speed > 0 && obs.X < -40) {
				furthestX := m.screenWidth
				for _, other := range m.Obstacles {
					if other != obs && other.X > furthestX {
						furthestX = other.X
					}
				}
				obs.X = furthestX + 140.0 + float64(rand.Intn(50))
				obs.Type = validTypes[rand.Intn(len(validTypes))]
				obs.Collided = false
				obs.Defeated = false
				obs.DefeatTicks = 0
				passedCount++
			}
			continue
		}

		// 1. Deslocamento com o scroll da câmera (todos os objetos do mundo acompanham o cenário)
		obs.X -= speed

		// 2. Movimentação autônoma suave dos predadores no chão/ar em direção ao herói
		switch obs.Type {
		case TypeJacare:
			// Jacaré-Açu rasteja sutilmente para a frente no chão
			obs.X -= 0.35
		case TypeSnake:
			// Cobra-Coral serpenteia com agilidade natural
			obs.X -= 0.50
		case TypeAir:
			// Urubu / ave sobrevoa os céus
			obs.X -= 0.70
		}

		if obs.X < -40 {
			furthestX := m.screenWidth
			for _, other := range m.Obstacles {
				if other != obs && other.X > furthestX {
					furthestX = other.X
				}
			}
			obs.X = furthestX + 140.0 + float64(rand.Intn(50))
			obs.Type = validTypes[rand.Intn(len(validTypes))]
			obs.Collided = false
			obs.Defeated = false
			obs.DefeatTicks = 0
			passedCount++
		} else if speed < 0 && obs.X > m.screenWidth+250.0 {
			obs.X = m.screenWidth + 250.0
		}
	}
	return passedCount
}

// CheckPlatformSupport verifica se o jogador está aterrissando ou em pé sobre o topo de um obstáculo sólido (Paneiro de Açaí ou Banco de Praça)
func (m *ObstacleManager) CheckPlatformSupport(playerX, playerY, playerW, playerH, playerVY, groundY float64) (bool, float64, *Obstacle) {
	var playerBottom float64
	if playerY > 50.0 {
		playerBottom = playerY + playerH
	} else {
		playerBottom = groundY + playerY
	}
	playerCenterX := playerX + playerW/2.0

	for _, obs := range m.Obstacles {
		if obs.Defeated {
			continue
		}
		// Apenas TypeGround (paneiro) e TypeBench (banco de praça) funcionam como plataformas sólidas
		if obs.Type != TypeGround && obs.Type != TypeBench {
			continue
		}

		obsX, obsY, obsW, obsH := obs.GetBounds()
		obsTop := obsY // groundY - obsH

		if playerCenterX >= obsX-4.0 && playerCenterX <= obsX+obsW+4.0 {
			dist := playerBottom - obsTop
			if dist >= -7.0 && dist <= 8.0 && playerVY >= -1.0 {
				return true, -obsH, obs
			}
		}
	}
	return false, 0, nil
}

// CheckStomp verifica se o jogador pulou em cima de um inimigo (Jacaré, Cobra ou Ave Aérea), derrotando-o instantaneamente e quicando no ar
func (m *ObstacleManager) CheckStomp(playerX, playerY, playerW, playerH, playerVY, groundY float64) (bool, float64, float64, ObstacleType) {
	var playerBottom float64
	if playerY > 50.0 {
		playerBottom = playerY + playerH
	} else {
		playerBottom = groundY + playerY
	}

	for _, obs := range m.Obstacles {
		if obs.Defeated || obs.Collided {
			continue
		}
		// Apenas alvos inimigos podem ser pisados: Jacaré, Cobra e Ave Aérea
		if obs.Type != TypeJacare && obs.Type != TypeSnake && obs.Type != TypeAir {
			continue
		}

		ox, oy, ow, oh := obs.GetBounds()
		obsTop := oy
		obsBottom := oy + oh

		// Sobreposição horizontal com margem confortável para gameplay fluida
		overlapX := playerX+playerW > ox-3.0 && playerX < ox+ow+3.0

		// O jogador está no ar (saltando) e seus pés tocam/estão no corpo do inimigo vindo de cima
		isJumping := playerBottom < groundY-1.0 || playerVY != 0
		isAboveBase := playerBottom <= obsBottom+4.0 && playerBottom >= obsTop-16.0

		if overlapX && isJumping && isAboveBase {
			obs.Defeated = true
			obs.DefeatTicks = 26
			return true, ox + ow/2.0, oy + oh/2.0, obs.Type
		}
	}
	return false, 0, 0, TypeGround
}

// CheckCollision realiza checagem de dano ignorando plataformas seguras (Banco de Praça e Paneiro de Açaí)
func (m *ObstacleManager) CheckCollision(playerX, playerY, playerW, playerH, playerVY, groundY float64, currentPlatform *Obstacle) (bool, ObstacleType) {
	var playerBottom, playerTop float64
	if playerY > 50.0 {
		playerTop = playerY
		playerBottom = playerY + playerH
	} else {
		playerBottom = groundY + playerY
		playerTop = playerBottom - playerH
	}

	for _, obs := range m.Obstacles {
		if obs.Collided || obs.Defeated {
			continue
		}
		// Se o jogador está apoiado neste obstáculo, não recebe dano
		if obs == currentPlatform {
			continue
		}

		// Banco de Praça e Paneiro de Açaí NUNCA causam dano ao herói!
		// O açaí é sagrado e nutritivo: concede energia em vez de machucar!
		if obs.Type == TypeBench || obs.Type == TypeGround {
			continue
		}

		ox, oy, ow, oh := obs.GetBounds()
		obsTop := oy
		obsBottom := oy + oh

		// SALVAGUARDA ABSOLUTA DE PULO:
		// Se for bicho inimigo (Jacaré, Cobra, Ave aérea) e o herói estiver saltando no ar com sobreposição,
		// ele NUNCA deve tomar dano ou perder vida; se atingir o bicho no pulo, o bicho é derrotado!
		if obs.Type == TypeJacare || obs.Type == TypeSnake || obs.Type == TypeAir {
			overlapX := playerX+playerW > ox-3.0 && playerX < ox+ow+3.0
			isJumping := playerBottom < groundY-1.0 || playerVY != 0
			if overlapX && isJumping && playerBottom <= obsBottom+4.0 {
				obs.Defeated = true
				obs.DefeatTicks = 26
				continue
			}
		}

		// Se for obstáculo no chão e o herói colide vindo de cima em movimento descendente, evita dano
		if playerBottom <= obsTop+5.0 && playerVY >= -0.5 {
			continue
		}

		if obs.CheckCollision(playerX, playerTop, playerW, playerH) {
			obs.Collided = true
			return true, obs.Type
		}
	}
	return false, TypeGround
}

// CheckEnergyCollection verifica se o jogador tocou ou passou pelo Paneiro de Açaí para absorver energia vital
func (m *ObstacleManager) CheckEnergyCollection(playerX, playerY, playerW, playerH float64) (bool, *Obstacle) {
	realPlayerY := playerY
	if playerY <= 50.0 {
		realPlayerY = m.groundY + playerY - playerH
	}
	for _, obs := range m.Obstacles {
		if obs.Collided || obs.Defeated || obs.Type != TypeGround {
			continue
		}
		if obs.CheckCollision(playerX, realPlayerY, playerW, playerH) {
			obs.Collided = true // Marcado como consumido para não pontuar/curar continuamente
			return true, obs
		}
	}
	return false, nil
}

func (m *ObstacleManager) CheckProjectileHit(projX, projY, projW, projH float64) (bool, float64, float64, ObstacleType) {
	for _, obs := range m.Obstacles {
		// Banco de praça e Paneiro de Açaí não são alvos inimigos destruídos por tiros
		if obs.Type == TypeBench || obs.Type == TypeGround {
			continue
		}
		if !obs.Defeated && obs.X > -20 && obs.X < m.screenWidth+20 {
			ox, oy, ow, oh := obs.GetBounds()
			overlapX := projX < ox+ow && projX+projW > ox
			overlapY := projY < oy+oh && projY+projH > oy
			if overlapX && overlapY {
				obs.Defeated = true
				obs.DefeatTicks = 26
				return true, ox + ow/2.0, oy + oh/2.0, obs.Type
			}
		}
	}
	return false, 0, 0, TypeGround
}

func (m *ObstacleManager) Draw(screen *ebiten.Image, ticks int, stage int) {
	for _, obs := range m.Obstacles {
		if obs.X > -40 && obs.X < m.screenWidth+40 {
			obs.Draw(screen, ticks, stage)
		}
	}
}

func (m *ObstacleManager) Reset() {
	spacing := 165.0
	types := []ObstacleType{TypeGround, TypeJacare, TypeBench, TypeSnake, TypeAir}
	for i, obs := range m.Obstacles {
		obs.X = m.screenWidth + 25.0 + float64(i)*spacing
		obs.Type = types[i%len(types)]
		obs.Collided = false
		obs.Defeated = false
		obs.DefeatTicks = 0
	}
}
