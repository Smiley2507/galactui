package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func renderTable(scores []int) string {
	header := lipgloss.JoinHorizontal(lipgloss.Left,
		styleTableHeader.Width(10).Align(lipgloss.Center).Render("RANK"),
		styleTableHeader.Width(14).Align(lipgloss.Center).Render("SCORE"),
	)
	var rows []string
	rows = append(rows, header)
	for i, s := range scores {
		rankStr := fmt.Sprintf("#%d", i+1)
		scoreStr := fmt.Sprintf("%d", s)
		style := styleTableRow
		if i == 0 {
			style = style.Foreground(lipgloss.Color("220")).Bold(true) // Gold
		} else if i == 1 {
			style = style.Foreground(lipgloss.Color("248")) // Silver
		} else if i == 2 {
			style = style.Foreground(lipgloss.Color("137")) // Bronze
		}
		row := lipgloss.JoinHorizontal(lipgloss.Left,
			style.Width(10).Align(lipgloss.Center).Render(rankStr),
			style.Width(14).Align(lipgloss.Center).Render(scoreStr),
		)
		rows = append(rows, row)
	}
	return lipgloss.JoinVertical(lipgloss.Center, rows...)
}

func (m model) View() string {
	if m.state == stateStart {
		lines := strings.Split(titleArt, "\n")
		var styledLines []string
		for i, line := range lines {
			style := lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
			if (m.titleTick/5+i)%2 == 0 {
				style = style.Foreground(lipgloss.Color("255")) // White pulse
			}
			styledLines = append(styledLines, style.Render(line))
		}
		art := lipgloss.JoinVertical(lipgloss.Center, styledLines...)

		normalText, bossText, scoresText, helpText, settingsText, creditsText, quitText := "Normal Mode", "Boss Rush", "High Scores", "How to Play", "Settings", "Credits", "Quit"
		normalStyle, bossStyle, scoresStyle, helpStyle, settingsStyle, creditsStyle, quitStyle := styleMenu, styleMenu, styleMenu, styleMenu, styleMenu, styleMenu, styleMenu

		if m.menuIndex == 0 {
			normalStyle = styleSelected
			normalText = "> " + normalText
		} else if m.menuIndex == 1 {
			bossStyle = styleSelected
			bossText = "> " + bossText
		} else if m.menuIndex == 2 {
			scoresStyle = styleSelected
			scoresText = "> " + scoresText
		} else if m.menuIndex == 3 {
			helpStyle = styleSelected
			helpText = "> " + helpText
		} else if m.menuIndex == 4 {
			settingsStyle = styleSelected
			settingsText = "> " + settingsText
		} else if m.menuIndex == 5 {
			creditsStyle = styleSelected
			creditsText = "> " + creditsText
		} else if m.menuIndex == 6 {
			quitStyle = styleSelected
			quitText = "> " + quitText
		}

		content := lipgloss.NewStyle().
			Padding(2).
			Border(lipgloss.ThickBorder()).
			BorderForeground(primaryColor).
			Render(
				lipgloss.JoinVertical(lipgloss.Center,
					art,
					normalStyle.Render(normalText),
					bossStyle.Render(bossText),
					scoresStyle.Render(scoresText),
					helpStyle.Render(helpText),
					settingsStyle.Render(settingsText),
					creditsStyle.Render(creditsText),
					quitStyle.Render(quitText),
				),
			)
		return lipgloss.Place(m.termWidth, m.termHeight, lipgloss.Center, lipgloss.Center, content)
	}

	if m.state == stateHowToPlay {
		content := lipgloss.NewStyle().
			Padding(2).
			Border(lipgloss.ThickBorder()).
			BorderForeground(primaryColor).
			Render(
				lipgloss.JoinVertical(lipgloss.Center,
					styleTitle.Render("HOW TO PLAY"),
					styleMenu.Render("CONTROLS"),
					"Move: Arrow Keys",
					"Shoot: Space",
					"Pause: P / ESC",
					"Quit: Q / Ctrl+C",
					styleMenu.Render("\nPOWERUPS"),
					fmt.Sprintf("%s Multi-Shot  %s Shield  %s Nuke", stylePowerUp.Render(spritePowerupMulti), stylePowerUp.Render(spritePowerupShield), stylePowerUp.Render(spritePowerupNuke)),
					styleMenu.Render("\nENEMIES"),
					fmt.Sprintf("%s Basic  %s ZigZag  %s Chaser", styleEnemy.Render(spriteEnemyBasic), styleEnemyZigZag.Render(spriteEnemyZigZag), styleEnemyChaser.Render(spriteEnemyChaser)),
					styleSelected.Render("\nPress ESC to Return"),
				),
			)
		return lipgloss.Place(m.termWidth, m.termHeight, lipgloss.Center, lipgloss.Center, content)
	}

	if m.state == stateCredits {
		content := lipgloss.NewStyle().
			Padding(2).
			Border(lipgloss.ThickBorder()).
			BorderForeground(primaryColor).
			Render(
				lipgloss.JoinVertical(lipgloss.Center,
					styleTitle.Render("CREDITS"),
					styleMenu.Render("DEVELOPED BY"),
					"Celse",
					"",
					styleMenu.Render("LIBRARIES"),
					"Bubble Tea (Charm)",
					"Lip Gloss (Charm)",
					"Beep (Gopxl)",
					styleSelected.Render("\nPress ESC to Return"),
				),
			)
		return lipgloss.Place(m.termWidth, m.termHeight, lipgloss.Center, lipgloss.Center, content)
	}

	if m.state == stateHighScores {
		modeText := "NORMAL MODE"
		if m.gameMode == 1 {
			modeText = "BOSS RUSH"
		}

		tableContent := renderTable(m.topScores)

		content := lipgloss.NewStyle().Padding(1, 4).Border(lipgloss.ThickBorder()).BorderForeground(primaryColor).Render(
			lipgloss.JoinVertical(lipgloss.Center,
				styleTitle.Render("HIGH SCORES"),
				styleMenu.Render("< "+modeText+" >"),
				"",
				tableContent,
				"",
				styleSelected.Render("Press ESC to Return"),
			),
		)
		return lipgloss.Place(m.termWidth, m.termHeight, lipgloss.Center, lipgloss.Center, content)
	}

	if m.state == stateSettings {
		// Ensure no side effects like playSound are here
		soundText := fmt.Sprintf("Sound: %v", m.soundEnabled)
		if m.soundEnabled {
			soundText = "Sound: ON "
		} else {
			soundText = "Sound: OFF"
		}
		themeText := fmt.Sprintf("Theme: %s", themeNames[m.themeIndex])
		backText := "Back"

		soundStyle, themeStyle, backStyle := styleMenu, styleMenu, styleMenu

		if m.settingsMenuIndex == 0 {
			soundStyle = styleSelected
			soundText = "> " + soundText
		} else if m.settingsMenuIndex == 1 {
			themeStyle = styleSelected
			themeText = "> " + themeText
		} else {
			backStyle = styleSelected
			backText = "> " + backText
		}

		content := lipgloss.NewStyle().Padding(2).Border(lipgloss.ThickBorder()).BorderForeground(primaryColor).Render(
			lipgloss.JoinVertical(lipgloss.Center, styleTitle.Render("SETTINGS"), soundStyle.Render(soundText), themeStyle.Render(themeText), backStyle.Render(backText)),
		)
		return lipgloss.Place(m.termWidth, m.termHeight, lipgloss.Center, lipgloss.Center, content)
	}

	if m.state == stateGameOver {
		lines := strings.Split(gameOverArt, "\n")
		var styledLines []string
		for _, line := range lines {
			styledLines = append(styledLines, styleBoss.Render(line))
		}
		art := lipgloss.JoinVertical(lipgloss.Center, styledLines...)

		restartText, menuText, quitText := "Restart", "Main Menu", "Quit"
		restartStyle, menuStyle, quitStyle := styleMenu, styleMenu, styleMenu

		if m.menuIndex == 0 {
			restartStyle = styleSelected
			restartText = "> " + restartText
		} else if m.menuIndex == 1 {
			menuStyle = styleSelected
			menuText = "> " + menuText
		} else if m.menuIndex == 2 {
			quitStyle = styleSelected
			quitText = "> " + quitText
		}

		scoresText := renderTable(m.topScores)

		content := lipgloss.NewStyle().
			Padding(2).
			Border(lipgloss.ThickBorder()).
			BorderForeground(primaryColor).
			Render(
				lipgloss.JoinVertical(lipgloss.Center,
					art,
					fmt.Sprintf("\nSCORE: %d", m.score),
					"",
					scoresText,
					restartStyle.Render(restartText),
					menuStyle.Render(menuText),
					quitStyle.Render(quitText),
				),
			)
		return lipgloss.Place(m.termWidth, m.termHeight, lipgloss.Center, lipgloss.Center, content)
	}

	if m.state == statePaused {
		resumeText, quitText := "Resume", "Quit to Title"
		resumeStyle, quitStyle := styleMenu, styleMenu

		if m.pauseMenuIndex == 0 {
			resumeStyle = styleSelected
			resumeText = "> " + resumeText
		} else {
			quitStyle = styleSelected
			quitText = "> " + quitText
		}

		stats := fmt.Sprintf("SCORE: %d\nLIVES: %d\nWAVE: %d\nHIGH: %d", m.score, m.lives, m.wave, m.topScores[0])

		content := lipgloss.NewStyle().
			Padding(2).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(primaryColor).
			Render(
				lipgloss.JoinVertical(lipgloss.Center,
					styleTitle.Render("PAUSED"),
					styleMenu.Render(stats),
					"",
					resumeStyle.Render(resumeText),
					quitStyle.Render(quitText),
				),
			)
		return lipgloss.Place(m.termWidth, m.termHeight, lipgloss.Center, lipgloss.Center, content)
	}

	return m.viewPlaying()
}

