package main

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) getTickDuration() time.Duration {
	return time.Millisecond * 50
}

func (m model) getScoreMultiplier() float64 {
	return 1.0
}

func (m model) updatePlaying() (tea.Model, tea.Cmd) {
	if m.levelUpTimer > 0 {
		m.levelUpTimer--
	}

	// Combo Decay
	if m.comboTimer > 0 {
		m.comboTimer--
	} else {
		m.combo = 0
	}

	// --- LOGIC START ---

	// Spawning Logic
	if len(m.bosses) == 0 && m.enemiesToSpawn > 0 && m.gameMode == 0 {
		m.spawnTimer--
		if m.spawnTimer <= 0 {
			// Spawn one enemy
			kind := 0
			if m.wave > 1 && rand.Intn(100) < 30 {
				kind = 1 // ZigZag
			}
			if m.wave > 2 && rand.Intn(100) < 20 {
				kind = 2 // Chaser
			}

			m.enemies = append(m.enemies, Enemy{
				x:     rand.Intn(m.width-8) + 4,
				y:     1,
				alive: true,
				kind:  kind,
				dir:   1,
			})
			m.enemiesToSpawn--
			// Reset timer (faster as wave increases)
			base := 30
			m.spawnTimer = max(10, base-(m.wave*2))
		}
	}

	// Boss Rush Spawning
	if m.gameMode == 1 {
		m.spawnTimer--
		if m.spawnTimer <= 0 {
			m.bossKillCount++
			isSuper := m.bossKillCount%5 == 0
			hp := 30
			if isSuper {
				hp = 150
			}

			m.bosses = append(m.bosses, &Boss{
				x:       rand.Intn(m.width - 10),
				y:       -5, // Start above screen
				hp:      hp,
				maxHp:   hp,
				width:   5,
				dir:     1,
				isSuper: isSuper,
			})

			// Spawn faster as we go
			m.spawnTimer = max(30, 80-(m.bossKillCount*2))
		}
	}

	// Particle Logic
	var activeParticles []Particle
	for _, p := range m.particles {
		p.x += p.vx
		p.y += p.vy
		p.life--
		if p.life > 0 && p.x >= 0 && p.x < float64(m.width) && p.y >= 0 && p.y < float64(m.height) {
			activeParticles = append(activeParticles, p)
		}
	}
	m.particles = activeParticles

	// Star Logic
	for i := range m.stars {
		m.stars[i].y += m.stars[i].speed
		if m.stars[i].y >= float64(m.height) {
			m.stars[i].y = 0
			m.stars[i].x = rand.Float64() * float64(m.width)
		}
	}

	// Powerup movement
	var activePowerups []PowerUp
	for _, p := range m.powerups {
		p.y++
		if p.y < m.height {
			activePowerups = append(activePowerups, p)
		}
	}
	m.powerups = activePowerups

	// Bullet movement
	var updated []Bullet
	for _, b := range m.bullets {
		if b.fromEnemy {
			b.y++
		} else {
			b.y--
		}

		if b.y > 0 && b.y < m.height {
			updated = append(updated, b)
		}
	}
	m.bullets = updated

	// Collision: Player Bullets -> Enemies
	for bi := range m.bullets {
		for ei := range m.enemies {
			if !m.enemies[ei].alive {
				continue
			}

			if m.bullets[bi].x == m.enemies[ei].x &&
				m.bullets[bi].y == m.enemies[ei].y {

				m.enemies[ei].alive = false
				m.bullets[bi].y = -1

				m.combo++
				m.comboTimer = 40 // ~3 seconds
				baseScore := 10 * (1 + m.combo/5)
				m.score += int(float64(baseScore) * m.getScoreMultiplier())
				playSound("explosion", m.soundEnabled)

				// Spawn Particles
				for i := 0; i < 6; i++ {
					m.particles = append(m.particles, Particle{
						x:    float64(m.enemies[ei].x),
						y:    float64(m.enemies[ei].y),
						vx:   (rand.Float64() - 0.5) * 2.0,
						vy:   (rand.Float64() - 0.5) * 2.0,
						life: 8,
						char: []string{"*", "•", "°", "·"}[rand.Intn(4)],
					})
				}

				// Chance to drop powerup
				if rand.Intn(100) < 5 { // 5% chance
					m.powerups = append(m.powerups, PowerUp{
						x:    m.enemies[ei].x,
						y:    m.enemies[ei].y,
						kind: rand.Intn(3),
					})
				}
			}
		}

		// Player Bullets -> Boss
		for i := range m.bosses {
			boss := m.bosses[i]
			b := m.bullets[bi]
			if !b.fromEnemy && b.x >= boss.x && b.x < boss.x+boss.width && b.y == boss.y {
				m.bullets[bi].y = -1
				boss.hp--
				if boss.hp <= 0 {
					// Boss Dead
					m.score += int(500 * m.getScoreMultiplier())
					playSound("explosion", m.soundEnabled)
					// Drop Powerup
					m.powerups = append(m.powerups, PowerUp{
						x:    boss.x + 2,
						y:    boss.y,
						kind: rand.Intn(3),
					})
					// Remove boss logic handled below in cleanup
				}
			}
		}
	}

	// Collision: Enemy Bullets -> Player
	isShielded := time.Now().Before(m.shieldUntil)
	for _, b := range m.bullets {
		if b.fromEnemy &&
			b.x == m.playerX &&
			b.y == m.playerY {

			if !isShielded {
				m.lives--
				if m.lives <= 0 {
					m.state = stateGameOver
					saveScore(m.score, m.gameMode)
					m.topScores = loadHighScores(m.gameMode)
					playSound("gameover", m.soundEnabled)
				}
			}
		}
	}

	// Collision: Powerups -> Player
	for i, p := range m.powerups {
		if p.x == m.playerX && p.y == m.playerY {
			// Activate powerup
			if p.kind == 0 {
				m.multiShotUntil = time.Now().Add(time.Second * 5)
			} else if p.kind == 1 {
				m.shieldUntil = time.Now().Add(time.Second * 5)
			} else {
				// Nuke
				for j := range m.enemies {
					if m.enemies[j].alive {
						m.enemies[j].alive = false
						m.score += int(10 * m.getScoreMultiplier())
						// Spawn particles
						for k := 0; k < 3; k++ {
							m.particles = append(m.particles, Particle{
								x:    float64(m.enemies[j].x),
								y:    float64(m.enemies[j].y),
								vx:   (rand.Float64() - 0.5) * 2.0,
								vy:   (rand.Float64() - 0.5) * 2.0,
								life: 8,
								char: []string{"*", "•"}[rand.Intn(2)],
							})
						}
					}
				}
				for _, b := range m.bosses {
					b.hp -= 50
				}
				playSound("explosion", m.soundEnabled)
			}
			// Remove powerup (move off screen)
			m.powerups[i].y = m.height + 1
		}
	}

	// Check Wave Cleared
	enemiesAlive := 0
	for _, e := range m.enemies {
		if e.alive {
			enemiesAlive++
		}
	}

	// Cleanup Dead Bosses
	var activeBosses []*Boss
	for _, b := range m.bosses {
		if b.hp > 0 {
			activeBosses = append(activeBosses, b)
		}
	}
	m.bosses = activeBosses

	if enemiesAlive == 0 && m.enemiesToSpawn == 0 && len(m.bosses) == 0 && m.gameMode == 0 {
		m.wave++
		m.levelUpTimer = 30
		m.bullets = []Bullet{}
		m.powerups = []PowerUp{}

		// Normal Mode
		// Boss Wave every 5 levels
		if m.wave%5 == 0 {
			m.bosses = append(m.bosses, &Boss{
				x:     m.width/2 - 2,
				y:     3,
				hp:    20 * (m.wave / 5),
				maxHp: 20 * (m.wave / 5),
				width: 5,
				dir:   1,
			})
		} else {
			m.enemiesToSpawn = 10 + (m.wave * 2)
		}
	}

	// Boss Logic
	for _, boss := range m.bosses {
		// Move Boss
		if boss.y < 2 {
			boss.y++ // Move into screen
		} else if m.gameMode == 1 && rand.Intn(100) < 2 {
			boss.y++ // Slowly move down in Boss Rush
		}

		if rand.Intn(10) < 3 {
			boss.x += boss.dir
			if boss.x <= 2 || boss.x >= m.width-7 {
				boss.dir *= -1
			}
		}
		// Boss Shoot
		if rand.Intn(100) < 5 {
			// Shoot from center
			m.bullets = append(m.bullets, Bullet{
				x:         boss.x + 2,
				y:         boss.y + 1,
				fromEnemy: true,
			})
			// Shoot from sides
			if rand.Intn(2) == 0 {
				m.bullets = append(m.bullets, Bullet{x: boss.x, y: boss.y + 1, fromEnemy: true})
				m.bullets = append(m.bullets, Bullet{x: boss.x + 4, y: boss.y + 1, fromEnemy: true})
			}
		}
	}

	var aliveBullets []Bullet
	for _, b := range m.bullets {
		if b.y >= 0 {
			aliveBullets = append(aliveBullets, b)
		}
	}
	m.bullets = aliveBullets

	// Move enemies
	edgeHit := false
	for i := range m.enemies {
		e := &m.enemies[i]
		if !e.alive {
			continue
		}

		switch e.kind {
		case 0: // Basic (Formation)
			e.x += m.dir
			if e.x <= 1 || e.x >= m.width-2 {
				edgeHit = true
			}
		case 1: // ZigZag
			e.x += e.dir
			if e.x <= 1 || e.x >= m.width-2 {
				e.dir *= -1
				e.x += e.dir
			}
			// Move down occasionally
			if rand.Intn(100) < 5 {
				e.y++
			}
		case 2: // Chaser
			e.tick++
			if e.tick%3 == 0 {
				e.y++
				if e.x < m.playerX {
					e.x++
				} else if e.x > m.playerX {
					e.x--
				}
			}
		}
	}

	// Enemy shooting
	for _, e := range m.enemies {
		chance := 5 + m.wave
		if e.alive && rand.Intn(1000) < chance { // Scale difficulty
			m.bullets = append(m.bullets, Bullet{
				x:         e.x,
				y:         e.y + 1,
				fromEnemy: true,
			})
		}
	}

	// Check enemy reached player zone
	for i := range m.enemies {
		if m.enemies[i].alive && m.enemies[i].y >= m.playerY {
			m.enemies[i].alive = false
			if !isShielded {
				m.lives--
			}
		}
	}

	if m.lives <= 0 {
		m.state = stateGameOver
		saveScore(m.score, m.gameMode)
		m.topScores = loadHighScores(m.gameMode)
		playSound("gameover", m.soundEnabled)
	}

	if edgeHit {
		m.dir *= -1
		for i := range m.enemies {
			if m.enemies[i].kind == 0 {
				m.enemies[i].y++
			}
		}
	}

	return m, tick(m.getTickDuration())
}
