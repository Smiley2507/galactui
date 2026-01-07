package main

import (
	"time"
)

// Game States
type gameState int

const (
	stateStart gameState = iota
	stateHowToPlay
	stateSettings
	stateCredits
	stateHighScores
	statePlaying
	statePaused
	stateGameOver
)

type Bullet struct {
	x         int
	y         int
	fromEnemy bool
}

type Enemy struct {
	x     int
	y     int
	alive bool
	kind  int // 0: Basic, 1: ZigZag, 2: Chaser
	dir   int
	tick  int
}

type Boss struct {
	x, y      int
	hp, maxHp int
	width     int
	dir       int
	isSuper   bool
}

type PowerUp struct {
	x, y int
	kind int // 0: MultiShot, 1: Shield
}

type Particle struct {
	x, y   float64
	vx, vy float64
	life   int
	char   string
}

type Star struct {
	x, y  float64
	speed float64
	char  string
}

type model struct {
	width             int
	height            int
	playerX           int
	playerY           int
	bullets           []Bullet
	enemies           []Enemy
	bosses            []*Boss
	powerups          []PowerUp
	particles         []Particle
	stars             []Star
	score             int
	topScores         []int
	dir               int
	lives             int
	state             gameState
	wave              int
	shieldUntil       time.Time
	multiShotUntil    time.Time
	termWidth         int
	termHeight        int
	spawnTimer        int
	enemiesToSpawn    int
	titleTick         int
	menuIndex         int
	levelUpTimer      int
	gameMode          int // 0: Normal, 1: Boss Rush
	pauseMenuIndex    int
	combo             int
	comboTimer        int
	bossKillCount     int
	settingsMenuIndex int
	soundEnabled      bool
	themeIndex        int
	grid              [][]string
}
