package config

import (
	internalconfig "pomodoro-cli/internal/config"
	"pomodoro-cli/internal/ui"
)

// Run はconfigコマンドを実行する
func Run(cfg *internalconfig.Config) {
	ui.ShowConfig(cfg)
}
