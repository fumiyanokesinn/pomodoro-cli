# pomodoro-cli

A simple, terminal-based Pomodoro timer written in Go.

## Features

- Customizable work and break durations
- Long breaks after configurable number of sessions
- System notifications and sound alerts
- Keyboard controls for pause, resume, skip, and reset
- Persistent configuration

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap fumiyanokesinn/tap
brew install pomodoro-cli
```

### From Source

Requires Go 1.21 or later.

```bash
git clone https://github.com/fumiyanokesinn/pomodoro-cli.git
cd pomodoro-cli
go build -o pomodoro ./cmd/pomodoro
```

## Usage

Start the timer with default settings (25 min work, 5 min short break, 15 min long break):

```bash
pomodoro
```

### Options

| Option | Description |
|--------|-------------|
| `-w`, `--work` | Work duration (e.g., `-w 25m`) |
| `-s`, `--short-break` | Short break duration (e.g., `-s 5m`) |
| `-l`, `--long-break` | Long break duration (e.g., `-l 15m`) |
| `-n`, `--sessions` | Sessions until long break (e.g., `-n 4`) |
| `--no-sound` | Disable notification sound |
| `--no-notify` | Disable system notifications |
| `--no-auto-break` | Disable auto-start breaks |
| `--no-auto-work` | Disable auto-start work sessions |
| `-v`, `--version` | Show version |
| `-h`, `--help` | Show help |

### Commands

| Command | Description |
|---------|-------------|
| `start` | Start the pomodoro timer (default) |
| `config` | Show current configuration |
| `init` | Initialize configuration file |

### Examples

```bash
# Start with custom durations
pomodoro -w 30m -s 10m -l 20m

# Start with 6 sessions before long break
pomodoro -n 6

# Start without sound
pomodoro --no-sound

# Show current configuration
pomodoro config
```

### Keyboard Shortcuts

During timer:

| Key | Action |
|-----|--------|
| `Space` | Pause/Resume |
| `q` | Quit |
| `s` | Skip to next session |
| `r` | Reset current session |

## Configuration

Configuration is stored at `~/.config/pomodoro/config.json`.

Run `pomodoro init` to create a default configuration file.

## License

[MIT](LICENSE)
