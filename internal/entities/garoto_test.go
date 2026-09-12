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
