package ui

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed title_screen.jpg
var titleScreenBytes []byte
var titleScreenImg *ebiten.Image

func getTitleScreenImage() *ebiten.Image {
	if titleScreenImg != nil {
		return titleScreenImg
	}
	if len(titleScreenBytes) == 0 {
		return nil
	}
	img, _, err := image.Decode(bytes.NewReader(titleScreenBytes))
	if err != nil {
		return nil
	}
	titleScreenImg = ebiten.NewImageFromImage(img)
	return titleScreenImg
}

func DrawHeart(screen *ebiten.Image, x, y float64, filled bool) {
	if filled {
		red := color.RGBA{R: 235, G: 40, B: 50, A: 255}
		ebitenutil.DrawRect(screen, x+1, y, 3, 2, red)
		ebitenutil.DrawRect(screen, x+5, y, 3, 2, red)
		ebitenutil.DrawRect(screen, x, y+2, 9, 3, red)
		ebitenutil.DrawRect(screen, x+1, y+5, 7, 2, red)
		ebitenutil.DrawRect(screen, x+2, y+7, 5, 1, red)
		ebitenutil.DrawRect(screen, x+3, y+8, 3, 1, red)
		ebitenutil.DrawRect(screen, x+4, y+9, 1, 1, red)
		ebitenutil.DrawRect(screen, x+2, y+1, 1, 1, color.RGBA{R: 255, G: 200, B: 200, A: 255})
	} else {
		outline := color.RGBA{R: 160, G: 70, B: 70, A: 255}
		ebitenutil.DrawRect(screen, x+1, y, 3, 1, outline)
		ebitenutil.DrawRect(screen, x+5, y, 3, 1, outline)
		ebitenutil.DrawRect(screen, x, y+1, 1, 4, outline)
		ebitenutil.DrawRect(screen, x+8, y+1, 1, 4, outline)
		ebitenutil.DrawRect(screen, x+4, y+1, 1, 2, outline)
		ebitenutil.DrawRect(screen, x+1, y+5, 1, 2, outline)
		ebitenutil.DrawRect(screen, x+7, y+5, 1, 2, outline)
		ebitenutil.DrawRect(screen, x+2, y+7, 1, 1, outline)
		ebitenutil.DrawRect(screen, x+6, y+7, 1, 1, outline)
		ebitenutil.DrawRect(screen, x+3, y+8, 1, 1, outline)
		ebitenutil.DrawRect(screen, x+5, y+8, 1, 1, outline)
		ebitenutil.DrawRect(screen, x+4, y+9, 1, 1, outline)
	}
}

func getCuriosityText(stage int) string {
	switch stage {
	case 2: // Estação das Docas
		return "★ ESTACAO DAS DOCAS: Armazens de ferro ingleses de 1897 restaurados na orla de Belem!  " +
			"★ GUINDASTES: Importados no sec. XX, hoje marcos historicos do porto!  " +
			"★ ILHA DO COMBU: Polo de cacau nativo e turismo ecologico a 10 min de barco!  " +
			"★ CARIMBO: Ritmo e danca tradicional paraense patrimonio cultural do Brasil!  " +
			"★ CHUVA DA TARDE: O belemense marca compromissos 'antes ou depois da chuva'!  " +
			"★ SORVETES TIPICOS: Saboreie acai, cupuacu, bacuri, tapereba e castanha nas Docas!  "
	case 3: // Theatro da Paz
		return "★ THEATRO DA PAZ: Fundado em 1878 na Belle Epoque, inspirado no Scala de Milao!  " +
			"★ CIDADE DAS MANGUEIRAS: Belem ganhou este titulo pelas arvores plantadas no sec. XIX!  " +
			"★ CIRIO DE NAZARE: O Natal dos paraenses reune mais de 2 milhoes de devotos!  " +
			"★ PRACA DA REPUBLICA: Grande praca colonial que abriga o Theatro no coracao da cidade!  " +
			"★ ONCA-PINTADA: Rainha da Amazonia com a mordida mais potente entre todos os felinos!  " +
			"★ PAIDEGUA: Expressao genuina do Para para algo excelente, sensacional e autentico!  "
	default: // Mercado do Ver-o-Peso
		return "★ VER-O-PESO: Fundado em 1627, e a maior feira a ceu aberto da America Latina!  " +
			"★ ACAI PURO: No Para, o acai e consumido tradicionalmente com peixe frito e farinha!  " +
			"★ TACACA: Servido na cuia quente com tucupi e folhas de jambu que amortecem os labios!  " +
			"★ JACARE-ACU: O gigante dos rios amazonicos alcanca mais de 4 metros de comprimento!  " +
			"★ BAIA DO GUAJARA: Aguas onde atracam diariamente os barcos com pescado e acai!  " +
			"★ PANEIRO: Cesto tipico de palha trancada usado para carregar o acai colhido!  "
	}
}

func DrawCityFooter(screen *ebiten.Image, screenWidth, screenHeight float64, stage int, ticks int) {
	footerH := 14.0
	footerY := screenHeight - footerH

	// Fundo escuro do rodapé com transparência suave
	ebitenutil.DrawRect(screen, 0, footerY, screenWidth, footerH, color.RGBA{R: 8, G: 12, B: 22, A: 235})
	// Filete dourado divisório superior
	ebitenutil.DrawRect(screen, 0, footerY, screenWidth, 1, color.RGBA{R: 250, G: 205, B: 55, A: 220})

	fullText := getCuriosityText(stage)
	textWidth := len(fullText) * 6
	if textWidth == 0 {
		return
	}

	// Rolagem contínua suave (1 pixel por frame)
	scrollOffset := ticks % textWidth
	baseX := 52 - scrollOffset

	// Desenha a primeira repetição
	ebitenutil.DebugPrintAt(screen, fullText, baseX, int(footerY)+1)
	// Desenha a segunda repetição para manter o fluxo contínuo
	if baseX+textWidth < int(screenWidth) {
		ebitenutil.DebugPrintAt(screen, fullText, baseX+textWidth, int(footerY)+1)
	}

	// Badge fixa que mascara o texto que passa sob ela
	ebitenutil.DrawRect(screen, 0, footerY, 48, footerH, color.RGBA{R: 16, G: 26, B: 46, A: 255})
	ebitenutil.DrawRect(screen, 48, footerY, 1, footerH, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DebugPrintAt(screen, "★BELEM", 4, int(footerY)+1)
}

func DrawHUD(screen *ebiten.Image, lives int, hearts int, score int, stage int, isDoubleJump bool, stageBannerTimer int, isMuted bool, ticks int, relicsCount int, heroName string, stageDistance float64) {
	for i := 0; i < 3; i++ {
		hx := 8.0 + float64(i*12)
		DrawHeart(screen, hx, 9, i < hearts)
	}

	livesText := fmt.Sprintf("x%d %s", lives, heroName)
	ebitenutil.DebugPrintAt(screen, livesText, 46, 9)

	stageName := "1. VER-O-PESO"
	if stage == 2 {
		stageName = "2. DOCAS"
	} else if stage == 3 {
		stageName = "3. THEATRO"
	}
	ebitenutil.DebugPrintAt(screen, stageName, 126, 9)

	scoreText := fmt.Sprintf("SCORE: %05d", score)
	ebitenutil.DebugPrintAt(screen, scoreText, 226, 9)

	// Contador de Relíquias / Tesouros Amazônicos coletados
	ebitenutil.DrawRect(screen, 8, 23, 7, 7, color.RGBA{R: 245, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, 9, 24, 5, 5, color.RGBA{R: 45, G: 195, B: 95, A: 255})
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("RELIQUIAS: %d", relicsCount), 18, 22)

	// Medidor de distância da corrida para a chegada da fase
	targetDist := 2800
	currDist := int(stageDistance)
	if currDist > targetDist {
		currDist = targetDist
	}
	distText := fmt.Sprintf("DIST: %dm/%dm", currDist, targetDist)
	ebitenutil.DebugPrintAt(screen, distText, 115, 22)

	if isMuted {
		ebitenutil.DebugPrintAt(screen, "[MUDO]", 226, 22)
	}

	if stageBannerTimer > 0 {
		var bannerTitle string
		switch stage {
		case 2:
			bannerTitle = "★ FASE 2: ESTACAO DAS DOCAS ★"
		case 3:
			bannerTitle = "★ FASE 3: THEATRO DA PAZ ★"
		default:
			bannerTitle = "★ FASE 1: MERCADO DO VER-O-PESO ★"
		}
		ebitenutil.DrawRect(screen, 20, 50, 280, 22, color.RGBA{R: 20, G: 20, B: 30, A: 210})
		ebitenutil.DrawRect(screen, 20, 50, 280, 2, color.RGBA{R: 250, G: 200, B: 50, A: 255})
		ebitenutil.DrawRect(screen, 20, 70, 280, 2, color.RGBA{R: 250, G: 200, B: 50, A: 255})
		ebitenutil.DebugPrintAt(screen, bannerTitle, 30, 55)
	}

	// Rodapé cultural com curiosidades dinâmicas de Belém do Pará
	DrawCityFooter(screen, float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy()), stage, ticks)
}

