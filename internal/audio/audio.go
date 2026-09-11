package audio

import (
	"bytes"
	"math"

	ebitenaudio "github.com/hajimehoshi/ebiten/v2/audio"
)

const sampleRate = 44100

type Manager struct {
	ctx           *ebitenaudio.Context
	sndRoar       []byte
	sndDoubleRoar []byte
	sndDuck       []byte
	sndHit        []byte
	sndEguaMano   []byte
	sndGameOver   []byte
	sndStageUp    []byte
	bgmPlayer     *ebitenaudio.Player
	introPlayer   *ebitenaudio.Player
	isMuted       bool
	isIntroActive bool
}

func createRoar(isDouble bool) []byte {
	durationMs := 170
	if isDouble {
		durationMs = 190
	}
	numSamples := sampleRate * durationMs / 1000
	buf := make([]byte, numSamples*4)

	baseFreq := 85.0
	if isDouble {
		baseFreq = 110.0
	}

	phase := 0.0
	amPhase := 0.0

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(numSamples)

		currentFreq := baseFreq*(1.0+0.25*math.Sin(t*math.Pi)) - t*18.0

		phase += 2.0 * math.Pi * currentFreq / float64(sampleRate)
		amPhase += 2.0 * math.Pi * 38.0 / float64(sampleRate)

		tone := math.Sin(phase) + 0.35*math.Sin(phase*2.0)
		throatTremolo := 0.6 + 0.4*math.Sin(amPhase)

		raw := tone * throatTremolo

		attack := math.Min(1.0, float64(i)/float64(sampleRate*0.025))
		decay := math.Exp(-3.2 * t)

		sample := int16(raw * attack * decay * 0.35 * 32767.0)

		idx := i * 4
		buf[idx] = byte(sample)
		buf[idx+1] = byte(sample >> 8)
		buf[idx+2] = byte(sample)
		buf[idx+3] = byte(sample >> 8)
	}
	return buf
}

func createHitPunch() []byte {
	durationMs := 200
	numSamples := sampleRate * durationMs / 1000
	buf := make([]byte, numSamples*4)
	phase := 0.0

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(numSamples)
		freq := 160.0 * math.Exp(-9.0*t) + 40.0

		phase += 2.0 * math.Pi * freq / float64(sampleRate)
		tone := math.Sin(phase)

		envelope := math.Exp(-4.5 * t)
		sample := int16(tone * envelope * 0.35 * 32767.0)

		idx := i * 4
		buf[idx] = byte(sample)
		buf[idx+1] = byte(sample >> 8)
		buf[idx+2] = byte(sample)
		buf[idx+3] = byte(sample >> 8)
	}
	return buf
}

func createEguaManoSfx() []byte {
	durationMs := 360
	numSamples := sampleRate * durationMs / 1000
	buf := make([]byte, numSamples*4)
	phase := 0.0
	vibPhase := 0.0

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(numSamples)

		freq := (380.0 - 230.0*math.Pow(t, 0.65)) + 32.0*math.Sin(vibPhase)
		phase += 2.0 * math.Pi * freq / float64(sampleRate)
		vibPhase += 2.0 * math.Pi * 18.0 / float64(sampleRate)

		tone := math.Sin(phase) + 0.35*math.Sin(phase*2.0)

		attack := math.Min(1.0, float64(i)/float64(sampleRate*0.012))
		decay := math.Exp(-2.8 * t)
		sample := int16(tone * attack * decay * 0.38 * 32767.0)

		idx := i * 4
		buf[idx] = byte(sample)
		buf[idx+1] = byte(sample >> 8)
		buf[idx+2] = byte(sample)
		buf[idx+3] = byte(sample >> 8)
	}
	return buf
}

func createSmoothTone(startFreq, endFreq float64, durationMs int, harmonic float64) []byte {
	numSamples := sampleRate * durationMs / 1000
	buf := make([]byte, numSamples*4)
	phase1 := 0.0
	phase2 := 0.0

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(numSamples)
		freq := startFreq + (endFreq-startFreq)*t

		phase1 += 2.0 * math.Pi * freq / float64(sampleRate)
		phase2 += 2.0 * math.Pi * (freq * 1.5) / float64(sampleRate)

		val1 := math.Sin(phase1)
		val2 := math.Sin(phase2) * harmonic
		mix := (val1 + val2) / (1.0 + harmonic)

		attack := math.Min(1.0, float64(i)/float64(sampleRate*0.012))
		decay := math.Exp(-3.0 * t)
		sample := int16(mix * attack * decay * 0.22 * 32767.0)

		idx := i * 4
		buf[idx] = byte(sample)
		buf[idx+1] = byte(sample >> 8)
		buf[idx+2] = byte(sample)
		buf[idx+3] = byte(sample >> 8)
	}
	return buf
}

