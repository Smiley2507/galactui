# Galactui
![galactui](/galactui.png)

Galactui is a retro-style, terminal-based space shooter game written in Go. It leverages the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework for the TUI (Terminal User Interface) and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for styling.

Experience classic arcade action directly in your terminal with smooth animations, sound effects, and intense boss battles.

## Features

- **Classic Arcade Action**: Fast-paced shooting mechanics.
- **Multiple Enemy Types**: Basic formations, ZigZag patterns, and Chasers.
- **Boss Battles**: Challenging bosses every 5 waves and a dedicated Boss Rush mode.
- **Power-ups**: Multi-Shot, Shield, and Nuke.
- **High Scores**: Local leaderboard tracking.
- **Cross-Platform**: Runs on Linux, Windows, and macOS.
- **Sound Effects**: Retro sound effects (can be toggled).

##  Download & Run

You can download the pre-built binaries for your operating system from the [GitHub Releases page](https://github.com/Smiley2507/galactui/releases/).

### Linux
1. Download `galactui-linux`.
2. Open a terminal in the download directory.
3. Make the binary executable and run it:
   ```bash
   chmod +x galactui-linux
   ./galactui-linux
   ```

### Windows
1. Download `galactui-windows.exe`.
2. Open **Command Prompt** or **PowerShell** in the download folder.
3. Run the game:
   ```powershell
   .\galactui-windows.exe
   ```

### macOS
1. Download `galactui-mac-intel` (Intel) or `galactui-mac-arm64` (Apple Silicon).
2. Open a terminal.
3. Make the binary executable and run it:
   ```bash
   chmod +x galactui-mac-arm64
   ./galactui-mac-arm64
   ```
   *Note: If macOS blocks the app, you may need to allow it in System Settings > Privacy & Security.*

## How to Play

![game-screen](/game-screen.png)


### Controls
| Key | Action |
| --- | --- |
| **Arrow Keys** | Move the ship |
| **Space** | Shoot |
| **P** / **ESC** | Pause game |
| **Q** / **Ctrl+C** | Quit |

### Game Modes
- **Normal Mode**: Survive endless waves of enemies. Difficulty increases over time. Defeat bosses every 5 levels.
- **Boss Rush**: Skip the minions and fight an endless stream of bosses with increasing difficulty.

## Building from Source

### Prerequisites
- Go (version 1.24 or higher)
- A terminal emulator with True Color support.

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Smiley2507/galactui.git
   cd galactui
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

### Running

```bash
# Using Make
make run
```

## Building

You can build the game for your current operating system using:

```bash
make build
```

### Cross-Platform Build

To build the game for multiple operating systems (Linux, Windows, macOS) so you can share it with friends:

```bash
make build-all
```

This will generate the following binaries in your directory:
- `galactui-linux` (Linux)
- `galactui-windows.exe` (Windows)
- `galactui-mac-intel` (macOS Intel)
- `galactui-mac-arm64` (macOS Apple Silicon)

To clean up built binaries:
```bash
make clean
```
#### **Happy Gaming!**
---
