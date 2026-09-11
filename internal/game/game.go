package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/luci-jr/paidegua-game/internal/audio"
	"github.com/luci-jr/paidegua-game/internal/entities"
	"github.com/luci-jr/paidegua-game/internal/scenery"
	"github.com/luci-jr/paidegua-game/internal/ui"
)

const (
	ScreenWidth  = 340.0
	ScreenHeight = 210.0
	GroundY      = 155.0
	BaseSpeed    = 2.3
)

type Engine struct {
	onca            *entities.Onca
	obstacles       *entities.ObstacleManager
	scenery         *scenery.Background
	audio           *audio.Manager

	score             int
	ticks             int
	lives             int
	hearts            int
	invincibleTicks   int
	shakeTimer        int
	hitDelayTimer     int
	speechBubbleTimer int
	speechBubbleText  string
	stage             int
	stageBannerTimer  int
	isSaoBrasIntro    bool
	saoBrasTimer      int
	isTitleScreen     bool
	isShowingCredits  bool
	isPaused          bool
	pauseMenuIndex    int
	isGameOver        bool
	isStageComplete   bool
}

func NewEngine() *Engine {
	return &Engine{
		onca:              entities.NewOnca(),
		obstacles:         entities.NewObstacleManager(ScreenWidth, GroundY),
		scenery:           scenery.NewBackground(),
		audio:             audio.NewManager(),
		score:             0,
		ticks:             0,
		lives:             3,
		hearts:            3,
		invincibleTicks:   0,
		shakeTimer:        0,
		hitDelayTimer:     0,
		speechBubbleTimer: 0,
		speechBubbleText:  "",
		stage:             1,
		stageBannerTimer:  120,
		isSaoBrasIntro:    true,
		saoBrasTimer:      900, // ~15 segundos de visualização (ou avança com toque/clique/tecla)
		isTitleScreen:     false,
		isShowingCredits:  false,
		isPaused:          false,
		pauseMenuIndex:    0,
		isGameOver:        false,
		isStageComplete:   false,
	}
}

func isPointerJustPressed() bool {
	if getVirtualKey("Enter") {
		resetVirtualKey("Enter")
		return true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) ||
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return true
	}
	touches := inpututil.AppendJustPressedTouchIDs(nil)
	if len(touches) > 0 {
		return true
	}
	allTouches := ebiten.AppendTouchIDs(nil)
	return len(allTouches) > 0
}

