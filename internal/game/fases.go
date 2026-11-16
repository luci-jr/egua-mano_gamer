package game

import (
	"log"

	"github.com/luci-jr/egua-mano_gamer/internal/entities"
	"github.com/luci-jr/egua-mano_gamer/internal/level"
	"github.com/luci-jr/egua-mano_gamer/internal/save"
)

// tipoPorNome liga os nomes usados em internal/level/fases.json aos tipos de
// obstáculo do jogo.
var tipoPorNome = map[string]entities.ObstacleType{
	level.Paneiro: entities.TypeGround,
	level.Urubu:   entities.TypeAir,
	level.Jacare:  entities.TypeJacare,
	level.Cobra:   entities.TypeSnake,
	level.Banco:   entities.TypeBench,
}

// carregarFases lê as fases embutidas. Se o arquivo estiver inválido o jogo
// continua com a regra original (nil), em vez de quebrar na abertura.
func carregarFases() []level.Fase {
	fases, err := level.Carregar()
	if err != nil {
		log.Printf("fases.json inválido, usando regras padrão: %v", err)
		return nil
	}
	return fases
}

// regraDoTrecho converte um trecho de fase na regra de spawn dos obstáculos.
func regraDoTrecho(t level.Trecho) entities.SpawnRule {
	nomes, pesos := t.TiposOrdenados()
	regra := entities.SpawnRule{
		GapMin:     t.EspacoMin,
		GapMax:     t.EspacoMax,
		EnemySpeed: t.VelInimigos,
	}
	for i, nome := range nomes {
		tipo, ok := tipoPorNome[nome]
		if !ok {
			continue
		}
		regra.Types = append(regra.Types, tipo)
		regra.Weights = append(regra.Weights, pesos[i])
	}
	return regra
}

// aplicarRegraDaFase entrega ao gerenciador de obstáculos a regra do trecho
// em que o jogador está. Só recalcula quando a fase ou o trecho mudam (ou
// quando forcar é true, na troca ou no reinício de fase).
func (e *Engine) aplicarRegraDaFase(metros float64, forcar bool) {
	if len(e.fases) == 0 {
		return
	}
	indice := e.stage - 1
	if indice < 0 {
		indice = 0
	}
	if indice >= len(e.fases) {
		indice = len(e.fases) - 1
	}
	trecho := e.fases[indice].TrechoEm(metros)
	if !forcar && e.regraFase == e.stage && e.regraAte == trecho.AteMetros {
		return
	}
	e.regraFase = e.stage
	e.regraAte = trecho.AteMetros
	e.obstacles.SetSpawnRule(regraDoTrecho(trecho))
}

// registrarRecorde salva o placar atual se ele superar o recorde.
func (e *Engine) registrarRecorde() {
	if save.Registrar(e.score) {
		e.recorde = e.score
		e.novoRecorde = true
	}
}
