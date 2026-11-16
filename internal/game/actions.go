package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Camada de ações de entrada.
//
// Cada tela do jogo perguntava a mesma coisa ("o jogador confirmou?") com a
// mesma lista de teclas, teclas virtuais do celular e botões do gamepad,
// copiada tela a tela. Estas funções dão nome à ação e concentram a lista num
// lugar só: para mapear uma tecla nova, muda-se aqui.

// consumeVirtual informa se alguma das teclas virtuais (botões de toque do
// index.html) está acionada e já as desliga, para que o mesmo toque não seja
// lido de novo no quadro seguinte.
func consumeVirtual(nomes ...string) bool {
	acionada := false
	for _, nome := range nomes {
		if getVirtualKey(nome) {
			resetVirtualKey(nome)
			acionada = true
		}
	}
	return acionada
}

// confirmKeyJustPressed: Enter, Enter do teclado numérico ou Espaço.
func confirmKeyJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

// confirmJustPressed é a ação "confirmar" dos menus e telas de transição:
// teclado, botões virtuais (Enter e Espaço; Start quando comStart é true) e
// botão de confirmação do gamepad.
func confirmJustPressed(comStart bool) bool {
	virtuais := []string{"Enter", "Space"}
	if comStart {
		virtuais = append(virtuais, "Start")
	}
	// As três fontes são sempre consultadas: consumeVirtual precisa rodar
	// mesmo que o teclado já tenha confirmado.
	teclado := confirmKeyJustPressed()
	virtual := consumeVirtual(virtuais...)
	gamepad := isGamepadConfirmJustPressed()
	return teclado || virtual || gamepad
}
