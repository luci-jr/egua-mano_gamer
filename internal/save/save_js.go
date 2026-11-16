//go:build js && wasm

package save

import "syscall/js"

// localStorage pode não existir ou lançar exceção (navegação privada, iframe
// sem permissão). Uma exceção de JavaScript vira panic em Go, por isso o
// recover.

func ler() (valor string) {
	defer func() {
		if recover() != nil {
			valor = ""
		}
	}()
	v := js.Global().Get("localStorage").Call("getItem", chave)
	if v.IsNull() || v.IsUndefined() {
		return ""
	}
	return v.String()
}

func gravar(valor string) {
	defer func() { _ = recover() }()
	js.Global().Get("localStorage").Call("setItem", chave, valor)
}