// drawFlyingUrubu desenha um urubu em pixel art com asas animadas voando no céu
func drawFlyingUrubu(screen *ebiten.Image, x, y float64, flapTick int) {
	cBody := color.RGBA{R: 20, G: 20, B: 26, A: 245}
	cBeak := color.RGBA{R: 75, G: 70, B: 65, A: 240}

	phase := (flapTick / 8) % 4
	switch phase {
	case 0: // Asas para cima \_/
		ebitenutil.DrawRect(screen, x+4, y+2, 4, 3, cBody)
		ebitenutil.DrawRect(screen, x+3, y+3, 1, 1, cBeak)
		ebitenutil.DrawRect(screen, x+2, y+1, 2, 2, cBody)
		ebitenutil.DrawRect(screen, x, y-1, 2, 2, cBody)
		ebitenutil.DrawRect(screen, x+8, y+1, 2, 2, cBody)
		ebitenutil.DrawRect(screen, x+10, y-1, 2, 2, cBody)
	case 1, 3: // Asas planas ---
		ebitenutil.DrawRect(screen, x+4, y+2, 4, 3, cBody)
		ebitenutil.DrawRect(screen, x+3, y+3, 1, 1, cBeak)
		ebitenutil.DrawRect(screen, x-1, y+2, 5, 2, cBody)
		ebitenutil.DrawRect(screen, x+8, y+2, 5, 2, cBody)
	case 2: // Asas para baixo / \
		ebitenutil.DrawRect(screen, x+4, y+2, 4, 3, cBody)
		ebitenutil.DrawRect(screen, x+3, y+3, 1, 1, cBeak)
		ebitenutil.DrawRect(screen, x+2, y+3, 2, 2, cBody)
		ebitenutil.DrawRect(screen, x, y+5, 2, 2, cBody)
		ebitenutil.DrawRect(screen, x+8, y+3, 2, 2, cBody)
		ebitenutil.DrawRect(screen, x+10, y+5, 2, 2, cBody)
	}
}

// drawFlyingGarca desenha uma garça branca esguia em pixel art cruzando a baía
func drawFlyingGarca(screen *ebiten.Image, x, y float64, flapTick int) {
	cWhite := color.RGBA{R: 250, G: 250, B: 255, A: 255}
	cBeak := color.RGBA{R: 255, G: 210, B: 40, A: 255}
	cLegs := color.RGBA{R: 50, G: 50, B: 55, A: 220}

	phase := (flapTick / 10) % 4
	switch phase {
	case 0: // Asas subindo
		ebitenutil.DrawRect(screen, x+5, y+2, 5, 2, cWhite)
		ebitenutil.DrawRect(screen, x+3, y+1, 2, 2, cWhite)
		ebitenutil.DrawRect(screen, x+1, y+2, 2, 1, cBeak)
		ebitenutil.DrawRect(screen, x+10, y+3, 3, 1, cLegs)
		ebitenutil.DrawRect(screen, x+4, y, 3, 2, cWhite)
		ebitenutil.DrawRect(screen, x+2, y-2, 2, 2, cWhite)
	case 1, 3: // Asas planas
		ebitenutil.DrawRect(screen, x+5, y+2, 5, 2, cWhite)
		ebitenutil.DrawRect(screen, x+3, y+2, 2, 2, cWhite)
		ebitenutil.DrawRect(screen, x+1, y+2, 2, 1, cBeak)
		ebitenutil.DrawRect(screen, x+10, y+3, 3, 1, cLegs)
		ebitenutil.DrawRect(screen, x+1, y+1, 6, 2, cWhite)
		ebitenutil.DrawRect(screen, x+8, y+1, 4, 2, cWhite)
	case 2: // Asas descendo
		ebitenutil.DrawRect(screen, x+5, y+2, 5, 2, cWhite)
		ebitenutil.DrawRect(screen, x+3, y+2, 2, 2, cWhite)
		ebitenutil.DrawRect(screen, x+1, y+2, 2, 1, cBeak)
		ebitenutil.DrawRect(screen, x+10, y+3, 3, 1, cLegs)
		ebitenutil.DrawRect(screen, x+4, y+3, 3, 2, cWhite)
		ebitenutil.DrawRect(screen, x+2, y+5, 2, 2, cWhite)
	}
}

// drawBelemSkyBirds desenha a revoada de urubus e garças sobre Belém
func drawBelemSkyBirds(screen *ebiten.Image, screenWidth float64, ticks int) {
	// 4 Urubus voando em círculos e planando no céu alto do Ver-o-Peso
	for i := 0; i < 4; i++ {
		baseSpeed := 0.7 + float64(i)*0.18
		rawX := float64(ticks)*baseSpeed + float64(i*85)
		ux := screenWidth + 20.0 - math.Mod(rawX, screenWidth+80.0)
		uy := 28.0 + float64(i*12) + math.Sin(float64(ticks+i*35)*0.03)*5.0
		drawFlyingUrubu(screen, ux, uy, ticks+i*17)
	}

	// 3 Garças Brancas voando sobre as águas da Baía do Guajará
	for i := 0; i < 3; i++ {
		baseSpeed := 1.1 + float64(i)*0.15
		rawX := float64(ticks)*baseSpeed + float64(i*115)
		gx := math.Mod(rawX, screenWidth+90.0) - 30.0
		gy := 105.0 + float64(i*11) + math.Sin(float64(ticks+i*40)*0.04)*4.0
		drawFlyingGarca(screen, gx, gy, ticks+i*23)
	}
}

// DrawTitleCoverScreen renderiza a tela de abertura oficial limpa (sem menu) estilo Pitfall / Super Metroid / Contra
func DrawTitleCoverScreen(screen *ebiten.Image, screenWidth, screenHeight float64, ticks int, version string) {
	// 1. Imagem de fundo de Belém (Mercado do Ver-o-Peso ao pôr do sol)
	tImg := getTitleScreenImage()
	if tImg != nil {
		op := &ebiten.DrawImageOptions{}
		bounds := tImg.Bounds()
		scaleX := screenWidth / float64(bounds.Dx())
		scaleY := screenHeight / float64(bounds.Dy())
		op.GeoM.Scale(scaleX, scaleY)
		screen.DrawImage(tImg, op)
	} else {
		ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{R: 10, G: 16, B: 24, A: 255})
	}

	// 2. Reflexos cintilantes do pôr do sol na água da Baía do Guajará
	for i := 0; i < 8; i++ {
		glX := 160.0 + float64(i*18) + math.Sin(float64(ticks+i*15)*0.08)*12.0
		glY := 140.0 + math.Sin(float64(ticks+i*22)*0.06)*16.0
		alpha := uint8(120 + math.Sin(float64(ticks+i*20)*0.1)*100)
		ebitenutil.DrawRect(screen, glX, glY, 3, 1, color.RGBA{R: 255, G: 220, B: 130, A: alpha})
	}

	// 3. Urubus e Garças voando pelo céu de Belém do Pará
	drawBelemSkyBirds(screen, screenWidth, ticks)

	// 4. Logo Principal do Jogo em destaque superior épico
	logoW := 304.0
	logoH := 30.0
	logoX := (screenWidth - logoW) / 2.0
	logoY := 10.0

	ebitenutil.DrawRect(screen, logoX+2, logoY+2, logoW, logoH, color.RGBA{R: 0, G: 0, B: 0, A: 170})
	ebitenutil.DrawRect(screen, logoX, logoY, logoW, logoH, color.RGBA{R: 12, G: 18, B: 30, A: 245})
	ebitenutil.DrawRect(screen, logoX, logoY, logoW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, logoX, logoY+logoH, logoW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, logoX, logoY, 2, logoH, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, logoX+logoW, logoY, 2, logoH+2, color.RGBA{R: 250, G: 205, B: 55, A: 255})

	ebitenutil.DebugPrintAt(screen, "★ P A I D ' E G U A   R U N N E R ★", int(logoX)+46, int(logoY)+5)
	ebitenutil.DebugPrintAt(screen, "UMA AVENTURA EM BELEM DO PARA", int(logoX)+54, int(logoY)+17)

	// Badge com a versão no canto superior direito
	verW := float64(len(version)*6 + 12)
	verX := screenWidth - verW - 6.0
	ebitenutil.DrawRect(screen, verX, 10.0, verW, 14, color.RGBA{R: 16, G: 26, B: 44, A: 235})
	ebitenutil.DrawRect(screen, verX, 10.0, verW, 1, color.RGBA{R: 250, G: 205, B: 55, A: 200})
	ebitenutil.DebugPrintAt(screen, version, int(verX)+6, 11)

	// 5. Chamada de ação estilo PRESS START BUTTON (em português paraense pulsante)
	startW := 250.0
	startX := (screenWidth - startW) / 2.0
	startY := 115.0

	if (ticks/24)%2 == 0 {
		ebitenutil.DrawRect(screen, startX, startY-2, startW, 18, color.RGBA{R: 8, G: 14, B: 24, A: 210})
		ebitenutil.DrawRect(screen, startX, startY-2, startW, 1, color.RGBA{R: 250, G: 205, B: 55, A: 255})
		ebitenutil.DrawRect(screen, startX, startY+16, startW, 1, color.RGBA{R: 250, G: 205, B: 55, A: 255})
		ebitenutil.DebugPrintAt(screen, ">> APERTE ENTER OU TOQUE NA TELA <<", int(startX)+18, int(startY)+2)
	} else {
		ebitenutil.DrawRect(screen, startX, startY-2, startW, 18, color.RGBA{R: 8, G: 14, B: 24, A: 140})
		ebitenutil.DebugPrintAt(screen, "   APERTE ENTER OU TOQUE NA TELA   ", int(startX)+18, int(startY)+2)
	}

	// 6. Rodapé clássico de abertura estilo Pitfall / Contra / Super Metroid
	botH := 26.0
	botY := screenHeight - botH
	ebitenutil.DrawRect(screen, 0, botY, screenWidth, botH, color.RGBA{R: 8, G: 12, B: 22, A: 245})
	ebitenutil.DrawRect(screen, 0, botY, screenWidth, 1, color.RGBA{R: 250, G: 205, B: 55, A: 255})

	ebitenutil.DebugPrintAt(screen, "(C) 2026 LUCIVALDO JUNIOR & NEXUS AI  |  [C] CREDITOS", 14, int(botY)+4)
	ebitenutil.DebugPrintAt(screen, "BELEM DO PARA - BRASIL  |  DESENVOLVIDO EM GO + EBITENGINE", 12, int(botY)+14)
}

