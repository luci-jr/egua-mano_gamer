package scenery

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed ver_o_peso_bg.jpg
var bgJpgBytes []byte

type Background struct {
	scrollOffset float64
	bgImage      *ebiten.Image
}

func NewBackground() *Background {
	b := &Background{scrollOffset: 0}

	if len(bgJpgBytes) > 0 {
		img, _, err := image.Decode(bytes.NewReader(bgJpgBytes))
		if err == nil {
			b.bgImage = ebiten.NewImageFromImage(img)
		}
	}

	return b
}

func (b *Background) Update(speed float64) {
	b.scrollOffset += speed * 0.7
	if b.bgImage != nil {
		bgWidth := 340.0
		if b.scrollOffset >= bgWidth {
			b.scrollOffset -= bgWidth
		}
	} else {
		if b.scrollOffset >= 40 {
			b.scrollOffset -= 40
		}
	}
}

func (b *Background) Draw(screen *ebiten.Image, screenWidth, groundY float64, ticks int, stage int) {
	switch stage {
	case 2:
		b.drawDocas(screen, screenWidth, groundY, ticks)
	case 3:
		b.drawTheatroDaPaz(screen, screenWidth, groundY, ticks)
	default:
		b.drawAmazonRuins(screen, screenWidth, groundY, ticks)
	}
}

func (b *Background) drawVerOPesoImage(screen *ebiten.Image, screenWidth, groundY float64) {
	bounds := b.bgImage.Bounds()
	imgW := float64(bounds.Dx())
	imgH := float64(bounds.Dy())

	scaleX := screenWidth / imgW
	scaleY := (groundY + 55.0) / imgH

	op1 := &ebiten.DrawImageOptions{}
	op1.GeoM.Scale(scaleX, scaleY)
	op1.GeoM.Translate(-b.scrollOffset, 0)
	screen.DrawImage(b.bgImage, op1)

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Scale(scaleX, scaleY)
	op2.GeoM.Translate(screenWidth-b.scrollOffset, 0)
	screen.DrawImage(b.bgImage, op2)

	ebitenutil.DrawRect(screen, 0, groundY+48, screenWidth, 12, color.RGBA{R: 45, G: 40, B: 35, A: 255})
}

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

func (b *Background) drawDocas(screen *ebiten.Image, screenWidth, groundY float64, ticks int) {
	screen.Fill(color.RGBA{R: 20, G: 35, B: 55, A: 255})

	ebitenutil.DrawRect(screen, 0, 70, screenWidth, 45, color.RGBA{R: 190, G: 80, B: 60, A: 255})
	ebitenutil.DrawRect(screen, 0, 105, screenWidth, 15, color.RGBA{R: 210, G: 130, B: 80, A: 255})
	ebitenutil.DrawRect(screen, 0, 115, screenWidth, 20, color.RGBA{R: 20, G: 60, B: 80, A: 255})

	docaRed := color.RGBA{R: 140, G: 50, B: 40, A: 255}
	docaDark := color.RGBA{R: 45, G: 45, B: 50, A: 255}
	glassYellow := color.RGBA{R: 245, G: 210, B: 110, A: 220}

	for arm := 0; arm < 2; arm++ {
		ax := float64(arm*170) + 15
		ebitenutil.DrawRect(screen, ax, groundY-52, 115, 52, docaRed)
		ebitenutil.DrawRect(screen, ax+4, groundY-50, 107, 48, color.RGBA{R: 165, G: 60, B: 50, A: 255})
		ebitenutil.DrawRect(screen, ax, groundY-58, 115, 6, docaDark)

		for w := 0; w < 5; w++ {
			wx := ax + 12 + float64(w*20)
			ebitenutil.DrawRect(screen, wx, groundY-38, 12, 22, glassYellow)
			ebitenutil.DrawRect(screen, wx+5, groundY-38, 2, 22, docaDark)
		}
	}

	craneX := float64(140)
	craneYellow := color.RGBA{R: 240, G: 185, B: 30, A: 255}
	ebitenutil.DrawRect(screen, craneX, groundY-68, 6, 68, craneYellow)
	ebitenutil.DrawRect(screen, craneX+16, groundY-68, 6, 68, craneYellow)
	ebitenutil.DrawRect(screen, craneX-4, groundY-76, 30, 8, craneYellow)
	ebitenutil.DrawRect(screen, craneX-8, groundY-88, 12, 12, color.RGBA{R: 60, G: 60, B: 65, A: 255})
	ebitenutil.DrawRect(screen, craneX+4, groundY-92, 42, 6, craneYellow)
	ebitenutil.DrawRect(screen, craneX+42, groundY-86, 2, 22, color.RGBA{R: 200, G: 200, B: 200, A: 255})

	ebitenutil.DrawRect(screen, 0, groundY, screenWidth, 60, color.RGBA{R: 90, G: 60, B: 40, A: 255})
	for x := -b.scrollOffset; x < screenWidth; x += 18 {
		ebitenutil.DrawRect(screen, x, groundY+2, 16, 55, color.RGBA{R: 110, G: 75, B: 50, A: 255})
		ebitenutil.DrawRect(screen, x+16, groundY+2, 2, 55, color.RGBA{R: 60, G: 40, B: 25, A: 255})
	}
}

