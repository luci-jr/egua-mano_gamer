package entities

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestGarotoBaladeiraDrawNoPanic(t *testing.T) {
	g := NewGaroto()
	img := ebiten.NewImage(320, 224)

	// 1. Idle
	g.Draw(img, 155.0, 60, 0)

	// 2. Correndo
	g.IsRunning = true
	for frame := 0; frame < 20; frame++ {
		g.RunTicks = frame
		g.Draw(img, 155.0, frame, 0)
	}

	// 3. Atacando no chão
	g.TriggerAttack()
	g.Draw(img, 155.0, 60, 0)
	g.IsAttacking = false

	// 4. Mirando para cima (anti-aéreo)
	g.SetAimUp(true)
	g.Draw(img, 155.0, 60, 0)
	g.TriggerAttack()
	g.Draw(img, 155.0, 60, 0)
	g.SetAimUp(false)
	g.IsAttacking = false

	// 5. Agachado
	g.SetCrouch(true)
	g.Draw(img, 155.0, 60, 0)
	g.SetCrouch(false)

	// 6. Pulando
	g.IsJumping = true
	g.VelocityY = -5.0
	g.Draw(img, 155.0, 60, 0)
	g.TriggerAttack()
	g.Draw(img, 155.0, 60, 0)
	g.SetAimUp(true)
	g.Draw(img, 155.0, 60, 0)
	g.IsJumping = false
	g.SetAimUp(false)
	g.IsAttacking = false

	// 7. Balançando no cipó
	vine := NewVine(160.0)
	g.GrabVine(vine)
	g.Draw(img, 155.0, 60, 0)
}

func TestPlayerAimDirections(t *testing.T) {
	heroes := []PlayerCharacter{NewGaroto(), NewOnca()}

	for _, p := range heroes {
		// 1. Tiro horizontal padrão (Facing Right)
		p.SetAimDirection(0, 0)
		_, _, vx, vy := p.GetShootOrigin(155.0)
		if vx <= 0 || vy != 0 {
			t.Errorf("%s: esperava vx > 0 e vy == 0, obteve vx=%f, vy=%f", p.GetName(), vx, vy)
		}

		// 2. Tiro diagonal cima-direita
		p.SetAimDirection(1.0, -1.0)
		_, _, vx, vy = p.GetShootOrigin(155.0)
		if vx <= 0 || vy >= 0 {
			t.Errorf("%s: tiro diagonal cima-direita falhou: vx=%f, vy=%f", p.GetName(), vx, vy)
		}

		// 3. Tiro diagonal cima-esquerda
		p.SetAimDirection(-1.0, -1.0)
		_, _, vx, vy = p.GetShootOrigin(155.0)
		if vx >= 0 || vy >= 0 {
			t.Errorf("%s: tiro diagonal cima-esquerda falhou: vx=%f, vy=%f", p.GetName(), vx, vy)
		}

		// 4. Tiro diagonal baixo-direita
		p.SetAimDirection(1.0, 1.0)
		_, _, vx, vy = p.GetShootOrigin(155.0)
		if vx <= 0 || vy <= 0 {
			t.Errorf("%s: tiro diagonal baixo-direita falhou: vx=%f, vy=%f", p.GetName(), vx, vy)
		}

		// 5. Tiro diagonal baixo-esquerda
		p.SetAimDirection(-1.0, 1.0)
		_, _, vx, vy = p.GetShootOrigin(155.0)
		if vx >= 0 || vy <= 0 {
			t.Errorf("%s: tiro diagonal baixo-esquerda falhou: vx=%f, vy=%f", p.GetName(), vx, vy)
		}

		// 6. Tiro vertical puro para cima
		p.SetAimDirection(0, -1.0)
		_, _, vx, vy = p.GetShootOrigin(155.0)
		if vx != 0 || vy >= 0 {
			t.Errorf("%s: tiro vertical cima falhou: vx=%f, vy=%f", p.GetName(), vx, vy)
		}

		// 7. Apontando puro para baixo (acocado / agachado): tiro frontal horizontal na altura agachada
		p.MoveForward(1.0) // Garante FacingRight
		p.SetAimDirection(0, 1.0)
		_, shootY, vx, vy := p.GetShootOrigin(155.0)
		if vx <= 0 || vy != 0 {
			t.Errorf("%s: tiro acocado deve sair horizontal para a frente, obteve vx=%f, vy=%f", p.GetName(), vx, vy)
		}
		if shootY != 155.0-8.0 {
			t.Errorf("%s: tiro acocado deve sair na altura de agachamento (y=147), obteve y=%f", p.GetName(), shootY)
		}
	}
}
