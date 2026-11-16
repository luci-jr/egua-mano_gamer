package entities

import "testing"

// avancar simula a fase rolando: move o mundo e devolve os tipos que
// reapareceram, na ordem.
func avancar(m *ObstacleManager, ticks int) []ObstacleType {
	var vistos []ObstacleType
	for i := 0; i < ticks; i++ {
		antes := make(map[*Obstacle]float64, len(m.Obstacles))
		for _, o := range m.Obstacles {
			antes[o] = o.X
		}
		m.Update(3.0)
		for _, o := range m.Obstacles {
			if o.X > antes[o] { // voltou para a frente da fila: reapareceu
				vistos = append(vistos, o.Type)
			}
		}
	}
	return vistos
}

func TestMesmaSementeMesmaFase(t *testing.T) {
	a := NewObstacleManager(340, 155)
	b := NewObstacleManager(340, 155)
	a.SetSeed(42)
	b.SetSeed(42)
	va, vb := avancar(a, 3000), avancar(b, 3000)
	if len(va) < 20 {
		t.Fatalf("poucos reaparecimentos para comparar: %d", len(va))
	}
	if len(va) != len(vb) {
		t.Fatalf("quantidades diferentes: %d e %d", len(va), len(vb))
	}
	for i := range va {
		if va[i] != vb[i] {
			t.Fatalf("divergiu no reaparecimento %d: %v e %v", i, va[i], vb[i])
		}
	}
	for i, o := range a.Obstacles {
		if o.X != b.Obstacles[i].X {
			t.Fatalf("posição do obstáculo %d divergiu: %v e %v", i, o.X, b.Obstacles[i].X)
		}
	}
}

func TestRegraLimitaOsTipos(t *testing.T) {
	m := NewObstacleManager(340, 155)
	m.SetSeed(7)
	m.SetSpawnRule(SpawnRule{Types: []ObstacleType{TypeGround, TypeBench}, GapMin: 150, GapMax: 200, EnemySpeed: 1})
	vistos := avancar(m, 4000)
	if len(vistos) < 20 {
		t.Fatalf("poucos reaparecimentos: %d", len(vistos))
	}
	for _, tipo := range vistos {
		if tipo != TypeGround && tipo != TypeBench {
			t.Fatalf("surgiu um tipo fora da regra: %v", tipo)
		}
	}
}

func TestPesoZeroNuncaAparecePesoAltoDomina(t *testing.T) {
	m := NewObstacleManager(340, 155)
	m.SetSeed(11)
	m.SetSpawnRule(SpawnRule{
		Types:      []ObstacleType{TypeGround, TypeSnake, TypeJacare},
		Weights:    []int{9, 1, 0},
		GapMin:     140,
		GapMax:     190,
		EnemySpeed: 1,
	})
	conta := map[ObstacleType]int{}
	for _, tipo := range avancar(m, 20000) {
		conta[tipo]++
	}
	if conta[TypeJacare] != 0 {
		t.Fatalf("peso zero apareceu %d vezes", conta[TypeJacare])
	}
	if conta[TypeGround] <= conta[TypeSnake]*3 {
		t.Fatalf("peso 9 deveria dominar o peso 1: %v", conta)
	}
}

func TestEspacamentoRespeitaARegra(t *testing.T) {
	m := NewObstacleManager(340, 155)
	m.SetSeed(3)
	// Só tipos parados, para medir o espaçamento sem a velocidade própria.
	m.SetSpawnRule(SpawnRule{Types: []ObstacleType{TypeGround, TypeBench}, GapMin: 160, GapMax: 180, EnemySpeed: 1})
	m.Reset()
	for i := 0; i < 5000; i++ {
		antes := make(map[*Obstacle]float64, len(m.Obstacles))
		for _, o := range m.Obstacles {
			antes[o] = o.X
		}
		m.Update(3.0)
		for _, o := range m.Obstacles {
			if o.X <= antes[o] {
				continue
			}
			maisAFrente := m.screenWidth
			for _, outro := range m.Obstacles {
				if outro != o && outro.X > maisAFrente {
					maisAFrente = outro.X
				}
			}
			if gap := o.X - maisAFrente; gap < 160 || gap > 180 {
				t.Fatalf("espaçamento %.1f fora de [160, 180]", gap)
			}
		}
	}
}

func TestRegraInvalidaEhIgnorada(t *testing.T) {
	m := NewObstacleManager(340, 155)
	padrao := DefaultSpawnRule()
	for _, ruim := range []SpawnRule{
		{},
		{Types: []ObstacleType{TypeGround}, GapMin: 0, GapMax: 10},
		{Types: []ObstacleType{TypeGround}, GapMin: 200, GapMax: 100},
	} {
		m.SetSpawnRule(ruim)
		if len(m.rule.Types) != len(padrao.Types) || m.rule.GapMin != padrao.GapMin {
			t.Fatalf("regra inválida %+v foi aceita", ruim)
		}
	}
}
