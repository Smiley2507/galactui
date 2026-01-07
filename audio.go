package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

// Sound System
var (
	sampleRate  = beep.SampleRate(44100)
	audioInited bool
)

func initAudio() {
	err := speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	if err == nil {
		audioInited = true
	}
}

type oscillator struct {
	freq      float64
	freqDecay float64
	pos       float64
	vol       float64
}

func (g *oscillator) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		if g.freq <= 20 && g.freqDecay > 0 {
			return i, false
		}
		v := math.Sin(g.pos*2*math.Pi) * g.vol
		samples[i][0], samples[i][1] = v, v
		g.pos += g.freq / float64(sampleRate)
		g.freq -= g.freqDecay
	}
	return len(samples), true
}
func (g *oscillator) Err() error { return nil }

type noiseGen struct{}

func (g noiseGen) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		v := (rand.Float64()*2 - 1) * 0.1
		samples[i][0], samples[i][1] = v, v
	}
	return len(samples), true
}
func (g noiseGen) Err() error { return nil }

func playSound(kind string, enabled bool) {
	if !enabled || !audioInited {
		return
	}
	var s beep.Streamer
	switch kind {
	case "shoot":
		s = &oscillator{freq: 800, freqDecay: 0.2, vol: 0.1}
	case "explosion":
		s = beep.Take(sampleRate.N(300*time.Millisecond), noiseGen{})
	case "menu":
		s = beep.Take(sampleRate.N(50*time.Millisecond), &oscillator{freq: 440, freqDecay: 0, vol: 0.1})
	case "start":
		s = beep.Take(sampleRate.N(500*time.Millisecond), &oscillator{freq: 200, freqDecay: -0.02, vol: 0.1})
	case "gameover":
		s = beep.Take(sampleRate.N(1000*time.Millisecond), &oscillator{freq: 400, freqDecay: 0.01, vol: 0.1})
	case "select":
		s = beep.Take(sampleRate.N(100*time.Millisecond), &oscillator{freq: 660, freqDecay: 0, vol: 0.1})
	}
	if s != nil {
		speaker.Play(s)
	}
}
