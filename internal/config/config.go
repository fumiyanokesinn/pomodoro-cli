package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Config はポモドーロタイマーの設定を保持する
type Config struct {
	WorkDuration       time.Duration `json:"work_duration"`
	ShortBreakDuration time.Duration `json:"short_break_duration"`
	LongBreakDuration  time.Duration `json:"long_break_duration"`
	SessionsUntilLong  int           `json:"sessions_until_long_break"`
	AutoStartBreaks    bool          `json:"auto_start_breaks"`
	AutoStartWork      bool          `json:"auto_start_work"`
	SoundEnabled       bool          `json:"sound_enabled"`
	NotifyEnabled      bool          `json:"notify_enabled"`
}

// Default はデフォルトの設定を返す
func Default() *Config {
	return &Config{
		WorkDuration:       25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
		SessionsUntilLong:  4,
		AutoStartBreaks:    true,
		AutoStartWork:      true,
		SoundEnabled:       true,
		NotifyEnabled:      true,
	}
}

// ConfigPath は設定ファイルのパスを返す
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "pomodoro", "config.json"), nil
}

// Load は設定ファイルから設定を読み込む
// ファイルが存在しない場合はデフォルト設定を返す
func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return Default(), err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Default(), err
	}

	cfg := Default()
	if err := json.Unmarshal(data, cfg); err != nil {
		return Default(), err
	}

	return cfg, nil
}

// Save は設定をファイルに保存する
func (c *Config) Save() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	// ディレクトリを作成
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