// DrawTitleIntro renderiza o menu de opções interativo que surge APÓS o jogador apertar Enter/clicar na tela de título
func DrawTitleIntro(screen *ebiten.Image, screenWidth, screenHeight float64, ticks int, isAudioPlaying bool, selectedIndex int, isMuted bool, speedLabel string, heroName string, version string) {
	// 1. Imagem de fundo de Belém (Ver-o-Peso ao pôr do sol)
	tImg := getTitleScreenImage()
	if tImg != nil {
		op := &ebiten.DrawImageOptions{}
		bounds := tImg.Bounds()
		scaleX := screenWidth / float64(bounds.Dx())
		scaleY := screenHeight / float64(bounds.Dy())
		op.GeoM.Scale(scaleX, scaleY)
		screen.DrawImage(tImg, op)
	} else {
		ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{R: 10, G: 16, B: 24, A: 255})
	}

	// 2. Reflexos cintilantes na água
	for i := 0; i < 8; i++ {
		glX := 160.0 + float64(i*18) + math.Sin(float64(ticks+i*15)*0.08)*12.0
		glY := 140.0 + math.Sin(float64(ticks+i*22)*0.06)*16.0
		alpha := uint8(120 + math.Sin(float64(ticks+i*20)*0.1)*100)
		ebitenutil.DrawRect(screen, glX, glY, 3, 1, color.RGBA{R: 255, G: 220, B: 130, A: alpha})
	}

	// 3. Urubus e Garças voando pelo céu de Belém
	drawBelemSkyBirds(screen, screenWidth, ticks)

	// 4. Logo Principal do Jogo
	logoW := 294.0
	logoH := 26.0
	logoX := (screenWidth - logoW) / 2.0
	logoY := 5.0

	ebitenutil.DrawRect(screen, logoX+2, logoY+2, logoW, logoH, color.RGBA{R: 0, G: 0, B: 0, A: 160})
	ebitenutil.DrawRect(screen, logoX, logoY, logoW, logoH, color.RGBA{R: 10, G: 18, B: 28, A: 245})
	ebitenutil.DrawRect(screen, logoX, logoY, logoW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, logoX, logoY+logoH, logoW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, logoX, logoY, 2, logoH, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, logoX+logoW, logoY, 2, logoH+2, color.RGBA{R: 250, G: 205, B: 55, A: 255})

	ebitenutil.DebugPrintAt(screen, "★ P A I D E G U A   G A M E ★", int(logoX)+52, int(logoY)+3)
	ebitenutil.DebugPrintAt(screen, "UMA AVENTURA PELA CIDADE DE BELEM DO PARA", int(logoX)+18, int(logoY)+14)

	// 5. Menu Principal Interativo centralizado
	boxW := 236.0
	boxH := 94.0
	boxX := (screenWidth - boxW) / 2.0
	boxY := 35.0

	ebitenutil.DrawRect(screen, boxX+2, boxY+2, boxW, boxH, color.RGBA{R: 0, G: 0, B: 0, A: 140})
	ebitenutil.DrawRect(screen, boxX, boxY, boxW, boxH, color.RGBA{R: 12, G: 18, B: 32, A: 245})
	ebitenutil.DrawRect(screen, boxX, boxY, boxW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY+boxH, boxW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY, 2, boxH, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX+boxW, boxY, 2, boxH+2, color.RGBA{R: 250, G: 205, B: 55, A: 255})

	bx := int(boxX)
	by := int(boxY)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("★ MENU DE AVENTURA (%s) ★", version), bx+28, by+5)
	ebitenutil.DrawRect(screen, boxX+10, boxY+17, boxW-20, 1, color.RGBA{R: 250, G: 205, B: 55, A: 160})

	soundStatus := "SOM: [ LIGADO ]"
	if isMuted {
		soundStatus = "SOM: [ MUDO ]  "
	}

	options := []string{
		"INICIAR AVENTURA",
		fmt.Sprintf("HEROI: [ %s ]", heroName),
		soundStatus,
		fmt.Sprintf("VELOCIDADE: [ %s ]", speedLabel),
		"VER CREDITOS",
	}

	startY := by + 21
	for i, opt := range options {
		itemY := startY + i*14
		if i == selectedIndex {
			ebitenutil.DrawRect(screen, boxX+6, float64(itemY-1), boxW-12, 13, color.RGBA{R: 245, G: 185, B: 45, A: 80})
			ebitenutil.DrawRect(screen, boxX+6, float64(itemY-1), 3, 13, color.RGBA{R: 250, G: 210, B: 50, A: 255})
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("> %s <", opt), bx+10, itemY)
		} else {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("  %s", opt), bx+10, itemY)
		}
	}

	// 6. Faixa inferior ampla com narrativa cultural de Belém do Pará
	bannerH := 45.0
	bannerY := screenHeight - bannerH
	ebitenutil.DrawRect(screen, 0, bannerY, screenWidth, bannerH, color.RGBA{R: 6, G: 12, B: 22, A: 245})
	ebitenutil.DrawRect(screen, 0, bannerY, screenWidth, 1, color.RGBA{R: 250, G: 205, B: 55, A: 255})

	if isAudioPlaying {
		ebitenutil.DebugPrintAt(screen, "★ TRILHA SONORA: PINDUCA - A DANCA DO CARIMBO ★", 24, int(bannerY)+3)
	} else {
		ebitenutil.DebugPrintAt(screen, "★ PAID'EGUA RUNNER: UMA AVENTURA EM BELEM DO PARA ★", 10, int(bannerY)+3)
	}

	// Letreiro de ação e cultura paraense
	loreTicker := "★ MISSAO: Explore o Mercado do Ver-o-Peso, o cais da Estacao das Docas e o Theatro da Paz!  " +
		"★ DOIS HEROIS: O destemido Garoto Curumim ou a guardiao Onca-Pintada!  " +
		"★ CULTURA: Saboreie acai, curta o carimbo de Belem e conquiste as reliquias sagradas!  "
	loreTextWidth := len(loreTicker) * 6
	loreOffset := ticks % loreTextWidth
	loreBaseX := 52 - loreOffset
	ebitenutil.DebugPrintAt(screen, loreTicker, loreBaseX, int(bannerY)+16)
	if loreBaseX+loreTextWidth < int(screenWidth) {
		ebitenutil.DebugPrintAt(screen, loreTicker, loreBaseX+loreTextWidth, int(bannerY)+16)
	}
	ebitenutil.DrawRect(screen, 0, bannerY+15, 48, 14, color.RGBA{R: 16, G: 32, B: 52, A: 255})
	ebitenutil.DrawRect(screen, 48, bannerY+15, 1, 14, color.RGBA{R: 45, G: 215, B: 175, A: 255})
	ebitenutil.DebugPrintAt(screen, "★BELEM", 4, int(bannerY)+16)

	// Dica clara de navegação
	ebitenutil.DebugPrintAt(screen, "[CIMA/BAIXO] Navegar | [ENTER/ESPACO] Escolher | [ESC] Voltar", 8, int(bannerY)+30)
}

