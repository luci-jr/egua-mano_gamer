package scenery

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed bg_fase1_8bit.png
var bgFase1Bytes []byte

//go:embed bg_fase2_8bit.png
var bgFase2Bytes []byte

//go:embed bg_fase3_8bit.png
var bgFase3Bytes []byte

type Background struct {
	scrollOffset float64
	imgFase1     *ebiten.Image
	imgFase2     *ebiten.Image
	imgFase3     *ebiten.Image
}

func decodeBgImage(data []byte) *ebiten.Image {
	if len(data) == 0 {
		return nil
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil
	}
	return ebiten.NewImageFromImage(img)
}

func NewBackground() *Background {
	return &Background{
		scrollOffset: 0,
		imgFase1:     decodeBgImage(bgFase1Bytes),
		imgFase2:     decodeBgImage(bgFase2Bytes),
		imgFase3:     decodeBgImage(bgFase3Bytes),
	}
}

func (b *Background) Update(speed float64) {
	b.scrollOffset += speed * 0.7
	// Reseta apenas em múltiplos gigantescos para evitar qualquer salto visível nos diferentes períodos de parallax
	if b.scrollOffset >= 1000000.0 {
		b.scrollOffset -= 1000000.0
	} else if b.scrollOffset < 0 {
		b.scrollOffset += 1000000.0
	}
}

func (b *Background) Draw(screen *ebiten.Image, screenWidth, groundY float64, ticks int, stage int) {
	switch stage {
	case 2:
		b.drawDocas(screen, screenWidth, groundY, ticks)
	case 3:
		b.drawTheatroDaPaz(screen, screenWidth, groundY, ticks)
	default:
		b.drawVerOPeso(screen, screenWidth, groundY, ticks)
	}
}

// drawLoopingImage desenha o background panorâmico 8-bit em looping horizontal contínuo (parallax infinito)
func (b *Background) drawLoopingImage(screen *ebiten.Image, bg *ebiten.Image, screenWidth, groundY float64, parallaxSpeed float64) {
	if bg == nil {
		return
	}
	bounds := bg.Bounds()
	imgW := float64(bounds.Dx())
	imgH := float64(bounds.Dy())

	scaleX := screenWidth / imgW
	// Escala para cobrir toda a tela vertical (210 pixels)
	scaleY := (groundY + 55.0) / imgH

	offset := math.Mod(b.scrollOffset*parallaxSpeed, screenWidth)
	if offset < 0 {
		offset += screenWidth
	}

	op1 := &ebiten.DrawImageOptions{}
	op1.GeoM.Scale(scaleX, scaleY)
	op1.GeoM.Translate(-offset, 0)
	screen.DrawImage(bg, op1)

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Scale(scaleX, scaleY)
	op2.GeoM.Translate(screenWidth-offset, 0)
	screen.DrawImage(bg, op2)

	// Sombra sutil de transição e ambient occlusion entre a paisagem e o calçadão/piso
	ebitenutil.DrawRect(screen, 0, groundY-1, screenWidth, 1, color.RGBA{R: 15, G: 20, B: 28, A: 70})
}

