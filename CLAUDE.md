# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CLI-based Pomodoro timer written in Go. Helps users manage work sessions using the Pomodoro Technique (typically 25-minute work sessions with 5-minute breaks).

## Build and Run Commands

```bash
# Build
go build -o pomodoro ./cmd/pomodoro

# Run
./pomodoro

# Run directly without building
go run ./cmd/pomodoro

# Run with options
./pomodoro -w 25m -s 5m -l 15m -n 4

# Run tests
go test ./...

# Run a single test
go test -run TestName ./path/to/package

# Format code
go fmt ./...

# Lint (requires golangci-lint)
golangci-lint run
```

## CLI Options

```
Commands:
  start              Start pomodoro timer (default)
  config             Show current configuration

Options:
  -w, --work         Work duration (e.g., -w 25m)
  -s, --short-break  Short break duration (e.g., -s 5m)
  -l, --long-break   Long break duration (e.g., -l 15m)
  -n, --sessions     Sessions until long break (e.g., -n 4)
      --no-sound     Disable notification sound
      --no-notify    Disable system notifications
      --no-auto-break  Disable auto-start breaks
      --no-auto-work   Disable auto-start work
  -v, --version      Show version
  -h, --help         Show help

Keyboard shortcuts (during timer):
  Space  - Pause/Resume
  q      - Quit
  s      - Skip to next session
  r      - Reset current session
```

## Architecture

```
cmd/pomodoro/        # Main entry point, event loop, key handling
internal/
  timer/             # Core timer logic
    timer.go         # Timer struct with Start/Pause/Resume/Stop
    state.go         # SessionType, TimerState, PomodoroState
  config/            # Configuration handling
    config.go        # Config struct, Load/Save, defaults
  ui/                # Terminal UI
    display.go       # Rendering, messages, progress bar
    input.go         # Raw mode input, KeyChan for select
    notification.go  # System notifications, sound
```

## Development Notes

- Use English for all user-facing messages and CLI help text
- Target Go 1.21+ for modern features
- Keep dependencies minimal for a CLI tool
- Use channel-based input (KeyChan) for integration with select statements
- Use time.Ticker for efficient event loops (not select/default + Sleep)
- Configuration stored at ~/.config/pomodoro/config.json
