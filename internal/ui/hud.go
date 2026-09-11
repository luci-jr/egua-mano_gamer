package ui

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed sao_bras.jpg
var saoBrasBytes []byte
var saoBrasImg *ebiten.Image

func getSaoBrasImage() *ebiten.Image {
	if saoBrasImg != nil {
		return saoBrasImg
	}
	if len(saoBrasBytes) == 0 {
		return nil
	}
	img, _, err := image.Decode(bytes.NewReader(saoBrasBytes))
	if err != nil {
		return nil
	}
	saoBrasImg = ebiten.NewImageFromImage(img)
	return saoBrasImg
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

func DrawHUD(screen *ebiten.Image, lives int, hearts int, score int, stage int, isDoubleJump bool, stageBannerTimer int, isMuted bool, ticks int) {
	for i := 0; i < 3; i++ {
		hx := 8.0 + float64(i*12)
		DrawHeart(screen, hx, 9, i < hearts)
	}

	livesText := fmt.Sprintf("x%d VIDAS", lives)
	ebitenutil.DebugPrintAt(screen, livesText, 46, 9)

	stageName := "1. VER-O-PESO"
	if stage == 2 {
		stageName = "2. DOCAS"
	} else if stage == 3 {
		stageName = "3. THEATRO"
	}
	ebitenutil.DebugPrintAt(screen, stageName, 116, 9)

	scoreText := fmt.Sprintf("SCORE: %05d", score)
	ebitenutil.DebugPrintAt(screen, scoreText, 226, 9)

	if isMuted {
		ebitenutil.DebugPrintAt(screen, "[MUDO]", 170, 24)
	}

	if stageBannerTimer > 0 {
		var bannerTitle string
		switch stage {
		case 2:
			bannerTitle = "★ FASE 2: ESTACAO DAS DOCAS ★"
		case 3:
			bannerTitle = "★ FASE 3: THEATRO DA PAZ ★"
		default:
			bannerTitle = "★ FASE 1: VER-O-PESO ★"
		}
		ebitenutil.DrawRect(screen, 30, 50, 260, 22, color.RGBA{R: 20, G: 20, B: 30, A: 210})
		ebitenutil.DrawRect(screen, 30, 50, 260, 2, color.RGBA{R: 250, G: 200, B: 50, A: 255})
		ebitenutil.DrawRect(screen, 30, 70, 260, 2, color.RGBA{R: 250, G: 200, B: 50, A: 255})
		ebitenutil.DebugPrintAt(screen, bannerTitle, 45, 55)
	}

	// Rodapé cultural com curiosidades dinâmicas de Belém do Pará
	DrawCityFooter(screen, float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy()), stage, ticks)
}

func DrawSaoBrasIntro(screen *ebiten.Image, screenWidth, screenHeight float64, ticks int, isAudioPlaying bool) {
	sImg := getSaoBrasImage()
	if sImg != nil {
		op := &ebiten.DrawImageOptions{}
		bounds := sImg.Bounds()
		scaleX := screenWidth / float64(bounds.Dx())
		scaleY := screenHeight / float64(bounds.Dy())
		op.GeoM.Scale(scaleX, scaleY)
		screen.DrawImage(sImg, op)
	} else {
		ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{R: 15, G: 20, B: 30, A: 255})
	}

	// Faixa inferior ampla com curiosidades históricas de São Brás e Belém
	bannerH := 45.0
	bannerY := screenHeight - bannerH
	ebitenutil.DrawRect(screen, 0, bannerY, screenWidth, bannerH, color.RGBA{R: 6, G: 10, B: 20, A: 240})
	ebitenutil.DrawRect(screen, 0, bannerY, screenWidth, 1, color.RGBA{R: 250, G: 205, B: 55, A: 255})

	if isAudioPlaying {
		ebitenutil.DebugPrintAt(screen, "★ MERCADO DE SAO BRAS (TRILHA DE APRESENTACAO) ★", 24, int(bannerY)+3)
	} else {
		ebitenutil.DebugPrintAt(screen, "★ MERCADO DE SAO BRAS ATUAL (BELEM - PA) ★", 34, int(bannerY)+3)
	}

	// Linha 2: Letreiro de curiosidade histórica de São Brás
	sbCuriosities := "★ SAO BRAS: Inaugurado em 1911 pelo arquiteto George Saint-Clair!  " +
		"★ ARQUITETURA: Fachada historica com ferro europeu e azulejos raros!  " +
		"★ MEMORIA: Erguido ao lado da antiga ferrovia Belem-Braganca!  " +
		"★ NOVO SAO BRAS: Polo cultural, gastronomico e historico revitalizado!  "
	sbTextWidth := len(sbCuriosities) * 6
	sbOffset := ticks % sbTextWidth
	sbBaseX := 52 - sbOffset
	ebitenutil.DebugPrintAt(screen, sbCuriosities, sbBaseX, int(bannerY)+16)
	if sbBaseX+sbTextWidth < int(screenWidth) {
		ebitenutil.DebugPrintAt(screen, sbCuriosities, sbBaseX+sbTextWidth, int(bannerY)+16)
	}
	ebitenutil.DrawRect(screen, 0, bannerY+15, 48, 14, color.RGBA{R: 16, G: 26, B: 46, A: 255})
	ebitenutil.DrawRect(screen, 48, bannerY+15, 1, 14, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DebugPrintAt(screen, "★FATO", 6, int(bannerY)+16)

	// Linha 3: Botão de largada piscante
	if isAudioPlaying {
		if (ticks/25)%2 == 0 {
			ebitenutil.DebugPrintAt(screen, ">> CLIQUE NA JANELA PARA COMEÇAR A CORRIDA <<", 30, int(bannerY)+30)
		} else {
			ebitenutil.DebugPrintAt(screen, "   CLIQUE NA JANELA PARA COMEÇAR A CORRIDA   ", 30, int(bannerY)+30)
		}
	} else {
		if (ticks/25)%2 == 0 {
			ebitenutil.DebugPrintAt(screen, ">> CLIQUE NA JANELA DO JOGO PARA INICIAR <<", 33, int(bannerY)+30)
		} else {
			ebitenutil.DebugPrintAt(screen, "   CLIQUE NA JANELA DO JOGO PARA INICIAR   ", 33, int(bannerY)+30)
		}
	}
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
	ebitenutil.DebugPrintAt(screen, "Mover: Setas / W-A-S-D", bx+28, by+41)
	ebitenutil.DebugPrintAt(screen, "Avancar/Voltar: Dir/Esq", bx+26, by+54)

	if (ticks/30)%2 == 0 {
		ebitenutil.DebugPrintAt(screen, ">> TOQUE OU APERTE ENTER <<", bx+20, by+72)
	} else {
		ebitenutil.DebugPrintAt(screen, "   TOQUE OU APERTE ENTER   ", bx+20, by+72)
	}

	ebitenutil.DebugPrintAt(screen, "[C] Creditos  |  [ESC] Sair", bx+18, by+91)
}

