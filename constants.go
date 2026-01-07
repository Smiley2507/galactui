package main

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	screenWidth  = 60
	screenHeight = 24
)

const titleArt = `
 ██████╗  █████╗ ██╗      █████╗  ██████╗████████╗██╗   ██╗██╗
██╔════╝ ██╔══██╗██║     ██╔══██╗██╔════╝╚══██╔══╝██║   ██║██║
██║  ███╗███████║██║     ███████║██║        ██║   ██║   ██║██║
██║   ██║██╔══██║██║     ██╔══██║██║        ██║   ██║   ██║██║
╚██████╔╝██║  ██║███████╗██║  ██║╚██████╗   ██║   ╚██████╔╝██║
 ╚═════╝ ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝ ╚═════╝   ╚═╝    ╚═════╝ ╚═╝
`

const gameOverArt = `
  ██████╗  █████╗ ███╗   ███╗███████╗     ██████╗ ██╗   ██╗███████╗██████╗ 
██╔════╝ ██╔══██╗████╗ ████║██╔════╝    ██╔═══██╗██║   ██║██╔════╝██╔══██╗
██║  ███╗███████║██╔████╔██║█████╗      ██║   ██║██║   ██║█████╗  ██████╔╝
██║   ██║██╔══██║██║╚██╔╝██║██╔══╝      ██║   ██║╚██╗ ██╔╝██╔══╝  ██╔══██╗
╚██████╔╝██║  ██║██║ ╚═╝ ██║███████╗    ╚██████╔╝ ╚████╔╝ ███████╗██║  ██║
 ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝     ╚═════╝   ╚═══╝  ╚══════╝╚═╝  ╚═╝  
`

// Styles
var (
	primaryColor     = lipgloss.Color("#FF5555") // Red
	styleTitle       = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true).Background(primaryColor).Padding(1, 2).MarginBottom(1)
	styleMenu        = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).MarginTop(1)
	styleSelected    = lipgloss.NewStyle().Foreground(primaryColor).Bold(true).MarginTop(1)
	stylePlayer      = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)  // Cyan
	styleEnemy       = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))            // Orange
	styleEnemyZigZag = lipgloss.NewStyle().Foreground(lipgloss.Color("213"))            // Pink
	styleEnemyChaser = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))             // Green
	styleBoss        = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true) // Red
	styleBullet      = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))            // Yellow
	styleShield      = lipgloss.NewStyle().Foreground(lipgloss.Color("51"))             // Light Blue
	stylePowerUp     = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)  // Green
	styleScore       = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	styleParticle    = lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // Orange/Gold
	styleStar        = lipgloss.NewStyle().Foreground(lipgloss.Color("238")) // Dark Gray
	styleBossBar     = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	styleLevelUp     = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
	styleCombo       = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true) // Gold
	styleTableHeader = lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Underline(true)
	styleTableRow    = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	// Themes
	themes     = []string{"#FF5555", "#50FA7B", "#8BE9FD", "#BD93F9"}
	themeNames = []string{"Red", "Green", "Blue", "Purple"}

	// Custom Sprites
	spritePlayer        = "🛦"
	spriteEnemyBasic    = "⮟"
	spriteEnemyZigZag   = "⛛"
	spriteEnemyChaser   = "𝈣"
	spriteBullet        = "⬝"
	spriteBoss          = []string{"╔", "𝩮", "🁢", "𝩮", "╗"}
	spriteSuperBoss     = []string{"◥", "▓", "☠", "▓", "◤"}
	spritePowerupMulti  = "󰜂"
	spritePowerupShield = "󰘜"
	spritePowerupNuke   = "☢"
)

func updateTheme(idx int) {
	if idx < 0 || idx >= len(themes) {
		return
	}
	c := themes[idx]
	primaryColor = lipgloss.Color(c)
	styleTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true).Background(primaryColor).Padding(1, 2).MarginBottom(1)
	styleSelected = lipgloss.NewStyle().Foreground(primaryColor).Bold(true).MarginTop(1)
	styleScore = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	styleTableHeader = lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Underline(true)
}