package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Vine struct {
	AnchorX       float64
	AnchorY       float64
	Length        float64
	Angle         float64
	AngleVelocity float64
	MaxAngle      float64
	SwingTimer    int
	Active        bool
	Grabbed       bool
}

func NewVine(anchorX float64) *Vine {
	return &Vine{
		AnchorX:       anchorX,
		AnchorY:       -4.0,
		Length:        112.0,
		Angle:         0,
		AngleVelocity: 0,
		MaxAngle:      0.68, // ~39 graus de arco pendular
		SwingTimer:    0,
		Active:        true,
		Grabbed:       false,
	}
}

func (v *Vine) Update(speed float64) {
	if !v.Active {
		return
	}
	v.AnchorX -= speed
	v.SwingTimer++

	// Movimento harmônico simples do pêndulo
	omega := 0.052
	v.Angle = v.MaxAngle * math.Sin(float64(v.SwingTimer)*omega)
	v.AngleVelocity = v.MaxAngle * omega * math.Cos(float64(v.SwingTimer)*omega)

	if v.AnchorX < -150.0 {
		v.Active = false
		v.Grabbed = false
	}
}

func (v *Vine) GetTipPosition() (float64, float64) {
	tipX := v.AnchorX + v.Length*math.Sin(v.Angle)
	tipY := v.AnchorY + v.Length*math.Cos(v.Angle)
	return tipX, tipY
}

func (v *Vine) CanGrab(playerX, playerY, playerW, playerH float64) bool {
	if !v.Active || v.Grabbed {
		return false
	}
	tipX, tipY := v.GetTipPosition()

	// Caixa de detecção ao redor da ponta/nó inferior do cipó
	grabRadius := 20.0
	playerCenterX := playerX + playerW/2.0
	playerCenterY := playerY + playerH/2.0

	dx := playerCenterX - tipX
	dy := playerCenterY - tipY
	distSq := dx*dx + dy*dy

	return distSq <= grabRadius*grabRadius
}

func (v *Vine) Draw(screen *ebiten.Image) {
	if !v.Active || v.AnchorX < -120 || v.AnchorX > 450 {
		return
	}

	cVineDark := color.RGBA{R: 55, G: 85, B: 30, A: 255}
	cVineLight := color.RGBA{R: 88, G: 140, B: 45, A: 255}
	cLeaf := color.RGBA{R: 40, G: 165, B: 55, A: 255}

	segments := 9
	prevX := v.AnchorX
	prevY := v.AnchorY

	// Efeito sutil de catenária/curvatura natural ao balançar
	for s := 1; s <= segments; s++ {
		t := float64(s) / float64(segments)
		segAngle := v.Angle * (0.85 + 0.15*t)
		currX := v.AnchorX + (v.Length*t)*math.Sin(segAngle)
		currY := v.AnchorY + (v.Length*t)*math.Cos(segAngle)

		// Desenha linha grossa do cipó (2px)
		ebitenutil.DrawLine(screen, prevX, prevY, currX, currY, cVineDark)
		ebitenutil.DrawLine(screen, prevX+1, prevY, currX+1, currY, cVineLight)

		// Folhas brotando ao longo do cipó
		if s%2 == 0 && s < segments {
			leafDir := 1.0
			if s%4 == 0 {
				leafDir = -1.0
			}
			lx := currX + leafDir*3.0
			ly := currY
			ebitenutil.DrawRect(screen, lx, ly, 3, 2, cLeaf)
		}

		prevX = currX
		prevY = currY
	}

	// Nó / Ponta reforçada onde o personagem se segura
	tipX, tipY := v.GetTipPosition()
	ebitenutil.DrawRect(screen, tipX-2, tipY-2, 5, 5, color.RGBA{R: 120, G: 80, B: 35, A: 255})
	ebitenutil.DrawRect(screen, tipX-1, tipY-1, 3, 3, cVineLight)
}

// VineManager gerencia os cipós que surgem na selva
type VineManager struct {
	Vines       []*Vine
	screenWidth float64
	spawnTimer  int
}

func NewVineManager(screenWidth float64) *VineManager {
	return &VineManager{
		Vines:       make([]*Vine, 0),
		screenWidth: screenWidth,
		spawnTimer:  0,
	}
}

func (vm *VineManager) SpawnAt(x float64) {
	// Desativado a pedido do jogador
}

func (vm *VineManager) Update(speed float64) {
	// Desativado
}

func (vm *VineManager) CheckGrab(playerX, playerY, playerW, playerH float64) *Vine {
	return nil
}

func (vm *VineManager) Draw(screen *ebiten.Image) {
	// Não desenha nada
}

func (vm *VineManager) Reset() {
	vm.Vines = vm.Vines[:0]
	vm.spawnTimer = 0
}