func createStageUpJingle() []byte {
	notes := []float64{261.63, 329.63, 392.00, 523.25}
	noteMs := 85
	totalSamples := (sampleRate * noteMs / 1000) * len(notes)
	buf := make([]byte, totalSamples*4)

	sampleIdx := 0
	for _, freq := range notes {
		stepSamples := sampleRate * noteMs / 1000
		phase := 0.0
		for s := 0; s < stepSamples; s++ {
			t := float64(s) / float64(stepSamples)
			phase += 2.0 * math.Pi * freq / float64(sampleRate)

			val := math.Sin(phase)
			attack := math.Min(1.0, float64(s)/float64(sampleRate*0.01))
			decay := math.Exp(-3.2 * t)
			sample := int16(val * attack * decay * 0.25 * 32767.0)

			idx := sampleIdx * 4
			buf[idx] = byte(sample)
			buf[idx+1] = byte(sample >> 8)
			buf[idx+2] = byte(sample)
			buf[idx+3] = byte(sample >> 8)
			sampleIdx++
		}
	}
	return buf
}

func createCarimboBGM() []byte {
	melodyNotes := []float64{
		261.63, 261.63, 329.63, 392.00,
		392.00, 329.63, 293.66, 261.63,
		293.66, 293.66, 329.63, 392.00,
		329.63, 293.66, 261.63, 220.00,
		261.63, 329.63, 392.00, 329.63,
		293.66, 261.63, 220.00, 196.00,
		220.00, 261.63, 293.66, 329.63,
		293.66, 261.63, 220.00, 261.63,
	}

	bassNotes := []float64{
		130.81, 130.81, 98.00, 130.81,
		130.81, 130.81, 98.00, 130.81,
		110.00, 110.00, 82.41, 110.00,
		110.00, 110.00, 82.41, 110.00,
		130.81, 130.81, 98.00, 130.81,
		130.81, 130.81, 98.00, 130.81,
		98.00, 98.00, 73.42, 98.00,
		130.81, 130.81, 98.00, 130.81,
	}

	noteMs := 150
	totalSamples := (sampleRate * noteMs / 1000) * len(melodyNotes)
	buf := make([]byte, totalSamples*4)

	sampleIdx := 0
	melPhase := 0.0
	bassPhase := 0.0

	for step := 0; step < len(melodyNotes); step++ {
		melFreq := melodyNotes[step]
		bassFreq := bassNotes[step]
		stepSamples := sampleRate * noteMs / 1000

		for s := 0; s < stepSamples; s++ {
			t := float64(s) / float64(stepSamples)

			melPhase += 2.0 * math.Pi * melFreq / float64(sampleRate)
			bassPhase += 2.0 * math.Pi * bassFreq / float64(sampleRate)

			melVal := math.Sin(melPhase)
			melAttack := math.Min(1.0, float64(s)/float64(sampleRate*0.012))
			melDecay := math.Exp(-2.8 * t)
			melEnvelope := melAttack * melDecay

			bassVal := math.Sin(bassPhase)
			bassAttack := math.Min(1.0, float64(s)/float64(sampleRate*0.015))
			bassDecay := math.Exp(-2.2 * t)
			bassEnvelope := bassAttack * bassDecay

			mix := (melVal*melEnvelope*0.11 + bassVal*bassEnvelope*0.18)

			sample := int16(mix * 32767.0)

			idx := sampleIdx * 4
			buf[idx] = byte(sample)
			buf[idx+1] = byte(sample >> 8)
			buf[idx+2] = byte(sample)
			buf[idx+3] = byte(sample >> 8)

			sampleIdx++
		}
	}
	return buf
}

