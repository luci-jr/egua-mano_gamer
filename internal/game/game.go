package game

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/luci-jr/egua-mano_gamer/internal/audio"
	"github.com/luci-jr/egua-mano_gamer/internal/entities"
	"github.com/luci-jr/egua-mano_gamer/internal/scenery"
	"github.com/luci-jr/egua-mano_gamer/internal/ui"
)

const (
	ScreenWidth         = 340.0
	ScreenHeight        = 210.0
	GroundY             = 155.0
	BaseSpeed           = 2.3
	StageTargetDistance = 2800.0
	GameVersion         = "VER. 2.4.0"

	CharSelectSourceTitle    = 0
	CharSelectSourcePause    = 1
	CharSelectSourceGameOver = 2
)

var (
	SpeedMultipliers = []float64{0.65, 0.85, 1.15}
	SpeedLabels      = []string{"LENTO", "NORMAL", "RAPIDO"}
)

type Engine struct {
	player            entities.PlayerCharacter
	garoto            *entities.Garoto
	onca              *entities.Onca
	selectedHero      int
	isCharSelect      bool
	charSelectSource  int
	projectiles       *entities.ProjectileManager
	obstacles         *entities.ObstacleManager
	currentPlatform   *entities.Obstacle
	vines             *entities.VineManager
	relics            *entities.RelicManager
	scenery           *scenery.Background
	audio             *audio.Manager

	score             int
	relicsCount       int
	ticks             int
	lives             int
	hearts            int
	attackCooldown    int
	chargeTimer       int
	isCharged         bool
	invincibleTicks   int
	starPowerTimer    int
	shakeTimer        int
	hitDelayTimer     int
	mudSinkTimer      int
	speechBubbleTimer int
	speechBubbleText  string
	heatSpeechTimer   int

	// Animação dramática/cômica de queda no rio / baía com splash
	waterFallActive          bool
	waterFallTimer           int
	waterFallHeroX           float64
	waterFallHeroY           float64
	waterFallHeroVX          float64
	waterFallHeroVY          float64
	waterFallSplashTriggered bool
	waterFallSplashTimer     int

	stage             int
	stageDistance     float64
	stageBannerTimer  int
	isLoadingStage    bool
	loadingTimer      int
	targetStage       int
	stageFadeTimer    int
	isTitleCover      bool
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
	wasFocused        bool
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
		charSelectSource:  CharSelectSourceTitle,
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
		attackCooldown:    0,
		chargeTimer:       0,
		isCharged:         false,
		invincibleTicks:   0,
		starPowerTimer:    0,
		shakeTimer:        0,
		hitDelayTimer:     0,
		mudSinkTimer:      0,
		speechBubbleTimer: 0,
		speechBubbleText:  "",
		heatSpeechTimer:   0,
		stage:             1,
		stageBannerTimer:  120,
		isTitleCover:      true,
		isSaoBrasIntro:    false,
		saoBrasTimer:      0,
		introMenuIndex:    0,
		speedIndex:        1,
		isTitleScreen:     false,
		isShowingCredits:  false,
		isPaused:          false,
		pauseMenuIndex:    0,
		isGameOver:        false,
		isStageComplete:   false,
		wasFocused:        true,
	}
}

func (e *Engine) applySelectedHero() {
	if e.selectedHero == entities.HeroGaroto {
		e.player = e.garoto
	} else {
		e.player = e.onca
	}
}

func (e *Engine) toggleSelectedHero() {
	if e.selectedHero == entities.HeroGaroto {
		e.selectedHero = entities.HeroOnca
	} else {
		e.selectedHero = entities.HeroGaroto
	}
	e.applySelectedHero()
}

func (e *Engine) restartStage() {
	e.applySelectedHero()
	e.player.Reset()
	e.currentPlatform = nil
	e.attackCooldown = 0
	e.chargeTimer = 0
	e.isCharged = false
	e.projectiles.Reset()
	e.obstacles.Reset()
	e.vines.Reset()
	e.relics.Reset()
	e.stageDistance = 0
	e.mudSinkTimer = 0
	e.lives = 3
	e.hearts = 3
	e.stageBannerTimer = 120
	e.stageFadeTimer = 20
	e.invincibleTicks = 60
	e.starPowerTimer = 0
	e.shakeTimer = 0
	e.speechBubbleTimer = 0
	e.speechBubbleText = ""
	e.waterFallActive = false
	e.isPaused = false
	e.isGameOver = false
	e.isStageComplete = false
	e.audio.RestartBGM()
}

func (e *Engine) resetGame() {
	e.stage = 1
	e.score = 0
	e.relicsCount = 0
	e.restartStage()
}

