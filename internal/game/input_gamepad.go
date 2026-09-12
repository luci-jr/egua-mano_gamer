package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const stickDeadZone = 0.35

// isGamepadJustPressed verifica se um botão padronizado acabou de ser pressionado em qualquer controle conectado
func isGamepadJustPressed(button ebiten.StandardGamepadButton) bool {
	ids := ebiten.AppendGamepadIDs(nil)
	for _, id := range ids {
		if inpututil.IsStandardGamepadButtonJustPressed(id, button) {
			return true
		}
	}
	return false
}

// isGamepadPressed verifica se um botão padronizado está sendo mantido pressionado em qualquer controle conectado
func isGamepadPressed(button ebiten.StandardGamepadButton) bool {
	ids := ebiten.AppendGamepadIDs(nil)
	for _, id := range ids {
		if ebiten.IsStandardGamepadButtonPressed(id, button) {
			return true
		}
	}
	return false
}

// isGamepadLeftPressed verifica se a direção Esquerda está ativada no D-Pad ou Analógico Esquerdo
func isGamepadLeftPressed() bool {
	if isGamepadPressed(ebiten.StandardGamepadButtonLeftLeft) {
		return true
	}
	ids := ebiten.AppendGamepadIDs(nil)
	for _, id := range ids {
		if ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal) < -stickDeadZone {
			return true
		}
	}
	return false
}

// isGamepadLeftJustPressed verifica acionamento rápido para Esquerda
func isGamepadLeftJustPressed() bool {
	return isGamepadJustPressed(ebiten.StandardGamepadButtonLeftLeft)
}

// isGamepadRightPressed verifica se a direção Direita está ativada no D-Pad ou Analógico Esquerdo
func isGamepadRightPressed() bool {
	if isGamepadPressed(ebiten.StandardGamepadButtonLeftRight) {
		return true
	}
	ids := ebiten.AppendGamepadIDs(nil)
	for _, id := range ids {
		if ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal) > stickDeadZone {
			return true
		}
	}
	return false
}

// isGamepadRightJustPressed verifica acionamento rápido para Direita
func isGamepadRightJustPressed() bool {
	return isGamepadJustPressed(ebiten.StandardGamepadButtonLeftRight)
}

// isGamepadUpPressed verifica se a direção Cima está ativada no D-Pad ou Analógico Esquerdo
func isGamepadUpPressed() bool {
	if isGamepadPressed(ebiten.StandardGamepadButtonLeftTop) {
		return true
	}
	ids := ebiten.AppendGamepadIDs(nil)
	for _, id := range ids {
		if ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical) < -stickDeadZone {
			return true
		}
	}
	return false
}

// isGamepadUpJustPressed verifica acionamento rápido para Cima
func isGamepadUpJustPressed() bool {
	return isGamepadJustPressed(ebiten.StandardGamepadButtonLeftTop)
}

// isGamepadDownPressed verifica se a direção Baixo está ativada no D-Pad ou Analógico Esquerdo
func isGamepadDownPressed() bool {
	if isGamepadPressed(ebiten.StandardGamepadButtonLeftBottom) {
		return true
	}
	ids := ebiten.AppendGamepadIDs(nil)
	for _, id := range ids {
		if ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical) > stickDeadZone {
			return true
		}
	}
	return false
}

// isGamepadDownJustPressed verifica acionamento rápido para Baixo
func isGamepadDownJustPressed() bool {
	return isGamepadJustPressed(ebiten.StandardGamepadButtonLeftBottom)
}

// isGamepadConfirmJustPressed verifica confirmação / aceitar (Botão A no Xbox, X no PlayStation, ou Start)
func isGamepadConfirmJustPressed() bool {
	return isGamepadJustPressed(ebiten.StandardGamepadButtonRightBottom) ||
		isGamepadJustPressed(ebiten.StandardGamepadButtonRightLeft) ||
		isGamepadJustPressed(ebiten.StandardGamepadButtonCenterRight)
}

// isGamepadCancelJustPressed verifica cancelamento / voltar (Botão B no Xbox ou Select/View)
func isGamepadCancelJustPressed() bool {
	return isGamepadJustPressed(ebiten.StandardGamepadButtonRightRight) ||
		isGamepadJustPressed(ebiten.StandardGamepadButtonRightTop) ||
		isGamepadJustPressed(ebiten.StandardGamepadButtonCenterLeft)
}

// isGamepadJumpJustPressed (Pular no primeiro frame: Botão A, D-Pad Cima)
func isGamepadJumpJustPressed() bool {
	return isGamepadJustPressed(ebiten.StandardGamepadButtonRightBottom) ||
		isGamepadJustPressed(ebiten.StandardGamepadButtonLeftTop)
}

// isGamepadJumpHolding (Manter pulo pressionado: Botão A, D-Pad Cima ou Analógico Cima)
func isGamepadJumpHolding() bool {
	return isGamepadPressed(ebiten.StandardGamepadButtonRightBottom) ||
		isGamepadUpPressed()
}

// isGamepadAttackHeld (Segurar disparo de açaí / rugido: Botões X, B, Y ou Triggers/Bumpers)
func isGamepadAttackHeld() bool {
	return isGamepadPressed(ebiten.StandardGamepadButtonRightLeft) ||
		isGamepadPressed(ebiten.StandardGamepadButtonRightRight) ||
		isGamepadPressed(ebiten.StandardGamepadButtonRightTop) ||
		isGamepadPressed(ebiten.StandardGamepadButtonFrontTopRight) ||
		isGamepadPressed(ebiten.StandardGamepadButtonFrontTopLeft) ||
		isGamepadPressed(ebiten.StandardGamepadButtonFrontBottomRight) ||
		isGamepadPressed(ebiten.StandardGamepadButtonFrontBottomLeft)
}

// isGamepadPauseJustPressed (Pausar: Botão Start / Menu)
func isGamepadPauseJustPressed() bool {
	return isGamepadJustPressed(ebiten.StandardGamepadButtonCenterRight)
}

// isGamepadCreditsJustPressed (Créditos: Botão Select / View)
func isGamepadCreditsJustPressed() bool {
	return isGamepadJustPressed(ebiten.StandardGamepadButtonCenterLeft)
}
