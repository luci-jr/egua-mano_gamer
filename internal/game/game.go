package game

import (
	"fmt"
	"image/color"

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

var (
	SpeedMultipliers = []float64{0.8, 1.0, 1.3, 1.6}
	SpeedLabels      = []string{"0.8x CALMO", "1.0x NORMAL", "1.3x RAPIDO", "1.6x TURBO"}
)

type Engine struct {
	player            entities.PlayerCharacter
	garoto            *entities.Garoto
	onca              *entities.Onca
	selectedHero      int
	isCharSelect      bool
	projectiles       *entities.ProjectileManager
	obstacles         *entities.ObstacleManager
	vines             *entities.VineManager
	relics            *entities.RelicManager
	scenery           *scenery.Background
	audio             *audio.Manager

	score             int
	relicsCount       int
	ticks             int
	lives             int
	hearts            int
	invincibleTicks   int
	shakeTimer        int
	hitDelayTimer     int
	mudSinkTimer      int
	speechBubbleTimer int
	speechBubbleText  string
	heatSpeechTimer   int
	stage             int
	stageDistance     float64
	stageBannerTimer  int
	isSaoBrasIntro    bool
	saoBrasTimer      int
	introMenuIndex    int
	speedIndex        int
	isTitleScreen     bool
	isShowingCredits  bool
	isPaused          bool
	pauseMenuIndex    int
	isGameOver        bool
	isStageComplete   bool
}

func NewEngine() *Engine {
	audioMgr := audio.NewManager()
	audioMgr.PlayIntroBGM()

	g := entities.NewGaroto()
	o := entities.NewOnca()

	return &Engine{
		player:            g,
		garoto:            g,
		onca:              o,
		selectedHero:      entities.HeroGaroto,
		isCharSelect:      false,
		projectiles:       entities.NewProjectileManager(ScreenWidth),
		obstacles:         entities.NewObstacleManager(ScreenWidth, GroundY),
		vines:             entities.NewVineManager(ScreenWidth),
		relics:            entities.NewRelicManager(ScreenWidth, GroundY),
		scenery:           scenery.NewBackground(),
		audio:             audioMgr,
		score:             0,
		relicsCount:       0,
		stageDistance:     0,
		ticks:             0,
		lives:             3,
		hearts:            3,
		invincibleTicks:   0,
		shakeTimer:        0,
		hitDelayTimer:     0,
		mudSinkTimer:      0,
		speechBubbleTimer: 0,
		speechBubbleText:  "",
		heatSpeechTimer:   0,
		stage:             1,
		stageBannerTimer:  120,
		isSaoBrasIntro:    true,
		saoBrasTimer:      0,
		introMenuIndex:    0,
		speedIndex:        1,
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
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		return true
	}
	touches := inpututil.AppendJustPressedTouchIDs(nil)
	return len(touches) > 0
}

func (e *Engine) Update() error {
	// 1. Tela de Créditos (visível sobre qualquer tela)
	if e.isShowingCredits {
		exitCredits := isPointerJustPressed() ||
			inpututil.IsKeyJustPressed(ebiten.KeyEscape) ||
			inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
			inpututil.IsKeyJustPressed(ebiten.KeyC) ||
			getVirtualKey("Escape") ||
			getVirtualKey("Enter") ||
			getVirtualKey("Space")

		if getVirtualKey("Escape") {
			resetVirtualKey("Escape")
		}
		if getVirtualKey("Enter") {
			resetVirtualKey("Enter")
		}
		if getVirtualKey("Space") {
			resetVirtualKey("Space")
		}

		if exitCredits {
			e.isShowingCredits = false
			if !e.isPaused && !e.isTitleScreen && !e.isSaoBrasIntro && !e.isCharSelect && !e.isGameOver && !e.isStageComplete && !e.audio.IsMuted() {
				e.audio.ResumeBGM()
			}
		}
		return nil
	}

	// 2. Tela de Seleção de Personagem (Garoto Curumim vs Onça-Pintada)
	if e.isCharSelect {
		e.ticks++
		if !e.audio.IsMuted() {
			e.audio.PlayIntroBGM()
		}

		leftPressed := inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) || getVirtualKey("ArrowLeft")
		if getVirtualKey("ArrowLeft") {
			resetVirtualKey("ArrowLeft")
		}
		rightPressed := inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) || getVirtualKey("ArrowRight")
		if getVirtualKey("ArrowRight") {
			resetVirtualKey("ArrowRight")
		}

		if leftPressed && e.selectedHero != entities.HeroGaroto {
			e.selectedHero = entities.HeroGaroto
			e.player = e.garoto
			e.audio.PlayShot()
		}
		if rightPressed && e.selectedHero != entities.HeroOnca {
			e.selectedHero = entities.HeroOnca
			e.player = e.onca
			e.audio.PlayRoar(false)
		}

		cardClicked := false
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			if float64(my) >= 32.0 && float64(my) <= 174.0 {
				if float64(mx) < ScreenWidth/2.0 {
					if e.selectedHero != entities.HeroGaroto {
						e.selectedHero = entities.HeroGaroto
						e.player = e.garoto
						e.audio.PlayShot()
					} else {
						cardClicked = true
					}
				} else {
					if e.selectedHero != entities.HeroOnca {
						e.selectedHero = entities.HeroOnca
						e.player = e.onca
						e.audio.PlayRoar(false)
					} else {
						cardClicked = true
					}
				}
			} else if float64(my) > 174.0 {
				cardClicked = true
			}
		}

		touches := inpututil.AppendJustPressedTouchIDs(nil)
		for _, id := range touches {
			tx, ty := ebiten.TouchPosition(id)
			if ty >= 32 && ty <= 174 {
				if float64(tx) < ScreenWidth/2.0 {
					if e.selectedHero != entities.HeroGaroto {
						e.selectedHero = entities.HeroGaroto
						e.player = e.garoto
						e.audio.PlayShot()
					} else {
						cardClicked = true
					}
				} else {
					if e.selectedHero != entities.HeroOnca {
						e.selectedHero = entities.HeroOnca
						e.player = e.onca
						e.audio.PlayRoar(false)
					} else {
						cardClicked = true
					}
				}
			} else if ty > 174 {
				cardClicked = true
			}
		}

		confirmPressed := cardClicked ||
			inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
			getVirtualKey("Enter") ||
			getVirtualKey("Space")

		if getVirtualKey("Enter") {
			resetVirtualKey("Enter")
		}
		if getVirtualKey("Space") {
			resetVirtualKey("Space")
		}

		if confirmPressed {
			if e.selectedHero == entities.HeroGaroto {
				e.player = e.garoto
			} else {
				e.player = e.onca
			}
			e.isCharSelect = false
			e.isSaoBrasIntro = false
			e.isTitleScreen = false
			e.stage = 1
			e.stageDistance = 0
			e.stageBannerTimer = 120
			e.audio.StopIntroBGM()
			e.audio.RestartBGM()
			return nil
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || getVirtualKey("Escape") {
			resetVirtualKey("Escape")
			e.isCharSelect = false
			e.isSaoBrasIntro = true
			return nil
		}

		return nil
	}

	// 3. Menu Principal (Abertura oficial)
	if e.isSaoBrasIntro {
		e.ticks++
		if !e.audio.IsMuted() {
			e.audio.PlayIntroBGM()
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			e.isShowingCredits = true
			return nil
		}

		// Navegação no Menu Principal (5 opções)
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) || getVirtualKey("ArrowUp") {
			resetVirtualKey("ArrowUp")
			e.introMenuIndex = (e.introMenuIndex - 1 + 5) % 5
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) || getVirtualKey("ArrowDown") {
			resetVirtualKey("ArrowDown")
			e.introMenuIndex = (e.introMenuIndex + 1) % 5
		}

		// Ajuste com Esquerda / Direita
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
			if e.introMenuIndex == 1 {
				if e.selectedHero == entities.HeroGaroto {
					e.selectedHero = entities.HeroOnca
					e.player = e.onca
					e.audio.PlayRoar(false)
				} else {
					e.selectedHero = entities.HeroGaroto
					e.player = e.garoto
					e.audio.PlayShot()
				}
			} else if e.introMenuIndex == 2 {
				e.audio.ToggleMute()
			} else if e.introMenuIndex == 3 {
				e.speedIndex = (e.speedIndex - 1 + len(SpeedMultipliers)) % len(SpeedMultipliers)
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
			if e.introMenuIndex == 1 {
				if e.selectedHero == entities.HeroGaroto {
					e.selectedHero = entities.HeroOnca
					e.player = e.onca
					e.audio.PlayRoar(false)
				} else {
					e.selectedHero = entities.HeroGaroto
					e.player = e.garoto
					e.audio.PlayShot()
				}
			} else if e.introMenuIndex == 2 {
				e.audio.ToggleMute()
			} else if e.introMenuIndex == 3 {
				e.speedIndex = (e.speedIndex + 1) % len(SpeedMultipliers)
			}
		}

		// Detecção de clique / toque no Menu
		mouseTriggered := false
		boxX := (ScreenWidth - 226.0) / 2.0
		boxY := 35.0
		startY := boxY + 21.0

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			fmx, fmy := float64(mx), float64(my)
			if fmx >= boxX && fmx <= boxX+226.0 && fmy >= startY && fmy <= startY+5*14.0 {
				clickedIdx := int((fmy - startY) / 14.0)
				if clickedIdx >= 0 && clickedIdx < 5 {
					e.introMenuIndex = clickedIdx
					mouseTriggered = true
				}
			} else if fmy >= boxY && fmy <= boxY+94.0 {
				mouseTriggered = true
			}
		}

		menuTouches := inpututil.AppendJustPressedTouchIDs(nil)
		for _, id := range menuTouches {
			tx, ty := ebiten.TouchPosition(id)
			ftx, fty := float64(tx), float64(ty)
			if ftx >= boxX && ftx <= boxX+226.0 && fty >= startY && fty <= startY+5*14.0 {
				clickedIdx := int((fty - startY) / 14.0)
				if clickedIdx >= 0 && clickedIdx < 5 {
					e.introMenuIndex = clickedIdx
					mouseTriggered = true
				}
			} else {
				mouseTriggered = true
			}
		}

		selectTriggered := mouseTriggered ||
			inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
			getVirtualKey("Enter") ||
			getVirtualKey("Space")

		if getVirtualKey("Enter") {
			resetVirtualKey("Enter")
		}
		if getVirtualKey("Space") {
			resetVirtualKey("Space")
		}

		if selectTriggered {
			switch e.introMenuIndex {
			case 0:
				e.isCharSelect = true
				e.isSaoBrasIntro = false
				return nil
			case 1:
				if e.selectedHero == entities.HeroGaroto {
					e.selectedHero = entities.HeroOnca
					e.player = e.onca
					e.audio.PlayRoar(false)
				} else {
					e.selectedHero = entities.HeroGaroto
					e.player = e.garoto
					e.audio.PlayShot()
				}
				e.isCharSelect = true
				e.isSaoBrasIntro = false
				return nil
			case 2:
				e.audio.ToggleMute()
			case 3:
				e.speedIndex = (e.speedIndex + 1) % len(SpeedMultipliers)
			case 4:
				e.isShowingCredits = true
			}
			return nil
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
			e.audio.RestartBGM()
			return nil
		}

		e.ticks++
		e.scenery.Update(0.6)
		e.player.Update()
		return nil
	}

	if e.isPaused {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) || getVirtualKey("ArrowUp") {
			resetVirtualKey("ArrowUp")
			e.pauseMenuIndex = (e.pauseMenuIndex - 1 + 6) % 6
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) || getVirtualKey("ArrowDown") {
			resetVirtualKey("ArrowDown")
			e.pauseMenuIndex = (e.pauseMenuIndex + 1) % 6
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
			if e.pauseMenuIndex == 1 {
				e.audio.ToggleMute()
			} else if e.pauseMenuIndex == 2 {
				e.speedIndex = (e.speedIndex - 1 + len(SpeedMultipliers)) % len(SpeedMultipliers)
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
			if e.pauseMenuIndex == 1 {
				e.audio.ToggleMute()
			} else if e.pauseMenuIndex == 2 {
				e.speedIndex = (e.speedIndex + 1) % len(SpeedMultipliers)
			}
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
				e.audio.ToggleMute()
			case 2:
				e.speedIndex = (e.speedIndex + 1) % len(SpeedMultipliers)
			case 3:
				e.player.Reset()
				e.projectiles.Reset()
				e.obstacles.Reset()
				e.vines.Reset()
				e.relics.Reset()
				e.score = 0
				e.relicsCount = 0
				e.mudSinkTimer = 0
				e.lives = 3
				e.hearts = 3
				e.stage = 1
				e.stageBannerTimer = 120
				e.speechBubbleTimer = 0
				e.speechBubbleText = ""
				e.isPaused = false
				e.audio.RestartBGM()
			case 4:
				e.isShowingCredits = true
			case 5:
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
			e.player.Reset()
			e.projectiles.Reset()
			e.obstacles.Reset()
			e.vines.Reset()
			e.relics.Reset()
			e.score = 0
			e.relicsCount = 0
			e.stageDistance = 0
			e.mudSinkTimer = 0
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
				e.stageDistance = 0
				e.stageBannerTimer = 130
				e.isStageComplete = false
				e.hearts = 3
				e.projectiles.Reset()
				e.obstacles.Reset()
				e.vines.Reset()
				e.relics.Reset()
				if !e.audio.IsMuted() {
					e.audio.ResumeBGM()
				}
			} else {
				e.player.Reset()
				e.projectiles.Reset()
				e.obstacles.Reset()
				e.vines.Reset()
				e.relics.Reset()
				e.score = 0
				e.relicsCount = 0
				e.stageDistance = 0
				e.mudSinkTimer = 0
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
			e.player.Reset()
			e.projectiles.Reset()
			e.obstacles.Reset()
			e.vines.Reset()
			e.relics.Reset()
			e.score = 0
			e.relicsCount = 0
			e.stageDistance = 0
			e.mudSinkTimer = 0
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
	} else {
		// A cada ~14 a 20 segundos correndo sob o sol de Belém, a onça reclama do calor
		e.heatSpeechTimer++
		if e.heatSpeechTimer >= 850 {
			if e.hitDelayTimer == 0 {
				e.heatSpeechTimer = 0
				e.speechBubbleText = "Egua da lua, um sol pra cada um!"
				e.speechBubbleTimer = 165
			}
		}
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

	// Movimentação horizontal do herói (Adiantar e Recuar com Teclado, Botões Virtuais ou Toque no Canvas)
	moveSpeed := 2.2
	if e.selectedHero == entities.HeroOnca {
		moveSpeed = 2.4 // Onça tem reflexos e velocidade felina ligeiramente superiores
	}
	moveForward := ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) || getVirtualKey("ArrowRight")
	moveBackward := ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) || getVirtualKey("ArrowLeft")

	duckKey := ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) || getVirtualKey("ArrowDown")
	duckJustPressed := inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS)

	jumpJustPressed := inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) ||
		inpututil.IsKeyJustPressed(ebiten.KeyW) ||
		getVirtualKey("JustJump")
	if getVirtualKey("JustJump") {
		resetVirtualKey("JustJump")
	}

	jumpHolding := ebiten.IsKeyPressed(ebiten.KeyArrowUp) ||
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

	aimUp := (ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) || getVirtualKey("ArrowUp")) && !e.player.IsPlayerJumping()
	e.player.SetAimUp(aimUp)

	attackJustPressed := inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyX) ||
		inpututil.IsKeyJustPressed(ebiten.KeyJ) ||
		getVirtualKey("Attack")
	if getVirtualKey("Attack") {
		resetVirtualKey("Attack")
	}
	if attackJustPressed {
		e.player.TriggerAttack()
		origX, origY, vx, vy := e.player.GetShootOrigin(GroundY)
		if e.selectedHero == entities.HeroOnca {
			e.projectiles.ShootRoar(origX, origY, vx, vy)
			e.audio.PlayRoar(false)
		} else {
			e.projectiles.ShootAcai(origX, origY, vx, vy)
			e.audio.PlayShot()
		}
	}

	if moveForward {
		e.player.MoveForward(moveSpeed)
	}
	if moveBackward {
		e.player.MoveBackward(moveSpeed)
	}
	if !moveForward && !moveBackward {
		e.player.StopRunning()
	}

	if duckJustPressed && !e.player.IsPlayerJumping() {
		e.audio.PlayDuck()
	}

	if duckKey {
		if e.player.IsPlayerJumping() {
			e.player.FastDrop()
		} else {
			e.player.SetCrouch(true)
		}
	} else {
		e.player.SetCrouch(false)
	}

	if jumpJustPressed {
		wasSwinging := e.player.IsPlayerSwinging()
		jumped, isDouble := e.player.Jump()
		if jumped {
			if wasSwinging {
				if e.selectedHero == entities.HeroOnca {
					e.audio.PlayRoar(true)
				} else {
					e.audio.PlayJungleYell()
				}
			} else {
				e.audio.PlayRoar(isDouble)
			}
		}
	}

	if !jumpHolding && e.player.GetJumpHolding() {
		e.player.ReleaseJump()
	}

	e.player.Update()
	e.projectiles.Update(ScreenWidth)

	playerX, playerY, playerW, playerH := e.player.GetBounds(GroundY)

	// Colisão de Projéteis (Sementes da Baladeira / Rugido Sônico) contra Inimigos e Obstáculos
	for _, proj := range e.projectiles.Projectiles {
		if !proj.Active {
			continue
		}
		px, py, pw, ph := proj.GetBounds()
		hit, hitX, hitY, obsType := e.obstacles.CheckProjectileHit(px, py, pw, ph)
		if hit {
			proj.Active = false
			e.audio.PlayDefeat()

			scoreBonus := 50
			burstColor := color.RGBA{R: 215, G: 50, B: 55, A: 255}
			switch obsType {
			case entities.TypeAir:
				scoreBonus = 100
				burstColor = color.RGBA{R: 240, G: 240, B: 250, A: 255} // Penas
			case entities.TypeJacare:
				scoreBonus = 120
				burstColor = color.RGBA{R: 45, G: 160, B: 55, A: 255} // Escamas
			case entities.TypeSnake:
				scoreBonus = 90
				burstColor = color.RGBA{R: 245, G: 180, B: 35, A: 255} // Coral
			default:
				scoreBonus = 50
				burstColor = color.RGBA{R: 120, G: 45, B: 140, A: 255} // Açaí
			}

			e.score += scoreBonus
			e.projectiles.SpawnHitBurst(hitX, hitY, burstColor, 10)
			e.projectiles.AddScorePopup(hitX-8, hitY-10, scoreBonus)
		}
	}

	stageTargetDist := 1200.0
	stageProgress := e.stageDistance / stageTargetDist
	if stageProgress > 1.0 {
		stageProgress = 1.0
	}
	currentSpeed := (BaseSpeed + float64(e.stage-1)*0.35 + (stageProgress * 0.6)) * SpeedMultipliers[e.speedIndex]
	e.stageDistance += currentSpeed * 0.45
	e.scenery.Update(currentSpeed)
	e.relics.Update(currentSpeed)

	// Conclusão de fase baseada EXCLUSIVAMENTE em distância percorrida da corrida
	if e.stageDistance >= stageTargetDist && !e.isStageComplete {
		e.isStageComplete = true
		e.stageDistance = 0
		e.audio.PauseBGM()
		e.audio.PlayStageUp()
		return nil
	}

	// Coleta de Relíquias e Tesouros Amazônicos (soma pontos e tesouros, SEM mudar de fase!)
	if collected, r := e.relics.CheckCollection(playerX, playerY, playerW, playerH); collected {
		e.relicsCount++
		e.score += r.Value
		e.audio.PlayTreasure()
		e.projectiles.SpawnHitBurst(r.X+8, r.Y+8, color.RGBA{R: 255, G: 220, B: 50, A: 255}, 12)
		e.projectiles.AddScorePopup(r.X, r.Y-8, r.Value)
	}

	if passed := e.obstacles.Update(currentSpeed); passed > 0 {
		e.score += passed * 25
	}
	if e.ticks%6 == 0 {
		e.score += 1
	}

	hit, hitType := e.obstacles.CheckCollision(playerX, playerY, playerW, playerH)
	if hit && e.invincibleTicks <= 0 {
		e.mudSinkTimer = 0
		e.hearts--
		e.shakeTimer = 14
			if e.hearts <= 0 {
				e.lives--
				if e.lives <= 0 {
					e.lives = 0
					e.hearts = 0
					e.isGameOver = true
					if e.selectedHero == entities.HeroOnca {
						e.speechBubbleText = "Arrgh! A floresta me chama..."
					} else {
						e.speechBubbleText = "Levei o farelo mano, mancada!"
					}
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
					if e.selectedHero == entities.HeroOnca {
						e.speechBubbleText = "EGUA DO JACARE FOFOQUEIRO!"
					} else {
						e.speechBubbleText = "EGUA DO JACARE!..."
					}
				} else if hitType == entities.TypeSnake {
					if e.selectedHero == entities.HeroOnca {
						e.speechBubbleText = "SAI PRA LA, COBRA TRAIDORA!"
					} else {
						e.speechBubbleText = "VALHA-ME! UMA COBRA!"
					}
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
	heroName := "GAROTO"
	heroFull := "GAROTO CURUMIM"
	if e.selectedHero == entities.HeroOnca {
		heroName = "ONCA"
		heroFull = "ONCA PINTADA"
	}

	if e.isShowingCredits {
		ui.DrawCreditsScreen(screen, ScreenWidth, ScreenHeight)
		return
	}

	if e.isCharSelect {
		ui.DrawCharacterSelectScreen(screen, ScreenWidth, ScreenHeight, e.ticks, e.selectedHero)
		return
	}

	if e.isSaoBrasIntro {
		ui.DrawTitleIntro(screen, ScreenWidth, ScreenHeight, e.ticks, e.audio.IsIntroPlaying(), e.introMenuIndex, e.audio.IsMuted(), SpeedLabels[e.speedIndex], heroFull)
		return
	}

	e.scenery.Draw(screen, ScreenWidth, GroundY, e.ticks, e.stage)
	e.vines.Draw(screen)
	e.relics.Draw(screen, e.ticks)

	if e.isTitleScreen {
		e.player.Draw(screen, GroundY, e.ticks, 0)
		ui.DrawCityFooter(screen, ScreenWidth, ScreenHeight, 1, e.ticks)
		if e.isShowingCredits {
			ui.DrawCreditsScreen(screen, ScreenWidth, ScreenHeight)
		} else {
			ui.DrawTitleScreen(screen, ScreenWidth, ScreenHeight, e.ticks)
		}
		return
	}

	e.player.Draw(screen, GroundY, e.ticks, e.invincibleTicks)

	if e.speechBubbleTimer > 0 && e.speechBubbleText != "" {
		pX, pY, _, _ := e.player.GetBounds(GroundY)
		bubbleX := pX - 10
		if bubbleX < 10 {
			bubbleX = 10
		}
		if bubbleX+float64(len(e.speechBubbleText)*6) > ScreenWidth-10 {
			bubbleX = pX - 25
		}
		ui.DrawSpeechBubble(screen, bubbleX, pY-24, e.speechBubbleText)
	}

	e.obstacles.Draw(screen, e.ticks, e.stage)
	e.projectiles.Draw(screen, e.ticks)

	isDoubleJump := e.player.GetJumpCount() == 2
	ui.DrawHUD(screen, e.lives, e.hearts, e.score, e.stage, isDoubleJump, e.stageBannerTimer, e.audio.IsMuted(), e.ticks, e.relicsCount, heroName, e.stageDistance)

	if e.isShowingCredits {
		ui.DrawCreditsScreen(screen, ScreenWidth, ScreenHeight)
		return
	}

	if e.isPaused {
		ui.DrawPauseMenu(screen, ScreenWidth, ScreenHeight, e.pauseMenuIndex, e.audio.IsMuted(), SpeedLabels[e.speedIndex])
	}

	if e.isStageComplete {
		ui.DrawStageCompleteScreen(screen, ScreenWidth, ScreenHeight, e.stage, e.score, e.lives, e.ticks)
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
	ebiten.SetWindowTitle("Pai D'Égua Game: Aventura Amazônica")

	engine := NewEngine()
	return ebiten.RunGame(engine)
}
