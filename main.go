package main

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func initialModel() model {
	var stars []Star
	for i := 0; i < 40; i++ {
		stars = append(stars, Star{
			x:     rand.Float64() * float64(screenWidth),
			y:     rand.Float64() * float64(screenHeight),
			speed: 0.1 + rand.Float64()*0.5,
			char:  ".",
		})
	}

	grid := make([][]string, screenHeight)
	for y := 0; y < screenHeight; y++ {
		grid[y] = make([]string, screenWidth)
	}

	return model{
		width:             screenWidth,
		height:            screenHeight,
		playerX:           screenWidth / 2,
		playerY:           screenHeight - 4,
		bullets:           []Bullet{},
		enemies:           []Enemy{},
		bosses:            []*Boss{},
		powerups:          []PowerUp{},
		particles:         []Particle{},
		stars:             stars,
		score:             0,
		topScores:         loadHighScores(0),
		dir:               1,
		lives:             3,
		state:             stateStart,
		wave:              1,
		levelUpTimer:      0,
		gameMode:          0,
		pauseMenuIndex:    0,
		combo:             0,
		comboTimer:        0,
		bossKillCount:     0,
		settingsMenuIndex: 0,
		soundEnabled:      true,
		themeIndex:        0,
		grid:              grid,
	}
}

func (m *model) startNewGame() {
	m.score = 0
	m.lives = 3
	m.bullets = []Bullet{}
	m.powerups = []PowerUp{}
	m.bosses = []*Boss{}
	m.playerX = m.width / 2
	m.playerY = m.height - 4
	m.enemies = []Enemy{}
	m.spawnTimer = 0
	m.state = statePlaying
	m.shieldUntil = time.Time{}
	m.multiShotUntil = time.Time{}
	m.levelUpTimer = 0
	m.combo = 0
	m.comboTimer = 0
	m.bossKillCount = 0

	if m.gameMode == 1 {
		m.wave = 0 // Will be incremented to 1 immediately by loop
		m.enemiesToSpawn = 0
	} else {
		m.wave = 1
		m.enemiesToSpawn = 10 + (m.wave * 2)
	}
	m.topScores = loadHighScores(m.gameMode)
	playSound("start", m.soundEnabled)
}

type tickMsg time.Time

func tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tick(m.getTickDuration())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "p", "esc":
			if m.state == statePlaying {
				m.state = statePaused
				m.pauseMenuIndex = 0
			} else if m.state == statePaused {
				m.state = statePlaying
			} else if m.state == stateHowToPlay {
				m.state = stateStart
			} else if m.state == stateSettings {
				m.state = stateStart
			} else if m.state == stateCredits {
				m.state = stateStart
			} else if m.state == stateHighScores {
				m.state = stateStart
			}
		case "up":
			if m.state == stateStart || m.state == stateGameOver {
				if m.menuIndex > 0 {
					m.menuIndex--
					if m.state == stateStart && m.menuIndex <= 1 {
						m.topScores = loadHighScores(m.menuIndex)
					}
					playSound("menu", m.soundEnabled)
				}
			} else if m.state == statePaused {
				if m.pauseMenuIndex > 0 {
					m.pauseMenuIndex--
				}
			} else if m.state == stateSettings {
				if m.settingsMenuIndex > 0 {
					m.settingsMenuIndex--
					playSound("menu", m.soundEnabled)
				}
			}
		case "down":
			if m.state == stateStart {
				if m.menuIndex < 6 {
					m.menuIndex++
					if m.menuIndex <= 1 {
						m.topScores = loadHighScores(m.menuIndex)
					}
					playSound("menu", m.soundEnabled)
				}
			} else if m.state == stateGameOver {
				if m.menuIndex < 2 {
					m.menuIndex++
					playSound("menu", m.soundEnabled)
				}
			} else if m.state == statePaused {
				if m.pauseMenuIndex < 1 {
					m.pauseMenuIndex++
					playSound("menu", m.soundEnabled)
				}
			} else if m.state == stateSettings {
				if m.settingsMenuIndex < 2 {
					m.settingsMenuIndex++
					playSound("menu", m.soundEnabled)
				}
			}
		case "left", "right":
			if m.state == stateSettings {
				if m.settingsMenuIndex == 0 {
					m.soundEnabled = !m.soundEnabled
				} else if m.settingsMenuIndex == 1 {
					if msg.String() == "left" {
						m.themeIndex--
						if m.themeIndex < 0 {
							m.themeIndex = len(themes) - 1
						}
					} else {
						m.themeIndex++
						if m.themeIndex >= len(themes) {
							m.themeIndex = 0
						}
					}
					updateTheme(m.themeIndex)
				}
			} else if m.state == stateHighScores {
				if m.gameMode == 0 {
					m.gameMode = 1
				} else {
					m.gameMode = 0
				}
				m.topScores = loadHighScores(m.gameMode)
				playSound("menu", m.soundEnabled)
			}
		case "enter":
			if m.state == stateStart {
				playSound("select", m.soundEnabled)
				if m.menuIndex == 0 {
					m.gameMode = 0
					m.startNewGame()
				} else if m.menuIndex == 1 {
					m.gameMode = 1
					m.startNewGame()
				} else if m.menuIndex == 2 {
					m.state = stateHowToPlay
					m.state = stateHighScores
					m.gameMode = 0 // Default to Normal mode view
					m.topScores = loadHighScores(m.gameMode)
				} else if m.menuIndex == 3 {
					m.state = stateHowToPlay
				} else if m.menuIndex == 4 {
					m.state = stateSettings
				} else if m.menuIndex == 5 {
					m.state = stateCredits
				} else {
					return m, tea.Quit
				}
			} else if m.state == stateGameOver {
				playSound("select", m.soundEnabled)
				if m.menuIndex == 0 {
					m.startNewGame()
				} else if m.menuIndex == 1 {
					m.state = stateStart
				} else {
					return m, tea.Quit
				}
			} else if m.state == statePaused {
				playSound("select", m.soundEnabled)
				if m.pauseMenuIndex == 0 {
					m.state = statePlaying
				} else {
					m.state = stateStart
				}
			} else if m.state == stateSettings {
				playSound("select", m.soundEnabled)
				if m.settingsMenuIndex == 0 {
					m.soundEnabled = !m.soundEnabled
				} else if m.settingsMenuIndex == 1 {
					m.themeIndex++
					if m.themeIndex >= len(themes) {
						m.themeIndex = 0
					}
					updateTheme(m.themeIndex)
				} else {
					m.state = stateStart
				}
			}
		}

	case tickMsg:
		if m.state != statePlaying {
			m.titleTick++
			return m, tick(m.getTickDuration())
		}
		return m.updatePlaying()
	}

	// Handle Player Input only if Playing
	if m.state == statePlaying {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case " ":
				m.bullets = append(m.bullets, Bullet{
					x:         m.playerX,
					y:         m.playerY - 1,
					fromEnemy: false,
				})
				if time.Now().Before(m.multiShotUntil) || m.gameMode == 1 {
					m.bullets = append(m.bullets, Bullet{x: m.playerX - 1, y: m.playerY - 1, fromEnemy: false})
					m.bullets = append(m.bullets, Bullet{x: m.playerX + 1, y: m.playerY - 1, fromEnemy: false})
				}
				playSound("shoot", m.soundEnabled)

			case "left":
				if m.playerX > 1 {
					m.playerX--
				}
			case "right":
				if m.playerX < m.width-2 {
					m.playerX++
				}
			case "up":
				if m.playerY > 1 {
					m.playerY--
				}
			case "down":
				if m.playerY < m.height-2 {
					m.playerY++
				}
			}
		}
	}

	return m, nil
}

func main() {
	initAudio()
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
