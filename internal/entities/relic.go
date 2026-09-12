package entities

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type RelicType int

const (
	RelicMuiraquita RelicType = 0 // Sapinho sagrado de jade verde (+100 pts)
	RelicUrna       RelicType = 1 // Vaso ancestral de cerâmica marajoara (+200 pts)
	RelicOuro       RelicType = 2 // Pepita de ouro do Tapajós (+50 pts)
	RelicAcaiBowl   RelicType = 3 // 🥣 Cuia de Tacacá / Tigela de Açaí Nutritiva (Recupera 1 Coração!)
)

type Relic struct {
	X        float64
	Y        float64
	BaseY    float64
	Type     RelicType
	Active   bool
	BobTimer int
	Value    int
}

func NewRelic(x, baseY float64, rType RelicType) *Relic {
	val := 100
	switch rType {
	case RelicUrna:
		val = 200
	case RelicOuro:
		val = 50
	case RelicAcaiBowl:
		val = 150
	}
	return &Relic{
		X:        x,
		Y:        baseY,
		BaseY:    baseY,
		Type:     rType,
		Active:   true,
		BobTimer: rand.Intn(100),
		Value:    val,
	}
}

func (r *Relic) Update(speed float64) {
	if !r.Active {
		return
	}
	r.X -= speed
	r.BobTimer++
	// Flutuação mística vertical
	r.Y = r.BaseY + math.Sin(float64(r.BobTimer)*0.08)*4.0

	if r.X < -50.0 {
		r.Active = false
	}
}

func (r *Relic) GetBounds() (float64, float64, float64, float64) {
	return r.X, r.Y, 16.0, 16.0
}

func (r *Relic) CheckCollection(playerX, playerY, playerW, playerH float64) bool {
	if !r.Active {
		return false
	}
	rx, ry, rw, rh := r.GetBounds()
	overlapX := rx < playerX+playerW && rx+rw > playerX
	overlapY := ry < playerY+playerH && ry+rh > playerY
	return overlapX && overlapY
}