func DrawCreditsScreen(screen *ebiten.Image, screenWidth, screenHeight float64) {
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{R: 8, G: 10, B: 18, A: 140})

	boxW := 210.0
	boxH := 122.0
	boxX := (screenWidth - boxW) / 2.0
	boxY := (screenHeight - boxH) / 2.0

	ebitenutil.DrawRect(screen, boxX, boxY, boxW, boxH, color.RGBA{R: 20, G: 26, B: 42, A: 245})
	ebitenutil.DrawRect(screen, boxX, boxY, boxW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY+boxH, boxW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY, 2, boxH, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX+boxW, boxY, 2, boxH+2, color.RGBA{R: 250, G: 205, B: 55, A: 255})

	bx := int(boxX)
	by := int(boxY)

	ebitenutil.DebugPrintAt(screen, "★ CREDITOS - PAIDEGUA ★", bx+32, by+8)
	ebitenutil.DrawRect(screen, boxX+10, boxY+21, boxW-20, 1, color.RGBA{R: 250, G: 205, B: 55, A: 160})

	ebitenutil.DebugPrintAt(screen, "Desenvolvedor: Luci Junior", bx+16, by+26)
	ebitenutil.DebugPrintAt(screen, "Co-criacao IA: Nexus AI (Lucy)", bx+16, by+39)
	ebitenutil.DebugPrintAt(screen, "Linguagem:     Go + Ebitengine", bx+16, by+52)
	ebitenutil.DebugPrintAt(screen, "Audio:         Carimbo 8-bit", bx+16, by+65)
	ebitenutil.DebugPrintAt(screen, "Cenarios:      Belem do Para", bx+16, by+78)

	ebitenutil.DebugPrintAt(screen, "[ENTER / ESC / C] Voltar", bx+28, by+102)
}

