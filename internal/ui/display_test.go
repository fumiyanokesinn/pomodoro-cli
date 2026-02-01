package ui

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"pomodoro-cli/internal/config"
	"pomodoro-cli/internal/timer"
)

// =============================================================================
// ProgressBar - プログレスバーの描画
// =============================================================================

func TestProgressBarShowsEmptyBarAtZeroPercent(t *testing.T) {
	got := progressBar(0.0, 10)
	if got != "░░░░░░░░░░" {
		t.Errorf("progressBar(0.0, 10) = %q, want empty bar", got)
	}
}

func TestProgressBarShowsHalfFilledAtFiftyPercent(t *testing.T) {
	got := progressBar(0.5, 10)
	if got != "█████░░░░░" {
		t.Errorf("progressBar(0.5, 10) = %q, want half filled", got)
	}
}

func TestProgressBarShowsFullBarAtHundredPercent(t *testing.T) {
	got := progressBar(1.0, 10)
	if got != "██████████" {
		t.Errorf("progressBar(1.0, 10) = %q, want full bar", got)
	}
}

// =============================================================================
// FormatDuration - 時間のフォーマット
// =============================================================================

func TestFormatDurationShowsMinutes(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{25 * time.Minute, "25m"},
		{5 * time.Minute, "5m"},
		{60 * time.Minute, "60m"},
		{0, "0m"},
	}
	for _, tt := range tests {
		got := FormatDuration(tt.duration)
		if got != tt.expected {
			t.Errorf("FormatDuration(%v) = %q, want %q", tt.duration, got, tt.expected)
		}
	}
}

func TestFormatDurationTruncatesSeconds(t *testing.T) {
	if got := FormatDuration(30 * time.Second); got != "0m" {
		t.Errorf("FormatDuration(30s) = %q, want 0m", got)
	}
	if got := FormatDuration(90 * time.Second); got != "1m" {
		t.Errorf("FormatDuration(90s) = %q, want 1m", got)
	}
}

// =============================================================================
// Init Command UI - 初期化コマンドの表示
// =============================================================================

func TestShowInitHeaderDisplaysSetupInstructions(t *testing.T) {
	output := captureStdout(t, func() {
		ShowInitHeader()
	})

	assertContains(t, output, "CONFIGURATION SETUP")
	assertContains(t, output, "Press Enter")
}

func TestShowConfigCreatedDisplaysFilePath(t *testing.T) {
	output := captureStdout(t, func() {
		ShowConfigCreated("/test/path/config.json")
	})

	assertContains(t, output, "Configuration saved")
	assertContains(t, output, "/test/path/config.json")
}

// =============================================================================
// Help and Usage - ヘルプとUsageの表示
// =============================================================================

func TestShowUsageDisplaysAllCommandsAndOptions(t *testing.T) {
	output := captureStderr(t, func() {
		ShowUsage()
	})

	expectedStrings := []string{
		"Usage: pomodoro",
		"Commands:",
		"start", "config", "init",
		"Options:",
		"-w, --work",
		"-s, --short-break",
		"-l, --long-break",
		"-n, --sessions",
		"--no-sound",
		"--no-notify",
		"-v, --version",
		"-h, --help",
	}

	for _, expected := range expectedStrings {
		assertContains(t, output, expected)
	}
}

func TestShowUnknownCommandDisplaysErrorMessage(t *testing.T) {
	output := captureStderr(t, func() {
		ShowUnknownCommand("unknown")
	})

	assertContains(t, output, "Unknown command: unknown")
}

func TestShowErrorDisplaysMessage(t *testing.T) {
	output := captureStderr(t, func() {
		ShowError("test error message")
	})

	assertContains(t, output, "test error message")
}

func TestShowVersionDisplaysVersionString(t *testing.T) {
	output := captureStdout(t, func() {
		ShowVersion("1.0.0")
	})

	assertContains(t, output, "pomodoro version 1.0.0")
}

// =============================================================================
// Config Display - 設定の表示
// =============================================================================