// DrawCharacterSelectScreen renderiza a tela dedicada de escolha de personagem (Garoto vs Onça)
func DrawCharacterSelectScreen(screen *ebiten.Image, screenWidth, screenHeight float64, ticks int, selectedHero int) {
	// 1. Fundo noturno amazônico com vaga-lumes flutuantes
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{R: 12, G: 16, B: 26, A: 255})

	for i := 0; i < 10; i++ {
		fx := float64((i*37 + ticks/2) % int(screenWidth))
		fy := 35.0 + math.Sin(float64(ticks+i*25)*0.07)*12.0 + float64(i*12)
		alpha := uint8(140 + math.Sin(float64(ticks+i*30)*0.1)*80)
		ebitenutil.DrawRect(screen, fx, fy, 2, 2, color.RGBA{R: 150, G: 255, B: 110, A: alpha})
	}

	// 2. Banner Superior
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, 27, color.RGBA{R: 8, G: 12, B: 22, A: 245})
	ebitenutil.DrawRect(screen, 0, 27, screenWidth, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DebugPrintAt(screen, "★ ESCOLHA SEU HEROI DA AMAZONIA ★", 56, 4)
	ebitenutil.DebugPrintAt(screen, "Selecione o protagonista para a jornada no Marajo e Belem", 14, 15)

	cardW := 150.0
	cardH := 142.0
	cardY := 32.0

	// ==========================================
	// CARD 0: GAROTO CURUMIM (Esquerda)
	// ==========================================
	c0X := 14.0
	is0Selected := selectedHero == 0

	bg0 := color.RGBA{R: 16, G: 24, B: 38, A: 245}
	border0 := color.RGBA{R: 80, G: 95, B: 120, A: 220}
	if is0Selected {
		bg0 = color.RGBA{R: 20, G: 34, B: 58, A: 250}
		border0 = color.RGBA{R: 255, G: 215, B: 50, A: 255}
	}
	ebitenutil.DrawRect(screen, c0X, cardY, cardW, cardH, bg0)
	ebitenutil.DrawRect(screen, c0X, cardY, cardW, 2, border0)
	ebitenutil.DrawRect(screen, c0X, cardY+cardH, cardW, 2, border0)
	ebitenutil.DrawRect(screen, c0X, cardY, 2, cardH, border0)
	ebitenutil.DrawRect(screen, c0X+cardW, cardY, 2, cardH+2, border0)

	// Cabeçalho do Card 0
	if is0Selected {
		ebitenutil.DrawRect(screen, c0X+4, cardY+4, cardW-8, 14, color.RGBA{R: 245, G: 190, B: 40, A: 220})
		ebitenutil.DebugPrintAt(screen, "► GAROTO CURUMIM ◄", int(c0X)+18, int(cardY)+5)
	} else {
		ebitenutil.DebugPrintAt(screen, "  GAROTO CURUMIM", int(c0X)+22, int(cardY)+5)
	}
	ebitenutil.DrawRect(screen, c0X+6, cardY+20, cardW-12, 1, color.RGBA{R: 200, G: 180, B: 80, A: 120})

	// Pedestal de pedra com musgo
	pedX := c0X + 8
	pedY := cardY + 54
	ebitenutil.DrawRect(screen, pedX, pedY+26, 44, 6, color.RGBA{R: 70, G: 65, B: 60, A: 255})
	ebitenutil.DrawRect(screen, pedX+2, pedY+24, 40, 3, color.RGBA{R: 90, G: 120, B: 50, A: 255})

	// Sprite Preview do Garoto no pedestal (com animação suave)
	breathe := float64((ticks / 20) % 2)
	gx := pedX + 12
	gy := pedY - 4 + breathe

	cSkin := color.RGBA{R: 228, G: 168, B: 125, A: 255}
	cHair := color.RGBA{R: 30, G: 25, B: 28, A: 255}
	cBandana := color.RGBA{R: 235, G: 50, B: 45, A: 255}
	cShirt := color.RGBA{R: 35, G: 130, B: 220, A: 255}
	cPants := color.RGBA{R: 175, G: 145, B: 100, A: 255}
	cWood := color.RGBA{R: 145, G: 85, B: 40, A: 255}

	ebitenutil.DrawRect(screen, gx+4, gy, 8, 8, cSkin)
	ebitenutil.DrawRect(screen, gx+2, gy-2, 12, 3, cHair)
	ebitenutil.DrawRect(screen, gx+3, gy+1, 10, 2, cBandana)
	ebitenutil.DrawRect(screen, gx+9, gy+3, 2, 2, color.RGBA{R: 10, G: 10, B: 10, A: 255})
	ebitenutil.DrawRect(screen, gx+3, gy+8, 9, 8, cShirt)
	ebitenutil.DrawRect(screen, gx+12, gy+7, 3, 5, cWood)
	ebitenutil.DrawRect(screen, gx+13, gy+6, 4, 2, color.RGBA{R: 230, G: 110, B: 50, A: 255})
	ebitenutil.DrawRect(screen, gx+4, gy+16, 3, 8, cPants)
	ebitenutil.DrawRect(screen, gx+9, gy+16, 3, 8, cPants)
	ebitenutil.DrawRect(screen, gx+4, gy+24, 4, 3, cSkin)
	ebitenutil.DrawRect(screen, gx+9, gy+24, 4, 3, cSkin)

	// Ficha de Atributos do Garoto (à direita no card)
	tx0 := int(c0X) + 54
	ebitenutil.DebugPrintAt(screen, "Arma: Baladeira", tx0, int(cardY)+25)
	ebitenutil.DebugPrintAt(screen, "Tiro: Acai veloz", tx0, int(cardY)+37)
	ebitenutil.DebugPrintAt(screen, "Pulo: Salto Duplo", tx0, int(cardY)+49)
	ebitenutil.DebugPrintAt(screen, "Esquiva: Agil no ar", tx0, int(cardY)+61)
	ebitenutil.DebugPrintAt(screen, "Modo: Aventureiro", tx0, int(cardY)+73)

	// Divisória inferior do card
	ebitenutil.DrawRect(screen, c0X+6, cardY+90, cardW-12, 1, color.RGBA{R: 150, G: 160, B: 180, A: 100})
	ebitenutil.DebugPrintAt(screen, "\"Bora la, maninho!\"", int(c0X)+14, int(cardY)+95)
	ebitenutil.DebugPrintAt(screen, "\"Nao tem jacare que pegue\"", int(c0X)+4, int(cardY)+107)

	if is0Selected {
		ebitenutil.DrawRect(screen, c0X+12, cardY+124, cardW-24, 13, color.RGBA{R: 245, G: 185, B: 45, A: 85})
		ebitenutil.DebugPrintAt(screen, "★ SELECIONADO ★", int(c0X)+24, int(cardY)+125)
	} else {
		ebitenutil.DebugPrintAt(screen, "[ ENTER P/ ESCOLHER ]", int(c0X)+12, int(cardY)+125)
	}

	// ==========================================
	// CARD 1: ONÇA PINTADA (Direita)
	// ==========================================
	c1X := 176.0
	is1Selected := selectedHero == 1

	bg1 := color.RGBA{R: 26, G: 20, B: 14, A: 245}
	border1 := color.RGBA{R: 110, G: 85, B: 60, A: 220}
	if is1Selected {
		bg1 = color.RGBA{R: 42, G: 30, B: 18, A: 250}
		border1 = color.RGBA{R: 255, G: 215, B: 50, A: 255}
	}
	ebitenutil.DrawRect(screen, c1X, cardY, cardW, cardH, bg1)
	ebitenutil.DrawRect(screen, c1X, cardY, cardW, 2, border1)
	ebitenutil.DrawRect(screen, c1X, cardY+cardH, cardW, 2, border1)
	ebitenutil.DrawRect(screen, c1X, cardY, 2, cardH, border1)
	ebitenutil.DrawRect(screen, c1X+cardW, cardY, 2, cardH+2, border1)

	// Cabeçalho do Card 1
	if is1Selected {
		ebitenutil.DrawRect(screen, c1X+4, cardY+4, cardW-8, 14, color.RGBA{R: 245, G: 190, B: 40, A: 220})
		ebitenutil.DebugPrintAt(screen, "► ONCA PINTADA ◄", int(c1X)+26, int(cardY)+5)
	} else {
		ebitenutil.DebugPrintAt(screen, "  ONCA PINTADA", int(c1X)+30, int(cardY)+5)
	}
	ebitenutil.DrawRect(screen, c1X+6, cardY+20, cardW-12, 1, color.RGBA{R: 200, G: 180, B: 80, A: 120})

	// Pedestal de pedra
	ped1X := c1X + 8
	ebitenutil.DrawRect(screen, ped1X, pedY+26, 46, 6, color.RGBA{R: 70, G: 65, B: 60, A: 255})
	ebitenutil.DrawRect(screen, ped1X+2, pedY+24, 42, 3, color.RGBA{R: 90, G: 120, B: 50, A: 255})

	// Sprite Preview da Onça no pedestal
	ox := ped1X + 4
	oy := pedY + 4

	cGold := color.RGBA{R: 235, G: 160, B: 35, A: 255}
	cGoldLight := color.RGBA{R: 248, G: 185, B: 65, A: 255}
	cSpotBlack := color.RGBA{R: 35, G: 18, B: 10, A: 255}
	cCream := color.RGBA{R: 245, G: 230, B: 195, A: 255}

	ebitenutil.DrawRect(screen, ox+2, oy+2, 24, 10, cGold)
	ebitenutil.DrawRect(screen, ox+4, oy+3, 20, 7, cGoldLight)
	ebitenutil.DrawRect(screen, ox+5, oy+10, 18, 3, cCream)
	ebitenutil.DrawRect(screen, ox+6, oy+4, 4, 3, cSpotBlack)
	ebitenutil.DrawRect(screen, ox+14, oy+5, 4, 3, cSpotBlack)
	ebitenutil.DrawRect(screen, ox+22, oy, 9, 9, cGold)
	ebitenutil.DrawRect(screen, ox+27, oy+3, 2, 2, color.RGBA{R: 85, G: 205, B: 110, A: 255})
	ebitenutil.DrawRect(screen, ox+29, oy+4, 2, 2, color.RGBA{R: 70, G: 30, B: 30, A: 255})
	ebitenutil.DrawRect(screen, ox+23, oy-3, 3, 3, cSpotBlack)
	ebitenutil.DrawRect(screen, ox+3, oy+12, 5, 8, cGold)
	ebitenutil.DrawRect(screen, ox+18, oy+12, 5, 8, cGold)
	ebitenutil.DrawRect(screen, ox+2, oy+18, 6, 2, cCream)
	ebitenutil.DrawRect(screen, ox+18, oy+18, 6, 2, cCream)
	oTail := math.Sin(float64(ticks)*0.2) * 2.0
	ebitenutil.DrawRect(screen, ox-3, oy+4+oTail, 5, 3, cGold)
	ebitenutil.DrawRect(screen, ox-5, oy+2+oTail, 3, 3, cSpotBlack)

	// Ficha de Atributos da Onça
	tx1 := int(c1X) + 54
	ebitenutil.DebugPrintAt(screen, "Arma: Rugido", tx1, int(cardY)+25)
	ebitenutil.DebugPrintAt(screen, "Tiro: Onda sonica", tx1, int(cardY)+37)
	ebitenutil.DebugPrintAt(screen, "Pulo: Bote Feroz", tx1, int(cardY)+49)
	ebitenutil.DebugPrintAt(screen, "Passada: 2.4x rapida", tx1, int(cardY)+61)
	ebitenutil.DebugPrintAt(screen, "Modo: Predadora", tx1, int(cardY)+73)

	// Divisória inferior do card
	ebitenutil.DrawRect(screen, c1X+6, cardY+90, cardW-12, 1, color.RGBA{R: 150, G: 160, B: 180, A: 100})
	ebitenutil.DebugPrintAt(screen, "\"RRRAUW! Rainha!\"", int(c1X)+20, int(cardY)+95)
	ebitenutil.DebugPrintAt(screen, "\"Ninguem segura a onca!\"", int(c1X)+6, int(cardY)+107)

	if is1Selected {
		ebitenutil.DrawRect(screen, c1X+12, cardY+124, cardW-24, 13, color.RGBA{R: 245, G: 185, B: 45, A: 85})
		ebitenutil.DebugPrintAt(screen, "★ SELECIONADO ★", int(c1X)+24, int(cardY)+125)
	} else {
		ebitenutil.DebugPrintAt(screen, "[ ENTER P/ ESCOLHER ]", int(c1X)+12, int(cardY)+125)
	}

	// ==========================================
	// BARRA INFERIOR DE INSTRUÇÃO E CONFIRMAÇÃO
	// ==========================================
	botY := 176.0
	botH := 34.0
	ebitenutil.DrawRect(screen, 0, botY, screenWidth, botH, color.RGBA{R: 8, G: 12, B: 20, A: 250})
	ebitenutil.DrawRect(screen, 0, botY, screenWidth, 1, color.RGBA{R: 250, G: 205, B: 55, A: 255})

	if (ticks/24)%2 == 0 {
		ebitenutil.DebugPrintAt(screen, ">> [ENTER / TOQUE] CONFIRMAR E INICIAR <<", 22, int(botY)+6)
	} else {
		ebitenutil.DebugPrintAt(screen, "   [ENTER / TOQUE] CONFIRMAR E INICIAR   ", 22, int(botY)+6)
	}
	ebitenutil.DebugPrintAt(screen, "[ESQ/DIR ou A/D] Alternar Heroi  |  [ESC] Voltar ao Menu", 16, int(botY)+19)
}

