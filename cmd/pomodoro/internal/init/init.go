package init

import (
	"fmt"

	"pomodoro-cli/internal/config"
	"pomodoro-cli/internal/ui"
)

// Run は対話的に設定ファイルを作成する
func Run() error {
	path, err := config.ConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// 現在の設定を読み込む（なければデフォルト）
	current, err := config.Load()
	if err != nil {
		// ファイルがない場合はデフォルト値を使用
		fmt.Printf("Note: %v (using defaults)\n", err)
	}
	defaults := config.Default()

	ui.ShowInitHeader()

	cfg := &config.Config{
		WorkDuration:       ui.PromptDuration("Work duration", current.WorkDuration, ui.FormatDuration(defaults.WorkDuration)),
		ShortBreakDuration: ui.PromptDuration("Short break duration", current.ShortBreakDuration, ui.FormatDuration(defaults.ShortBreakDuration)),
		LongBreakDuration:  ui.PromptDuration("Long break duration", current.LongBreakDuration, ui.FormatDuration(defaults.LongBreakDuration)),
		SessionsUntilLong:  ui.PromptInt("Sessions until long break", current.SessionsUntilLong, defaults.SessionsUntilLong),
		AutoStartBreaks:    ui.PromptBool("Auto-start breaks", current.AutoStartBreaks, defaults.AutoStartBreaks),
		AutoStartWork:      ui.PromptBool("Auto-start work", current.AutoStartWork, defaults.AutoStartWork),
		SoundEnabled:       ui.PromptBool("Enable sound", current.SoundEnabled, defaults.SoundEnabled),
		NotifyEnabled:      ui.PromptBool("Enable notifications", current.NotifyEnabled, defaults.NotifyEnabled),
	}

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	ui.ShowConfigCreated(path)
	return nil
}