func TestShowConfigDisplaysAllSettings(t *testing.T) {
	cfg := &config.Config{
		WorkDuration:       25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
		SessionsUntilLong:  4,
		AutoStartBreaks:    true,
		AutoStartWork:      false,
		SoundEnabled:       true,
		NotifyEnabled:      false,
	}

	output := captureStdout(t, func() {
		ShowConfig(cfg)
	})

	expectedStrings := []string{
		"CURRENT CONFIGURATION",
		"Work duration:",
		"Short break:",
		"Long break:",
		"Sessions until long:",
		"Auto-start breaks:",
		"Auto-start work:",
		"Sound enabled:",
		"Notify enabled:",
		"Yes",
		"No",
	}

	for _, expected := range expectedStrings {
		assertContains(t, output, expected)
	}
}

// =============================================================================
// Timer Display - タイマーの表示
// =============================================================================

func TestRenderTimerDoesNotPanicWithNilSession(t *testing.T) {
	RenderTimer(nil, timer.StateRunning)
}

func TestRenderTimerDisplaysSessionTypeAndTime(t *testing.T) {
	session := &timer.Session{
		Type:      timer.SessionWork,
		Duration:  25 * time.Minute,
		Remaining: 24 * time.Minute,
	}

	output := captureStdout(t, func() {
		RenderTimer(session, timer.StateRunning)
	})

	assertContains(t, output, "Work")
	assertContains(t, output, "24:00")
}

func TestRenderTimerShowsPauseIconWhenPaused(t *testing.T) {
	session := &timer.Session{
		Type:      timer.SessionShortBreak,
		Duration:  5 * time.Minute,
		Remaining: 3 * time.Minute,
	}

	output := captureStdout(t, func() {
		RenderTimer(session, timer.StatePaused)
	})

	assertContains(t, output, "⏸")
}

// =============================================================================
// Welcome Screen - ウェルカム画面
// =============================================================================

func TestShowWelcomeDisplaysSettingsAndShortcuts(t *testing.T) {
	output := captureStdout(t, func() {
		ShowWelcome(25*time.Minute, 5*time.Minute, 15*time.Minute)
	})

	expectedStrings := []string{
		"Work:",
		"Short Break:",
		"Long Break:",
		"Keyboard Shortcuts",
		"[Space]",
		"Pause/Resume",
		"[q] Quit",
	}

	for _, expected := range expectedStrings {
		assertContains(t, output, expected)
	}
}

// =============================================================================
// Session Messages - セッションメッセージ
// =============================================================================

func TestShowSessionCompleteDisplaysWorkComplete(t *testing.T) {
	output := captureStdout(t, func() {
		ShowSessionComplete(timer.SessionWork)
	})
	assertContains(t, output, "Work session complete")
}

func TestShowSessionCompleteDisplaysBreakOver(t *testing.T) {
	output := captureStdout(t, func() {
		ShowSessionComplete(timer.SessionShortBreak)
	})
	assertContains(t, output, "Break over")
}

func TestShowStartSessionDisplaysSessionType(t *testing.T) {
	output := captureStdout(t, func() {
		ShowStartSession(timer.SessionWork)
	})
	assertContains(t, output, "Starting Work")
}

func TestShowPausedDisplaysPausedMessage(t *testing.T) {
	output := captureStdout(t, func() {
		ShowPaused()
	})
	assertContains(t, output, "Paused")
}

func TestShowResumedDisplaysResumedMessage(t *testing.T) {
	output := captureStdout(t, func() {
		ShowResumed()
	})
	assertContains(t, output, "Resumed")
}

func TestShowSkippedDisplaysSessionType(t *testing.T) {
	output := captureStdout(t, func() {
		ShowSkipped(timer.SessionShortBreak)
	})
	assertContains(t, output, "Skipped")
	assertContains(t, output, "Short Break")
}

func TestShowResetDisplaysResetMessage(t *testing.T) {
	output := captureStdout(t, func() {
		ShowReset()
	})
	assertContains(t, output, "Reset")
}

func TestShowExitDisplaysExitMessage(t *testing.T) {
	output := captureStdout(t, func() {
		ShowExit()
	})
	assertContains(t, output, "Thanks for using Pomodoro")
}

// =============================================================================
// Prompt - ユーザー入力プロンプト
// =============================================================================

func TestPromptDurationReturnsCurrentValueOnEmptyInput(t *testing.T) {
	var got time.Duration
	withInput("\n", func() {
		got = PromptDuration("Test", 25*time.Minute, "25m")
	})
	if got != 25*time.Minute {
		t.Errorf("PromptDuration() = %v, want 25m", got)
	}
}