func DrawTitleScreen(screen *ebiten.Image, screenWidth, screenHeight float64, ticks int) {
	// Fundo com escurecimento suave para destacar a onça correndo e o cenário ao redor
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{R: 10, G: 12, B: 20, A: 110})

	boxW := 195.0
	boxH := 108.0
	boxX := (screenWidth - boxW) / 2.0
	boxY := (screenHeight - boxH) / 2.0

	ebitenutil.DrawRect(screen, boxX, boxY, boxW, boxH, color.RGBA{R: 18, G: 24, B: 38, A: 245})
	ebitenutil.DrawRect(screen, boxX, boxY, boxW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY+boxH, boxW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY, 2, boxH, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX+boxW, boxY, 2, boxH+2, color.RGBA{R: 250, G: 205, B: 55, A: 255})

	bx := int(boxX)
	by := int(boxY)

	ebitenutil.DebugPrintAt(screen, "★ PAIDEGUA GAME ★", bx+45, by+8)
	ebitenutil.DrawRect(screen, boxX+10, boxY+22, boxW-20, 1, color.RGBA{R: 250, G: 205, B: 55, A: 160})

	ebitenutil.DebugPrintAt(screen, "Dev: Luci Junior & Nexus AI", bx+14, by+26)
	ebitenutil.DebugPrintAt(screen, "Mover: A/D ou Setas | Pulo: W/Cima", bx+8, by+41)
	ebitenutil.DebugPrintAt(screen, "Ataque: Espaco/X | Baixo: S", bx+10, by+54)

	if (ticks/30)%2 == 0 {
		ebitenutil.DebugPrintAt(screen, ">> TOQUE OU APERTE ENTER <<", bx+20, by+72)
	} else {
		ebitenutil.DebugPrintAt(screen, "   TOQUE OU APERTE ENTER   ", bx+20, by+72)
	}

	ebitenutil.DebugPrintAt(screen, "[C] Creditos  |  [ESC] Sair", bx+18, by+91)
}

func DrawCreditsScreen(screen *ebiten.Image, screenWidth, screenHeight float64) {
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{R: 8, G: 10, B: 18, A: 140})

	boxW := 216.0
	boxH := 134.0
	boxX := (screenWidth - boxW) / 2.0
	boxY := (screenHeight - boxH) / 2.0

	ebitenutil.DrawRect(screen, boxX, boxY, boxW, boxH, color.RGBA{R: 20, G: 26, B: 42, A: 245})
	ebitenutil.DrawRect(screen, boxX, boxY, boxW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY+boxH, boxW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY, 2, boxH, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX+boxW, boxY, 2, boxH+2, color.RGBA{R: 250, G: 205, B: 55, A: 255})

	bx := int(boxX)
	by := int(boxY)

	ebitenutil.DebugPrintAt(screen, "★ CREDITOS - PAID'EGUA RUNNER ★", bx+12, by+8)
	ebitenutil.DrawRect(screen, boxX+10, boxY+21, boxW-20, 1, color.RGBA{R: 250, G: 205, B: 55, A: 160})

	ebitenutil.DebugPrintAt(screen, "Desenvolvedor: Luci Junior", bx+14, by+26)
	ebitenutil.DebugPrintAt(screen, "Co-criacao IA: Nexus AI (Lucy)", bx+14, by+38)
	ebitenutil.DebugPrintAt(screen, "Linguagem:     Go (Golang)", bx+14, by+50)
	ebitenutil.DebugPrintAt(screen, "Game Engine:   Ebitengine v2", bx+14, by+62)
	ebitenutil.DebugPrintAt(screen, "Trilha Sonora: Pinduca (Carimbo)", bx+14, by+74)
	ebitenutil.DebugPrintAt(screen, "Cenarios:      Belem do Para", bx+14, by+86)

	ebitenutil.DebugPrintAt(screen, "[ENTER / ESC / C] Voltar", bx+30, by+112)
}

func DrawPauseMenu(screen *ebiten.Image, screenWidth, screenHeight float64, selectedIndex int, isMuted bool, speedLabel string) {
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{R: 10, G: 12, B: 20, A: 140})

	boxW := 214.0
	boxH := 136.0
	boxX := (screenWidth - boxW) / 2.0
	boxY := (screenHeight - boxH) / 2.0

	ebitenutil.DrawRect(screen, boxX, boxY, boxW, boxH, color.RGBA{R: 22, G: 28, B: 44, A: 245})
	ebitenutil.DrawRect(screen, boxX, boxY, boxW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY+boxH, boxW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY, 2, boxH, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX+boxW, boxY, 2, boxH+2, color.RGBA{R: 250, G: 205, B: 55, A: 255})

	bx := int(boxX)
	by := int(boxY)

	ebitenutil.DebugPrintAt(screen, "PAUSA - PAID'EGUA RUNNER", bx+34, by+8)
	ebitenutil.DrawRect(screen, boxX+10, boxY+22, boxW-20, 1, color.RGBA{R: 250, G: 205, B: 55, A: 160})

	soundStatus := "Som: [ LIGADO ]"
	if isMuted {
		soundStatus = "Som: [ MUDO ]  "
	}

	options := []string{
		"Continuar",
		soundStatus,
		fmt.Sprintf("Velocidade: [ %s ]", speedLabel),
		"Reiniciar Aventura",
		"Ver Creditos",
		"Fechar o Jogo",
	}

	startY := by + 28
	for i, opt := range options {
		itemY := startY + i*15
		if i == selectedIndex {
			ebitenutil.DrawRect(screen, boxX+8, float64(itemY-2), boxW-16, 13, color.RGBA{R: 245, G: 185, B: 45, A: 70})
			ebitenutil.DrawRect(screen, boxX+8, float64(itemY-2), 3, 13, color.RGBA{R: 250, G: 210, B: 50, A: 255})
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("> %s <", opt), bx+16, itemY)
		} else {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("  %s", opt), bx+16, itemY)
		}
	}

	ebitenutil.DebugPrintAt(screen, "[CIMA/BAIXO] | [ENTER/DIR] | [ESC]", bx+8, by+int(boxH)-12)
}

