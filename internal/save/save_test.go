//go:build !(js && wasm)

package save

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func usarPastaTemporaria(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	antes := dirBase
	dirBase = func() (string, error) { return dir, nil }
	t.Cleanup(func() { dirBase = antes })
	return dir
}

func TestRecordeComecaEmZero(t *testing.T) {
	usarPastaTemporaria(t)
	if got := Recorde(); got != 0 {
		t.Fatalf("Recorde() = %d, esperava 0", got)
	}
}

func TestRegistrarSoGravaQuandoSupera(t *testing.T) {
	usarPastaTemporaria(t)
	if !Registrar(500) {
		t.Fatal("o primeiro placar deveria virar recorde")
	}
	if Registrar(500) || Registrar(120) {
		t.Fatal("placar igual ou menor não é recorde")
	}
	if got := Recorde(); got != 500 {
		t.Fatalf("Recorde() = %d, esperava 500", got)
	}
	if !Registrar(900) || Recorde() != 900 {
		t.Fatal("placar maior deveria substituir o recorde")
	}
}

func TestArquivoCorrompidoViraZero(t *testing.T) {
	dir := usarPastaTemporaria(t)
	p := filepath.Join(dir, "egua-mano-gamer", "recorde.txt")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, lixo := range []string{"abc", "-50", ""} {
		if err := os.WriteFile(p, []byte(lixo), 0o644); err != nil {
			t.Fatal(err)
		}
		if got := Recorde(); got != 0 {
			t.Errorf("conteúdo %q: Recorde() = %d, esperava 0", lixo, got)
		}
	}
}

func TestSemPastaDeConfiguracaoNaoQuebra(t *testing.T) {
	antes := dirBase
	dirBase = func() (string, error) { return "", errors.New("sem pasta") }
	t.Cleanup(func() { dirBase = antes })
	if Recorde() != 0 {
		t.Fatal("sem pasta, o recorde é 0")
	}
	Registrar(100) // não pode dar panic
	if Recorde() != 0 {
		t.Fatal("sem pasta, nada é gravado")
	}
}
