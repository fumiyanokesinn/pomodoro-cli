<p align="center">
  <img src="https://raw.githubusercontent.com/fumiyanokesinn/pomodoro-cli/main/assets/logo.svg" alt="pomodoro-cli logo" width="120">
</p>

<h1 align="center">pomodoro-cli</h1>

<p align="center">
  <strong>A beautiful, interactive Pomodoro timer for your terminal</strong>
</p>

<p align="center">
  <a href="https://github.com/fumiyanokesinn/pomodoro-cli/releases"><img src="https://img.shields.io/github/v/release/fumiyanokesinn/pomodoro-cli?style=flat-square&color=blue" alt="Release"></a>
  <a href="https://github.com/fumiyanokesinn/pomodoro-cli/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-green?style=flat-square" alt="License"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go"></a>
</p>

<p align="center">
  <img src="https://raw.githubusercontent.com/fumiyanokesinn/pomodoro-cli/main/assets/demo.gif" alt="Demo" width="600">
</p>

---

## Features

- **Interactive TUI** — Real-time progress bar and countdown display
- **Keyboard Controls** — Pause, resume, skip, and reset with single keystrokes
- **Smart Sessions** — Automatic short/long break rotation
- **System Notifications** — Desktop alerts when sessions complete
- **Sound Alerts** — Audio notifications (can be disabled)
- **Fully Configurable** — Customize durations, sessions, and behavior
- **Persistent Config** — Settings saved to `~/.config/pomodoro/`

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap fumiyanokesinn/tap
brew install pomodoro
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/fumiyanokesinn/pomodoro-cli/releases).

### From Source

```bash
git clone https://github.com/fumiyanokesinn/pomodoro-cli.git
cd pomodoro-cli
go build -o pomodoro ./cmd/pomodoro
```

## Quick Start

```bash
# Initialize config file
pomodoro init

# Start with default settings (25min work / 5min break)
pomodoro

# Custom work duration
pomodoro -w 30m

# Custom everything
pomodoro -w 25m -s 5m -l 15m -n 4
```

## Keyboard Shortcuts

| Key | Action |
|:---:|--------|
| `Space` | Pause / Resume |
| `s` | Skip to next session |
| `r` | Reset current session |
| `q` | Quit |

## Options

```
Options:
  -w, --work          Work duration (default: 25m)
  -s, --short-break   Short break duration (default: 5m)
  -l, --long-break    Long break duration (default: 15m)
  -n, --sessions      Sessions until long break (default: 4)
      --no-sound      Disable notification sound
      --no-notify     Disable system notifications
      --no-auto-break Disable auto-start breaks
      --no-auto-work  Disable auto-start work sessions
  -v, --version       Show version
  -h, --help          Show help

Commands:
  start               Start pomodoro timer (default)
  config              Show current configuration
  init                Initialize configuration file
```

## Configuration

Settings are stored at `~/.config/pomodoro/config.json`.

```bash
# Create default config
pomodoro init

# View current config
pomodoro config
```

Example config:

```json
{
  "work_duration": 1500000000000,
  "short_break_duration": 300000000000,
  "long_break_duration": 900000000000,
  "sessions_until_long_break": 4,
  "auto_start_breaks": true,
  "auto_start_work": true,
  "sound_enabled": true,
  "notify_enabled": true
}
```

## License

[MIT](LICENSE)

---

<p align="center">
  Made with ❤️ for focused work
</p>