func (b *Background) drawTheatroDaPaz(screen *ebiten.Image, screenWidth, groundY float64, ticks int) {
	screen.Fill(color.RGBA{R: 35, G: 45, B: 60, A: 255})

	ebitenutil.DrawRect(screen, 0, 80, screenWidth, 35, color.RGBA{R: 180, G: 120, B: 90, A: 255})

	for m := 0; m < 3; m++ {
		mx := float64(m*120) + 10
		ebitenutil.DrawRect(screen, mx+18, groundY-65, 8, 65, color.RGBA{R: 65, G: 45, B: 30, A: 255})
		ebitenutil.DrawRect(screen, mx-10, groundY-105, 65, 50, color.RGBA{R: 25, G: 75, B: 35, A: 255})
		ebitenutil.DrawRect(screen, mx-5, groundY-110, 55, 45, color.RGBA{R: 35, G: 95, B: 45, A: 255})
		ebitenutil.DrawRect(screen, mx+5, groundY-115, 35, 35, color.RGBA{R: 45, G: 115, B: 55, A: 255})

		ebitenutil.DrawRect(screen, mx+10, groundY-85, 4, 5, color.RGBA{R: 240, G: 160, B: 30, A: 255})
		ebitenutil.DrawRect(screen, mx+32, groundY-90, 4, 5, color.RGBA{R: 240, G: 160, B: 30, A: 255})
	}

	tX := float64(75)
	theaterWhite := color.RGBA{R: 240, G: 235, B: 225, A: 255}
	theaterGold := color.RGBA{R: 220, G: 180, B: 90, A: 255}

	ebitenutil.DrawRect(screen, tX, groundY-75, 140, 75, theaterWhite)
	ebitenutil.DrawRect(screen, tX+10, groundY-85, 120, 10, theaterGold)
	for p := 0; p < 6; p++ {
		px := tX + 22 + float64(p*18)
		ebitenutil.DrawRect(screen, px, groundY-75, 6, 75, color.RGBA{R: 215, G: 210, B: 200, A: 255})
	}

	ebitenutil.DrawRect(screen, tX+20, groundY-100, 100, 15, theaterWhite)
	ebitenutil.DrawRect(screen, tX+35, groundY-108, 70, 8, theaterGold)

	ebitenutil.DrawRect(screen, 0, groundY, screenWidth, 60, color.RGBA{R: 75, G: 75, B: 80, A: 255})
	for x := -b.scrollOffset; x < screenWidth; x += 16 {
		ebitenutil.DrawRect(screen, x, groundY+3, 7, 7, color.RGBA{R: 220, G: 220, B: 220, A: 255})
		ebitenutil.DrawRect(screen, x+8, groundY+3, 7, 7, color.RGBA{R: 40, G: 40, B: 45, A: 255})
		ebitenutil.DrawRect(screen, x, groundY+12, 7, 7, color.RGBA{R: 40, G: 40, B: 45, A: 255})
		ebitenutil.DrawRect(screen, x+8, groundY+12, 7, 7, color.RGBA{R: 220, G: 220, B: 220, A: 255})
	}
}