func createSaoBrasIntroBGM() []byte {
	// Melodia de apresentação inspirada no Mercado de São Brás e Carimbó/Guitarrada Paraense
	// Frases harmônicas com balanço característico de Belém do Pará
	melodyNotes := []float64{
		// Frase 1: Entrada majestosa de boas-vindas ao Mercado de São Brás (C - G - C - E)
		261.63, 329.63, 392.00, 523.25, 659.25, 587.33, 523.25, 392.00,
		// Frase 2: Guitarrada Paraense nostálgica em Am (Lá Menor)
		440.00, 523.25, 659.25, 880.00, 783.99, 659.25, 523.25, 440.00,
		// Frase 3: Brisa do Ver-o-Rio e Baía do Guajará em F (Fá Maior)
		349.23, 440.00, 523.25, 698.46, 659.25, 587.33, 523.25, 440.00,
		// Frase 4: Cadência rítmica e chamada de Carimbó em G7
		392.00, 493.88, 587.33, 698.46, 783.99, 698.46, 587.33, 493.88,

		// Frase 5: Tema alegre e vibrante da Onça-Pintada
		659.25, 783.99, 659.25, 587.33, 523.25, 392.00, 440.00, 523.25,
		// Frase 6: Balanço sincopado de carimbó tradicional
		587.33, 523.25, 440.00, 523.25, 587.33, 698.46, 659.25, 587.33,
		// Frase 7: Arpejo ascendente de guitarrada paraense
		783.99, 587.33, 493.88, 392.00, 493.88, 587.33, 659.25, 587.33,
		// Frase 8: Resolução harmônica convidativa para o loop contínuo
		523.25, 659.25, 783.99, 1046.50, 783.99, 659.25, 587.33, 523.25,
	}

	bassNotes := []float64{
		// C Maior
		130.81, 130.81, 98.00, 130.81, 130.81, 130.81, 98.00, 130.81,
		// Am
		110.00, 110.00, 82.41, 110.00, 110.00, 110.00, 82.41, 110.00,
		// F Maior
		87.31, 87.31, 130.81, 87.31, 87.31, 87.31, 130.81, 87.31,
		// G Maior
		98.00, 98.00, 146.83, 98.00, 98.00, 98.00, 146.83, 98.00,

		// C Maior
		130.81, 130.81, 98.00, 130.81, 130.81, 130.81, 98.00, 130.81,
		// F Maior
		87.31, 87.31, 130.81, 87.31, 87.31, 87.31, 130.81, 87.31,
		// G Maior
		98.00, 98.00, 146.83, 98.00, 98.00, 98.00, 146.83, 98.00,
		// C Maior cadência
		130.81, 98.00, 130.81, 98.00, 130.81, 130.81, 98.00, 130.81,
	}

	noteMs := 160 // Andamento musical (~160ms por nota, ~10s ciclo completo)
	totalSamples := (sampleRate * noteMs / 1000) * len(melodyNotes)
	buf := make([]byte, totalSamples*4)

	sampleIdx := 0
	melPhase := 0.0
	bassPhase := 0.0

	for step := 0; step < len(melodyNotes); step++ {
		melFreq := melodyNotes[step]
		bassFreq := bassNotes[step]
		stepSamples := sampleRate * noteMs / 1000

		for s := 0; s < stepSamples; s++ {
			t := float64(s) / float64(stepSamples)

			// Leve vibrato acústico de guitarrada
			vibrato := 1.0 + 0.006*math.Sin(2.0*math.Pi*5.0*t)
			melPhase += 2.0 * math.Pi * (melFreq * vibrato) / float64(sampleRate)
			bassPhase += 2.0 * math.Pi * bassFreq / float64(sampleRate)

			// Timbre rico: fundamental + harmônico 2x (oitava) + harmônico 3x com ataque rápido
			melVal := math.Sin(melPhase) + 0.28*math.Sin(melPhase*2.0) + 0.08*math.Sin(melPhase*3.0)
			melAttack := math.Min(1.0, float64(s)/float64(sampleRate*0.012))
			melDecay := math.Exp(-2.5 * t)
			melEnvelope := melAttack * melDecay

			// Baixo acústico encorpado com harmônico de 2ª oitava para alto-falantes de celular
			bassVal := math.Sin(bassPhase) + 0.32*math.Sin(bassPhase*2.0)
			bassAttack := math.Min(1.0, float64(s)/float64(sampleRate*0.018))
			bassDecay := math.Exp(-2.0 * t)
			bassEnvelope := bassAttack * bassDecay

			// Toque sutil de maraca/percussão de carimbó sintetizado
			maraca := math.Sin(float64(s)*0.75) * math.Exp(-35.0*t) * 0.035

			mix := (melVal*melEnvelope*0.13 + bassVal*bassEnvelope*0.16 + maraca)
			sample := int16(mix * 32767.0)

			idx := sampleIdx * 4
			buf[idx] = byte(sample)
			buf[idx+1] = byte(sample >> 8)
			buf[idx+2] = byte(sample)
			buf[idx+3] = byte(sample >> 8)

			sampleIdx++
		}
	}
	return buf
}

