// Package level descreve as fases do jogo em dados (fases.json), em vez de
// regras espalhadas pelo código. Mudar a dificuldade, o ritmo ou os inimigos
// de uma fase passa a ser editar um arquivo, não recompilar lógica.
//
// O pacote não depende da Ebitengine: dá para testar com "go test" em
// qualquer máquina, sem janela nem placa de vídeo.
package level

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
)

//go:embed fases.json
var fasesJSON []byte

// Tipos de obstáculo aceitos em fases.json.
const (
	Paneiro = "paneiro"
	Urubu   = "urubu"
	Jacare  = "jacare"
	Cobra   = "cobra"
	Banco   = "banco"
)

var tiposValidos = map[string]bool{Paneiro: true, Urubu: true, Jacare: true, Cobra: true, Banco: true}

// Trecho é um pedaço da fase com as mesmas regras de obstáculos.
type Trecho struct {
	AteMetros   float64        `json:"ate_metros"`
	Tipos       map[string]int `json:"tipos"`
	EspacoMin   float64        `json:"espaco_min"`
	EspacoMax   float64        `json:"espaco_max"`
	VelInimigos float64        `json:"vel_inimigos"`
}

// Fase é uma fase completa, do início até a meta.
type Fase struct {
	Nome       string   `json:"nome"`
	MetaMetros float64  `json:"meta_metros"`
	Trechos    []Trecho `json:"trechos"`
}

type arquivo struct {
	Fases []Fase `json:"fases"`
}

// Carregar lê e valida as fases embutidas no binário.
func Carregar() ([]Fase, error) {
	return Ler(fasesJSON)
}

// Ler interpreta e valida um JSON de fases.
func Ler(dados []byte) ([]Fase, error) {
	var a arquivo
	if err := json.Unmarshal(dados, &a); err != nil {
		return nil, fmt.Errorf("fases: JSON inválido: %w", err)
	}
	if len(a.Fases) == 0 {
		return nil, fmt.Errorf("fases: nenhuma fase definida")
	}
	for i, f := range a.Fases {
		if err := f.validar(); err != nil {
			return nil, fmt.Errorf("fase %d (%s): %w", i+1, f.Nome, err)
		}
	}
	return a.Fases, nil
}

func (f Fase) validar() error {
	if f.Nome == "" {
		return fmt.Errorf("sem nome")
	}
	if f.MetaMetros <= 0 {
		return fmt.Errorf("meta_metros deve ser maior que zero")
	}
	if len(f.Trechos) == 0 {
		return fmt.Errorf("sem trechos")
	}
	anterior := 0.0
	for i, t := range f.Trechos {
		n := i + 1
		if t.AteMetros <= anterior {
			return fmt.Errorf("trecho %d: ate_metros (%.0f) deve ser maior que o do trecho anterior (%.0f)", n, t.AteMetros, anterior)
		}
		anterior = t.AteMetros
		if t.EspacoMin <= 0 || t.EspacoMax < t.EspacoMin {
			return fmt.Errorf("trecho %d: espaco_min/espaco_max inválidos (%.0f/%.0f)", n, t.EspacoMin, t.EspacoMax)
		}
		if t.VelInimigos <= 0 {
			return fmt.Errorf("trecho %d: vel_inimigos deve ser maior que zero", n)
		}
		soma := 0
		for tipo, peso := range t.Tipos {
			if !tiposValidos[tipo] {
				return fmt.Errorf("trecho %d: tipo desconhecido %q", n, tipo)
			}
			if peso < 0 {
				return fmt.Errorf("trecho %d: peso negativo para %q", n, tipo)
			}
			soma += peso
		}
		if soma == 0 {
			return fmt.Errorf("trecho %d: nenhum tipo com peso maior que zero", n)
		}
	}
	if anterior < f.MetaMetros {
		return fmt.Errorf("o último trecho termina em %.0f m, antes da meta (%.0f m)", anterior, f.MetaMetros)
	}
	return nil
}

// TrechoEm devolve o trecho que vale na distância informada. Depois do último
// limite, continua valendo o último trecho.
func (f Fase) TrechoEm(metros float64) Trecho {
	for _, t := range f.Trechos {
		if metros < t.AteMetros {
			return t
		}
	}
	return f.Trechos[len(f.Trechos)-1]
}

// TiposOrdenados devolve os tipos do trecho com peso maior que zero e seus
// pesos, em ordem alfabética, para que o sorteio seja reproduzível (a ordem
// de iteração de um map em Go muda a cada execução).
func (t Trecho) TiposOrdenados() (tipos []string, pesos []int) {
	for tipo, peso := range t.Tipos {
		if peso > 0 {
			tipos = append(tipos, tipo)
		}
	}
	sort.Strings(tipos)
	for _, tipo := range tipos {
		pesos = append(pesos, t.Tipos[tipo])
	}
	return tipos, pesos
}
