package ui

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestRenderTitleCoverNoPanic(t *testing.T) {
	img := ebiten.NewImage(320, 224)
	DrawTitleCoverScreen(img, 320, 224, 60, "VER. 2.4.0")
	DrawTitleIntro(img, 320, 224, 60, true, 0, false, "1.0x", "Curumim", "VER. 2.4.0")
}
