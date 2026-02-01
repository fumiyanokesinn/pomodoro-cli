package init

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"pomodoro-cli/internal/config"
	"pomodoro-cli/internal/ui"
)

// =============================================================================
// Run - init コマンドの実行
// =============================================================================

func TestRunはデフォルト値で設定ファイルを作成する(t *testing.T) {
	withTempHome(t, func(tmpHome string) {
		withStdinInput(t, "\n\n\n\n\n\n\n\n", func() {
			err := Run()
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			configPath := filepath.Join(tmpHome, ".config", "pomodoro", "config.json")
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				t.Error("設定ファイルが作成されていない")
			}

			defaults := config.Default()
			loaded, err := config.Load()
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}

			if loaded.WorkDuration != defaults.WorkDuration {
				t.Errorf("WorkDuration = %v, want %v", loaded.WorkDuration, defaults.WorkDuration)
			}
			if loaded.ShortBreakDuration != defaults.ShortBreakDuration {
				t.Errorf("ShortBreakDuration = %v, want %v", loaded.ShortBreakDuration, defaults.ShortBreakDuration)
			}
		})
	})
}

func TestRunはカスタム値を保存する(t *testing.T) {
	withTempHome(t, func(tmpHome string) {
		input := "30m\n10m\n20m\n6\nn\nn\ny\ny\n"
		withStdinInput(t, input, func() {
			err := Run()
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			loaded, err := config.Load()
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}

			if loaded.WorkDuration != 30*time.Minute {
				t.Errorf("WorkDuration = %v, want 30m", loaded.WorkDuration)
			}
			if loaded.ShortBreakDuration != 10*time.Minute {
				t.Errorf("ShortBreakDuration = %v, want 10m", loaded.ShortBreakDuration)
			}
			if loaded.LongBreakDuration != 20*time.Minute {
				t.Errorf("LongBreakDuration = %v, want 20m", loaded.LongBreakDuration)
			}
			if loaded.SessionsUntilLong != 6 {
				t.Errorf("SessionsUntilLong = %d, want 6", loaded.SessionsUntilLong)
			}
			if loaded.AutoStartBreaks {
				t.Error("AutoStartBreaks = true, want false")
			}
			if loaded.AutoStartWork {
				t.Error("AutoStartWork = true, want false")
			}
		})
	})
}

func TestRunは既存の設定値をEnterで維持する(t *testing.T) {
	withTempHome(t, func(tmpHome string) {
		existingCfg := &config.Config{
			WorkDuration:       45 * time.Minute,
			ShortBreakDuration: 10 * time.Minute,
			LongBreakDuration:  30 * time.Minute,
			SessionsUntilLong:  3,
			AutoStartBreaks:    false,
			AutoStartWork:      false,
			SoundEnabled:       false,
			NotifyEnabled:      false,
		}
		if err := existingCfg.Save(); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		withStdinInput(t, "\n\n\n\n\n\n\n\n", func() {
			err := Run()
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			loaded, err := config.Load()
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}

			if loaded.WorkDuration != 45*time.Minute {
				t.Errorf("WorkDuration = %v, want 45m", loaded.WorkDuration)
			}
			if loaded.SessionsUntilLong != 3 {
				t.Errorf("SessionsUntilLong = %d, want 3", loaded.SessionsUntilLong)
			}
			if loaded.AutoStartBreaks {
				t.Error("AutoStartBreaks = true, want false")
			}
		})
	})
}

// =============================================================================
// Config Format - 設定ファイルのフォーマット
// =============================================================================

func TestRunは正しいJSON形式で設定を保存する(t *testing.T) {
	withTempHome(t, func(tmpHome string) {
		withStdinInput(t, "\n\n\n\n\n\n\n\n", func() {
			err := Run()
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			configPath := filepath.Join(tmpHome, ".config", "pomodoro", "config.json")
			data, err := os.ReadFile(configPath)
			if err != nil {
				t.Fatalf("ReadFile() error = %v", err)
			}

			var raw map[string]interface{}
			if err := json.Unmarshal(data, &raw); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			expectedKeys := []string{
				"work_duration",
				"short_break_duration",
				"long_break_duration",
				"sessions_until_long_break",
				"auto_start_breaks",
				"auto_start_work",
				"sound_enabled",
				"notify_enabled",
			}

			for _, key := range expectedKeys {
				if _, ok := raw[key]; !ok {
					t.Errorf("config.json にキー %q が存在しない", key)
				}
			}
		})
	})
}

// =============================================================================
// Test Helpers
// =============================================================================

func withTempHome(t *testing.T, fn func(tmpHome string)) {
	t.Helper()
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	defer func() {
		if err := os.Setenv("HOME", origHome); err != nil {
			t.Errorf("failed to restore HOME: %v", err)
		}
	}()
	fn(tmpDir)
}

func withStdinInput(t *testing.T, input string, fn func()) {
	t.Helper()
	oldScanner := ui.GetScanner()
	ui.SetScanner(bufio.NewScanner(strings.NewReader(input)))
	defer ui.SetScanner(oldScanner)
	fn()
}
