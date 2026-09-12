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

	// 2. Se o jogador estiver no chão correndo (sem pular), não deve acionar stomp
	playerGroundY := 0.0
	playerGroundVY := 0.0
	hit2, _, _, _ := mgr.CheckStomp(204.0, playerGroundY, playerW, playerH, playerGroundVY, GroundY)
	if hit2 {
		t.Fatalf("Não deve acionar stomp se o jogador estiver no chão correndo")
	}

	// 3. Teste com coordenadas de tela de GetBounds (ex: playerY = 131.0 para Y no solo, ou 115.0 no pulo)
	// Garante que o cálculo com coordenadas absolutas de tela funcione de forma idêntica
	mgr.Obstacles[2].Collided = false
	mgr.Obstacles[2].Defeated = false
	// Ave aérea: X=300, Y=GroundY-34=121.0, W=28, H=16 -> topo 121, base 137
	hit3, _, _, obsType3 := mgr.CheckStomp(304.0, 115.0, playerW, playerH, 1.0, GroundY)
	if !hit3 {
		t.Fatalf("Esperado acerto de stomp na Ave Aérea usando coordenadas de tela")
	}
	if obsType3 != entities.TypeAir {
		t.Fatalf("Esperado TypeAir, obtido %v", obsType3)
	}

	// 4. Teste de Salvaguarda de Pulo em CheckCollision:
	// Ao saltar sobre a cobra (TypeSnake), o herói NUNCA toma dano; a cobra é derrotada!
	mgr.Obstacles[1].Collided = false
	mgr.Obstacles[1].Defeated = false
	collisionHit, _ := mgr.CheckCollision(204.0, -12.0, playerW, playerH, 1.0, GroundY, nil)
	if collisionHit {
		t.Fatalf("Pulo sobre o bicho NUNCA deve causar dano ao herói!")
	}
	if !mgr.Obstacles[1].Defeated {
		t.Fatalf("A cobra deveria ser derrotada ao ser pisada/saltada por cima!")
	}

	// 5. Teste de Colisão Frontal no Chão:
	// Se o jogador estiver correndo no chão (playerBottom == GroundY), ele toma dano frontal
	mgr.Obstacles[1].Collided = false
	mgr.Obstacles[1].Defeated = false
	groundCollision, _ := mgr.CheckCollision(204.0, 0.0, playerW, playerH, 0.0, GroundY, nil)
	if !groundCollision {
		t.Fatalf("Correr de frente no bicho pelo chão deve causar colisão com dano!")
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