func (r *Relic) Draw(screen *ebiten.Image, ticks int) {
	if !r.Active || r.X < -20 || r.X > 360 {
		return
	}

	rx := int(r.X)
	ry := int(r.Y)

	// Brilho pulsante em volta da relíquia
	glowAlpha := uint8(80 + int(math.Sin(float64(ticks)*0.15)*40.0))
	glowColor := color.RGBA{R: 255, G: 240, B: 150, A: glowAlpha}
	ebitenutil.DrawRect(screen, r.X-2, r.Y-2, 20, 20, glowColor)

	switch r.Type {
	case RelicMuiraquita:
		// 🐸 Muiraquitã Sagrado (Sapinho de Jade verde polida)
		cJade := color.RGBA{R: 35, G: 195, B: 95, A: 255}
		cJadeDark := color.RGBA{R: 20, G: 130, B: 60, A: 255}
		cJadeLight := color.RGBA{R: 120, G: 245, B: 160, A: 255}
		cGoldEye := color.RGBA{R: 250, G: 220, B: 50, A: 255}

		// Corpo do sapinho
		ebitenutil.DrawRect(screen, float64(rx+3), float64(ry+4), 10, 8, cJade)
		ebitenutil.DrawRect(screen, float64(rx+4), float64(ry+5), 8, 6, cJadeDark)
		ebitenutil.DrawRect(screen, float64(rx+5), float64(ry+5), 3, 2, cJadeLight)

		// Patas arqueadas
		ebitenutil.DrawRect(screen, float64(rx+1), float64(ry+8), 3, 4, cJade)
		ebitenutil.DrawRect(screen, float64(rx+12), float64(ry+8), 3, 4, cJade)
		ebitenutil.DrawRect(screen, float64(rx+2), float64(ry+3), 3, 3, cJade)
		ebitenutil.DrawRect(screen, float64(rx+11), float64(ry+3), 3, 3, cJade)

		// Olhos dourados sagrados
		ebitenutil.DrawRect(screen, float64(rx+4), float64(ry+3), 2, 2, cGoldEye)
		ebitenutil.DrawRect(screen, float64(rx+10), float64(ry+3), 2, 2, cGoldEye)

	case RelicUrna:
		// 🏺 Urna de Cerâmica Marajoara (grafismos vermelhos e pretos na terracota)
		cClay := color.RGBA{R: 215, G: 135, B: 80, A: 255}
		cClayDark := color.RGBA{R: 165, G: 95, B: 50, A: 255}
		cPattern := color.RGBA{R: 30, G: 25, B: 30, A: 255}
		cRed := color.RGBA{R: 210, G: 45, B: 35, A: 255}

		// Gargalo e boca
		ebitenutil.DrawRect(screen, float64(rx+4), float64(ry+1), 8, 3, cClayDark)
		ebitenutil.DrawRect(screen, float64(rx+5), float64(ry+2), 6, 2, cRed)

		// Bojo arredondado
		ebitenutil.DrawRect(screen, float64(rx+2), float64(ry+4), 12, 10, cClay)
		ebitenutil.DrawRect(screen, float64(rx+3), float64(ry+5), 10, 8, cClayDark)

		// Grafismos labirínticos marajoaras
		ebitenutil.DrawRect(screen, float64(rx+4), float64(ry+6), 3, 2, cPattern)
		ebitenutil.DrawRect(screen, float64(rx+9), float64(ry+6), 3, 2, cPattern)
		ebitenutil.DrawRect(screen, float64(rx+6), float64(ry+9), 4, 2, cRed)
		ebitenutil.DrawRect(screen, float64(rx+4), float64(ry+11), 8, 2, cPattern)

		// Base
		ebitenutil.DrawRect(screen, float64(rx+5), float64(ry+14), 6, 2, cClayDark)

	case RelicOuro:
		// 🪙 Pepita de Ouro do Rio Tapajós (facetas douradas cintilantes)
		cGoldBright := color.RGBA{R: 255, G: 235, B: 85, A: 255}
		cGoldMid := color.RGBA{R: 240, G: 185, B: 30, A: 255}
		cGoldShade := color.RGBA{R: 180, G: 125, B: 15, A: 255}
		cSparkle := color.RGBA{R: 255, G: 255, B: 240, A: 255}

		ebitenutil.DrawRect(screen, float64(rx+3), float64(ry+3), 10, 10, cGoldMid)
		ebitenutil.DrawRect(screen, float64(rx+5), float64(ry+2), 6, 12, cGoldBright)
		ebitenutil.DrawRect(screen, float64(rx+2), float64(ry+5), 12, 6, cGoldBright)
		ebitenutil.DrawRect(screen, float64(rx+6), float64(ry+7), 6, 5, cGoldShade)

		// Ponto de brilho estelar (cintilação)
		if (ticks/8)%2 == 0 {
			ebitenutil.DrawRect(screen, float64(rx+4), float64(ry+4), 2, 2, cSparkle)
		}

	case RelicAcaiBowl:
		// 🥣 Cuia de Tacacá Fumegante com Jambu & Camarão Seco (Recupera 1 Coração!)
		cCuia := color.RGBA{R: 110, G: 65, B: 30, A: 255}
		cCuiaDark := color.RGBA{R: 75, G: 40, B: 18, A: 255}
		cTucupi := color.RGBA{R: 245, G: 205, B: 40, A: 255}
		cJambu := color.RGBA{R: 45, G: 165, B: 55, A: 255}
		cCamarao := color.RGBA{R: 235, G: 95, B: 45, A: 255}
		cHeartRed := color.RGBA{R: 255, G: 45, B: 65, A: 255}
		cHeartLight := color.RGBA{R: 255, G: 175, B: 190, A: 255}

		// Cuia arredondada
		ebitenutil.DrawRect(screen, float64(rx+2), float64(ry+8), 12, 7, cCuia)
		ebitenutil.DrawRect(screen, float64(rx+4), float64(ry+14), 8, 2, cCuiaDark)
		ebitenutil.DrawRect(screen, float64(rx+1), float64(ry+7), 14, 2, cCuiaDark)

		// Caldo de tucupi borbulhante e folhas
		ebitenutil.DrawRect(screen, float64(rx+3), float64(ry+8), 10, 3, cTucupi)
		ebitenutil.DrawRect(screen, float64(rx+4), float64(ry+9), 3, 2, cJambu)
		ebitenutil.DrawRect(screen, float64(rx+8), float64(ry+8), 4, 2, cCamarao)

		// Coração flutuante pulsante sobre a cuia
		pulse := (ticks / 8) % 2
		hy := float64(ry + 1)
		if pulse == 0 {
			hy -= 1.0
		}
		ebitenutil.DrawRect(screen, float64(rx+4), hy, 3, 2, cHeartRed)
		ebitenutil.DrawRect(screen, float64(rx+8), hy, 3, 2, cHeartRed)
		ebitenutil.DrawRect(screen, float64(rx+3), hy+2, 9, 2, cHeartRed)
		ebitenutil.DrawRect(screen, float64(rx+4), hy+4, 7, 2, cHeartRed)
		ebitenutil.DrawRect(screen, float64(rx+6), hy+6, 3, 1, cHeartRed)
		ebitenutil.DrawRect(screen, float64(rx+5), hy+1, 1, 1, cHeartLight)
	}
}