// ==========================================
// FASE 1: MERCADO DO VER-O-PESO (8-BIT)
// ==========================================
func (b *Background) drawVerOPeso(screen *ebiten.Image, screenWidth, groundY float64, ticks int) {
	if b.imgFase1 != nil {
		// 1. Cenário panorâmico real 8-bit do Ver-o-Peso (Mercado de Ferro e Baía do Guajará ao pôr do sol)
		b.drawLoopingImage(screen, b.imgFase1, screenWidth, groundY, 0.45)

		// 2. Revoada característica de urubus e garças de Belém cortando o céu dourado
		cBird := color.RGBA{R: 28, G: 20, B: 24, A: 210}
		for bird := 0; bird < 4; bird++ {
			bx := math.Mod(float64(bird*90)+float64(ticks)*0.45, screenWidth+40) - 20
			by := 18.0 + math.Sin(float64(ticks+bird*35)*0.04)*8.0 + float64(bird*5)
			flap := (ticks + bird*7) / 7
			drawSkyBird(screen, bx, by, flap, cBird)
		}

		// 3. Postes coloniais vitorianos de ferro fundido do cais (Belle Époque) em primeiro plano
		lampIron := color.RGBA{R: 28, G: 32, B: 42, A: 255}
		lampGlow := color.RGBA{R: 255, G: 215, B: 90, A: 220}
		lampLightBeam := color.RGBA{R: 255, G: 220, B: 120, A: 25}

		for p := 0; p < 3; p++ {
			px := float64(p*130) + 50.0 - math.Mod(b.scrollOffset*0.75, 130.0)
			ebitenutil.DrawRect(screen, px, groundY-46, 2, 46, lampIron)
			ebitenutil.DrawRect(screen, px-2, groundY-48, 6, 2, lampIron)
			ebitenutil.DrawRect(screen, px-1, groundY-50, 4, 3, lampIron)
			ebitenutil.DrawRect(screen, px-1, groundY-47, 4, 3, lampGlow)
			ebitenutil.DrawRect(screen, px-5, groundY-44, 12, 44, lampLightBeam)
		}

		// 4. Chão: Cais de pedra de cantaria histórica do Ver-o-Peso
		cStoneDark := color.RGBA{R: 52, G: 54, B: 58, A: 255}
		cStoneMid := color.RGBA{R: 72, G: 76, B: 82, A: 255}
		cStoneLight := color.RGBA{R: 98, G: 104, B: 112, A: 255}
		cWaterEdge := color.RGBA{R: 35, G: 70, B: 85, A: 255}

		ebitenutil.DrawRect(screen, 0, groundY, screenWidth, 55, cStoneDark)
		ebitenutil.DrawRect(screen, 0, groundY, screenWidth, 3, cStoneLight)

		// Lajes de pedra do cais em scroll
		for x := -math.Mod(b.scrollOffset, 24.0); x < screenWidth; x += 24.0 {
			ebitenutil.DrawRect(screen, x, groundY+4, 22, 10, cStoneMid)
			ebitenutil.DrawRect(screen, x+2, groundY+5, 18, 2, cStoneLight)
			ebitenutil.DrawRect(screen, x+12, groundY+18, 22, 11, cStoneMid)
			ebitenutil.DrawRect(screen, x+14, groundY+19, 18, 2, cStoneLight)
			ebitenutil.DrawRect(screen, x, groundY+33, 22, 11, cStoneMid)
		}

		// Borda da Baía com água no limite inferior
		ebitenutil.DrawRect(screen, 0, groundY+49, screenWidth, 6, cWaterEdge)
		ebitenutil.DrawRect(screen, 0, groundY+48, screenWidth, 1, color.RGBA{R: 90, G: 145, B: 170, A: 200})
		return
	}

	// Fallback procedural
	b.drawVeroPesoProcedural(screen, screenWidth, groundY, ticks)
}

