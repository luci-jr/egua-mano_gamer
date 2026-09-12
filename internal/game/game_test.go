package game

import (
	"testing"

	"github.com/luci-jr/egua-mano_gamer/internal/entities"
)

func TestHeroToggleAndApply(t *testing.T) {
	e := NewEngine()

	// Inicialmente herói é Garoto
	if e.selectedHero != entities.HeroGaroto {
		t.Fatalf("Esperado HeroGaroto inicialmente, obtido %d", e.selectedHero)
	}
	if e.player != e.garoto {
		t.Fatalf("Esperado player apontar para garoto")
	}

	// Alterna para Onça
	e.toggleSelectedHero()
	if e.selectedHero != entities.HeroOnca {
		t.Fatalf("Esperado HeroOnca após toggle, obtido %d", e.selectedHero)
	}
	if e.player != e.onca {
		t.Fatalf("Esperado player apontar para onca")
	}

	// Alterna de volta para Garoto
	e.toggleSelectedHero()
	if e.selectedHero != entities.HeroGaroto {
		t.Fatalf("Esperado HeroGaroto após segundo toggle, obtido %d", e.selectedHero)
	}
	if e.player != e.garoto {
		t.Fatalf("Esperado player apontar para garoto")
	}
}

func TestRestartStagePreservesStageAndHero(t *testing.T) {
	e := NewEngine()
	e.stage = 2
	e.score = 1500
	e.stageDistance = 800
	e.lives = 1
	e.hearts = 1
	e.isGameOver = true
	e.isPaused = true

	// Seleciona Onça
	e.selectedHero = entities.HeroOnca

	// Reinicia a fase
	e.restartStage()

	if e.stage != 2 {
		t.Fatalf("Esperado stage continuar 2, obtido %d", e.stage)
	}
	if e.lives != 3 {
		t.Fatalf("Esperado vidas restauradas para 3, obtido %d", e.lives)
	}
	if e.hearts != 3 {
		t.Fatalf("Esperado corações restaurados para 3, obtido %d", e.hearts)
	}
	if e.stageDistance != 0 {
		t.Fatalf("Esperado stageDistance zerada, obtido %f", e.stageDistance)
	}
	if e.isGameOver {
		t.Fatalf("Esperado isGameOver ser falso")
	}
	if e.isPaused {
		t.Fatalf("Esperado isPaused ser falso")
	}
	if e.player != e.onca {
		t.Fatalf("Esperado player ser a onca após reiniciar com HeroOnca selecionado")
	}
}

func TestResetGameResetsAllProgress(t *testing.T) {
	e := NewEngine()
	e.stage = 3
	e.score = 9999
	e.relicsCount = 8
	e.stageDistance = 2500
	e.lives = 1
	e.hearts = 1
	e.isGameOver = true

	// Seleciona Garoto
	e.selectedHero = entities.HeroGaroto

	// Reseta o jogo completo
	e.resetGame()

	if e.stage != 1 {
		t.Fatalf("Esperado stage resetado para 1, obtido %d", e.stage)
	}
	if e.score != 0 {
		t.Fatalf("Esperado score resetado para 0, obtido %d", e.score)
	}
	if e.relicsCount != 0 {
		t.Fatalf("Esperado relicsCount resetado para 0, obtido %d", e.relicsCount)
	}
	if e.stageDistance != 0 {
		t.Fatalf("Esperado stageDistance zerada, obtido %f", e.stageDistance)
	}
	if e.lives != 3 || e.hearts != 3 {
		t.Fatalf("Esperado 3 vidas e 3 corações, obtido vidas=%d coracoes=%d", e.lives, e.hearts)
	}
	if e.player != e.garoto {
		t.Fatalf("Esperado player ser garoto após reset com HeroGaroto selecionado")
	}
}

func TestReturnToTitle(t *testing.T) {
	e := NewEngine()
	e.isPaused = true
	e.isGameOver = true
	e.pauseMenuIndex = 3

	e.returnToTitle()

	if !e.isTitleCover {
		t.Fatalf("Esperado isTitleCover ser true")
	}
	if e.isPaused || e.isGameOver {
		t.Fatalf("Esperado isPaused e isGameOver serem falsos")
	}
	if e.pauseMenuIndex != 0 {
		t.Fatalf("Esperado pauseMenuIndex ser 0")
	}
}
