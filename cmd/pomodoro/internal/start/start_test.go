package start

import (
	"testing"
	"time"

	"pomodoro-cli/internal/config"
	"pomodoro-cli/internal/timer"
	"pomodoro-cli/internal/ui"
)

// =============================================================================
// ShouldAutoStart - 自動開始の判定
// =============================================================================

func TestShouldAutoStartはWorkセッションで設定に従う(t *testing.T) {
	cfg := &config.Config{AutoStartWork: true}
	if !ShouldAutoStart(cfg, timer.SessionWork) {
		t.Error("ShouldAutoStart(Work) = false, want true")
	}

	cfg.AutoStartWork = false
	if ShouldAutoStart(cfg, timer.SessionWork) {
		t.Error("ShouldAutoStart(Work) = true, want false")
	}
}

func TestShouldAutoStartはBreakセッションで設定に従う(t *testing.T) {
	cfg := &config.Config{AutoStartBreaks: true}

	if !ShouldAutoStart(cfg, timer.SessionShortBreak) {
		t.Error("ShouldAutoStart(ShortBreak) = false, want true")
	}
	if !ShouldAutoStart(cfg, timer.SessionLongBreak) {
		t.Error("ShouldAutoStart(LongBreak) = false, want true")
	}

	cfg.AutoStartBreaks = false
	if ShouldAutoStart(cfg, timer.SessionShortBreak) {
		t.Error("ShouldAutoStart(ShortBreak) = true, want false")
	}
}

// =============================================================================
// handleKeyInput - キー入力の処理
// =============================================================================

func TestQキーでタイマーを終了する(t *testing.T) {
	cfg := config.Default()
	tmr := timer.New(cfg)
	tmr.Start(timer.SessionWork)

	shouldExit := handleKeyInput(tmr, cfg, ui.KeyQ)
	if !shouldExit {
		t.Error("handleKeyInput(KeyQ) = false, want true")
	}
}

func TestSpaceキーで一時停止と再開を切り替える(t *testing.T) {
	cfg := config.Default()
	tmr := timer.New(cfg)
	tmr.Start(timer.SessionWork)

	// 一時停止
	handleKeyInput(tmr, cfg, ui.KeySpace)
	if tmr.State().TimerState != timer.StatePaused {
		t.Error("Spaceキー後に一時停止状態にならない")
	}

	// 再開
	handleKeyInput(tmr, cfg, ui.KeySpace)
	if tmr.State().TimerState != timer.StateRunning {
		t.Error("Spaceキー後に実行状態にならない")
	}
}

func TestSキーで次のセッションにスキップする(t *testing.T) {
	cfg := &config.Config{
		WorkDuration:       1 * time.Minute,
		ShortBreakDuration: 1 * time.Minute,
		SessionsUntilLong:  4,
	}
	tmr := timer.New(cfg)
	tmr.Start(timer.SessionWork)

	handleKeyInput(tmr, cfg, ui.KeyS)

	if tmr.State().CurrentSession.Type != timer.SessionShortBreak {
		t.Errorf("スキップ後のセッション = %v, want SessionShortBreak", tmr.State().CurrentSession.Type)
	}
}

func TestRキーで現在のセッションをリセットする(t *testing.T) {
	cfg := &config.Config{WorkDuration: 5 * time.Minute}
	tmr := timer.New(cfg)
	tmr.Start(timer.SessionWork)

	time.Sleep(100 * time.Millisecond)
	handleKeyInput(tmr, cfg, ui.KeyR)

	if tmr.State().CurrentSession.Remaining != cfg.WorkDuration {
		t.Errorf("リセット後の残り時間 = %v, want %v", tmr.State().CurrentSession.Remaining, cfg.WorkDuration)
	}
}

// =============================================================================
// handleSessionComplete - セッション完了時の処理
// =============================================================================

func TestWork完了後に自動開始が有効なら次のBreakが始まる(t *testing.T) {
	cfg := &config.Config{
		WorkDuration:       1 * time.Second,
		ShortBreakDuration: 1 * time.Second,
		SessionsUntilLong:  4,
		AutoStartBreaks:    true,
	}
	tmr := timer.New(cfg)
	tmr.Start(timer.SessionWork)

	time.Sleep(2 * time.Second)

	state := tmr.State()
	if state.TimerState != timer.StateCompleted {
		t.Fatalf("state = %v, want StateCompleted", state.TimerState)
	}

	var prevState timer.TimerState
	handleSessionComplete(tmr, cfg, state, &prevState)

	if tmr.State().CurrentSession.Type != timer.SessionShortBreak {
		t.Errorf("次のセッション = %v, want SessionShortBreak", tmr.State().CurrentSession.Type)
	}
}

func TestWork完了後に自動開始が無効なら待機状態のまま(t *testing.T) {
	cfg := &config.Config{
		WorkDuration:       1 * time.Second,
		ShortBreakDuration: 1 * time.Second,
		SessionsUntilLong:  4,
		AutoStartBreaks:    false,
	}
	tmr := timer.New(cfg)
	tmr.Start(timer.SessionWork)

	time.Sleep(2 * time.Second)

	state := tmr.State()
	if state.TimerState != timer.StateCompleted {
		t.Fatalf("state = %v, want StateCompleted", state.TimerState)
	}

	var prevState timer.TimerState
	handleSessionComplete(tmr, cfg, state, &prevState)

	if tmr.State().TimerState != timer.StateCompleted {
		t.Errorf("state = %v, want StateCompleted（自動開始無効）", tmr.State().TimerState)
	}
}

func TestBreak完了後に自動開始が有効なら次のWorkが始まる(t *testing.T) {
	cfg := &config.Config{
		WorkDuration:       1 * time.Second,
		ShortBreakDuration: 1 * time.Second,
		SessionsUntilLong:  4,
		AutoStartWork:      true,
	}
	tmr := timer.New(cfg)
	tmr.Start(timer.SessionShortBreak)

	time.Sleep(2 * time.Second)

	state := tmr.State()
	if state.TimerState != timer.StateCompleted {
		t.Fatalf("state = %v, want StateCompleted", state.TimerState)
	}

	var prevState timer.TimerState
	handleSessionComplete(tmr, cfg, state, &prevState)

	if tmr.State().CurrentSession.Type != timer.SessionWork {
		t.Errorf("次のセッション = %v, want SessionWork", tmr.State().CurrentSession.Type)
	}
}