// ==========================================
// FASE 2: ESTAÇÃO DAS DOCAS (8-BIT)
// ==========================================
func (b *Background) drawDocas(screen *ebiten.Image, screenWidth, groundY float64, ticks int) {
	if b.imgFase2 != nil {
		// 1. Cenário panorâmico real 8-bit das Docas (Galpões ingleses vermelhos, guindaste amarelo e cais)
		b.drawLoopingImage(screen, b.imgFase2, screenWidth, groundY, 0.45)

		// 2. Luminárias industriais de ferro do cais das Docas com luz âmbar em primeiro plano
		cDocasLamp := color.RGBA{R: 35, G: 38, B: 44, A: 255}
		cAmberGlow := color.RGBA{R: 255, G: 185, B: 65, A: 210}
		cAmberBeam := color.RGBA{R: 255, G: 190, B: 70, A: 25}

		for p := 0; p < 3; p++ {
			px := float64(p*140) + 60.0 - math.Mod(b.scrollOffset*0.75, 140.0)
			ebitenutil.DrawRect(screen, px, groundY-32, 2, 32, cDocasLamp)
			ebitenutil.DrawRect(screen, px-2, groundY-35, 6, 3, cDocasLamp)
			ebitenutil.DrawRect(screen, px-1, groundY-34, 4, 3, cAmberGlow)
			ebitenutil.DrawRect(screen, px-4, groundY-31, 10, 31, cAmberBeam)
		}

		// 3. Chão: Famoso Deck de Madeira de Lei (Ipê/Itaúba) das Docas e Trilhos dos Guindastes
		cWoodDark := color.RGBA{R: 78, G: 46, B: 30, A: 255}
		cWoodMid := color.RGBA{R: 110, G: 68, B: 44, A: 255}
		cWoodLight := color.RGBA{R: 138, G: 86, B: 56, A: 255}
		cRailMetal := color.RGBA{R: 45, G: 50, B: 55, A: 255}
		cRailShine := color.RGBA{R: 165, G: 175, B: 185, A: 255}

		ebitenutil.DrawRect(screen, 0, groundY, screenWidth, 55, cWoodDark)

		// Réguas do deck de madeira com juntas
		for y := groundY; y < groundY+55; y += 6.0 {
			ebitenutil.DrawRect(screen, 0, y, screenWidth, 5, cWoodMid)
			ebitenutil.DrawRect(screen, 0, y, screenWidth, 1, cWoodLight)
			ebitenutil.DrawRect(screen, 0, y+5, screenWidth, 1, cWoodDark)
		}

		// Parafusos e nós da madeira em movimento
		for x := -math.Mod(b.scrollOffset, 30.0); x < screenWidth; x += 30.0 {
			ebitenutil.DrawRect(screen, x, groundY+2, 2, 2, cWoodDark)
			ebitenutil.DrawRect(screen, x+14, groundY+14, 2, 2, cWoodDark)
			ebitenutil.DrawRect(screen, x+4, groundY+26, 2, 2, cWoodDark)
			ebitenutil.DrawRect(screen, x+18, groundY+38, 2, 2, cWoodDark)
		}

		// Trilho de ferro fundido do guindaste inglês
		ebitenutil.DrawRect(screen, 0, groundY+8, screenWidth, 4, cRailMetal)
		ebitenutil.DrawRect(screen, 0, groundY+8, screenWidth, 1, cRailShine)
		ebitenutil.DrawRect(screen, 0, groundY+20, screenWidth, 4, cRailMetal)
		ebitenutil.DrawRect(screen, 0, groundY+20, screenWidth, 1, cRailShine)
		return
	}

	// Fallback procedural
	screen.Fill(color.RGBA{R: 20, G: 35, B: 55, A: 255})
	ebitenutil.DrawRect(screen, 0, 70, screenWidth, 45, color.RGBA{R: 190, G: 80, B: 60, A: 255})
	ebitenutil.DrawRect(screen, 0, groundY, screenWidth, 55, color.RGBA{R: 90, G: 60, B: 40, A: 255})
}

