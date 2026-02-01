package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// =============================================================================
// Default - デフォルト設定の生成
// =============================================================================

func TestDefaultは25分の作業時間を返す(t *testing.T) {
	cfg := Default()
	if cfg.WorkDuration != 25*time.Minute {
		t.Errorf("WorkDuration = %v, want 25m", cfg.WorkDuration)
	}
}

func TestDefaultは5分の短い休憩を返す(t *testing.T) {
	cfg := Default()
	if cfg.ShortBreakDuration != 5*time.Minute {
		t.Errorf("ShortBreakDuration = %v, want 5m", cfg.ShortBreakDuration)
	}
}

func TestDefaultは15分の長い休憩を返す(t *testing.T) {
	cfg := Default()
	if cfg.LongBreakDuration != 15*time.Minute {
		t.Errorf("LongBreakDuration = %v, want 15m", cfg.LongBreakDuration)
	}
}

func TestDefaultは4セッション後に長い休憩を設定する(t *testing.T) {
	cfg := Default()
	if cfg.SessionsUntilLong != 4 {
		t.Errorf("SessionsUntilLong = %d, want 4", cfg.SessionsUntilLong)
	}
}

func TestDefaultは自動開始を有効にする(t *testing.T) {
	cfg := Default()
	if !cfg.AutoStartBreaks {
		t.Error("AutoStartBreaks should be true by default")
	}
	if !cfg.AutoStartWork {
		t.Error("AutoStartWork should be true by default")
	}
}

func TestDefaultは通知と音声を有効にする(t *testing.T) {
	cfg := Default()
	if !cfg.SoundEnabled {
		t.Error("SoundEnabled should be true by default")
	}
	if !cfg.NotifyEnabled {
		t.Error("NotifyEnabled should be true by default")
	}
}

// =============================================================================
// Save/Load - 設定の保存と読み込み
// =============================================================================

func TestSaveAndLoadで設定を永続化できる(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := Default()
	cfg.WorkDuration = 30 * time.Minute
	cfg.ShortBreakDuration = 10 * time.Minute
	cfg.SessionsUntilLong = 6

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("json.MarshalIndent() error = %v", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	loaded, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}

	var loadedCfg Config
	if err := json.Unmarshal(loaded, &loadedCfg); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if loadedCfg.WorkDuration != 30*time.Minute {
		t.Errorf("loaded WorkDuration = %v, want 30m", loadedCfg.WorkDuration)
	}
	if loadedCfg.ShortBreakDuration != 10*time.Minute {
		t.Errorf("loaded ShortBreakDuration = %v, want 10m", loadedCfg.ShortBreakDuration)
	}
	if loadedCfg.SessionsUntilLong != 6 {
		t.Errorf("loaded SessionsUntilLong = %d, want 6", loadedCfg.SessionsUntilLong)
	}
}

func TestLoadはファイルから設定を読み込む(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".config", "pomodoro")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}
	configPath := filepath.Join(configDir, "config.json")

	cfg := &Config{
		WorkDuration:       30 * time.Minute,
		ShortBreakDuration: 10 * time.Minute,
		LongBreakDuration:  20 * time.Minute,
		SessionsUntilLong:  6,
		AutoStartBreaks:    false,
		AutoStartWork:      false,
		SoundEnabled:       false,
		NotifyEnabled:      false,
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("json.MarshalIndent() error = %v", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	origHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("os.Setenv() error = %v", err)
	}
	defer func() {
		if err := os.Setenv("HOME", origHome); err != nil {
			t.Errorf("os.Setenv() restore error = %v", err)
		}
	}()

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.WorkDuration != 30*time.Minute {
		t.Errorf("WorkDuration = %v, want 30m", loaded.WorkDuration)
	}
	if loaded.SessionsUntilLong != 6 {
		t.Errorf("SessionsUntilLong = %d, want 6", loaded.SessionsUntilLong)
	}
	if loaded.AutoStartBreaks != false {
		t.Error("AutoStartBreaks should be false")
	}
}

func TestSaveは設定をファイルに保存する(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("os.Setenv() error = %v", err)
	}
	defer func() {
		if err := os.Setenv("HOME", origHome); err != nil {
			t.Errorf("os.Setenv() restore error = %v", err)
		}
	}()

	cfg := Default()
	cfg.WorkDuration = 45 * time.Minute
	cfg.SessionsUntilLong = 3

	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.WorkDuration != 45*time.Minute {
		t.Errorf("WorkDuration = %v, want 45m", loaded.WorkDuration)
	}
	if loaded.SessionsUntilLong != 3 {
		t.Errorf("SessionsUntilLong = %d, want 3", loaded.SessionsUntilLong)
	}
}

// =============================================================================
// Load Fallback - 設定ファイルがない場合のフォールバック
// =============================================================================

func TestLoadはファイルがない場合デフォルト値を返す(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("os.Setenv() error = %v", err)
	}
	defer func() {
		if err := os.Setenv("HOME", origHome); err != nil {
			t.Errorf("os.Setenv() restore error = %v", err)
		}
	}()

	cfg, _ := Load()

	if cfg.WorkDuration != 25*time.Minute {
		t.Errorf("WorkDuration = %v, want 25m (default)", cfg.WorkDuration)
	}
	if cfg.ShortBreakDuration != 5*time.Minute {
		t.Errorf("ShortBreakDuration = %v, want 5m (default)", cfg.ShortBreakDuration)
	}
}

func TestLoadは破損したファイルの場合デフォルト値を返す(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".config", "pomodoro")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}
	configPath := filepath.Join(configDir, "config.json")

	if err := os.WriteFile(configPath, []byte("invalid json {{{"), 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	origHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("os.Setenv() error = %v", err)
	}
	defer func() {
		if err := os.Setenv("HOME", origHome); err != nil {
			t.Errorf("os.Setenv() restore error = %v", err)
		}
	}()

	cfg, _ := Load()

	if cfg.WorkDuration != 25*time.Minute {
		t.Errorf("WorkDuration = %v, want 25m (default)", cfg.WorkDuration)
	}
}

// =============================================================================
// Directory Creation - ディレクトリの自動作成
// =============================================================================

func TestSaveはディレクトリを自動作成する(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("os.Setenv() error = %v", err)
	}
	defer func() {
		if err := os.Setenv("HOME", origHome); err != nil {
			t.Errorf("os.Setenv() restore error = %v", err)
		}
	}()

	cfg := Default()
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	configPath := filepath.Join(tmpDir, ".config", "pomodoro", "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file was not created")
	}
}