func DrawPauseMenu(screen *ebiten.Image, screenWidth, screenHeight float64, selectedIndex int, isMuted bool) {
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{R: 10, G: 12, B: 20, A: 140})

	boxW := 196.0
	boxH := 126.0
	boxX := (screenWidth - boxW) / 2.0
	boxY := (screenHeight - boxH) / 2.0

	ebitenutil.DrawRect(screen, boxX, boxY, boxW, boxH, color.RGBA{R: 22, G: 28, B: 44, A: 245})
	ebitenutil.DrawRect(screen, boxX, boxY, boxW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY+boxH, boxW, 2, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY, 2, boxH, color.RGBA{R: 250, G: 205, B: 55, A: 255})
	ebitenutil.DrawRect(screen, boxX+boxW, boxY, 2, boxH+2, color.RGBA{R: 250, G: 205, B: 55, A: 255})

	bx := int(boxX)
	by := int(boxY)

	ebitenutil.DebugPrintAt(screen, "PAUSA - PAIDEGUA", bx+48, by+8)
	ebitenutil.DrawRect(screen, boxX+10, boxY+22, boxW-20, 1, color.RGBA{R: 250, G: 205, B: 55, A: 160})

	soundStatus := "[ SOM: LIGADO ]"
	if isMuted {
		soundStatus = "[ SOM: MUDO ]  "
	}

	options := []string{
		"Continuar",
		"Reiniciar Aventura",
		soundStatus,
		"Ver Creditos",
		"Fechar o Jogo",
	}

	startY := by + 28
	for i, opt := range options {
		itemY := startY + i*16
		if i == selectedIndex {
			ebitenutil.DrawRect(screen, boxX+8, float64(itemY-2), boxW-16, 14, color.RGBA{R: 245, G: 185, B: 45, A: 70})
			ebitenutil.DrawRect(screen, boxX+8, float64(itemY-2), 3, 14, color.RGBA{R: 250, G: 210, B: 50, A: 255})
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("> %s <", opt), bx+22, itemY)
		} else {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("  %s", opt), bx+22, itemY)
		}
	}

	ebitenutil.DebugPrintAt(screen, "[CIMA/BAIXO] | [ENTER] | [ESC]", bx+14, by+int(boxH)-12)
}

func DrawStageCompleteScreen(screen *ebiten.Image, screenWidth, screenHeight float64, stage int, score int) {
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{R: 10, G: 15, B: 25, A: 210})

	boxX := 20.0
	boxY := 12.0
	boxW := screenWidth - 40.0
	boxH := screenHeight - 24.0

	ebitenutil.DrawRect(screen, boxX, boxY, boxW, boxH, color.RGBA{R: 25, G: 30, B: 45, A: 240})
	ebitenutil.DrawRect(screen, boxX, boxY, boxW, 2, color.RGBA{R: 250, G: 210, B: 50, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY+boxH, boxW, 2, color.RGBA{R: 250, G: 210, B: 50, A: 255})
	ebitenutil.DrawRect(screen, boxX, boxY, 2, boxH, color.RGBA{R: 250, G: 210, B: 50, A: 255})
	ebitenutil.DrawRect(screen, boxX+boxW, boxY, 2, boxH+2, color.RGBA{R: 250, G: 210, B: 50, A: 255})

	if stage < 3 {
		completedName := "MERCADO DO VER-O-PESO"
		nextName := "ESTACAO DAS DOCAS"
		if stage == 2 {
			completedName = "ESTACAO DAS DOCAS"
			nextName = "THEATRO DA PAZ"
		}

		ebitenutil.DebugPrintAt(screen, "===================================", 45, 26)
		ebitenutil.DebugPrintAt(screen, "     PARABENS! FASE CONCLUIDA!     ", 45, 40)
		ebitenutil.DebugPrintAt(screen, "===================================", 45, 54)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Concluido: %s", completedName), 40, 72)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Proxima:   %s", nextName), 40, 88)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Pontuacao Atual: %05d pts", score), 40, 104)
		ebitenutil.DebugPrintAt(screen, "[ENTER] Proxima Fase", 40, 126)
		ebitenutil.DebugPrintAt(screen, "[R]     Recomecar do Inicio", 40, 142)
		ebitenutil.DebugPrintAt(screen, "[ESC]   Finalizar e Fechar o Jogo", 40, 158)
	} else {
		ebitenutil.DebugPrintAt(screen, "========================================", 42, 18)
		ebitenutil.DebugPrintAt(screen, "    PARABENS! EXPEDICAO CONCLUIDA!      ", 42, 29)
		ebitenutil.DebugPrintAt(screen, "========================================", 42, 40)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Pontuacao Maxima: %05d pts", score), 42, 54)
		ebitenutil.DebugPrintAt(screen, "----------- CREDITOS FINAIS ------------", 42, 67)
		ebitenutil.DebugPrintAt(screen, "Desenvolvedor: Lucivaldo Junior (Luci)", 42, 80)
		ebitenutil.DebugPrintAt(screen, "Co-criacao IA: Nexus Squad (Time de IA)", 42, 93)
		ebitenutil.DebugPrintAt(screen, "Linguagem:     Go (Golang 1.22+)", 42, 106)
		ebitenutil.DebugPrintAt(screen, "Game Engine:   Ebitengine (v2)", 42, 119)
		ebitenutil.DebugPrintAt(screen, "Trilha Sonora: Carimbo Procedural 8-bit", 42, 132)
		ebitenutil.DebugPrintAt(screen, "Localizacao:   Belem do Para - Amazonia", 42, 145)
		ebitenutil.DebugPrintAt(screen, "[ENTER / R] Jogar Novamente | [ESC] Sair", 32, 168)
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

	// Rabicho do balao apontando para baixo (cabeca da onca)
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