// ==========================================
// FASE 3: THEATRO DA PAZ (8-BIT)
// ==========================================
func (b *Background) drawTheatroDaPaz(screen *ebiten.Image, screenWidth, groundY float64, ticks int) {
	if b.imgFase3 != nil {
		// 1. Cenário panorâmico real 8-bit do Theatro da Paz (Fachada neoclássica, mangueiras e pedras portuguesas)
		b.drawLoopingImage(screen, b.imgFase3, screenWidth, groundY, 0.45)

		// 2. Postes republicanos de ferro com globo de iluminação da Praça da República em primeiro plano
		cRepIron := color.RGBA{R: 30, G: 34, B: 38, A: 255}
		cGlobeWhite := color.RGBA{R: 255, G: 250, B: 230, A: 240}
		cGlobeBeam := color.RGBA{R: 255, G: 245, B: 200, A: 20}

		for p := 0; p < 3; p++ {
			px := float64(p*135) + 45.0 - math.Mod(b.scrollOffset*0.75, 135.0)
			ebitenutil.DrawRect(screen, px, groundY-46, 2, 46, cRepIron)
			ebitenutil.DrawRect(screen, px-2, groundY-48, 6, 2, cRepIron)
			ebitenutil.DrawRect(screen, px-2, groundY-54, 6, 6, cGlobeWhite)
			ebitenutil.DrawRect(screen, px-6, groundY-46, 14, 46, cGlobeBeam)
		}

		// 3. Chão: Calçadão clássico de pedras portuguesas preto e branco (Mosaico da Praça da República)
		cStoneWhite := color.RGBA{R: 228, G: 228, B: 222, A: 255}
		cStoneBlack := color.RGBA{R: 42, G: 44, B: 48, A: 255}
		cBaseGrey := color.RGBA{R: 90, G: 92, B: 96, A: 255}
		cGrassGreen := color.RGBA{R: 35, G: 95, B: 45, A: 255}

		ebitenutil.DrawRect(screen, 0, groundY, screenWidth, 55, cBaseGrey)

		// Mosaico geométrico ondulado em pixel art
		step := 16.0
		for x := -math.Mod(b.scrollOffset, step*2); x < screenWidth; x += step {
			ebitenutil.DrawRect(screen, x, groundY+2, 7, 7, cStoneWhite)
			ebitenutil.DrawRect(screen, x+8, groundY+2, 7, 7, cStoneBlack)
			ebitenutil.DrawRect(screen, x, groundY+10, 7, 7, cStoneBlack)
			ebitenutil.DrawRect(screen, x+8, groundY+10, 7, 7, cStoneWhite)
			ebitenutil.DrawRect(screen, x, groundY+18, 7, 7, cStoneWhite)
			ebitenutil.DrawRect(screen, x+8, groundY+18, 7, 7, cStoneBlack)
			ebitenutil.DrawRect(screen, x, groundY+26, 7, 7, cStoneBlack)
			ebitenutil.DrawRect(screen, x+8, groundY+26, 7, 7, cStoneWhite)
		}

		// Meio-fio e gramado da praça no rodapé
		ebitenutil.DrawRect(screen, 0, groundY+36, screenWidth, 4, color.RGBA{R: 160, G: 160, B: 165, A: 255})
		ebitenutil.DrawRect(screen, 0, groundY+40, screenWidth, 15, cGrassGreen)
		ebitenutil.DrawRect(screen, 0, groundY+40, screenWidth, 2, color.RGBA{R: 55, G: 135, B: 65, A: 255})
		return
	}

	// Fallback procedural
	screen.Fill(color.RGBA{R: 35, G: 45, B: 60, A: 255})
	ebitenutil.DrawRect(screen, 0, 80, screenWidth, 35, color.RGBA{R: 180, G: 120, B: 90, A: 255})
	ebitenutil.DrawRect(screen, 0, groundY, screenWidth, 55, color.RGBA{R: 75, G: 75, B: 80, A: 255})
}

// ==========================================
// PÁSSAROS E DETALHES DE BELÉM
// ==========================================
func drawSkyBird(screen *ebiten.Image, x, y float64, flap int, dark color.RGBA) {
	switch flap % 3 {
	case 0: // Asas erguidas
		ebitenutil.DrawRect(screen, x, y, 2, 2, dark)
		ebitenutil.DrawRect(screen, x-2, y-1, 2, 2, dark)
		ebitenutil.DrawRect(screen, x+2, y-1, 2, 2, dark)
		ebitenutil.DrawRect(screen, x-4, y-2, 2, 2, dark)
		ebitenutil.DrawRect(screen, x+4, y-2, 2, 2, dark)
	case 1: // Asas retas planando no ar
		ebitenutil.DrawRect(screen, x, y, 2, 2, dark)
		ebitenutil.DrawRect(screen, x-3, y, 3, 1, dark)
		ebitenutil.DrawRect(screen, x+2, y, 3, 1, dark)
	case 2: // Asas baixadas
		ebitenutil.DrawRect(screen, x, y, 2, 2, dark)
		ebitenutil.DrawRect(screen, x-2, y+1, 2, 2, dark)
		ebitenutil.DrawRect(screen, x+2, y+1, 2, 2, dark)
		ebitenutil.DrawRect(screen, x-4, y+2, 2, 2, dark)
		ebitenutil.DrawRect(screen, x+4, y+2, 2, 2, dark)
	}
}

