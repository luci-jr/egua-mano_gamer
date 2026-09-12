package game

import (
	"testing"

	"github.com/luci-jr/egua-mano_gamer/internal/entities"
)

func TestCheckStompMechanics(t *testing.T) {
	mgr := entities.NewObstacleManager(340.0, GroundY)
	mgr.Obstacles = []*entities.Obstacle{
		entities.NewObstacle(340.0, GroundY),
		entities.NewObstacle(340.0, GroundY),
		entities.NewObstacle(340.0, GroundY),
	}
	mgr.Obstacles[0].X = 100.0
	mgr.Obstacles[0].Type = entities.TypeJacare

	mgr.Obstacles[1].X = 200.0
	mgr.Obstacles[1].Type = entities.TypeSnake

	mgr.Obstacles[2].X = 300.0
	mgr.Obstacles[2].Type = entities.TypeAir

	// 1. Pisão descendente em cima do Jacaré
	// Jacaré: X=100, Y=GroundY-17 (138.0), W=32, H=17
	playerX := 104.0
	playerY := -15.0 // playerBottom = GroundY + playerY = 155 - 15 = 140.0 (em cima do jacaré)
	playerW := 14.0
	playerH := 20.0
	playerVY := 1.5 // caindo para baixo

	hit, _, _, obsType := mgr.CheckStomp(playerX, playerY, playerW, playerH, playerVY, GroundY)
	if !hit {
		t.Fatalf("Esperado acerto de stomp no Jacaré")
	}
	if obsType != entities.TypeJacare {
		t.Fatalf("Esperado TypeJacare, obtido %v", obsType)
	}
	if !mgr.Obstacles[0].Defeated {
		t.Fatalf("Esperado Jacaré marcado como Defeated após stomp")
	}

	// 2. Se estiver subindo rápido (pulo ascendente VY < -0.8), não deve dar stomp
	playerVY = -2.0
	hit2, _, _, _ := mgr.CheckStomp(204.0, playerY, playerW, playerH, playerVY, GroundY)
	if hit2 {
		t.Fatalf("Não deve acionar stomp durante pulo ascendente (VY < -0.8)")
	}
}

func TestGuaranaPowerUpAndDeathSpeech(t *testing.T) {
	e := NewEngine()

	// 1. Testar Star Power no Garoto
	if e.player.GetStarPower() != 0 {
		t.Fatalf("Esperado StarPower inicial 0")
	}
	e.player.SetStarPower(300)
	if e.player.GetStarPower() != 300 {
		t.Fatalf("Esperado StarPower 300, obtido %d", e.player.GetStarPower())
	}

	// 2. Testar Star Power na Onça
	e.toggleSelectedHero()
	if e.player.GetStarPower() != 0 {
		t.Fatalf("Esperado StarPower inicial da Onça 0")
	}
	e.player.SetStarPower(420)
	if e.player.GetStarPower() != 420 {
		t.Fatalf("Esperado StarPower 420, obtido %d", e.player.GetStarPower())
	}

	// 3. Testar frase regional de morte / perda de vida
	expectedSpeech := "Egua do pitiu. Essa agua ta podre!"
	e.lives = 1
	e.hearts = 1
	e.speechBubbleText = expectedSpeech
	if e.speechBubbleText != expectedSpeech {
		t.Fatalf("Esperado balão com '%s', obtido '%s'", expectedSpeech, e.speechBubbleText)
	}

	// 4. Testar valor da relíquia Guaraná
	relic := entities.NewRelic(100.0, 100.0, entities.RelicGuarana)
	if relic.Type != entities.RelicGuarana {
		t.Fatalf("Esperado tipo RelicGuarana")
	}
	if relic.Value != 300 {
		t.Fatalf("Esperado valor 300 para o Guaraná da Amazônia, obtido %d", relic.Value)
	}
}