func (e *Engine) Update() error {
	if e.isSaoBrasIntro {
		e.ticks++
		e.saoBrasTimer--
		anyInput := isPointerJustPressed()
		if !anyInput {
			for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
				if inpututil.IsKeyJustPressed(k) {
					anyInput = true
					break
				}
			}
		}
		if e.saoBrasTimer <= 0 || anyInput {
			e.isSaoBrasIntro = false
			e.isTitleScreen = true
		}
		return nil
	}

	if e.isShowingCredits {
		exitCredits := isPointerJustPressed() ||
			inpututil.IsKeyJustPressed(ebiten.KeyEscape) ||
			inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
			inpututil.IsKeyJustPressed(ebiten.KeyC)
		if exitCredits {
			e.isShowingCredits = false
			if !e.isPaused && !e.isTitleScreen && !e.isGameOver && !e.isStageComplete && !e.audio.IsMuted() {
				e.audio.ResumeBGM()
			}
		}
		return nil
	}

	if e.isTitleScreen {
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			e.isShowingCredits = true
			return nil
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			return ebiten.Termination
		}

		startPressed := isPointerJustPressed() ||
			inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeySpace)

		if startPressed {
			e.isTitleScreen = false
			e.stageBannerTimer = 120
			return nil
		}

		e.ticks++
		e.scenery.Update(0.6)
		e.onca.Update()
		return nil
	}

	if e.isPaused {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || getVirtualKey("ArrowUp") {
			resetVirtualKey("ArrowUp")
			e.pauseMenuIndex = (e.pauseMenuIndex - 1 + 5) % 5
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || getVirtualKey("ArrowDown") {
			resetVirtualKey("ArrowDown")
			e.pauseMenuIndex = (e.pauseMenuIndex + 1) % 5
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || getVirtualKey("Escape") {
			resetVirtualKey("Escape")
			e.isPaused = false
			if !e.audio.IsMuted() {
				e.audio.ResumeBGM()
			}
			return nil
		}

		selectPressed := isPointerJustPressed() ||
			inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeySpace)

		if selectPressed {
			switch e.pauseMenuIndex {
			case 0:
				e.isPaused = false
				if !e.audio.IsMuted() {
					e.audio.ResumeBGM()
				}
			case 1:
				e.onca.Reset()
				e.obstacles.Reset()
				e.score = 0
				e.lives = 3
				e.hearts = 3
				e.stage = 1
				e.stageBannerTimer = 120
				e.speechBubbleTimer = 0
				e.speechBubbleText = ""
				e.isPaused = false
				e.audio.RestartBGM()
			case 2:
				e.audio.ToggleMute()
			case 3:
				e.isShowingCredits = true
			case 4:
				return ebiten.Termination
			}
		}
		return nil
	}

	if e.isStageComplete {
		exitPressed := inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ)
		if exitPressed {
			return ebiten.Termination
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			e.onca.Reset()
			e.obstacles.Reset()
			e.score = 0
			e.lives = 3
			e.hearts = 3
			e.stage = 1
			e.stageBannerTimer = 120
			e.isStageComplete = false
			e.audio.RestartBGM()
			return nil
		}

		continuePressed := isPointerJustPressed() ||
			inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeySpace)

		if continuePressed {
			if e.stage < 3 {
				e.stage++
				e.stageBannerTimer = 130
				e.isStageComplete = false
				e.hearts = 3
				e.obstacles.Reset()
				if !e.audio.IsMuted() {
					e.audio.ResumeBGM()
				}
			} else {
				e.onca.Reset()
				e.obstacles.Reset()
				e.score = 0
				e.lives = 3
				e.hearts = 3
				e.stage = 1
				e.stageBannerTimer = 120
				e.isStageComplete = false
				e.audio.RestartBGM()
			}
			return nil
		}
		return nil
	}

	if e.isGameOver {
		exitPressed := inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ)
		if exitPressed {
			return ebiten.Termination
		}

		restartPressed := isPointerJustPressed() ||
			inpututil.IsKeyJustPressed(ebiten.KeyR) ||
			inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeySpace)

		if restartPressed {
			e.onca.Reset()
			e.obstacles.Reset()
			e.score = 0
			e.lives = 3
			e.hearts = 3
			e.stage = 1
			e.stageBannerTimer = 120
			e.invincibleTicks = 0
			e.shakeTimer = 0
			e.speechBubbleTimer = 0
			e.speechBubbleText = ""
			e.isPaused = false
			e.isGameOver = false
			e.audio.RestartBGM()
		}
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		e.isShowingCredits = true
		e.audio.PauseBGM()
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || getVirtualKey("Escape") {
		resetVirtualKey("Escape")
		e.isPaused = true
		e.pauseMenuIndex = 0
		e.audio.PauseBGM()
		return nil
	}

	e.ticks++

	if e.hitDelayTimer > 0 {
		e.hitDelayTimer--
		if e.shakeTimer > 0 {
			e.shakeTimer--
		}
		if e.speechBubbleTimer > 0 {
			e.speechBubbleTimer--
		}
		return nil
	}

	if e.speechBubbleTimer > 0 {
		e.speechBubbleTimer--
	}

	if e.invincibleTicks > 0 {
		e.invincibleTicks--
	}

	if e.shakeTimer > 0 {
		e.shakeTimer--
	}

	if e.stageBannerTimer > 0 {
		e.stageBannerTimer--
	}

	stageTarget := 1000
	if e.stage == 2 {
		stageTarget = 2000
	} else if e.stage == 3 {
		stageTarget = 3000
	}

	if e.score >= stageTarget && !e.isStageComplete {
		e.isStageComplete = true
		e.audio.PauseBGM()
		e.audio.PlayStageUp()
		return nil
	}

	// Movimentação horizontal da Onça (Adiantar e Recuar com Teclado, Botões Virtuais ou Toque no Canvas)
	moveSpeed := 2.2
	moveForward := ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) || getVirtualKey("ArrowRight")
	moveBackward := ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) || getVirtualKey("ArrowLeft")

	duckKey := ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) || getVirtualKey("ArrowDown")
	duckJustPressed := inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS)

	jumpJustPressed := inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyW) ||
		getVirtualKey("JustJump")
	if getVirtualKey("JustJump") {
		resetVirtualKey("JustJump")
	}

	jumpHolding := ebiten.IsKeyPressed(ebiten.KeyArrowUp) ||
		ebiten.IsKeyPressed(ebiten.KeySpace) ||
		ebiten.IsKeyPressed(ebiten.KeyW) ||
		getVirtualKey("Jump")

	// Mapeamento de toques nativos direto na tela do celular
	touches := ebiten.AppendTouchIDs(nil)
	for _, id := range touches {
		tx, ty := ebiten.TouchPosition(id)
		if tx < 80 {
			moveBackward = true
		} else if tx >= 80 && tx < 170 {
			moveForward = true
		} else if tx >= 170 {
			if ty < 135 {
				jumpHolding = true
			} else {
				duckKey = true
			}
		}
	}
	justTouches := inpututil.AppendJustPressedTouchIDs(nil)
	for _, id := range justTouches {
		tx, ty := ebiten.TouchPosition(id)
		if tx >= 170 && ty < 135 {
			jumpJustPressed = true
		} else if tx >= 170 && ty >= 135 {
			duckJustPressed = true
		}
	}

	if moveForward {
		e.onca.MoveForward(moveSpeed)
	}
	if moveBackward {
		e.onca.MoveBackward(moveSpeed)
	}

	if duckJustPressed && !e.onca.IsJumping {
		e.audio.PlayDuck()
	}

	if duckKey {
		if e.onca.IsJumping {
			e.onca.FastDrop()
		} else {
			e.onca.SetCrouch(true)
		}
	} else {
		e.onca.SetCrouch(false)
	}

	if jumpJustPressed {
		jumped, isDouble := e.onca.Jump()
		if jumped {
			e.audio.PlayRoar(isDouble)
		}
	}

	if !jumpHolding && e.onca.JumpHolding {
		e.onca.ReleaseJump()
	}

	e.onca.Update()

	stageProgress := float64(e.score % 1000)
	currentSpeed := BaseSpeed + float64(e.stage-1)*0.35 + (stageProgress / 2500.0)
	e.scenery.Update(currentSpeed)

	if passed := e.obstacles.Update(currentSpeed); passed > 0 {
		e.score += passed * 25
	}
	if e.ticks%6 == 0 {
		e.score += 1
	}

	oncaX, oncaY, oncaW, oncaH := e.onca.GetBounds(GroundY)
	hit, hitType := e.obstacles.CheckCollision(oncaX, oncaY, oncaW, oncaH)
	if hit && e.invincibleTicks <= 0 {
		e.hearts--
		e.shakeTimer = 14
		if e.hearts <= 0 {
			e.lives--
			if e.lives <= 0 {
				e.lives = 0
				e.hearts = 0
				e.isGameOver = true
				e.speechBubbleText = "Levei o farelo mano, mancada!"
				e.speechBubbleTimer = 999999
				e.audio.PauseBGM()
				e.audio.PlayGameOver()
			} else {
				e.hearts = 3 // Restaura os 3 corações para a próxima vida
				e.speechBubbleText = fmt.Sprintf("PERDEU 1 VIDA! RESTAM %d", e.lives)
				e.speechBubbleTimer = 85
				e.hitDelayTimer = 25
				e.invincibleTicks = 90
				e.audio.PlayHit()
			}
		} else {
			if hitType == entities.TypeJacare {
				e.speechBubbleText = "EGUA DO JACARE!..."
			} else if hitType == entities.TypeSnake {
				e.speechBubbleText = "VALHA-ME! UMA COBRA!"
			} else {
				e.speechBubbleText = "EGUA MANO!..."
			}
			e.speechBubbleTimer = 65
			e.hitDelayTimer = 22
			e.invincibleTicks = 75
			e.audio.PlayHit()
		}
	}

	return nil
}