// RelicManager orquestra o surgimento das relíquias no cenário
type RelicManager struct {
	Relics      []*Relic
	screenWidth float64
	groundY     float64
	spawnTimer  int
	TotalFound  int
}

func NewRelicManager(screenWidth, groundY float64) *RelicManager {
	rm := &RelicManager{
		Relics:      make([]*Relic, 0),
		screenWidth: screenWidth,
		groundY:     groundY,
		spawnTimer:  80,
		TotalFound:  0,
	}
	// Relíquia inicial suspensa
	rm.SpawnRelic(screenWidth+50.0, groundY-65.0, RelicMuiraquita)
	return rm
}

func (rm *RelicManager) SpawnRelic(x, y float64, rType RelicType) {
	rm.Relics = append(rm.Relics, NewRelic(x, y, rType))
}

func (rm *RelicManager) Update(speed float64) {
	rm.spawnTimer++
	// Spawna uma relíquia/item de cura a cada ~230 ticks
	if rm.spawnTimer >= 230 {
		rm.spawnTimer = 0
		rType := RelicType(rand.Intn(4)) // Sorteia entre Muiraquitã, Urna, Ouro e Cuia de Tacacá/Açaí
		spawnY := rm.groundY - 55.0 - float64(rand.Intn(35))
		rm.SpawnRelic(rm.screenWidth+40.0, spawnY, rType)
	}

	alive := rm.Relics[:0]
	for _, r := range rm.Relics {
		r.Update(speed)
		if r.Active {
			alive = append(alive, r)
		}
	}
	rm.Relics = alive
}

func (rm *RelicManager) CheckCollection(playerX, playerY, playerW, playerH float64) (bool, *Relic) {
	for _, r := range rm.Relics {
		if r.CheckCollection(playerX, playerY, playerW, playerH) {
			r.Active = false
			rm.TotalFound++
			return true, r
		}
	}
	return false, nil
}

func (rm *RelicManager) Draw(screen *ebiten.Image, ticks int) {
	for _, r := range rm.Relics {
		r.Draw(screen, ticks)
	}
}

func (rm *RelicManager) Reset() {
	rm.Relics = rm.Relics[:0]
	rm.spawnTimer = 80
	rm.TotalFound = 0
	rm.SpawnRelic(rm.screenWidth+60.0, rm.groundY-65.0, RelicMuiraquita)
}
