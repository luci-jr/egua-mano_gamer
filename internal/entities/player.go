package entities

import "github.com/hajimehoshi/ebiten/v2"

const (
	HeroGaroto = 0
	HeroOnca   = 1
)

// PlayerCharacter define a interface unificada para os heróis jogáveis (Garoto Curumim e Onça-Pintada)
type PlayerCharacter interface {
	MoveForward(speed float64)
	MoveBackward(speed float64)
	StopRunning()
	SetPositionX(x float64)
	GetPosition() (x, y float64)
	SetGroundOffset(offset float64)
	GetGroundOffset() float64
	GetVelocityY() float64
	FallFromPlatform()
	SetAimUp(aim bool)
	TriggerAttack()
	GetShootOrigin(groundY float64) (x, y, vx, vy float64)
	Jump() (jumped bool, isDouble bool)
	ReleaseJump()
	SetCrouch(crouch bool)
	FastDrop()
	GrabVine(v *Vine)
	ReleaseVine() bool
	Update()
	Reset()
	GetBounds(groundY float64) (x, y, w, h float64)
	Draw(screen *ebiten.Image, groundY float64, ticks int, invincibleTicks int)
	IsPlayerJumping() bool
	IsPlayerCrouching() bool
	IsPlayerSwinging() bool
	GetJumpHolding() bool
	GetJumpCount() int
	GetHeroKind() int
	GetName() string
}
