package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/luci-jr/paidegua-runner/internal/game"
)

func main() {
	if err := game.Start(); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
