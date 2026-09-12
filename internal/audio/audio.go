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
	sndShot       []byte
	sndDefeat     []byte
	sndJungleYell []byte
	sndTreasure   []byte
	sndEguaMano   []byte
	sndGameOver   []byte
	sndStageUp    []byte
	sndSplash     []byte
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

func createSplashSound() []byte {
	durationMs := 340
	numSamples := sampleRate * durationMs / 1000
	buf := make([]byte, numSamples*4)
	phase := 0.0
	noiseSeed := uint32(987654321)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(numSamples)

		// Impacto aquático: queda de frequência de 240Hz para 50Hz com modulação de borbulhas
		freq := 240.0*math.Exp(-9.0*t) + 50.0 + 28.0*math.Sin(t*math.Pi*16.0)
		phase += 2.0 * math.Pi * freq / float64(sampleRate)

		tone := math.Sin(phase) * 0.55

		// Borrifo líquido (ruído filtrado de splash d'água)
		noiseSeed = noiseSeed*1664525 + 1013904223
		noise := (float64((noiseSeed>>16)&0xFF)/128.0 - 1.0) * math.Exp(-4.5*t) * 0.45

		envelope := math.Exp(-3.2 * t)
		sample := int16((tone + noise) * envelope * 0.42 * 32767.0)

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

func createSlingshotShot() []byte {
	durationMs := 95
	numSamples := sampleRate * durationMs / 1000
	buf := make([]byte, numSamples*4)
	phase := 0.0

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(numSamples)
		freq := 780.0*math.Exp(-12.0*t) + 180.0
		phase += 2.0 * math.Pi * freq / float64(sampleRate)

		tone := math.Sin(phase) + 0.3*math.Sin(phase*2.0)
		attack := math.Min(1.0, float64(i)/float64(sampleRate*0.005))
		decay := math.Exp(-8.0 * t)
		sample := int16(tone * attack * decay * 0.32 * 32767.0)

		idx := i * 4
		buf[idx] = byte(sample)
		buf[idx+1] = byte(sample >> 8)
		buf[idx+2] = byte(sample)
		buf[idx+3] = byte(sample >> 8)
	}
	return buf
}

func createEnemyDefeat() []byte {
	durationMs := 130
	numSamples := sampleRate * durationMs / 1000
	buf := make([]byte, numSamples*4)
	phase := 0.0

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(numSamples)
		freq := 280.0*math.Sin(t*math.Pi*2.0) + 160.0*math.Exp(-5.0*t) + 80.0
		phase += 2.0 * math.Pi * freq / float64(sampleRate)

		noise := (float64((i*1103515245+12345)%32768)/16384.0 - 1.0) * 0.25
		tone := math.Sin(phase)*0.75 + noise

		attack := math.Min(1.0, float64(i)/float64(sampleRate*0.008))
		decay := math.Exp(-6.5 * t)
		sample := int16(tone * attack * decay * 0.35 * 32767.0)

		idx := i * 4
		buf[idx] = byte(sample)
		buf[idx+1] = byte(sample >> 8)
		buf[idx+2] = byte(sample)
		buf[idx+3] = byte(sample >> 8)
	}
	return buf
}

func createJungleYell() []byte {
	durationMs := 480
	numSamples := sampleRate * durationMs / 1000
	buf := make([]byte, numSamples*4)
	phase := 0.0

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(numSamples)

		var freq float64
		if t < 0.18 {
			freq = 360.0 + (t/0.18)*200.0
		} else if t < 0.55 {
			yodelPhase := (t - 0.18) / 0.37
			yodel := math.Sin(yodelPhase * math.Pi * 12.0)
			freq = 520.0 + yodel*65.0
		} else if t < 0.82 {
			yodelPhase := (t - 0.55) / 0.27
			yodel := math.Sin(yodelPhase * math.Pi * 8.0)
			freq = 430.0 + yodel*50.0
		} else {
			dropPhase := (t - 0.82) / 0.18
			freq = 380.0 - dropPhase*90.0
		}

		phase += 2.0 * math.Pi * freq / float64(sampleRate)
		tone := math.Sin(phase) + 0.45*math.Sin(phase*2.0) + 0.22*math.Sin(phase*3.0)

		attack := math.Min(1.0, float64(i)/float64(sampleRate*0.015))
		decay := 1.0
		if t > 0.75 {
			decay = math.Exp(-7.0 * (t - 0.75) / 0.25)
		}

		sample := int16(tone * attack * decay * 0.38 * 32767.0)
		idx := i * 4
		buf[idx] = byte(sample)
		buf[idx+1] = byte(sample >> 8)
		buf[idx+2] = byte(sample)
		buf[idx+3] = byte(sample >> 8)
	}
	return buf
}