func (b *Background) drawAmazonRuins(screen *ebiten.Image, screenWidth, groundY float64, ticks int) {
	// Céu denso da selva amazônica (Estilo Pitfall 16-bit)
	screen.Fill(color.RGBA{R: 16, G: 34, B: 26, A: 255})

	// Raios solares translúcidos (God Rays) filtrando pelas copas
	godRayColor := color.RGBA{R: 245, G: 235, B: 160, A: 16}
	for ray := 0; ray < 3; ray++ {
		rx := float64(ray*115) + 25.0
		ebitenutil.DrawRect(screen, rx, 0, 36, groundY, godRayColor)
		ebitenutil.DrawRect(screen, rx+8, 0, 16, groundY, godRayColor)
	}

	// Copas densas de árvores no topo da tela (Dossel da Selva)
	canopyDark := color.RGBA{R: 14, G: 42, B: 20, A: 255}
	canopyMid := color.RGBA{R: 24, G: 68, B: 32, A: 255}
	canopyLight := color.RGBA{R: 42, G: 105, B: 50, A: 255}
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, 22, canopyDark)
	for leafX := 0.0; leafX < screenWidth; leafX += 16.0 {
		ebitenutil.DrawRect(screen, leafX, 20, 15, 8, canopyMid)
		ebitenutil.DrawRect(screen, leafX+3, 26, 8, 4, canopyLight)
	}

	// Cipós pendurados balançando suavemente no fundo (decorativos)
	cVineBg := color.RGBA{R: 32, G: 60, B: 24, A: 210}
	for v := 0; v < 4; v++ {
		vx := float64(v*90) + 30.0 - math.Mod(b.scrollOffset*0.3, 90.0)
		swing := math.Sin(float64(ticks+v*30)*0.04) * 4.0
		ebitenutil.DrawLine(screen, vx, 14, vx+swing, groundY-40, cVineBg)
		ebitenutil.DrawLine(screen, vx+1, 14, vx+swing+1, groundY-40, cVineBg)
	}

	// Sumaúmas gigantes no horizonte (Parallax distante)
	treeDark := color.RGBA{R: 38, G: 30, B: 24, A: 255}
	treeLight := color.RGBA{R: 55, G: 44, B: 32, A: 255}
	treeMoss := color.RGBA{R: 48, G: 90, B: 38, A: 255}

	for t := 0; t < 3; t++ {
		tx := float64(t*130) + 15.0 - math.Mod(b.scrollOffset*0.4, 130.0)
		// Tronco massivo
		ebitenutil.DrawRect(screen, tx, groundY-85, 24, 85, treeDark)
		ebitenutil.DrawRect(screen, tx+4, groundY-82, 16, 82, treeLight)
		// Sapopembas (raízes tubulares na base)
		ebitenutil.DrawRect(screen, tx-10, groundY-25, 12, 25, treeDark)
		ebitenutil.DrawRect(screen, tx+22, groundY-25, 12, 25, treeDark)
		// Musgo no tronco
		ebitenutil.DrawRect(screen, tx+6, groundY-50, 6, 12, treeMoss)
		ebitenutil.DrawRect(screen, tx+8, groundY-30, 8, 8, treeMoss)
	}

	// Ruínas Arqueológicas de Pedra Marajoara com Tochas Crepitantes
	stoneBase := color.RGBA{R: 68, G: 76, B: 72, A: 255}
	stoneLight := color.RGBA{R: 92, G: 102, B: 96, A: 255}
	stoneGlyph := color.RGBA{R: 42, G: 48, B: 45, A: 255}
	mossGreen := color.RGBA{R: 42, G: 110, B: 48, A: 255}

	for r := 0; r < 2; r++ {
		rx := float64(r*175) + 60.0 - math.Mod(b.scrollOffset*0.7, 175.0)

		// Pilar / Altar de pedra escalonada
		ebitenutil.DrawRect(screen, rx, groundY-55, 45, 55, stoneBase)
		ebitenutil.DrawRect(screen, rx+3, groundY-52, 39, 52, stoneLight)
		ebitenutil.DrawRect(screen, rx-4, groundY-58, 53, 5, stoneBase)
		ebitenutil.DrawRect(screen, rx-2, groundY-56, 49, 3, stoneLight)

		// Glifos geométricos esculpidos na pedra
		ebitenutil.DrawRect(screen, rx+8, groundY-44, 10, 10, stoneGlyph)
		ebitenutil.DrawRect(screen, rx+10, groundY-42, 6, 6, stoneLight)
		ebitenutil.DrawRect(screen, rx+26, groundY-44, 10, 10, stoneGlyph)
		ebitenutil.DrawRect(screen, rx+28, groundY-42, 6, 6, stoneLight)
		ebitenutil.DrawRect(screen, rx+12, groundY-26, 20, 8, stoneGlyph)

		// Musgo trepando nas frestas da rocha
		ebitenutil.DrawRect(screen, rx+2, groundY-48, 5, 8, mossGreen)
		ebitenutil.DrawRect(screen, rx+38, groundY-32, 6, 12, mossGreen)

		// Tocha de Pedra Crepitante montada na ruína
		torchX := rx + 20.0
		torchY := groundY - 68.0
		// Suporte da tocha
		ebitenutil.DrawRect(screen, torchX, torchY+7, 5, 7, color.RGBA{R: 48, G: 42, B: 38, A: 255})
		ebitenutil.DrawRect(screen, torchX-2, torchY+4, 9, 3, color.RGBA{R: 85, G: 72, B: 52, A: 255})

		// Chama animada com tremulação de 3 camadas em pixel art
		flicker := (ticks / 4) % 3
		flameOffset := float64(flicker)
		// Brilho da tocha no ambiente
		ebitenutil.DrawRect(screen, torchX-6, torchY-8, 17, 17, color.RGBA{R: 255, G: 195, B: 50, A: 45})

		// Camada externa vermelha da chama
		ebitenutil.DrawRect(screen, torchX-1, torchY-4+flameOffset*0.5, 7, 8, color.RGBA{R: 230, G: 55, B: 40, A: 255})
		// Camada média laranja
		ebitenutil.DrawRect(screen, torchX, torchY-6+flameOffset*0.5, 5, 7, color.RGBA{R: 245, G: 140, B: 25, A: 255})
		// Núcleo dourado brilhante
		ebitenutil.DrawRect(screen, torchX+1, torchY-7+flameOffset*0.5, 3, 5, color.RGBA{R: 255, G: 225, B: 85, A: 255})
		// Fagulhas subindo no ar
		if flicker == 1 {
			ebitenutil.DrawRect(screen, torchX+2, torchY-11, 2, 2, color.RGBA{R: 255, G: 200, B: 60, A: 220})
		}
	}

	// Solo terroso da selva com musgo, raízes e folhas caídas
	cEarthDark := color.RGBA{R: 36, G: 24, B: 15, A: 255}
	cEarthMid := color.RGBA{R: 52, G: 36, B: 22, A: 255}
	cForestMoss := color.RGBA{R: 48, G: 128, B: 42, A: 255}
	cLeafBrown := color.RGBA{R: 138, G: 80, B: 32, A: 255}

	ebitenutil.DrawRect(screen, 0, groundY, screenWidth, 60, cEarthDark)
	ebitenutil.DrawRect(screen, 0, groundY, screenWidth, 4, cEarthMid)

	for x := -math.Mod(b.scrollOffset, 16.0); x < screenWidth; x += 16.0 {
		// Grama e musgo no topo do solo
		ebitenutil.DrawRect(screen, x, groundY-1, 11, 3, cForestMoss)
		ebitenutil.DrawRect(screen, x+2, groundY-2, 5, 2, cForestMoss)

		// Raízes e folhas na terra
		ebitenutil.DrawRect(screen, x+6, groundY+6, 6, 2, cEarthMid)
		ebitenutil.DrawRect(screen, x+12, groundY+12, 3, 2, cLeafBrown)
	}
}
