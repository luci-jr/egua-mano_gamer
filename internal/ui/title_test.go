package ui

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestRenderTitleCoverNoPanic(t *testing.T) {
	img := ebiten.NewImage(320, 224)
	DrawTitleCoverScreen(img, 320, 224, 60, "VER. 2.4.0")
	DrawTitleIntro(img, 320, 224, 60, true, 0, false, "NORMAL", "Curumim", "VER. 2.4.0")
	DrawPauseMenu(img, 320, 224, 0, false, "NORMAL")
	DrawCharacterSelectScreen(img, 320, 224, 60, 0)
	DrawCharacterSelectScreen(img, 320, 224, 60, 1)
}

func TestRenderHeroWaterFallNoPanic(t *testing.T) {
	img := ebiten.NewImage(320, 224)
	// Garoto caindo no ar
	DrawHeroWaterFall(img, 80, 100, 0, -2.5, false)
	// Garoto na água
	DrawHeroWaterFall(img, 80, 160, 0, 1.0, true)
	// Onça caindo no ar
	DrawHeroWaterFall(img, 80, 100, 1, -2.5, false)
	// Onça na água
	DrawHeroWaterFall(img, 80, 160, 1, 1.0, true)
	// Splash de água
	DrawWaterSplash(img, 80, 160, 30)
	DrawWaterSplash(img, 80, 160, 5)
}