func DrawIndigenousWarrior(screen *ebiten.Image, x, y float64, ticks int) {
	cSkin := color.RGBA{R: 195, G: 125, B: 85, A: 255}
	cSkinDark := color.RGBA{R: 165, G: 98, B: 65, A: 255}
	cHair := color.RGBA{R: 18, G: 16, B: 20, A: 255}
	cFeatherRed := color.RGBA{R: 235, G: 45, B: 40, A: 255}
	cFeatherBlue := color.RGBA{R: 35, G: 125, B: 240, A: 255}
	cFeatherYellow := color.RGBA{R: 250, G: 215, B: 45, A: 255}
	cFeatherWhite := color.RGBA{R: 245, G: 245, B: 250, A: 255}
	cHeadband := color.RGBA{R: 205, G: 165, B: 90, A: 255}
	cUrucum := color.RGBA{R: 215, G: 30, B: 30, A: 255}
	cJenipapo := color.RGBA{R: 25, G: 20, B: 25, A: 255}
	cSkirt := color.RGBA{R: 180, G: 140, B: 75, A: 255}
	cNecklace := color.RGBA{R: 235, G: 225, B: 210, A: 255}

	sway := math.Sin(float64(ticks)*0.09) * 1.0

	// 1. Cocar Amazônico de Penas Coloridas em Leque
	ebitenutil.DrawRect(screen, x+21, y+2+sway, 5, 18, cFeatherYellow)
	ebitenutil.DrawRect(screen, x+22, y-2+sway, 3, 5, cFeatherRed)
	ebitenutil.DrawRect(screen, x+15, y+5+sway, 5, 16, cFeatherBlue)
	ebitenutil.DrawRect(screen, x+27, y+5+sway, 5, 16, cFeatherBlue)
	ebitenutil.DrawRect(screen, x+10, y+9+sway, 4, 13, cFeatherRed)
	ebitenutil.DrawRect(screen, x+33, y+9+sway, 4, 13, cFeatherRed)
	ebitenutil.DrawRect(screen, x+6, y+14+sway, 4, 10, cFeatherWhite)
	ebitenutil.DrawRect(screen, x+38, y+14+sway, 4, 10, cFeatherWhite)

	// Testeira de palha trançada com grafismos
	ebitenutil.DrawRect(screen, x+9, y+21+sway, 30, 4, cHeadband)
	for d := 0; d < 6; d++ {
		dx := x + 11 + float64(d*5)
		ebitenutil.DrawRect(screen, dx, y+22+sway, 2, 2, cUrucum)
	}

	// 2. Cabeça e Cabelo Preto Liso
	ebitenutil.DrawRect(screen, x+10, y+24+sway, 4, 25, cHair)
	ebitenutil.DrawRect(screen, x+34, y+24+sway, 4, 25, cHair)
	ebitenutil.DrawRect(screen, x+14, y+24+sway, 20, 18, cSkin)
	ebitenutil.DrawRect(screen, x+17, y+41+sway, 14, 3, cSkin)

	// Pintura facial de Urucum e Olhos
	ebitenutil.DrawRect(screen, x+15, y+32+sway, 6, 3, cUrucum)
	ebitenutil.DrawRect(screen, x+27, y+32+sway, 6, 3, cUrucum)
	ebitenutil.DrawRect(screen, x+17, y+29+sway, 4, 2, color.RGBA{R: 250, G: 250, B: 250, A: 255})
	ebitenutil.DrawRect(screen, x+19, y+29+sway, 2, 2, cJenipapo)
	ebitenutil.DrawRect(screen, x+27, y+29+sway, 4, 2, color.RGBA{R: 250, G: 250, B: 250, A: 255})
	ebitenutil.DrawRect(screen, x+27, y+29+sway, 2, 2, cJenipapo)

	// Nariz e Sorriso de comemoração
	ebitenutil.DrawRect(screen, x+23, y+33+sway, 2, 4, cSkinDark)
	ebitenutil.DrawRect(screen, x+20, y+38+sway, 8, 3, color.RGBA{R: 80, G: 20, B: 20, A: 255})
	ebitenutil.DrawRect(screen, x+21, y+38+sway, 6, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255})

	// 3. Pescoço e Colar de Sementes de Açaí
	ebitenutil.DrawRect(screen, x+20, y+43+sway, 8, 4, cSkin)
	ebitenutil.DrawRect(screen, x+16, y+45+sway, 16, 2, cNecklace)
	ebitenutil.DrawRect(screen, x+23, y+47+sway, 2, 3, cFeatherWhite)

	// 4. Peitoral com Grafismos de Jenipapo
	ebitenutil.DrawRect(screen, x+15, y+47+sway, 18, 18, cSkin)
	ebitenutil.DrawRect(screen, x+17, y+52+sway, 14, 2, cJenipapo)
	ebitenutil.DrawRect(screen, x+19, y+56+sway, 10, 2, cJenipapo)

	// 5. Braço Esquerdo (apoiado, com braçadeira)
	ebitenutil.DrawRect(screen, x+9, y+48+sway, 6, 16, cSkin)
	ebitenutil.DrawRect(screen, x+9, y+52+sway, 6, 3, cFeatherRed)
	ebitenutil.DrawRect(screen, x+11, y+64+sway, 5, 7, cSkin)

	// 6. Braço Direito Erguido em Sinal de Vitória / Saudação
	ebitenutil.DrawRect(screen, x+33, y+48+sway, 6, 6, cSkin)
	ebitenutil.DrawRect(screen, x+37, y+34+sway, 6, 15, cSkin)
	ebitenutil.DrawRect(screen, x+37, y+40+sway, 6, 3, cFeatherYellow)
	ebitenutil.DrawRect(screen, x+40, y+20+sway, 5, 15, cSkin)
	ebitenutil.DrawRect(screen, x+39, y+13+sway, 8, 8, cSkin)
	ebitenutil.DrawRect(screen, x+39, y+10+sway, 2, 4, cSkin)
	ebitenutil.DrawRect(screen, x+42, y+9+sway, 2, 5, cSkin)
	ebitenutil.DrawRect(screen, x+45, y+10+sway, 2, 4, cSkin)

	// 7. Saiote de Palha e Fibras
	ebitenutil.DrawRect(screen, x+13, y+65+sway, 22, 12, cSkirt)
	for f := 0; f < 5; f++ {
		fx := x + 15 + float64(f*4)
		ebitenutil.DrawRect(screen, fx, y+66+sway, 2, 10, color.RGBA{R: 150, G: 110, B: 55, A: 255})
	}

	// 8. Pernas e Tornozeleiras
	ebitenutil.DrawRect(screen, x+16, y+77+sway, 6, 14, cSkin)
	ebitenutil.DrawRect(screen, x+26, y+77+sway, 6, 14, cSkin)
	ebitenutil.DrawRect(screen, x+16, y+88+sway, 6, 2, cHeadband)
	ebitenutil.DrawRect(screen, x+26, y+88+sway, 6, 2, cHeadband)
	ebitenutil.DrawRect(screen, x+14, y+91+sway, 8, 3, cSkinDark)
	ebitenutil.DrawRect(screen, x+26, y+91+sway, 8, 3, cSkinDark)

	// 9. Plaqueta Identificadora do Tuxaua
	ebitenutil.DrawRect(screen, x+6, y+96, 36, 12, color.RGBA{R: 14, G: 20, B: 34, A: 245})
	ebitenutil.DrawRect(screen, x+6, y+96, 36, 1, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, x+6, y+107, 36, 1, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DebugPrintAt(screen, "TUXAUA", int(x)+9, int(y)+97)
}

func DrawStageCompleteScreen(screen *ebiten.Image, screenWidth, screenHeight float64, stage int, score int, lives int, ticks int) {
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{R: 10, G: 15, B: 25, A: 220})

	boxX := 12.0
	boxY := 8.0
	boxW := screenWidth - 24.0
	boxH := screenHeight - 16.0

	ebitenutil.DrawRect(screen, boxX, boxY, boxW, boxH, color.RGBA{R: 22, G: 28, B: 44, A: 245})
	ebitenutil.DrawRect(screen, boxX, boxY, boxW, 2, color.RGBA{R: 250, G: 210, B: 50, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY+boxH, boxW, 2, color.RGBA{R: 250, G: 210, B: 50, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY, 2, boxH, color.RGBA{R: 250, G: 210, B: 50, A: 255})
	ebitenutil.DrawRect(screen, boxX+boxW, boxY, 2, boxH+2, color.RGBA{R: 250, G: 210, B: 50, A: 255})

	// Desenha o Guerreiro Indígena Tuxaua no lado esquerdo
	DrawIndigenousWarrior(screen, boxX+6, boxY+22, ticks)

	// Painel de Conteúdo e Felicitações à Direita
	contentX := int(boxX + 66)

	if stage < 3 {
		completedName := "MERCADO DO VER-O-PESO"
		nextName := "ESTACAO DAS DOCAS"
		tuxauaLine1 := "Warana! Mandou paidegua, guerreiro!"
		tuxauaLine2 := "Superou os perigos do Ver-o-Peso!"
		if stage == 2 {
			completedName = "ESTACAO DAS DOCAS"
			nextName = "THEATRO DA PAZ"
			tuxauaLine1 = "Egua mano, que corrida veloz!"
			tuxauaLine2 = "As Docas foram vencidas com bravura!"
		}

		// Faixa do Título Superior
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("★ PARABENS! FASE %d CONCLUIDA! ★", stage), contentX+10, int(boxY)+8)
		ebitenutil.DrawRect(screen, float64(contentX), boxY+22, boxW-72, 1, color.RGBA{R: 250, G: 205, B: 55, A: 160})

		// Balão de Fala do Índio Tuxaua
		balloonX := float64(contentX)
		balloonY := boxY + 28
		balloonW := boxW - 72
		balloonH := 46.0

		// Sombra e fundo do balão
		ebitenutil.DrawRect(screen, balloonX+2, balloonY+2, balloonW, balloonH, color.RGBA{R: 0, G: 0, B: 0, A: 120})
		ebitenutil.DrawRect(screen, balloonX, balloonY, balloonW, balloonH, color.RGBA{R: 16, G: 24, B: 40, A: 245})
		ebitenutil.DrawRect(screen, balloonX, balloonY, balloonW, 1, color.RGBA{R: 250, G: 205, B: 55, A: 240})
		ebitenutil.DrawRect(screen, balloonX, balloonY+balloonH, balloonW, 1, color.RGBA{R: 250, G: 205, B: 55, A: 240})
		ebitenutil.DrawRect(screen, balloonX+balloonW, balloonY, 1, balloonH, color.RGBA{R: 250, G: 205, B: 55, A: 240})

		// Rabicho apontando para o índio
		ebitenutil.DrawRect(screen, balloonX-4, balloonY+12, 4, 3, color.RGBA{R: 250, G: 205, B: 55, A: 240})
		ebitenutil.DrawRect(screen, balloonX-6, balloonY+13, 3, 2, color.RGBA{R: 250, G: 205, B: 55, A: 240})

		ebitenutil.DebugPrintAt(screen, "TUXAUA DA AMAZONIA:", contentX+8, int(balloonY)+5)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("\"%s\"", tuxauaLine1), contentX+8, int(balloonY)+18)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf(" \"%s\"", tuxauaLine2), contentX+8, int(balloonY)+30)

		// Dados da Etapa
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Vencido: %s", completedName), contentX+6, int(boxY)+80)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Proximo: %s", nextName), contentX+6, int(boxY)+95)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Pontos:  %05d pts  |  x%d VIDAS", score, lives), contentX+6, int(boxY)+110)

		// Botões de Ação Piscantes
		if (ticks/25)%2 == 0 {
			ebitenutil.DebugPrintAt(screen, ">> [CLIQUE / ENTER] PROXIMA FASE <<", contentX+2, int(boxY)+130)
		} else {
			ebitenutil.DebugPrintAt(screen, "   [CLIQUE / ENTER] PROXIMA FASE   ", contentX+2, int(boxY)+130)
		}
		ebitenutil.DebugPrintAt(screen, "[R] Recomecar   |   [ESC] Sair", contentX+16, int(boxY)+145)

	} else {
		// Vitória Final da Expedição (Theatro da Paz)
		ebitenutil.DebugPrintAt(screen, "★ VITORIA TOTAL! EXPEDICAO CONCLUIDA! ★", contentX-8, int(boxY)+8)
		ebitenutil.DrawRect(screen, float64(contentX), boxY+22, boxW-72, 1, color.RGBA{R: 250, G: 205, B: 55, A: 160})

		balloonX := float64(contentX)
		balloonY := boxY + 26
		balloonW := boxW - 72
		balloonH := 46.0

		ebitenutil.DrawRect(screen, balloonX+2, balloonY+2, balloonW, balloonH, color.RGBA{R: 0, G: 0, B: 0, A: 120})
		ebitenutil.DrawRect(screen, balloonX, balloonY, balloonW, balloonH, color.RGBA{R: 16, G: 24, B: 40, A: 245})
		ebitenutil.DrawRect(screen, balloonX, balloonY, balloonW, 1, color.RGBA{R: 250, G: 205, B: 55, A: 240})
		ebitenutil.DrawRect(screen, balloonX, balloonY+balloonH, balloonW, 1, color.RGBA{R: 250, G: 205, B: 55, A: 240})
		ebitenutil.DrawRect(screen, balloonX+balloonW, balloonY, 1, balloonH, color.RGBA{R: 250, G: 205, B: 55, A: 240})

		ebitenutil.DrawRect(screen, balloonX-4, balloonY+12, 4, 3, color.RGBA{R: 250, G: 205, B: 55, A: 240})
		ebitenutil.DrawRect(screen, balloonX-6, balloonY+13, 3, 2, color.RGBA{R: 250, G: 205, B: 55, A: 240})

		ebitenutil.DebugPrintAt(screen, "TUXAUA DA AMAZONIA:", contentX+8, int(balloonY)+5)
		ebitenutil.DebugPrintAt(screen, "\"TRIUNFO TOTAL! Voce e a lenda viva!\"", contentX+8, int(balloonY)+18)
		ebitenutil.DebugPrintAt(screen, "\"A Onca Paidegua reina em Belem!\"", contentX+8, int(balloonY)+30)

		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Pontuacao Maxima: %05d pts", score), contentX+6, int(boxY)+78)
		ebitenutil.DebugPrintAt(screen, "Desenvolvedor: Lucivaldo Junior (Luci)", contentX+6, int(boxY)+92)
		ebitenutil.DebugPrintAt(screen, "Co-criacao IA: Nexus Squad (AI Team)", contentX+6, int(boxY)+106)
		ebitenutil.DebugPrintAt(screen, "Localizacao:   Belem do Para - Brasil", contentX+6, int(boxY)+120)

		if (ticks/25)%2 == 0 {
			ebitenutil.DebugPrintAt(screen, ">> [CLIQUE / ENTER / R] JOGAR NOVAMENTE <<", contentX-6, int(boxY)+136)
		} else {
			ebitenutil.DebugPrintAt(screen, "   [CLIQUE / ENTER / R] JOGAR NOVAMENTE   ", contentX-6, int(boxY)+136)
		}
		ebitenutil.DebugPrintAt(screen, "[ESC] Fechar o Jogo", contentX+38, int(boxY)+149)
	}
}

