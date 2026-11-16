//go:build !(js && wasm)

package save

import (
	"os"
	"path/filepath"
	"strings"
)

// dirBase permite que os testes troquem a pasta de configuração.
var dirBase = os.UserConfigDir

func caminho() (string, bool) {
	dir, err := dirBase()
	if err != nil || dir == "" {
		return "", false
	}
	return filepath.Join(dir, "egua-mano-gamer", "recorde.txt"), true
}

func ler() string {
	p, ok := caminho()
	if !ok {
		return ""
	}
	dados, err := os.ReadFile(p)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(dados))
}

func gravar(valor string) {
	p, ok := caminho()
	if !ok {
		return
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return
	}
	_ = os.WriteFile(p, []byte(valor), 0o644)
}