func createTreasureJingle() []byte {
	durationMs := 340
	numSamples := sampleRate * durationMs / 1000
	buf := make([]byte, numSamples*4)
	phase := 0.0

	notes := []float64{1318.5, 1661.2, 1975.5, 2637.0}
	samplesPerNote := numSamples / len(notes)

	for i := 0; i < numSamples; i++ {
		noteIdx := i / samplesPerNote
		if noteIdx >= len(notes) {
			noteIdx = len(notes) - 1
		}
		freq := notes[noteIdx]
		phase += 2.0 * math.Pi * freq / float64(sampleRate)

		noteT := float64(i%samplesPerNote) / float64(samplesPerNote)
		envelope := math.Exp(-5.5 * noteT)

		tone := math.Sin(phase) + 0.35*math.Sin(phase*2.75)
		sample := int16(tone * envelope * 0.32 * 32767.0)

		idx := i * 4
		buf[idx] = byte(sample)
		buf[idx+1] = byte(sample >> 8)
		buf[idx+2] = byte(sample)
		buf[idx+3] = byte(sample >> 8)
	}
	return buf
}

// nesSquareWave gera onda retangular com duty cycle (estilo Canais Pulse 1 e Pulse 2 do NES)
func nesSquareWave(phase float64, duty float64) float64 {
	p := math.Mod(phase, 2.0*math.Pi)
	if p < 0 {
		p += 2.0 * math.Pi
	}
	if p/(2.0*math.Pi) < duty {
		return 1.0
	}
	return -1.0
}

// nesTriangleWave gera onda triangular pura (estilo Canal Triangle do NES)
func nesTriangleWave(phase float64) float64 {
	p := math.Mod(phase, 2.0*math.Pi)
	if p < 0 {
		p += 2.0 * math.Pi
	}
	norm := p / (2.0 * math.Pi)
	if norm < 0.5 {
		return 4.0*norm - 1.0
	}
	return 3.0 - 4.0*norm
}