func DrawGameOverScreen(screen *ebiten.Image, screenWidth, screenHeight float64, finalScore int, stage int) {
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{R: 5, G: 5, B: 12, A: 140})

	boxW := 216.0
	boxH := 76.0
	boxX := (screenWidth - boxW) / 2.0
	boxY := 18.0

	ebitenutil.DrawRect(screen, boxX, boxY, boxW, boxH, color.RGBA{R: 28, G: 16, B: 24, A: 245})
	ebitenutil.DrawRect(screen, boxX, boxY, boxW, 2, color.RGBA{R: 240, G: 65, B: 65, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY+boxH, boxW, 2, color.RGBA{R: 240, G: 65, B: 65, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY, 2, boxH, color.RGBA{R: 240, G: 65, B: 65, A: 255})
	ebitenutil.DrawRect(screen, boxX+boxW, boxY, 2, boxH+2, color.RGBA{R: 240, G: 65, B: 65, A: 255})

	bx := int(boxX)
	by := int(boxY)

	stgStr := "VER-O-PESO"
	if stage == 2 {
		stgStr = "ESTACAO DAS DOCAS"
	} else if stage == 3 {
		stgStr = "THEATRO DA PAZ"
	}

	ebitenutil.DebugPrintAt(screen, "★ GAME OVER ★", bx+64, by+8)
	ebitenutil.DrawRect(screen, boxX+10, boxY+21, boxW-20, 1, color.RGBA{R: 240, G: 65, B: 65, A: 160})
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Pontos: %05d  |  %s", finalScore, stgStr), bx+14, by+27)
	ebitenutil.DebugPrintAt(screen, "[ENTER / R] Jogar Novamente", bx+24, by+44)
	ebitenutil.DebugPrintAt(screen, "[ESC] Voltar ao Menu", bx+42, by+58)
}

func DrawSpeechBubble(screen *ebiten.Image, x, y float64, text string) {
	w := float64(len(text)*6 + 14)
	h := 18.0

	origX := x
	if x+w > 334.0 {
		x = 334.0 - w
	}
	if x < 6.0 {
		x = 6.0
	}

	// Sombra suave do balao
	ebitenutil.DrawRect(screen, x+2, y+2, w, h, color.RGBA{R: 0, G: 0, B: 0, A: 140})

	// Fundo do balao (estilo quadrinhos retro paraense)
	ebitenutil.DrawRect(screen, x, y, w, h, color.RGBA{R: 20, G: 25, B: 40, A: 245})

	// Borda dourada de destaque
	ebitenutil.DrawRect(screen, x, y, w, 1, color.RGBA{R: 250, G: 210, B: 50, A: 255})
	ebitenutil.DrawRect(screen, x, y+h-1, w, 1, color.RGBA{R: 250, G: 210, B: 50, A: 255})
	ebitenutil.DrawRect(screen, x, y, 1, h, color.RGBA{R: 250, G: 210, B: 50, A: 255})
	ebitenutil.DrawRect(screen, x+w-1, y, 1, h, color.RGBA{R: 250, G: 210, B: 50, A: 255})

	// Rabicho do balao apontando para baixo (cabeca do garoto)
	tailX := origX + 14.0
	if tailX < x+8 {
		tailX = x + 8
	}
	if tailX > x+w-12 {
		tailX = x + w - 12
	}
	tailY := y + h
	ebitenutil.DrawRect(screen, tailX, tailY, 5, 2, color.RGBA{R: 250, G: 210, B: 50, A: 255})
	ebitenutil.DrawRect(screen, tailX+1, tailY+2, 3, 2, color.RGBA{R: 250, G: 210, B: 50, A: 255})
	ebitenutil.DrawRect(screen, tailX+2, tailY+4, 1, 2, color.RGBA{R: 250, G: 210, B: 50, A: 255})

	// Texto da giria
	ebitenutil.DebugPrintAt(screen, text, int(x)+7, int(y)+3)
}

// drawPopopoBoat renderiza o tradicional barco paraense de madeira (popopó) em pixel art
func drawPopopoBoat(screen *ebiten.Image, x, y float64, ticks int) {
	hullColor := color.RGBA{R: 28, G: 74, B: 140, A: 255}   // Azul amazônico
	hullStripe := color.RGBA{R: 245, G: 195, B: 35, A: 255} // Listra amarela
	cabinColor := color.RGBA{R: 240, G: 240, B: 245, A: 255} // Branco
	roofColor := color.RGBA{R: 190, G: 35, B: 35, A: 255}    // Toldo vermelho
	windowColor := color.RGBA{R: 20, G: 30, B: 50, A: 255}

	// Casco de madeira
	ebitenutil.DrawRect(screen, x+4, y+10, 28, 8, hullColor)
	ebitenutil.DrawRect(screen, x+2, y+12, 32, 6, hullColor)
	ebitenutil.DrawRect(screen, x+32, y+10, 4, 4, hullColor) // Proa elevada
	ebitenutil.DrawRect(screen, x+3, y+13, 30, 2, hullStripe)

	// Cabine com janelinhas
	ebitenutil.DrawRect(screen, x+8, y+2, 18, 9, cabinColor)
	ebitenutil.DrawRect(screen, x+6, y, 22, 3, roofColor)
	ebitenutil.DrawRect(screen, x+10, y+4, 3, 4, windowColor)
	ebitenutil.DrawRect(screen, x+15, y+4, 3, 4, windowColor)
	ebitenutil.DrawRect(screen, x+20, y+4, 3, 4, windowColor)

	// Chaminé na popa soltando fumaça "po-po-pó"
	ebitenutil.DrawRect(screen, x+6, y-3, 2, 4, color.RGBA{R: 50, G: 50, B: 55, A: 255})
	for f := 0; f < 3; f++ {
		smokeTimer := (ticks*2 + f*22) % 60
		sX := x + 6 - float64(smokeTimer)*0.45
		sY := y - 4 - float64(smokeTimer)*0.35
		alpha := uint8(220 - smokeTimer*3)
		if alpha > 0 && smokeTimer < 52 {
			sz := 2.0 + float64(smokeTimer/18)
			ebitenutil.DrawRect(screen, sX, sY, sz, sz, color.RGBA{R: 230, G: 235, B: 245, A: alpha})
		}
	}

	// Mastro e Bandeirinha do Pará balançando
	ebitenutil.DrawRect(screen, x+26, y-6, 1, 8, color.RGBA{R: 200, G: 200, B: 200, A: 255})
	flagWave := int(math.Sin(float64(ticks)*0.2) * 1.5)
	ebitenutil.DrawRect(screen, x+27, y-6+float64(flagWave), 4, 3, color.RGBA{R: 220, G: 40, B: 40, A: 255})

	// Espuma de água na proa cortando as ondas
	foamW := 2.0 + math.Sin(float64(ticks)*0.3)*1.5
	ebitenutil.DrawRect(screen, x+34, y+15, foamW+3, 2, color.RGBA{R: 240, G: 250, B: 255, A: 220})
}

// DrawStageTransitionLoadingScreen renderiza a tela de transição náutica e cultural de Belém entre as fases
func DrawStageTransitionLoadingScreen(screen *ebiten.Image, screenWidth, screenHeight float64, ticks int, nextStage int, progress float64) {
	if progress < 0.0 {
		progress = 0.0
	}
	if progress > 1.0 {
		progress = 1.0
	}

	// 1. Céu crepuscular sobre a Baía do Guajará (degradê roxo-açaí para dourado-pôr-do-sol)
	skyH := screenHeight * 0.58
	for y := 0.0; y < skyH; y += 2.0 {
		t := y / skyH
		r := uint8(20*(1.0-t) + 215*t)
		g := uint8(24*(1.0-t) + 115*t)
		b := uint8(65*(1.0-t) + 40*t)
		ebitenutil.DrawRect(screen, 0, y, screenWidth, 2.0, color.RGBA{R: r, G: g, B: b, A: 255})
	}

	// Sol poente no horizonte de Belém
	sunY := skyH - 14.0
	ebitenutil.DrawRect(screen, 155.0, sunY, 32, 16, color.RGBA{R: 255, G: 215, B: 90, A: 210})
	ebitenutil.DrawRect(screen, 159.0, sunY-3, 24, 22, color.RGBA{R: 255, G: 235, B: 140, A: 230})

	// Silhueta distante de Belém no horizonte (torres da Sé, casario histórico, mangueiras e guindastes)
	horizonY := skyH - 10.0
	cSil := color.RGBA{R: 22, G: 16, B: 28, A: 255}
	ebitenutil.DrawRect(screen, 26, horizonY-14, 8, 14, cSil)   // Torre da Sé
	ebitenutil.DrawRect(screen, 28, horizonY-18, 4, 4, cSil)    // Cúpula
	ebitenutil.DrawRect(screen, 50, horizonY-8, 24, 8, cSil)    // Casario
	ebitenutil.DrawRect(screen, 105, horizonY-13, 16, 13, cSil) // Copa de mangueira centenária
	ebitenutil.DrawRect(screen, 250, horizonY-16, 3, 16, cSil)  // Guindaste inglês
	ebitenutil.DrawRect(screen, 246, horizonY-16, 12, 3, cSil)

	// Urubus e garças distantes voando no horizonte
	drawFlyingUrubu(screen, 70.0+math.Mod(float64(ticks)*0.4, 280.0), 22.0, ticks)
	drawFlyingGarca(screen, 280.0-math.Mod(float64(ticks)*0.6, 300.0), 32.0, ticks)

	// 2. Águas da Baía do Guajará com reflexos dourados e ondulação
	waterY := skyH
	waterH := screenHeight - waterY
	ebitenutil.DrawRect(screen, 0, waterY, screenWidth, waterH, color.RGBA{R: 16, G: 36, B: 68, A: 255})
	for wy := waterY; wy < screenHeight; wy += 4.0 {
		wt := (wy - waterY) / waterH
		r := uint8(20*(1.0-wt) + 8*wt)
		g := uint8(50*(1.0-wt) + 20*wt)
		b := uint8(95*(1.0-wt) + 38*wt)
		ebitenutil.DrawRect(screen, 0, wy, screenWidth, 4.0, color.RGBA{R: r, G: g, B: b, A: 255})

		// Brilho cintilante da água
		waveOffset := math.Sin(float64(ticks)*0.08 + wy*0.5) * 8.0
		ebitenutil.DrawRect(screen, 150.0+waveOffset, wy, 42.0*(1.0-wt*0.5), 1.0, color.RGBA{R: 255, G: 215, B: 120, A: uint8(140 - wt*100)})
	}

	// 3. Barquinho Popopó tradicional navegando pela Baía
	boatBaseX := 20.0 + progress*(screenWidth-85.0)
	boatBob := math.Sin(float64(ticks)*0.12) * 2.5
	boatY := waterY - 14.0 + boatBob
	drawPopopoBoat(screen, boatBaseX, boatY, ticks)

	// Som do motor "popopo~"
	if (ticks/20)%3 == 0 {
		ebitenutil.DebugPrintAt(screen, "popopo~", int(boatBaseX)+36, int(boatY)-8)
	}

	// 4. Placa Superior de Destino & Lore Cultural
	cardW := 294.0
	cardH := 52.0
	cardX := (screenWidth - cardW) / 2.0
	cardY := 10.0

	ebitenutil.DrawRect(screen, cardX+2, cardY+2, cardW, cardH, color.RGBA{R: 0, G: 0, B: 0, A: 160})
	ebitenutil.DrawRect(screen, cardX, cardY, cardW, cardH, color.RGBA{R: 14, G: 20, B: 34, A: 245})
	ebitenutil.DrawRect(screen, cardX, cardY, cardW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, cardX, cardY+cardH, cardW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, cardX, cardY, 2, cardH, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, cardX+cardW, cardY, 2, cardH+2, color.RGBA{R: 250, G: 205, B: 55, A: 255})

	var stageTitle, stageRoute, stageLore string
	switch nextStage {
	case 2:
		stageTitle = "★ PROXIMA PARADA: ESTACAO DAS DOCAS ★"
		stageRoute = "Travessia fluvial: Mercado do Ver-o-Peso -> Docas"
		stageLore = "Dica: Guindastes historicos europeus e jacares no cais!"
	case 3:
		stageTitle = "★ PROXIMA PARADA: THEATRO DA PAZ & MANGUEIRAS ★"
		stageRoute = "Subindo a Av. Presidente Vargas sob as mangueiras"
		stageLore = "Dica: Fundado em 1878, joia neoclassica da Amazonia!"
	default:
		stageTitle = "★ RETORNANDO AO MERCADO DO VER-O-PESO ★"
		stageRoute = "A maior feira a ceu aberto da America Latina"
		stageLore = "Dica: Acai fresquinho e peixe frito com farinha d'agua!"
	}

	ebitenutil.DebugPrintAt(screen, stageTitle, int(cardX)+18, int(cardY)+6)
	ebitenutil.DrawRect(screen, cardX+10, cardY+18, cardW-20, 1, color.RGBA{R: 250, G: 205, B: 55, A: 140})
	ebitenutil.DebugPrintAt(screen, stageRoute, int(cardX)+12, int(cardY)+23)
	ebitenutil.DebugPrintAt(screen, stageLore, int(cardX)+12, int(cardY)+36)

	// 5. Barra de Carregamento Náutica Inferior
	barBoxW := 270.0
	barBoxH := 28.0
	barBoxX := (screenWidth - barBoxW) / 2.0
	barBoxY := screenHeight - barBoxH - 8.0

	ebitenutil.DrawRect(screen, barBoxX+2, barBoxY+2, barBoxW, barBoxH, color.RGBA{R: 0, G: 0, B: 0, A: 160})
	ebitenutil.DrawRect(screen, barBoxX, barBoxY, barBoxW, barBoxH, color.RGBA{R: 12, G: 18, B: 30, A: 245})
	ebitenutil.DrawRect(screen, barBoxX, barBoxY, barBoxW, 1, color.RGBA{R: 250, G: 205, B: 55, A: 220})
	ebitenutil.DrawRect(screen, barBoxX, barBoxY+barBoxH, barBoxW, 1, color.RGBA{R: 250, G: 205, B: 55, A: 220})

	// Trilho da barra
	fillTrackW := barBoxW - 20.0
	fillTrackH := 8.0
	fillX := barBoxX + 10.0
	fillY := barBoxY + 6.0
	ebitenutil.DrawRect(screen, fillX, fillY, fillTrackW, fillTrackH, color.RGBA{R: 25, G: 32, B: 48, A: 255})

	// Preenchimento em verde floresta / dourado
	fillW := fillTrackW * progress
	ebitenutil.DrawRect(screen, fillX, fillY, fillW, fillTrackH, color.RGBA{R: 42, G: 165, B: 85, A: 255})
	ebitenutil.DrawRect(screen, fillX, fillY, fillW, 2, color.RGBA{R: 110, G: 235, B: 140, A: 255})

	// Texto de porcentagem e pontinhos animados
	dots := strings.Repeat(".", (ticks/18)%4)
	pctText := fmt.Sprintf("VIAJANDO POR BELEM: %3d%%%s", int(progress*100), dots)
	ebitenutil.DebugPrintAt(screen, pctText, int(barBoxX)+14, int(barBoxY)+16)

	// Atalho rápido
	ebitenutil.DebugPrintAt(screen, "[ENTER] Pular >>", int(barBoxX)+176, int(barBoxY)+16)

	// 6. Transição cinematográfica: Fade-in suave ao entrar e Fade-out ao sair da tela náutica
	if progress < 0.15 {
		fadeAlpha := uint8(255.0 * (1.0 - progress/0.15))
		ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{R: 12, G: 14, B: 20, A: fadeAlpha})
	} else if progress > 0.85 {
		fadeAlpha := uint8(255.0 * (progress - 0.85) / 0.15)
		ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{R: 12, G: 14, B: 20, A: fadeAlpha})
	}
}