func (e *Engine) Draw(screen *ebiten.Image) {
	if e.isSaoBrasIntro {
		ui.DrawSaoBrasIntro(screen, ScreenWidth, ScreenHeight, e.ticks)
		return
	}

	e.scenery.Draw(screen, ScreenWidth, GroundY, e.ticks, e.stage)

	if e.isTitleScreen {
		e.onca.Draw(screen, GroundY, e.ticks, 0)
		if e.isShowingCredits {
			ui.DrawCreditsScreen(screen, ScreenWidth, ScreenHeight)
		} else {
			ui.DrawTitleScreen(screen, ScreenWidth, ScreenHeight, e.ticks)
		}
		return
	}

	e.onca.Draw(screen, GroundY, e.ticks, e.invincibleTicks)

	if e.speechBubbleTimer > 0 {
		oncaX, oncaY, _, _ := e.onca.GetBounds(GroundY)
		bubbleX := oncaX + 8
		if e.isGameOver {
			bubbleX = oncaX - 25
		}
		ui.DrawSpeechBubble(screen, bubbleX, oncaY-24, e.speechBubbleText)
	}

	e.obstacles.Draw(screen, e.ticks, e.stage)

	isDoubleJump := e.onca.JumpCount == 2
	ui.DrawHUD(screen, e.lives, e.hearts, e.score, e.stage, isDoubleJump, e.stageBannerTimer, e.audio.IsMuted())

	if e.isShowingCredits {
		ui.DrawCreditsScreen(screen, ScreenWidth, ScreenHeight)
		return
	}

	if e.isPaused {
		ui.DrawPauseMenu(screen, ScreenWidth, ScreenHeight, e.pauseMenuIndex, e.audio.IsMuted())
	}

	if e.isStageComplete {
		ui.DrawStageCompleteScreen(screen, ScreenWidth, ScreenHeight, e.stage, e.score)
	}

	if e.isGameOver {
		ui.DrawGameOverScreen(screen, ScreenWidth, ScreenHeight, e.score, e.stage)
	}
}

func (e *Engine) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(ScreenWidth), int(ScreenHeight)
}

func Start() error {
	ebiten.SetWindowSize(680, 420)
	ebiten.SetWindowTitle("Pai D'Égua Game: A Aventura da Onça em Belém do Pará")

	engine := NewEngine()
	return ebiten.RunGame(engine)
}