// ==========================================
// FALLBACK PROCEDURAL DO VER-O-PESO
// ==========================================
func (b *Background) drawVeroPesoProcedural(screen *ebiten.Image, screenWidth, groundY float64, ticks int) {
	screen.Fill(color.RGBA{R: 24, G: 48, B: 68, A: 255})

	ebitenutil.DrawRect(screen, 0, 75, screenWidth, 40, color.RGBA{R: 200, G: 110, B: 55, A: 255})
	ebitenutil.DrawRect(screen, 0, 105, screenWidth, 15, color.RGBA{R: 220, G: 150, B: 75, A: 255})
	ebitenutil.DrawRect(screen, 265, 30, 26, 26, color.RGBA{R: 255, G: 215, B: 110, A: 255})
	ebitenutil.DrawRect(screen, 0, 115, screenWidth, 20, color.RGBA{R: 30, G: 75, B: 95, A: 255})

	b.drawMercadoDeFerro(screen, 30, groundY)
	b.drawMercadoDeFerro(screen, 185, groundY)

	ebitenutil.DrawRect(screen, 0, groundY, screenWidth, 60, color.RGBA{R: 85, G: 85, B: 90, A: 255})
	for x := -b.scrollOffset; x < screenWidth; x += 20 {
		ebitenutil.DrawRect(screen, x, groundY+4, 18, 8, color.RGBA{R: 70, G: 70, B: 75, A: 255})
	}
}

func (b *Background) drawMercadoDeFerro(screen *ebiten.Image, startX, groundY float64) {
	metalBase := color.RGBA{R: 55, G: 65, B: 75, A: 255}
	metalLight := color.RGBA{R: 80, G: 95, B: 110, A: 255}
	cupulaBlue := color.RGBA{R: 35, G: 95, B: 165, A: 255}
	cupulaDark := color.RGBA{R: 22, G: 60, B: 110, A: 255}

	ebitenutil.DrawRect(screen, startX+12, groundY-48, 80, 48, metalBase)
	ebitenutil.DrawRect(screen, startX+16, groundY-44, 72, 44, metalLight)

	for win := 0; win < 4; win++ {
		wx := startX + 22 + float64(win*17)
		ebitenutil.DrawRect(screen, wx, groundY-36, 10, 18, color.RGBA{R: 30, G: 35, B: 45, A: 255})
	}

	b.drawTorre(screen, startX, groundY, metalBase, cupulaBlue, cupulaDark)
	b.drawTorre(screen, startX+88, groundY, metalBase, cupulaBlue, cupulaDark)
}

func (b *Background) drawTorre(screen *ebiten.Image, tx, groundY float64, base, cupula, cupulaDark color.RGBA) {
	ebitenutil.DrawRect(screen, tx, groundY-56, 14, 56, base)
	ebitenutil.DrawRect(screen, tx+2, groundY-54, 10, 54, color.RGBA{R: 70, G: 80, B: 90, A: 255})
	ebitenutil.DrawRect(screen, tx-1, groundY-62, 16, 6, cupulaDark)
	ebitenutil.DrawRect(screen, tx, groundY-68, 14, 7, cupula)
	ebitenutil.DrawRect(screen, tx+2, groundY-74, 10, 6, cupula)
	ebitenutil.DrawRect(screen, tx+4, groundY-79, 6, 5, cupula)
	ebitenutil.DrawRect(screen, tx+6, groundY-86, 2, 7, color.RGBA{R: 215, G: 215, B: 220, A: 255})
}
