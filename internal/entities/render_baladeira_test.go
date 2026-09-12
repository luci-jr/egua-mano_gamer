package entities

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func TestRenderBaladeiraPreview(t *testing.T) {
	width, height := 160, 100
	screen := ebiten.NewImage(width, height)
	ebitenutil.DrawRect(screen, 0, 0, float64(width), float64(height), color.RGBA{R: 15, G: 23, B: 42, A: 255})
	ebitenutil.DrawRect(screen, 0, 75, float64(width), 25, color.RGBA{R: 50, G: 65, B: 85, A: 255})

	gIdle := NewGaroto()
	gIdle.X = 20
	gIdle.Draw(screen, 75, 60, 0)

	gAttack := NewGaroto()
	gAttack.X = 70
	gAttack.TriggerAttack()
	gAttack.Draw(screen, 75, 60, 0)

	gAim := NewGaroto()
	gAim.X = 120
	gAim.SetAimUp(true)
	gAim.Draw(screen, 75, 60, 0)
}