func (m model) viewPlaying() string {
	// Clear Grid
	// Clear grid
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			m.grid[y][x] = " "
		}
	}

	// Stars
	for _, s := range m.stars {
		sx, sy := int(s.x), int(s.y)
		if sx >= 0 && sx < m.width && sy >= 0 && sy < m.height {
			m.grid[sy][sx] = styleStar.Render(s.char)
		}
	}

	// Player ship
	playerChar := spritePlayer
	if time.Now().Before(m.shieldUntil) {
		m.grid[m.playerY][m.playerX] = styleShield.Render(playerChar)
		if m.playerX > 0 {
			m.grid[m.playerY][m.playerX-1] = styleShield.Render("(")
		}
		if m.playerX < m.width-1 {
			m.grid[m.playerY][m.playerX+1] = styleShield.Render(")")
		}
	} else {
		m.grid[m.playerY][m.playerX] = stylePlayer.Render(playerChar)
	}

	// Bullets
	for _, b := range m.bullets {
		if b.y >= 0 && b.y < m.height && b.x >= 0 && b.x < m.width {
			if b.fromEnemy {
				m.grid[b.y][b.x] = styleBullet.Render(spriteBullet)
			} else {
				m.grid[b.y][b.x] = styleBullet.Render(spriteBullet)
			}
		}
	}

	// Enemies
	for _, e := range m.enemies {
		if e.alive && e.y >= 0 && e.y < m.height && e.x >= 0 && e.x < m.width {
			switch e.kind {
			case 1:
				m.grid[e.y][e.x] = styleEnemyZigZag.Render(spriteEnemyZigZag)
			case 2:
				m.grid[e.y][e.x] = styleEnemyChaser.Render(spriteEnemyChaser)
			default:
				m.grid[e.y][e.x] = styleEnemy.Render(spriteEnemyBasic)
			}
		}
	}

	// Boss
	for _, boss := range m.bosses {
		bossSprite := spriteBoss
		if boss.isSuper {
			bossSprite = spriteSuperBoss
		}
		for i, char := range bossSprite {
			bx := boss.x + i
			if bx >= 0 && bx < m.width && boss.y >= 0 && boss.y < m.height {
				m.grid[boss.y][bx] = styleBoss.Render(char)
			}
		}
	}

	// Powerups
	for _, p := range m.powerups {
		if p.y >= 0 && p.y < m.height && p.x >= 0 && p.x < m.width {
			if p.kind == 0 {
				m.grid[p.y][p.x] = stylePowerUp.Render(spritePowerupMulti) // Multi
			} else if p.kind == 1 {
				m.grid[p.y][p.x] = stylePowerUp.Render(spritePowerupShield) // Shield
			} else {
				m.grid[p.y][p.x] = stylePowerUp.Render(spritePowerupNuke) // Nuke
			}
		}
	}

	// Particles
	for _, p := range m.particles {
		px, py := int(p.x), int(p.y)
		if px >= 0 && px < m.width && py >= 0 && py < m.height {
			m.grid[py][px] = styleParticle.Render(p.char)
		}
	}

	// Level Up Notification
	if m.levelUpTimer > 0 && m.levelUpTimer%6 < 4 {
		text := " LEVEL UP! "
		row := m.height / 2
		col := (m.width - len(text)) / 2
		for i, r := range text {
			if col+i >= 0 && col+i < m.width {
				m.grid[row][col+i] = styleLevelUp.Render(string(r))
			}
		}
	}

	// Convert grid to rows
	var rows []string
	for _, row := range m.grid {
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Left, row...))
	}

	gameArea := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(primaryColor).
		Render(
			lipgloss.JoinVertical(lipgloss.Left, rows...),
		)

	lives := m.lives
	if lives < 0 {
		lives = 0
	}

	comboStr := ""
	if m.combo > 1 {
		comboStr = styleCombo.Render(fmt.Sprintf(" COMBO x%d ", m.combo))
	}

	header := styleScore.Render(
		fmt.Sprintf(" SCORE: %d   %s   WAVE: %d   HIGH: %d %s",
			m.score, strings.Repeat("❤️ ", lives), m.wave, m.topScores[0], comboStr),
	)

	var uiParts []string
	uiParts = append(uiParts, header)

	var bossesForBar []*Boss
	if m.gameMode == 1 {
		for _, b := range m.bosses {
			if b.isSuper {
				bossesForBar = append(bossesForBar, b)
			}
		}
	} else {
		bossesForBar = m.bosses
	}

	if len(bossesForBar) > 0 {
		totalHp, totalMax := 0, 0
		for _, b := range bossesForBar {
			totalHp += b.hp
			totalMax += b.maxHp
		}

		barWidth := 40
		percent := float64(totalHp) / float64(totalMax)
		if percent < 0 {
			percent = 0
		}
		filled := int(percent * float64(barWidth))

		barStr := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
		label := "BOSS"
		if m.gameMode == 1 {
			label = "SUPER BOSS"
		}
		bossBar := styleBossBar.Render(fmt.Sprintf("%s: [%s]", label, barStr))
		// Center the boss bar relative to the game width
		bossBar = lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(bossBar)
		uiParts = append(uiParts, bossBar)
	}

	uiParts = append(uiParts, gameArea)
	ui := lipgloss.JoinVertical(lipgloss.Left, uiParts...)

	return lipgloss.Place(m.termWidth, m.termHeight, lipgloss.Center, lipgloss.Center, ui)
}