func NewManager() *Manager {
	ctx := ebitenaudio.NewContext(sampleRate)

	m := &Manager{
		ctx:           ctx,
		sndRoar:       createRoar(false),
		sndDoubleRoar: createRoar(true),
		sndDuck:       createSmoothTone(260, 160, 70, 0.1),
		sndHit:        createHitPunch(),
		sndEguaMano:   createEguaManoSfx(),
		sndGameOver:   createSmoothTone(320, 95, 450, 0.3),
		sndStageUp:    createStageUpJingle(),
		isMuted:       false,
		isIntroActive: false,
	}

	bgmBytes := createCarimboBGM()
	bgmLoop := ebitenaudio.NewInfiniteLoop(bytes.NewReader(bgmBytes), int64(len(bgmBytes)))
	bgmPlayer, err := ctx.NewPlayer(bgmLoop)
	if err == nil {
		bgmPlayer.SetVolume(0.32)
		m.bgmPlayer = bgmPlayer
	}

	introBytes := createSaoBrasIntroBGM()
	introLoop := ebitenaudio.NewInfiniteLoop(bytes.NewReader(introBytes), int64(len(introBytes)))
	introPlayer, err := ctx.NewPlayer(introLoop)
	if err == nil {
		introPlayer.SetVolume(0.50) // Volume presente e audível em smartphones e notebooks
		m.introPlayer = introPlayer
	}

	return m
}

func (m *Manager) PlayIntroBGM() {
	m.isIntroActive = true
	if m.bgmPlayer != nil && m.bgmPlayer.IsPlaying() {
		m.bgmPlayer.Pause()
	}
	if !m.isMuted && m.introPlayer != nil && !m.introPlayer.IsPlaying() {
		m.introPlayer.Play()
	}
}

func (m *Manager) StopIntroBGM() {
	m.isIntroActive = false
	if m.introPlayer != nil && m.introPlayer.IsPlaying() {
		m.introPlayer.Pause()
	}
}

func (m *Manager) IsIntroPlaying() bool {
	return m.introPlayer != nil && m.introPlayer.IsPlaying()
}

func (m *Manager) ToggleMute() bool {
	m.isMuted = !m.isMuted
	if m.isMuted {
		m.PauseBGM()
	} else {
		m.ResumeBGM()
	}
	return !m.isMuted
}

func (m *Manager) IsMuted() bool {
	return m.isMuted
}

func (m *Manager) PlayRoar(isDouble bool) {
	if m.isMuted || m.ctx == nil {
		return
	}
	if isDouble {
		m.ctx.NewPlayerFromBytes(m.sndDoubleRoar).Play()
	} else {
		m.ctx.NewPlayerFromBytes(m.sndRoar).Play()
	}
}

func (m *Manager) PlayDuck() {
	if m.isMuted || m.ctx == nil {
		return
	}
	m.ctx.NewPlayerFromBytes(m.sndDuck).Play()
}

func (m *Manager) PlayHit() {
	if m.isMuted || m.ctx == nil {
		return
	}
	if len(m.sndEguaMano) > 0 {
		m.ctx.NewPlayerFromBytes(m.sndEguaMano).Play()
	}
	if len(m.sndHit) > 0 {
		m.ctx.NewPlayerFromBytes(m.sndHit).Play()
	}
}

func (m *Manager) PlayGameOver() {
	if m.isMuted || m.ctx == nil {
		return
	}
	m.ctx.NewPlayerFromBytes(m.sndGameOver).Play()
}

func (m *Manager) PlayStageUp() {
	if m.isMuted || m.ctx == nil {
		return
	}
	m.ctx.NewPlayerFromBytes(m.sndStageUp).Play()
}

func (m *Manager) PauseBGM() {
	if m.bgmPlayer != nil && m.bgmPlayer.IsPlaying() {
		m.bgmPlayer.Pause()
	}
	if m.introPlayer != nil && m.introPlayer.IsPlaying() {
		m.introPlayer.Pause()
	}
}

func (m *Manager) ResumeBGM() {
	if m.isMuted {
		return
	}
	if m.isIntroActive {
		if m.introPlayer != nil && !m.introPlayer.IsPlaying() {
			m.introPlayer.Play()
		}
	} else {
		if m.bgmPlayer != nil && !m.bgmPlayer.IsPlaying() {
			m.bgmPlayer.Play()
		}
	}
}

func (m *Manager) RestartBGM() {
	m.StopIntroBGM()
	if m.bgmPlayer != nil {
		_ = m.bgmPlayer.Rewind()
		if !m.isMuted {
			m.bgmPlayer.Play()
		}
	}
}
