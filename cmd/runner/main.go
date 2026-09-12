package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/luci-jr/egua-mano_gamer/internal/game"
)

func main() {
	if err := game.Start(); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