func TestPromptDurationParsesValidInput(t *testing.T) {
	var got time.Duration
	withInput("30m\n", func() {
		got = PromptDuration("Test", 25*time.Minute, "25m")
	})
	if got != 30*time.Minute {
		t.Errorf("PromptDuration() = %v, want 30m", got)
	}
}

func TestPromptDurationParsesComplexDuration(t *testing.T) {
	var got time.Duration
	withInput("1h30m\n", func() {
		got = PromptDuration("Test", 25*time.Minute, "25m")
	})
	if got != 90*time.Minute {
		t.Errorf("PromptDuration() = %v, want 1h30m", got)
	}
}

func TestPromptDurationReturnsCurrentValueOnInvalidInput(t *testing.T) {
	var got time.Duration
	withInput("invalid\n", func() {
		got = PromptDuration("Test", 25*time.Minute, "25m")
	})
	if got != 25*time.Minute {
		t.Errorf("PromptDuration() = %v, want 25m (fallback)", got)
	}
}

func TestPromptIntReturnsCurrentValueOnEmptyInput(t *testing.T) {
	var got int
	withInput("\n", func() {
		got = PromptInt("Test", 4, 4)
	})
	if got != 4 {
		t.Errorf("PromptInt() = %d, want 4", got)
	}
}

func TestPromptIntParsesValidInput(t *testing.T) {
	var got int
	withInput("6\n", func() {
		got = PromptInt("Test", 4, 4)
	})
	if got != 6 {
		t.Errorf("PromptInt() = %d, want 6", got)
	}
}

func TestPromptIntAcceptsZero(t *testing.T) {
	var got int
	withInput("0\n", func() {
		got = PromptInt("Test", 4, 4)
	})
	if got != 0 {
		t.Errorf("PromptInt() = %d, want 0", got)
	}
}

func TestPromptIntReturnsCurrentValueOnInvalidInput(t *testing.T) {
	var got int
	withInput("abc\n", func() {
		got = PromptInt("Test", 4, 4)
	})
	if got != 4 {
		t.Errorf("PromptInt() = %d, want 4 (fallback)", got)
	}
}

func TestPromptBoolReturnsCurrentValueOnEmptyInput(t *testing.T) {
	var got bool
	withInput("\n", func() {
		got = PromptBool("Test", true, true)
	})
	if got != true {
		t.Error("PromptBool() = false, want true")
	}
}

func TestPromptBoolReturnsTrueForYesInput(t *testing.T) {
	inputs := []string{"y\n", "yes\n", "Y\n", "YES\n"}
	for _, input := range inputs {
		var got bool
		withInput(input, func() {
			got = PromptBool("Test", false, true)
		})
		if got != true {
			t.Errorf("PromptBool(%q) = false, want true", input)
		}
	}
}

func TestPromptBoolReturnsFalseForNoInput(t *testing.T) {
	inputs := []string{"n\n", "no\n", "N\n", "NO\n"}
	for _, input := range inputs {
		var got bool
		withInput(input, func() {
			got = PromptBool("Test", true, true)
		})
		if got != false {
			t.Errorf("PromptBool(%q) = true, want false", input)
		}
	}
}

func TestPromptBoolReturnsFalseForUnknownInput(t *testing.T) {
	var got bool
	withInput("other\n", func() {
		got = PromptBool("Test", true, true)
	})
	if got != false {
		t.Error("PromptBool(other) = true, want false")
	}
}

// =============================================================================
// Test Helpers
// =============================================================================

func withInput(input string, fn func()) {
	oldScanner := scanner
	scanner = bufio.NewScanner(strings.NewReader(input))
	defer func() { scanner = oldScanner }()
	fn()
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = w

	fn()

	if err := w.Close(); err != nil {
		t.Errorf("w.Close() error = %v", err)
	}
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Errorf("buf.ReadFrom() error = %v", err)
	}
	return buf.String()
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stderr = w

	fn()

	if err := w.Close(); err != nil {
		t.Errorf("w.Close() error = %v", err)
	}
	os.Stderr = oldStderr

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Errorf("buf.ReadFrom() error = %v", err)
	}
	return buf.String()
}

func assertContains(t *testing.T, output, expected string) {
	t.Helper()
	if !strings.Contains(output, expected) {
		t.Errorf("output missing %q", expected)
	}
}