func (e *Engine) returnToTitle() {
	e.isPaused = false
	e.isGameOver = false
	e.isStageComplete = false
	e.isTitleScreen = false
	e.isCharSelect = false
	e.isSaoBrasIntro = false
	e.isTitleCover = true
	e.introMenuIndex = 0
	e.pauseMenuIndex = 0
	e.audio.PauseBGM()
	if !e.audio.IsMuted() {
		e.audio.PlayIntroBGM()
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
	// Se a tela não estiver em gameplay ativa, consome evento de unfocus pendente
	if e.isShowingCredits || e.isTitleCover || e.isSaoBrasIntro || e.isCharSelect || e.isTitleScreen || e.isPaused || e.isStageComplete || e.isGameOver {
		resetVirtualKey("JustUnfocused")
	}

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

	// 1.5. Tela de Abertura Limpa Arcade (Pitfall / Super Metroid / Contra) - Urubus e Garças Voando
	if e.isTitleCover {
		e.ticks++
		if !e.audio.IsMuted() {
			e.audio.PlayIntroBGM()
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			e.isShowingCredits = true
			return nil
		}

		startTriggered := isPointerJustPressed() ||
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

		if startTriggered {
			e.isTitleCover = false
			e.isSaoBrasIntro = true
			e.audio.PlayShot()
			return nil
		}
		return nil
	}

	// 1.8. Tela de Carregamento Náutica / Viagem Cultural por Belém (Transição entre Fases)
	if e.isLoadingStage {
		e.ticks++
		e.loadingTimer--

		skipLoading := isPointerJustPressed() ||
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

		if skipLoading || e.loadingTimer <= 0 {
			e.isLoadingStage = false
			e.stage = e.targetStage
			e.stageDistance = 0
			e.stageBannerTimer = 130
			e.stageFadeTimer = 25
			e.hearts = 3
			e.currentPlatform = nil
			e.attackCooldown = 0
			e.chargeTimer = 0
			e.isCharged = false
			e.projectiles.Reset()
			e.obstacles.Reset()
			e.relics.Reset()
			if !e.audio.IsMuted() {
				e.audio.ResumeBGM()
			}
			return nil
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

			if e.charSelectSource == CharSelectSourcePause {
				e.applySelectedHero()
				e.isPaused = false
				if !e.audio.IsMuted() {
					e.audio.ResumeBGM()
				}
				return nil
			}

			if e.charSelectSource == CharSelectSourceGameOver {
				e.applySelectedHero()
				e.isGameOver = true
				return nil
			}

			// Veio da abertura/título (CharSelectSourceTitle): Inicia a partida!
			e.isSaoBrasIntro = false
			e.isTitleScreen = false
			e.currentPlatform = nil
			e.stage = 1
			e.stageDistance = 0
			e.stageBannerTimer = 120
			e.stageFadeTimer = 25
			e.audio.StopIntroBGM()
			e.audio.RestartBGM()
			return nil
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || getVirtualKey("Escape") {
			resetVirtualKey("Escape")
			e.isCharSelect = false
			if e.charSelectSource == CharSelectSourcePause {
				e.isPaused = true
				return nil
			}
			if e.charSelectSource == CharSelectSourceGameOver {
				e.isGameOver = true
				return nil
			}
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

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || getVirtualKey("Escape") {
			resetVirtualKey("Escape")
			e.isSaoBrasIntro = false
			e.isTitleCover = true
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

		// Ajuste com Esquerda / Direita (apenas para SOM e VELOCIDADE)
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
			if e.introMenuIndex == 2 {
				e.audio.ToggleMute()
			} else if e.introMenuIndex == 3 {
				e.speedIndex = (e.speedIndex - 1 + len(SpeedMultipliers)) % len(SpeedMultipliers)
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
			if e.introMenuIndex == 2 {
				e.audio.ToggleMute()
			} else if e.introMenuIndex == 3 {
				e.speedIndex = (e.speedIndex + 1) % len(SpeedMultipliers)
			}
		}

		// Detecção de clique / toque no Menu
		mouseTriggered := false
		boxW := 226.0
		boxH := 88.0
		boxX := (ScreenWidth - boxW) / 2.0
		boxY := 62.0
		startY := boxY + 20.0

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			fmx, fmy := float64(mx), float64(my)
			if fmx >= boxX && fmx <= boxX+boxW && fmy >= startY && fmy <= startY+5*13.0 {
				clickedIdx := int((fmy - startY) / 13.0)
				if clickedIdx >= 0 && clickedIdx < 5 {
					e.introMenuIndex = clickedIdx
					mouseTriggered = true
				}
			} else if fmy >= boxY && fmy <= boxY+boxH {
				mouseTriggered = true
			}
		}

		menuTouches := inpututil.AppendJustPressedTouchIDs(nil)
		for _, id := range menuTouches {
			tx, ty := ebiten.TouchPosition(id)
			ftx, fty := float64(tx), float64(ty)
			if ftx >= boxX && ftx <= boxX+boxW && fty >= startY && fty <= startY+5*13.0 {
				clickedIdx := int((fty - startY) / 13.0)
				if clickedIdx >= 0 && clickedIdx < 5 {
					e.introMenuIndex = clickedIdx
					mouseTriggered = true
				}
			} else if fty >= boxY && fty <= boxY+boxH {
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
				e.charSelectSource = CharSelectSourceTitle
				e.isCharSelect = true
				e.isSaoBrasIntro = false
				return nil
			case 1:
				e.charSelectSource = CharSelectSourceTitle
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
			e.pauseMenuIndex = (e.pauseMenuIndex - 1 + 7) % 7
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) || getVirtualKey("ArrowDown") {
			resetVirtualKey("ArrowDown")
			e.pauseMenuIndex = (e.pauseMenuIndex + 1) % 7
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) || getVirtualKey("ArrowLeft") {
			resetVirtualKey("ArrowLeft")
			if e.pauseMenuIndex == 4 {
				e.audio.ToggleMute()
			} else if e.pauseMenuIndex == 5 {
				e.speedIndex = (e.speedIndex - 1 + len(SpeedMultipliers)) % len(SpeedMultipliers)
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) || getVirtualKey("ArrowRight") {
			resetVirtualKey("ArrowRight")
			if e.pauseMenuIndex == 4 {
				e.audio.ToggleMute()
			} else if e.pauseMenuIndex == 5 {
				e.speedIndex = (e.speedIndex + 1) % len(SpeedMultipliers)
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || getVirtualKey("Escape") {
			resetVirtualKey("Escape")
			e.applySelectedHero()
			e.isPaused = false
			if !e.audio.IsMuted() {
				e.audio.ResumeBGM()
			}
			return nil
		}

		selectPressed := isPointerJustPressed() ||
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

		if selectPressed {
			switch e.pauseMenuIndex {
			case 0: // CONTINUAR
				e.applySelectedHero()
				e.isPaused = false
				if !e.audio.IsMuted() {
					e.audio.ResumeBGM()
				}
			case 1: // SELECIONAR JOGADOR
				e.charSelectSource = CharSelectSourcePause
				e.isCharSelect = true
				e.isPaused = false
				return nil
			case 2: // REINICIAR FASE ATUAL
				e.restartStage()
			case 3: // RESETAR JOGO (DO ZERO)
				e.resetGame()
			case 4: // SOM
				e.audio.ToggleMute()
			case 5: // VELOCIDADE
				e.speedIndex = (e.speedIndex + 1) % len(SpeedMultipliers)
			case 6: // MENU INICIAL
				e.returnToTitle()
			}
		}
		return nil
	}

	if e.isStageComplete {
		exitPressed := inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) || getVirtualKey("Escape")
		if exitPressed {
			resetVirtualKey("Escape")
			e.returnToTitle()
			return nil
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			e.resetGame()
			return nil
		}

		continuePressed := isPointerJustPressed() ||
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

		if continuePressed {
			if e.stage < 3 {
				e.isStageComplete = false
				e.isLoadingStage = true
				e.targetStage = e.stage + 1
				e.loadingTimer = 180 // ~3 segundos navegando no barco Popopó pela Baía do Guajará
				return nil
			} else {
				e.resetGame()
			}
			return nil
		}
		return nil
	}

	if e.isGameOver {
		exitPressed := inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) || getVirtualKey("Escape")
		if exitPressed {
			resetVirtualKey("Escape")
			e.returnToTitle()
			return nil
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyJ) || getVirtualKey("KeyJ") {
			resetVirtualKey("KeyJ")
			e.charSelectSource = CharSelectSourceGameOver
			e.isCharSelect = true
			e.isGameOver = false
			return nil
		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			_, my := ebiten.CursorPosition()
			if float64(my) >= 46.0 && float64(my) <= 66.0 {
				e.charSelectSource = CharSelectSourceGameOver
				e.isCharSelect = true
				e.isGameOver = false
				return nil
			}
		}
		for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
			_, ty := ebiten.TouchPosition(id)
			if ty >= 46 && ty <= 66 {
				e.charSelectSource = CharSelectSourceGameOver
				e.isCharSelect = true
				e.isGameOver = false
				return nil
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyN) || getVirtualKey("KeyN") {
			resetVirtualKey("KeyN")
			e.resetGame()
			return nil
		}

		restartPressed := isPointerJustPressed() ||
			inpututil.IsKeyJustPressed(ebiten.KeyR) ||
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

		if restartPressed {
			e.restartStage()
			return nil
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

	// Auto-pausa inteligente e silenciamento se a tela perder foco ou for para segundo plano (ex: alternar para WhatsApp / outro app)
	currentFocused := ebiten.IsFocused()
	focusLost := !currentFocused && e.wasFocused
	e.wasFocused = currentFocused

	if focusLost || !currentFocused || getVirtualKey("JustUnfocused") || getVirtualKey("PageHidden") {
		resetVirtualKey("JustUnfocused")
		e.isPaused = true
		e.pauseMenuIndex = 0
		e.audio.PauseBGM()
		return nil
	}

	e.ticks++

	// Se a animação de queda no rio / baía com splash estiver ativa, processa a cinemática de perda de vida
	if e.waterFallActive {
		e.waterFallTimer--
		if !e.waterFallSplashTriggered {
			e.waterFallHeroX += e.waterFallHeroVX
			e.waterFallHeroY += e.waterFallHeroVY
			e.waterFallHeroVY += 0.32 // Gravidade acelerando a queda
			if e.waterFallHeroY >= GroundY+6.0 {
				e.waterFallSplashTriggered = true
				e.waterFallSplashTimer = 45
				e.audio.PlaySplash()
				e.shakeTimer = 10
				e.speechBubbleText = "Egua do pitiu. Essa agua ta podre!"
				e.speechBubbleTimer = 85
			}
		} else {
			if e.waterFallSplashTimer > 0 {
				e.waterFallSplashTimer--
			}
		}

		if e.waterFallTimer <= 0 {
			e.waterFallActive = false
			if e.lives <= 0 {
				e.lives = 0
				e.hearts = 0
				e.isGameOver = true
				e.speechBubbleText = "Egua do pitiu. Essa agua ta podre!"
				e.speechBubbleTimer = 999999
				e.audio.PauseBGM()
				e.audio.PlayGameOver()
			} else {
				e.hearts = 3 // Restaura os 3 corações para a próxima vida
				e.player.Reset()
				e.player.SetPositionX(45.0)
				e.player.SetGroundOffset(0)
				e.invincibleTicks = 100 // Proteção temporária após respawn
				e.speechBubbleText = fmt.Sprintf("AGORA VAI! RESTAM %d VIDAS", e.lives)
				e.speechBubbleTimer = 75
			}
		}
		return nil
	}

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

	if e.stageFadeTimer > 0 {
		e.stageFadeTimer--
	}

	if e.speechBubbleTimer > 0 {
		e.speechBubbleTimer--
	} else {
		// A cada ~14 a 20 segundos correndo pela cidade de Belém, o herói solta um brado de aventura regional
		e.heatSpeechTimer++
		if e.heatSpeechTimer >= 850 {
			if e.hitDelayTimer == 0 {
				e.heatSpeechTimer = 0
				if e.selectedHero == entities.HeroOnca {
					oncaPhrases := []string{
						"RRRAUW! Ninguem segura a onca!",
						"A cidade e as matas me pertencem!",
						"Sentiram o poder do rugido?!",
						"Bote certeiro, mano!",
					}
					e.speechBubbleText = oncaPhrases[(e.ticks/60)%len(oncaPhrases)]
				} else {
					garotoPhrases := []string{
						"Bora voando, maninho!",
						"O carimbo de Belem ta paidegua!",
						"Vou passar o rodo nesses bichos!",
						"Acai na tigela da forca!",
					}
					e.speechBubbleText = garotoPhrases[(e.ticks/60)%len(garotoPhrases)]
				}
				e.speechBubbleTimer = 165
			}
		}
	}

	if e.invincibleTicks > 0 {
		e.invincibleTicks--
	}

	if e.starPowerTimer > 0 {
		e.starPowerTimer--
		// Spawna faíscas estelares coloridas atrás do herói em movimento estilo Starman
		if e.ticks%3 == 0 {
			starColors := []color.RGBA{
				{R: 255, G: 235, B: 60, A: 255},
				{R: 60, G: 240, B: 255, A: 255},
				{R: 255, G: 90, B: 220, A: 255},
				{R: 255, G: 255, B: 255, A: 255},
			}
			sCol := starColors[(e.ticks/3)%len(starColors)]
			pBoundsX, pBoundsY, pBoundsW, pBoundsH := e.player.GetBounds(GroundY)
			e.projectiles.SpawnHitBurst(pBoundsX+pBoundsW/2.0, pBoundsY+pBoundsH/2.0, sCol, 2)
		}
	}
	e.player.SetStarPower(e.starPowerTimer)

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
	if e.starPowerTimer > 0 {
		moveSpeed *= 1.22 // Boost de agilidade e velocidade com o Guaraná Power
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

	if e.attackCooldown > 0 {
		e.attackCooldown--
	}

	attackHeld := ebiten.IsKeyPressed(ebiten.KeySpace) ||
		ebiten.IsKeyPressed(ebiten.KeyX) ||
		ebiten.IsKeyPressed(ebiten.KeyJ) ||
		getVirtualKey("Attack")

	if attackHeld {
		e.chargeTimer++
		// Efeito visual de carregamento de energia (partículas de açaí roxo escuro para o Garoto, ondas translúcidas para a Onça; zero bolinhas amarelas)
		if e.chargeTimer >= 18 {
			pX, pY, pW, pH := e.player.GetBounds(GroundY)
			var sparkColor color.RGBA
			if e.selectedHero == entities.HeroOnca {
				sparkColor = color.RGBA{R: 240, G: 110, B: 40, A: 160} // Onda sônica âmbar/ar translúcida
				if e.chargeTimer >= 42 {
					sparkColor = color.RGBA{R: 255, G: 255, B: 255, A: 200} // Halo branco de ar
				}
			} else {
				sparkColor = color.RGBA{R: 25, G: 6, B: 32, A: 200} // Açaí roxo bem escuro
				if e.chargeTimer >= 42 {
					sparkColor = color.RGBA{R: 42, G: 10, B: 52, A: 220} // Açaí encorpado
				}
			}
			if e.chargeTimer%2 == 0 {
				e.projectiles.Particles = append(e.projectiles.Particles, &entities.Particle{
					X:     pX + pW/2.0 + float64(e.ticks%9-4)*2.0,
					Y:     pY + pH/2.0 + float64(e.ticks%7-3)*2.0,
					VX:    float64(e.ticks%5-2) * 0.4,
					VY:    -1.2 - float64(e.ticks%3)*0.4,
					Life:  0,
					Max:   12,
					Size:  2.5,
					Color: sparkColor,
				})
			}
			if e.chargeTimer == 42 {
				// Halo anunciando carga máxima pronta (sem amarelo)
				burstCol := color.RGBA{R: 255, G: 255, B: 255, A: 200}
				if e.selectedHero == entities.HeroGaroto {
					burstCol = color.RGBA{R: 35, G: 8, B: 44, A: 220}
				}
				e.projectiles.SpawnHitBurst(pX+pW/2.0, pY+pH/2.0, burstCol, 10)
			}
		}
		e.isCharged = (e.chargeTimer >= 42)
	} else if e.chargeTimer > 0 {
		// Soltou o botão de ataque!
		if e.chargeTimer >= 42 {
			// ★ DISPARO DO TIRO ESPECIAL CARREGADO!
			e.player.TriggerAttack()
			origX, origY, vx, vy := e.player.GetShootOrigin(GroundY)
			if e.selectedHero == entities.HeroOnca {
				e.projectiles.ShootSpecialRoar(origX, origY, vx, vy)
				e.audio.PlayRoar(true)
			} else {
				e.projectiles.ShootSpecialAcai(origX, origY, vx, vy)
				e.audio.PlayShot()
			}
			e.attackCooldown = 26
		} else {
			// Toque rápido: disparo de tiro comum com cooldown cadenciado anti-spam
			if e.attackCooldown == 0 {
				e.player.TriggerAttack()
				origX, origY, vx, vy := e.player.GetShootOrigin(GroundY)
				if e.selectedHero == entities.HeroOnca {
					e.projectiles.ShootRoar(origX, origY, vx, vy)
					e.audio.PlayRoar(false)
				} else {
					e.projectiles.ShootAcai(origX, origY, vx, vy)
					e.audio.PlayShot()
				}
				e.attackCooldown = 18 // Cadência equilibrada de disparo (~0.3s)
			}
		}
		e.chargeTimer = 0
		e.isCharged = false
		if getVirtualKey("Attack") {
			resetVirtualKey("Attack")
		}
	}

	// Câmera dinâmica de plataforma estilo Pitfall:
	// O herói se desloca livremente pela tela.
	// Ao correr à frente ultrapassando o limiar (135px), a câmera avança com ele (scroll progressivo).
	// Ao recuar aquém do limiar esquerdo (45px) e havendo distância percorrida, o mundo retrocede suavemente.
	// Quando parado, o mundo estabiliza/pausa (scroll zero).
	cameraForwardLimit := 135.0
	cameraBackLimit := 45.0
	worldScrollSpeed := 0.0

	curX, _ := e.player.GetPosition()

	if moveForward {
		e.player.MoveForward(moveSpeed)
		curX, _ = e.player.GetPosition()
		if curX >= cameraForwardLimit {
			excess := curX - cameraForwardLimit
			e.player.SetPositionX(cameraForwardLimit)
			worldScrollSpeed = (moveSpeed + excess) * SpeedMultipliers[e.speedIndex]
		}
	} else if moveBackward {
		e.player.MoveBackward(moveSpeed)
		curX, _ = e.player.GetPosition()
		if curX <= cameraBackLimit && e.stageDistance > 0 {
			excess := cameraBackLimit - curX
			e.player.SetPositionX(cameraBackLimit)
			worldScrollSpeed = -(moveSpeed*0.75 + excess) * SpeedMultipliers[e.speedIndex]
		}
	} else {
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
		e.currentPlatform = nil // Desprende da plataforma sólida para novo salto
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

	// Lógica de Plataformas Sólidas estilo Pitfall (Pousar e subir no Paneiro de Açaí, Dorso do Jacaré ou Banco)
	if e.currentPlatform != nil {
		if !e.currentPlatform.IsPlayerOnTop(playerX, playerW) {
			// O jogador caminhou para fora do obstáculo ou o obstáculo se moveu: queda suave
			e.currentPlatform = nil
			e.player.FallFromPlatform()
		}
	}
	if e.currentPlatform == nil {
		supported, platformOffset, obs := e.obstacles.CheckPlatformSupport(playerX, playerY, playerW, playerH, e.player.GetVelocityY(), GroundY)
		if supported {
			e.currentPlatform = obs
			e.player.SetGroundOffset(platformOffset)
			e.audio.PlayDuck()
		}
	}

	// Colisão de Projéteis (Sementes da Baladeira / Rugido Sônico) contra Inimigos e Obstáculos
	for _, proj := range e.projectiles.Projectiles {
		if !proj.Active {
			continue
		}
		px, py, pw, ph := proj.GetBounds()
		hit, hitX, hitY, obsType := e.obstacles.CheckProjectileHit(px, py, pw, ph)
		if hit {
			proj.HitsLeft--
			if proj.HitsLeft <= 0 {
				proj.Active = false
			}
			e.audio.PlayDefeat()

			scoreBonus := 50
			burstColor := color.RGBA{R: 215, G: 50, B: 55, A: 255}
			switch obsType {
			case entities.TypeAir:
				scoreBonus = 50
				burstColor = color.RGBA{R: 240, G: 240, B: 250, A: 255} // Penas
			case entities.TypeJacare:
				scoreBonus = 60
				burstColor = color.RGBA{R: 45, G: 160, B: 55, A: 255} // Escamas
			case entities.TypeSnake:
				scoreBonus = 40
				burstColor = color.RGBA{R: 245, G: 180, B: 35, A: 255} // Coral
			default:
				scoreBonus = 20
				burstColor = color.RGBA{R: 120, G: 45, B: 140, A: 255} // Açaí
			}

			if proj.IsSpecial {
				scoreBonus += 30
			}

			e.score += scoreBonus
			e.projectiles.SpawnHitBurst(hitX, hitY, burstColor, 12)
			e.projectiles.AddScorePopup(hitX-8, hitY-10, scoreBonus)
		}
	}

	stageTargetDist := StageTargetDistance
	// Atualiza distância e cenário apenas de acordo com a movimentação real do jogador
	if worldScrollSpeed > 0 {
		e.stageDistance += worldScrollSpeed * 0.30
	} else if worldScrollSpeed < 0 {
		e.stageDistance += worldScrollSpeed * 0.30
		if e.stageDistance < 0 {
			e.stageDistance = 0
		}
	}
	e.scenery.Update(worldScrollSpeed)
	e.relics.Update(worldScrollSpeed)

	// Conclusão de fase baseada EXCLUSIVAMENTE em distância percorrida da corrida
	if e.stageDistance >= stageTargetDist && !e.isStageComplete {
		e.isStageComplete = true
		e.stageDistance = 0
		e.currentPlatform = nil
		e.audio.PauseBGM()
		e.audio.PlayStageUp()
		return nil
	}

	// Coleta de Relíquias e Tesouros Amazônicos (incluindo Cuia de Tacacá que recupera 1 coração!)
	// Coleta de Relíquias e Tesouros Amazônicos (Muiraquitã, Urna, Ouro, Açaí e Guaraná Power!)
	if collected, r := e.relics.CheckCollection(playerX, playerY, playerW, playerH); collected {
		e.relicsCount++
		if r.Type == entities.RelicGuarana {
			// 🌟 Fruto do Guaraná da Amazônia: Super Invencibilidade Estilo Mario!
			e.starPowerTimer = 420 // ~7 segundos de invencibilidade estelar
			e.score += r.Value
			e.audio.PlayEguaMano()
			e.projectiles.AddTextPopup(r.X-22, r.Y-20, "GUARANA POWER!", color.RGBA{R: 255, G: 220, B: 50, A: 255})
			e.projectiles.SpawnHitBurst(r.X+8, r.Y+8, color.RGBA{R: 255, G: 215, B: 40, A: 255}, 18)
			if e.selectedHero == entities.HeroOnca {
				e.speechBubbleText = "RROAAR! NINGUEM ME SEGURA!"
			} else {
				e.speechBubbleText = "EGUA MANO! TO INVENCIVEL!"
			}
			e.speechBubbleTimer = 75
		} else if r.Type == entities.RelicAcaiBowl {
			// 🥣 Cuia de Tacacá / Tigela de Açaí: item raro de cura
			if e.hearts < 3 {
				e.hearts++
				e.projectiles.AddTextPopup(r.X-20, r.Y-14, "+1 ENERGIA!", color.RGBA{R: 215, G: 65, B: 245, A: 255})
				e.projectiles.SpawnHitBurst(r.X+8, r.Y+8, color.RGBA{R: 190, G: 45, B: 230, A: 255}, 16)
			} else {
				bonus := 200
				e.score += bonus
				e.projectiles.AddScorePopup(r.X, r.Y-8, bonus)
				e.projectiles.AddTextPopup(r.X-20, r.Y-24, "ACAI POWER!", color.RGBA{R: 215, G: 65, B: 245, A: 255})
				e.projectiles.SpawnHitBurst(r.X+8, r.Y+8, color.RGBA{R: 45, G: 10, B: 58, A: 255}, 14)
			}
			e.audio.PlayTreasure()
		} else {
			e.score += r.Value
			e.audio.PlayTreasure()
			e.projectiles.SpawnHitBurst(r.X+8, r.Y+8, color.RGBA{R: 255, G: 220, B: 50, A: 255}, 12)
			e.projectiles.AddScorePopup(r.X, r.Y-8, r.Value)
		}
	}

	if passed := e.obstacles.Update(worldScrollSpeed); passed > 0 && worldScrollSpeed > 0 {
		e.score += passed * 10
	}
	// Pontuação por avanço ativo: a cada 60 ticks em deslocamento para frente
	if worldScrollSpeed > 0 && e.ticks%60 == 0 {
		e.score += 1
	}

	// Coleta do Paneiro de Açaí no solo: não vem coração toda hora (apenas 15% de chance de cura se ferido)
	if collected, obs := e.obstacles.CheckEnergyCollection(playerX, playerY, playerW, playerH); collected {
		ox, oy, ow, _ := obs.GetBounds()
		centerX := ox + ow/2.0
		centerY := oy + 4.0

		giveHeart := e.hearts < 3 && rand.Intn(100) < 15
		if giveHeart {
			e.hearts++
			e.projectiles.AddTextPopup(centerX-24, centerY-14, "+1 ENERGIA!", color.RGBA{R: 215, G: 65, B: 245, A: 255})
			e.projectiles.SpawnHitBurst(centerX, centerY, color.RGBA{R: 190, G: 45, B: 230, A: 255}, 16)
			e.audio.PlayTreasure()
		} else {
			bonus := 50
			e.score += bonus
			e.projectiles.AddScorePopup(centerX-12, centerY-14, bonus)
			e.projectiles.AddTextPopup(centerX-20, centerY-26, "ACAI PURO!", color.RGBA{R: 180, G: 55, B: 220, A: 255})
			e.projectiles.SpawnHitBurst(centerX, centerY, color.RGBA{R: 45, G: 10, B: 58, A: 255}, 10)
		}
	}

	// 1. Efeito do Guaraná Power: Destrói instantaneamente qualquer inimigo que encostar (Super Mario Starman)
	if e.starPowerTimer > 0 {
		for _, obs := range e.obstacles.Obstacles {
			if obs.Defeated || obs.Collided || obs.Type == entities.TypeBench || obs.Type == entities.TypeGround {
				continue
			}
			if obs.CheckCollision(playerX, playerY, playerW, playerH) {
				obs.Defeated = true
				obs.DefeatTicks = 26
				e.audio.PlayDefeat()
				scoreBonus := 150
				e.score += scoreBonus
				ox, oy, ow, oh := obs.GetBounds()
				e.projectiles.SpawnHitBurst(ox+ow/2.0, oy+oh/2.0, color.RGBA{R: 255, G: 225, B: 55, A: 255}, 16)
				e.projectiles.AddTextPopup(ox-8, oy-14, "SMASH! +150", color.RGBA{R: 255, G: 240, B: 80, A: 255})
			}
		}
	}

	// 2. Mecânica de Pisão na Cabeça (Stomp): Pular em cima do Jacaré, Cobra ou Garça/Ave mata o inimigo e quica no ar!
	stompHit, stompX, stompY, stompType := e.obstacles.CheckStomp(playerX, playerY, playerW, playerH, e.player.GetVelocityY(), GroundY)
	if stompHit {
		e.player.Bounce(-6.2) // Herói quica no ar estilo Mario
		e.audio.PlayDefeat()
		stompBonus := 100
		burstCol := color.RGBA{R: 255, G: 220, B: 60, A: 255}
		switch stompType {
		case entities.TypeAir:
			stompBonus = 80
			burstCol = color.RGBA{R: 245, G: 245, B: 255, A: 255}
		case entities.TypeJacare:
			stompBonus = 120
			burstCol = color.RGBA{R: 45, G: 180, B: 65, A: 255}
		case entities.TypeSnake:
			stompBonus = 90
			burstCol = color.RGBA{R: 245, G: 190, B: 40, A: 255}
		}
		e.score += stompBonus
		e.projectiles.SpawnHitBurst(stompX, stompY, burstCol, 14)
		e.projectiles.AddTextPopup(stompX-16, stompY-16, fmt.Sprintf("PISAO! +%d", stompBonus), color.RGBA{R: 255, G: 235, B: 70, A: 255})
	}

	// 3. Colisão de Dano Normal (se não estiver invencível pelo Guaraná Power nem por dano recente)
	hit, hitType := e.obstacles.CheckCollision(playerX, playerY, playerW, playerH, e.player.GetVelocityY(), GroundY, e.currentPlatform)
	if hit && e.invincibleTicks <= 0 && e.starPowerTimer <= 0 {
		e.currentPlatform = nil
		e.mudSinkTimer = 0
		e.hearts--
		e.shakeTimer = 14
		if e.hearts <= 0 {
			e.lives--
			e.hearts = 0
			e.audio.PlayHit()

			// Dispara a animação dramática de queda no rio / splash na baía do Guajará
			e.waterFallActive = true
			e.waterFallTimer = 75
			e.waterFallHeroX = playerX
			e.waterFallHeroY = playerY
			e.waterFallHeroVX = -1.2
			e.waterFallHeroVY = -3.8
			e.waterFallSplashTriggered = false
			e.waterFallSplashTimer = 0
			e.hitDelayTimer = 0
			e.invincibleTicks = 0
			e.speechBubbleText = "Egua do pitiu. Essa agua ta podre!"
			e.speechBubbleTimer = 90
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

	if e.isTitleCover {
		ui.DrawTitleCoverScreen(screen, ScreenWidth, ScreenHeight, e.ticks, GameVersion)
		return
	}

	if e.isCharSelect {
		ui.DrawCharacterSelectScreen(screen, ScreenWidth, ScreenHeight, e.ticks, e.selectedHero)
		return
	}

	if e.isSaoBrasIntro {
		ui.DrawTitleIntro(screen, ScreenWidth, ScreenHeight, e.ticks, e.audio.IsIntroPlaying(), e.introMenuIndex, e.audio.IsMuted(), SpeedLabels[e.speedIndex], heroFull, GameVersion)
		return
	}

	if e.isLoadingStage {
		progress := 1.0 - float64(e.loadingTimer)/180.0
		ui.DrawStageTransitionLoadingScreen(screen, ScreenWidth, ScreenHeight, e.ticks, e.targetStage, progress)
		return
	}

	e.scenery.Draw(screen, ScreenWidth, GroundY, e.ticks, e.stage)
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

	if e.waterFallActive {
		ui.DrawHeroWaterFall(screen, e.waterFallHeroX, e.waterFallHeroY, e.selectedHero, e.waterFallHeroVY, e.waterFallSplashTriggered)
		if e.waterFallSplashTriggered {
			ui.DrawWaterSplash(screen, e.waterFallHeroX, GroundY+8.0, e.waterFallSplashTimer)
		}
	} else {
		e.player.Draw(screen, GroundY, e.ticks, e.invincibleTicks)
	}

	if e.speechBubbleTimer > 0 && e.speechBubbleText != "" {
		pX, pY := 0.0, 0.0
		if e.waterFallActive {
			pX = e.waterFallHeroX
			pY = e.waterFallHeroY
		} else {
			pX, pY, _, _ = e.player.GetBounds(GroundY)
		}
		bubbleX := pX - 10
		if bubbleX < 10 {
			bubbleX = 10
		}
		if bubbleX+float64(len(e.speechBubbleText)*6) > ScreenWidth-10 {
			bubbleX = ScreenWidth - 10 - float64(len(e.speechBubbleText)*6)
		}
		ui.DrawSpeechBubble(screen, bubbleX, pY-24, e.speechBubbleText)
	}

	e.obstacles.Draw(screen, e.ticks, e.stage)
	e.projectiles.Draw(screen, e.ticks)

	if e.stageFadeTimer > 0 {
		fadeAlpha := uint8(255.0 * float64(e.stageFadeTimer) / 25.0)
		ebitenutil.DrawRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{R: 12, G: 14, B: 20, A: fadeAlpha})
	}

	isDoubleJump := e.player.GetJumpCount() == 2
	ui.DrawHUD(screen, e.lives, e.hearts, e.score, e.stage, isDoubleJump, e.stageBannerTimer, e.audio.IsMuted(), e.ticks, e.relicsCount, heroName, e.stageDistance, e.starPowerTimer)

	if e.isShowingCredits {
		ui.DrawCreditsScreen(screen, ScreenWidth, ScreenHeight)
		return
	}

	if e.isPaused {
		ui.DrawPauseMenu(screen, ScreenWidth, ScreenHeight, e.pauseMenuIndex, e.audio.IsMuted(), SpeedLabels[e.speedIndex], e.selectedHero)
	}

	if e.isStageComplete {
		ui.DrawStageCompleteScreen(screen, ScreenWidth, ScreenHeight, e.stage, e.score, e.lives, e.ticks)
	}

	if e.isGameOver {
		ui.DrawGameOverScreen(screen, ScreenWidth, ScreenHeight, e.score, e.stage, e.selectedHero)
	}
}

func (e *Engine) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(ScreenWidth), int(ScreenHeight)
}

func Start() error {
	ebiten.SetWindowSize(680, 420)
	ebiten.SetWindowTitle("Égua Mano! Gamer - Uma Aventura em Belém do Pará")
	ebiten.SetRunnableOnUnfocused(false)

	engine := NewEngine()
	return ebiten.RunGame(engine)
}