// createCarimboBGM sintetiza um autêntico Carimbó Paraense em puro estilo 8-Bit Chiptune (NES/Arcade).
// 100% autoral e procedural, eliminando qualquer risco de direitos autorais de gravações comerciais.
func createCarimboBGM() []byte {
	const (
		r = 0.0 // Silêncio

		// Baixo (Oitava 2)
		c2 = 65.41
		d2 = 73.42
		e2 = 82.41
		f2 = 87.31
		g2 = 98.00
		a2 = 110.00
		b2 = 123.47

		// Baixo / Harmonia (Oitava 3)
		c3 = 130.81
		d3 = 146.83
		e3 = 164.81
		f3 = 174.61
		g3 = 196.00
		a3 = 220.00
		b3 = 246.94

		// Harmonia / Guitarrada (Oitava 4)
		c4 = 261.63
		d4 = 293.66
		e4 = 329.63
		f4 = 349.23
		g4 = 392.00
		a4 = 440.00
		b4 = 493.88

		// Melodia Principal (Oitava 5)
		c5 = 523.25
		d5 = 587.33
		e5 = 659.25
		f5 = 698.46
		g5 = 783.99
		a5 = 880.00
		b5 = 987.77

		// Agudos de Guitarrada (Oitava 6)
		c6 = 1046.50
		d6 = 1174.66
		e6 = 1318.51
		f6 = 1396.91
		g6 = 1567.98
	)

	// 1. Canal Pulse 1 (Square 50%): Melodia contagiante do Carimbó Paraense
	melodyNotes := []float64{
		// Seção A: Tema Principal Alegre (C -> G7)
		c5, r, e5, g5, g5, e5, d5, c5,
		d5, e5, d5, c5, a4, c5, d5, e5,
		d5, r, f5, a5, a5, f5, e5, d5,
		b4, c5, d5, b4, g4, b4, d5, f5,

		// Seção B: Balanço do Rio e Ver-o-Peso (F -> C -> G7 -> C)
		a5, r, c6, a5, g5, f5, e5, f5,
		g5, f5, e5, d5, c5, e5, g5, a5,
		g5, e5, c5, e5, d5, c5, b4, c5,
		d5, e5, f5, g5, a5, b5, c6, r,

		// Seção C: Solo Virtuoso de Guitarrada Paraense 8-Bit (Am -> Em -> F -> G7)
		a5, c6, b5, a5, e5, a5, b5, c6,
		b5, a5, g5, e5, g5, b5, d6, b5,
		c6, b5, a5, f5, a5, c6, e6, d6,
		d6, c6, b5, g5, b5, d6, f6, g6,

		// Seção D: Dança das Saias Rodadas e Cadência Triunfal (C -> F -> G7 -> C)
		e6, r, d6, c6, g5, c6, d6, e6,
		f6, e6, d6, c6, a5, c6, d6, f6,
		g6, f6, e6, d6, c6, b5, a5, g5,
		c5, e5, g5, c6, b5, g5, c6, r,
	}

	// 2. Canal Pulse 2 (Square 25%): Guitarrada Sincopada e Arpejos de Contraponto
	harmonyNotes := []float64{
		// Seção A (C / G)
		e4, g4, c5, g4, e4, g4, c5, g4,
		f4, a4, d5, a4, f4, a4, d5, a4,
		d4, g4, b4, g4, d4, g4, b4, g4,
		d4, g4, b4, g4, f4, g4, b4, g4,

		// Seção B (F / C / G / C)
		c4, f4, a4, f4, c4, f4, a4, f4,
		e4, g4, c5, g4, e4, g4, c5, g4,
		d4, g4, b4, g4, d4, g4, b4, g4,
		e4, g4, c5, g4, f4, a4, c5, r,

		// Seção C (Am / Em / F / G)
		c5, e5, a5, e5, c5, e5, a5, e5,
		b4, e5, g5, e5, b4, e5, g5, e5,
		a4, c5, f5, c5, a4, c5, f5, c5,
		b4, d5, g5, d5, b4, d5, g5, d5,

		// Seção D (C / F / G / C)
		c5, e5, g5, e5, c5, e5, g5, e5,
		d5, f5, a5, f5, d5, f5, a5, f5,
		b4, d5, f5, d5, b4, d5, f5, d5,
		c5, e5, g5, e5, d5, g5, c5, r,
	}

	// 3. Canal Triangle: Baixo Curimbó Tumbao Sincopado Tradicional
	bassNotes := []float64{
		// Seção A (C / G)
		c3, r, g2, c3, r, g2, c3, g2,
		c3, r, g2, c3, r, g2, c3, g2,
		g2, r, d2, g2, r, d2, g2, d2,
		g2, r, d2, g2, r, d2, g2, d2,

		// Seção B (F / C / G / C)
		f2, r, c2, f2, r, c2, f2, c2,
		c3, r, g2, c3, r, g2, c3, g2,
		g2, r, d2, g2, r, d2, g2, d2,
		c3, r, g2, c3, r, c3, g2, c3,

		// Seção C (Am / Em / F / G)
		a2, r, e2, a2, r, e2, a2, e2,
		e2, r, b2, e2, r, b2, e2, b2,
		f2, r, c2, f2, r, c2, f2, c2,
		g2, r, d2, g2, r, d2, g2, d2,

		// Seção D (C / F / G / C)
		c3, r, g2, c3, r, g2, c3, g2,
		f2, r, c2, f2, r, c2, f2, c2,
		g2, r, d2, g2, r, d2, g2, d2,
		c3, g2, c3, g2, c3, r, c3, r,
	}

	noteMs := 112 // ~135 BPM (semicolcheia rápida de Carimbó)
	stepSamples := sampleRate * noteMs / 1000
	totalSamples := stepSamples * len(melodyNotes)
	buf := make([]byte, totalSamples*4)

	sampleIdx := 0
	melPhase := 0.0
	harmPhase := 0.0
	bassPhase := 0.0
	kickPhase := 0.0
	noiseSeed := uint32(2147483647)

	nextNoise := func() float64 {
		noiseSeed = noiseSeed*1664525 + 1013904223
		val := int((noiseSeed >> 16) & 0x0F)
		return float64(val-8) / 8.0
	}

	for step := 0; step < len(melodyNotes); step++ {
		melFreq := melodyNotes[step]
		harmFreq := harmonyNotes[step]
		bassFreq := bassNotes[step]

		beatInBar := step % 8
		isTurnaround := (step == 63 || step == 127)

		for s := 0; s < stepSamples; s++ {
			t := float64(s) / float64(stepSamples)

			// 1. Canal 1 - Melodia Square 50% com vibrato sutil
			var melVal float64
			if melFreq > 0 {
				vibrato := 1.0
				if t > 0.45 {
					vibrato += 0.007 * math.Sin(2.0*math.Pi*5.5*(t-0.45)/0.55)
				}
				melPhase += 2.0 * math.Pi * (melFreq * vibrato) / float64(sampleRate)
				square := nesSquareWave(melPhase, 0.50)
				attack := math.Min(1.0, float64(s)/float64(sampleRate*0.006))
				decay := math.Exp(-1.7 * t)
				envelope := attack * decay
				if t > 0.86 {
					envelope *= math.Exp(-18.0 * (t - 0.86))
				}
				melVal = square * envelope
			}

			// 2. Canal 2 - Harmonia / Guitarrada Square 25% (staccato sincopado)
			var harmVal float64
			if harmFreq > 0 {
				harmPhase += 2.0 * math.Pi * harmFreq / float64(sampleRate)
				square := nesSquareWave(harmPhase, 0.25)
				attack := math.Min(1.0, float64(s)/float64(sampleRate*0.005))
				decay := math.Exp(-4.2 * t)
				harmVal = square * attack * decay
			}

			// 3. Canal 3 - Baixo Triangle (Tumbao sincopado de Curimbó)
			var bassVal float64
			if bassFreq > 0 {
				bassPhase += 2.0 * math.Pi * bassFreq / float64(sampleRate)
				tri := nesTriangleWave(bassPhase)
				attack := math.Min(1.0, float64(s)/float64(sampleRate*0.008))
				decay := math.Exp(-2.2 * t)
				bassVal = tri * attack * decay
			}

			// 4. Canal 4 - Percussão Curimbó (Kick de pitch sweep + Snare + Maraca)
			var percVal float64

			// Curimbó Kick grave (TUM)
			if beatInBar == 0 || beatInBar == 3 || beatInBar == 6 || isTurnaround {
				kickFreq := 155.0*math.Exp(-26.0*t) + 45.0
				kickPhase += 2.0 * math.Pi * kickFreq / float64(sampleRate)
				kickTone := nesTriangleWave(kickPhase)
				kickNoise := nextNoise() * math.Exp(-55.0*t) * 0.35
				kickEnv := math.Exp(-11.0 * t)
				percVal += (kickTone*0.85 + kickNoise) * kickEnv * 0.42
			}

			// Curimbó Snare estalado (TÁ)
			if beatInBar == 2 || beatInBar == 5 || beatInBar == 7 || isTurnaround {
				snareNoise := nextNoise()
				snareEnv := math.Exp(-26.0 * t)
				percVal += snareNoise * snareEnv * 0.32
			}

			// Maraca / Ganzá nos contratempos
			maracaNoise := nextNoise()
			maracaEnv := math.Exp(-65.0 * t)
			percVal += maracaNoise * maracaEnv * 0.08

			// Mixagem equilibrada estilo console 8-bit
			mix := melVal*0.19 + harmVal*0.12 + bassVal*0.22 + percVal*0.20

			// Soft limiter
			if mix > 0.95 {
				mix = 0.95
			} else if mix < -0.95 {
				mix = -0.95
			}

			// Quantização DAC estilo chiptune clássico
			sample := int16(mix * 32767.0)
			sample = (sample / 64) * 64

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
		sndShot:       createSlingshotShot(),
		sndDefeat:     createEnemyDefeat(),
		sndJungleYell: createJungleYell(),
		sndTreasure:   createTreasureJingle(),
		sndEguaMano:   createEguaManoSfx(),
		sndGameOver:   createSmoothTone(320, 95, 450, 0.3),
		sndStageUp:    createStageUpJingle(),
		sndSplash:     createSplashSound(),
		isMuted:       false,
		isIntroActive: false,
	}

	// 1. Trilha Sonora Principal: Autêntico Carimbó 8-Bit Chiptune Paraense (100% Autoral & Livre de Direitos Autorais)
	bgmBytes := createCarimboBGM()
	bgmLoop := ebitenaudio.NewInfiniteLoop(bytes.NewReader(bgmBytes), int64(len(bgmBytes)))
	bgmPlayer, err := ctx.NewPlayer(bgmLoop)
	if err == nil {
		bgmPlayer.SetVolume(0.20) // Volume agradável, suave e perfeitamente calibrado
		m.bgmPlayer = bgmPlayer
	}

	introBytes := createSaoBrasIntroBGM()
	introLoop := ebitenaudio.NewInfiniteLoop(bytes.NewReader(introBytes), int64(len(introBytes)))
	introPlayer, err := ctx.NewPlayer(introLoop)
	if err == nil {
		introPlayer.SetVolume(0.22) // Volume confortável e agradável para a tela de introdução e menus
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

func (m *Manager) PlayShot() {
	if m.isMuted || m.ctx == nil || len(m.sndShot) == 0 {
		return
	}
	m.ctx.NewPlayerFromBytes(m.sndShot).Play()
}

func (m *Manager) PlayDefeat() {
	if m.isMuted || m.ctx == nil || len(m.sndDefeat) == 0 {
		return
	}
	m.ctx.NewPlayerFromBytes(m.sndDefeat).Play()
}

func (m *Manager) PlayJungleYell() {
	if m.isMuted || m.ctx == nil || len(m.sndJungleYell) == 0 {
		return
	}
	m.ctx.NewPlayerFromBytes(m.sndJungleYell).Play()
}

func (m *Manager) PlayTreasure() {
	if m.isMuted || m.ctx == nil || len(m.sndTreasure) == 0 {
		return
	}
	m.ctx.NewPlayerFromBytes(m.sndTreasure).Play()
}

func (m *Manager) PlaySplash() {
	if m.isMuted || m.ctx == nil || len(m.sndSplash) == 0 {
		return
	}
	m.ctx.NewPlayerFromBytes(m.sndSplash).Play()
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
