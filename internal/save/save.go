// Package save guarda o recorde do jogador entre uma partida e outra.
//
// No navegador (WebAssembly) usa o localStorage; nas outras plataformas, um
// arquivo pequeno na pasta de configuração do usuário. Qualquer falha ao ler
// ou gravar é ignorada: o recorde é um conforto, nunca um motivo para o jogo
// travar.
package save

import "strconv"

const chave = "eguamano.recorde"

// Recorde devolve o maior placar salvo, ou 0 se não houver nenhum.
func Recorde() int {
	n, err := strconv.Atoi(ler())
	if err != nil || n < 0 {
		return 0
	}
	return n
}

// Registrar grava o placar se ele for maior que o recorde atual e informa se
// houve novo recorde.
func Registrar(placar int) bool {
	if placar <= Recorde() {
		return false
	}
	gravar(strconv.Itoa(placar))
	return true
}
