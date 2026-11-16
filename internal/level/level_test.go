package level

import (
	"strings"
	"testing"
)

func TestFasesEmbutidasSaoValidas(t *testing.T) {
	fases, err := Carregar()
	if err != nil {
		t.Fatal(err)
	}
	if len(fases) != 3 {
		t.Fatalf("esperava 3 fases, vieram %d", len(fases))
	}
}

// A dificuldade não pode cair de um trecho para o seguinte, nem de uma fase
// para a próxima: o espaço entre obstáculos só diminui e os inimigos só
// aceleram. Este teste é a "curva de dificuldade" escrita como regra.
func TestDificuldadeNuncaDiminui(t *testing.T) {
	fases, err := Carregar()
	if err != nil {
		t.Fatal(err)
	}
	var anterior *Trecho
	for _, f := range fases {
		for i := range f.Trechos {
			tr := f.Trechos[i]
			if anterior != nil {
				if tr.EspacoMin > anterior.EspacoMin {
					t.Errorf("%s, trecho %d: espaco_min subiu de %.0f para %.0f", f.Nome, i+1, anterior.EspacoMin, tr.EspacoMin)
				}
				if tr.VelInimigos < anterior.VelInimigos {
					t.Errorf("%s, trecho %d: vel_inimigos caiu de %.2f para %.2f", f.Nome, i+1, anterior.VelInimigos, tr.VelInimigos)
				}
			}
			anterior = &f.Trechos[i]
		}
	}
}

func TestTrechoEm(t *testing.T) {
	f := Fase{Nome: "x", MetaMetros: 300, Trechos: []Trecho{
		{AteMetros: 100, EspacoMin: 200},
		{AteMetros: 300, EspacoMin: 150},
	}}
	casos := map[float64]float64{0: 200, 99.9: 200, 100: 150, 299: 150, 300: 150, 9999: 150}
	for metros, quer := range casos {
		if got := f.TrechoEm(metros).EspacoMin; got != quer {
			t.Errorf("TrechoEm(%.1f).EspacoMin = %.0f, esperava %.0f", metros, got, quer)
		}
	}
}

func TestTiposOrdenadosIgnoraPesoZeroEEhEstavel(t *testing.T) {
	tr := Trecho{Tipos: map[string]int{Urubu: 3, Banco: 0, Cobra: 1, Jacare: 2}}
	for i := 0; i < 20; i++ {
		tipos, pesos := tr.TiposOrdenados()
		if strings.Join(tipos, ",") != "cobra,jacare,urubu" {
			t.Fatalf("tipos = %v", tipos)
		}
		if len(pesos) != 3 || pesos[0] != 1 || pesos[1] != 2 || pesos[2] != 3 {
			t.Fatalf("pesos = %v", pesos)
		}
	}
}

func TestLerRecusaDadosRuins(t *testing.T) {
	base := `{"fases":[{"nome":"A","meta_metros":100,"trechos":[%s]}]}`
	trechoOK := `{"ate_metros":100,"tipos":{"paneiro":1},"espaco_min":140,"espaco_max":190,"vel_inimigos":1}`
	casos := map[string]string{
		"JSON quebrado":         `{"fases":[`,
		"sem fases":             `{"fases":[]}`,
		"sem trechos":           strings.Replace(base, "%s", "", 1),
		"tipo desconhecido":     strings.Replace(base, "%s", strings.Replace(trechoOK, "paneiro", "dragao", 1), 1),
		"espaço invertido":      strings.Replace(base, "%s", strings.Replace(trechoOK, `"espaco_max":190`, `"espaco_max":100`, 1), 1),
		"velocidade zero":       strings.Replace(base, "%s", strings.Replace(trechoOK, `"vel_inimigos":1`, `"vel_inimigos":0`, 1), 1),
		"pesos todos zero":      strings.Replace(base, "%s", strings.Replace(trechoOK, `"paneiro":1`, `"paneiro":0`, 1), 1),
		"termina antes da meta": strings.Replace(base, "%s", strings.Replace(trechoOK, `"ate_metros":100`, `"ate_metros":50`, 1), 1),
		"trechos fora de ordem": strings.Replace(base, "%s", trechoOK+","+strings.Replace(trechoOK, `"ate_metros":100`, `"ate_metros":80`, 1), 1),
	}
	for nome, dados := range casos {
		if _, err := Ler([]byte(dados)); err == nil {
			t.Errorf("%s: deveria dar erro", nome)
		}
	}
	if _, err := Ler([]byte(strings.Replace(base, "%s", trechoOK, 1))); err != nil {
		t.Errorf("caso válido deu erro: %v", err)
	}
}
