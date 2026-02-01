package main

import (
	"flag"
	"os"
	"time"

	configcmd "pomodoro-cli/cmd/pomodoro/internal/config"
	initcmd "pomodoro-cli/cmd/pomodoro/internal/init"
	"pomodoro-cli/cmd/pomodoro/internal/start"
	"pomodoro-cli/internal/config"
	"pomodoro-cli/internal/ui"
)

const version = "0.1.0"

func main() {
	// -h, --help, -v, --version を先に処理（flag.Parse前に）
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-h", "--help", "-help":
			ui.ShowUsage()
			os.Exit(0)
		case "-v", "--version", "-version":
			ui.ShowVersion(version)
			os.Exit(0)
		}
	}

	// Flag definitions
	var workDuration time.Duration
	var shortBreak time.Duration
	var longBreak time.Duration
	var sessions int
	var noSound, noNotify, noAutoBreak, noAutoWork bool
	var showVersion, showHelp bool

	// Duration flags
	flag.DurationVar(&workDuration, "w", 0, "")
	flag.DurationVar(&workDuration, "work", 0, "Work duration (e.g., 25m)")
	flag.DurationVar(&shortBreak, "s", 0, "")
	flag.DurationVar(&shortBreak, "short-break", 0, "Short break duration (e.g., 5m)")
	flag.DurationVar(&longBreak, "l", 0, "")
	flag.DurationVar(&longBreak, "long-break", 0, "Long break duration (e.g., 15m)")
	flag.IntVar(&sessions, "n", 0, "")
	flag.IntVar(&sessions, "sessions", 0, "Sessions until long break (e.g., 4)")

	// Disable flags
	flag.BoolVar(&noSound, "no-sound", false, "Disable notification sound")
	flag.BoolVar(&noNotify, "no-notify", false, "Disable system notifications")
	flag.BoolVar(&noAutoBreak, "no-auto-break", false, "Disable auto-start breaks")
	flag.BoolVar(&noAutoWork, "no-auto-work", false, "Disable auto-start work")

	// Other flags
	flag.BoolVar(&showVersion, "v", false, "")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&showHelp, "h", false, "")
	flag.BoolVar(&showHelp, "help", false, "Show help")

	flag.Usage = ui.ShowUsage
	flag.Parse()

	// 設定の読み込み
	cfg, loadErr := config.Load()
	if loadErr != nil {
		cfg = config.Default()
	}

	// フラグによる上書き
	if workDuration > 0 {
		cfg.WorkDuration = workDuration
	}
	if shortBreak > 0 {
		cfg.ShortBreakDuration = shortBreak
	}
	if longBreak > 0 {
		cfg.LongBreakDuration = longBreak
	}
	if sessions > 0 {
		cfg.SessionsUntilLong = sessions
	}
	if noSound {
		cfg.SoundEnabled = false
	}
	if noNotify {
		cfg.NotifyEnabled = false
	}
	if noAutoBreak {
		cfg.AutoStartBreaks = false
	}
	if noAutoWork {
		cfg.AutoStartWork = false
	}

	// コマンドの取得
	args := flag.Args()
	command := "start"
	if len(args) > 0 {
		command = args[0]
	}

	var err error
	switch command {
	case "start":
		err = start.Run(cfg)
	case "config":
		configcmd.Run(cfg)
	case "init":
		err = initcmd.Run()
	default:
		ui.ShowUnknownCommand(command)
		flag.Usage()
		os.Exit(1)
	}

	if err != nil {
		ui.ShowError(err.Error())
		os.Exit(1)
	}
}
