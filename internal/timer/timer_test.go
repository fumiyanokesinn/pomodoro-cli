package timer

import (
	"testing"
	"time"

	"pomodoro-cli/internal/config"
)

// =============================================================================
// New - タイマーの初期化
// =============================================================================

func TestNewはアイドル状態のタイマーを返す(t *testing.T) {
	cfg := config.Default()
	tmr := New(cfg)

	if tmr == nil {
		t.Fatal("New() returned nil")
	}
	state := tmr.State()
	if state.TimerState != StateIdle {
		t.Errorf("initial state = %v, want StateIdle", state.TimerState)
	}
	if state.CompletedWork != 0 {
		t.Errorf("initial CompletedWork = %d, want 0", state.CompletedWork)
	}
}

// =============================================================================
// Start - タイマーの開始
// =============================================================================

func TestStartはタイマーを実行状態にする(t *testing.T) {
	cfg := config.Default()
	cfg.WorkDuration = 1 * time.Hour
	tmr := New(cfg)

	tmr.Start(SessionWork)
	time.Sleep(10 * time.Millisecond)

	state := tmr.State()
	if state.TimerState != StateRunning {
		t.Errorf("after Start() state = %v, want StateRunning", state.TimerState)
	}
	if state.CurrentSession == nil {
		t.Fatal("CurrentSession is nil after Start()")
	}
	if state.CurrentSession.Type != SessionWork {
		t.Errorf("session type = %v, want SessionWork", state.CurrentSession.Type)
	}

	tmr.Stop()
}

// =============================================================================
// Pause/Resume - 一時停止と再開
// =============================================================================

func TestPauseはタイマーを一時停止する(t *testing.T) {
	cfg := config.Default()
	cfg.WorkDuration = 1 * time.Hour
	tmr := New(cfg)

	tmr.Start(SessionWork)
	time.Sleep(10 * time.Millisecond)

	tmr.Pause()
	state := tmr.State()
	if state.TimerState != StatePaused {
		t.Errorf("after Pause() state = %v, want StatePaused", state.TimerState)
	}

	tmr.Stop()
}

func TestResumeはタイマーを再開する(t *testing.T) {
	cfg := config.Default()
	cfg.WorkDuration = 1 * time.Hour
	tmr := New(cfg)

	tmr.Start(SessionWork)
	tmr.Pause()
	tmr.Resume()

	state := tmr.State()
	if state.TimerState != StateRunning {
		t.Errorf("after Resume() state = %v, want StateRunning", state.TimerState)
	}

	tmr.Stop()
}

// =============================================================================
// Stop - タイマーの停止
// =============================================================================

func TestStopはタイマーをアイドル状態に戻す(t *testing.T) {
	cfg := config.Default()
	cfg.WorkDuration = 1 * time.Hour
	tmr := New(cfg)

	tmr.Start(SessionWork)
	time.Sleep(10 * time.Millisecond)

	tmr.Stop()
	state := tmr.State()
	if state.TimerState != StateIdle {
		t.Errorf("after Stop() state = %v, want StateIdle", state.TimerState)
	}
}

// =============================================================================
// Completion - タイマー完了
// =============================================================================

func TestタイマーはWorkセッション完了時にCompletedWorkをインクリメントする(t *testing.T) {
	cfg := config.Default()
	cfg.WorkDuration = 2 * time.Second
	tmr := New(cfg)

	tmr.Start(SessionWork)
	time.Sleep(3 * time.Second)

	state := tmr.State()
	if state.TimerState != StateCompleted {
		t.Errorf("after completion state = %v, want StateCompleted", state.TimerState)
	}
	if state.CompletedWork != 1 {
		t.Errorf("CompletedWork = %d, want 1", state.CompletedWork)
	}
}

// =============================================================================
// getDuration - セッション別の時間取得
// =============================================================================

func TestGetDurationはセッションタイプに応じた時間を返す(t *testing.T) {
	cfg := config.Default()
	cfg.WorkDuration = 25 * time.Minute
	cfg.ShortBreakDuration = 5 * time.Minute
	cfg.LongBreakDuration = 15 * time.Minute
	tmr := New(cfg)

	tests := []struct {
		sessionType SessionType
		expected    time.Duration
	}{
		{SessionWork, 25 * time.Minute},
		{SessionShortBreak, 5 * time.Minute},
		{SessionLongBreak, 15 * time.Minute},
	}

	for _, tt := range tests {
		got := tmr.getDuration(tt.sessionType)
		if got != tt.expected {
			t.Errorf("getDuration(%v) = %v, want %v", tt.sessionType, got, tt.expected)
		}
	}
}

// =============================================================================
// Session Types - 各セッションタイプの動作
// =============================================================================

func TestStartは各セッションタイプを正しく設定する(t *testing.T) {
	cfg := config.Default()
	cfg.WorkDuration = 100 * time.Millisecond
	cfg.ShortBreakDuration = 100 * time.Millisecond
	tmr := New(cfg)

	tmr.Start(SessionWork)
	if tmr.State().CurrentSession.Type != SessionWork {
		t.Errorf("session type = %v, want SessionWork", tmr.State().CurrentSession.Type)
	}
	tmr.Stop()

	tmr.Start(SessionShortBreak)
	if tmr.State().CurrentSession.Type != SessionShortBreak {
		t.Errorf("session type = %v, want SessionShortBreak", tmr.State().CurrentSession.Type)
	}
	tmr.Stop()

	tmr.Start(SessionLongBreak)
	if tmr.State().CurrentSession.Type != SessionLongBreak {
		t.Errorf("session type = %v, want SessionLongBreak", tmr.State().CurrentSession.Type)
	}
	tmr.Stop()
}

// =============================================================================
// NextSessionType - 次のセッションタイプの決定
// =============================================================================

func TestNextSessionTypeはWork後にBreakを返す(t *testing.T) {
	state := &PomodoroState{
		CurrentSession: &Session{Type: SessionWork},
		CompletedWork:  0,
	}

	if got := state.NextSessionType(4); got != SessionShortBreak {
		t.Errorf("after work 1: got %v, want ShortBreak", got)
	}
}

func TestNextSessionTypeは4回目のWork後にLongBreakを返す(t *testing.T) {
	state := &PomodoroState{
		CurrentSession: &Session{Type: SessionWork},
		CompletedWork:  3,
	}

	if got := state.NextSessionType(4); got != SessionLongBreak {
		t.Errorf("after work 4: got %v, want LongBreak", got)
	}
}

func TestNextSessionTypeはBreak後にWorkを返す(t *testing.T) {
	state := &PomodoroState{
		CurrentSession: &Session{Type: SessionShortBreak},
		CompletedWork:  2,
	}

	if got := state.NextSessionType(4); got != SessionWork {
		t.Errorf("after short break: got %v, want Work", got)
	}

	state.CurrentSession.Type = SessionLongBreak
	if got := state.NextSessionType(4); got != SessionWork {
		t.Errorf("after long break: got %v, want Work", got)
	}
}

func TestNextSessionTypeはサイクルを繰り返す(t *testing.T) {
	state := &PomodoroState{
		CurrentSession: &Session{Type: SessionWork},
		CompletedWork:  4,
	}

	// 5回目のWork後 → ShortBreak（サイクルリセット）
	if got := state.NextSessionType(4); got != SessionShortBreak {
		t.Errorf("after work 5: got %v, want ShortBreak", got)
	}

	// 8回目のWork後 → LongBreak
	state.CompletedWork = 7
	if got := state.NextSessionType(4); got != SessionLongBreak {
		t.Errorf("after work 8: got %v, want LongBreak", got)
	}
}

// =============================================================================
// Slow Tests - 時間のかかるテスト（-shortでスキップ）
// =============================================================================

func Test複数セッションでLongBreakが正しく発生する(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test")
	}

	cfg := config.Default()
	cfg.WorkDuration = 1 * time.Second
	cfg.ShortBreakDuration = 1 * time.Second
	cfg.SessionsUntilLong = 2
	tmr := New(cfg)

	// Work 1完了
	tmr.Start(SessionWork)
	time.Sleep(2 * time.Second)
	state := tmr.State()
	if state.CompletedWork != 1 {
		t.Fatalf("after work 1: CompletedWork = %d, want 1", state.CompletedWork)
	}

	// SessionsUntilLong=2なので、1回目の後はLongBreak
	nextType := state.NextSessionType(cfg.SessionsUntilLong)
	if nextType != SessionLongBreak {
		t.Errorf("after work 1 (2-session cycle): nextType = %v, want SessionLongBreak", nextType)
	}

	// Work 2完了
	tmr.Start(SessionWork)
	time.Sleep(2 * time.Second)
	state = tmr.State()
	if state.CompletedWork != 2 {
		t.Fatalf("after work 2: CompletedWork = %d, want 2", state.CompletedWork)
	}

	nextType = state.NextSessionType(cfg.SessionsUntilLong)
	if nextType != SessionShortBreak {
		t.Errorf("after work 2: nextType = %v, want SessionShortBreak", nextType)
	}
}

func Test一時停止中は残り時間が減らない(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test")
	}

	cfg := config.Default()
	cfg.WorkDuration = 2 * time.Second
	tmr := New(cfg)

	tmr.Start(SessionWork)
	time.Sleep(500 * time.Millisecond)

	tmr.Pause()
	state := tmr.State()
	if state.TimerState != StatePaused {
		t.Fatalf("after Pause: state = %v, want StatePaused", state.TimerState)
	}
	remainingAtPause := state.CurrentSession.Remaining

	time.Sleep(500 * time.Millisecond)
	state = tmr.State()
	if state.CurrentSession.Remaining != remainingAtPause {
		t.Error("remaining time changed during pause")
	}

	tmr.Resume()
	time.Sleep(2 * time.Second)

	state = tmr.State()
	if state.TimerState != StateCompleted {
		t.Errorf("after completion: state = %v, want StateCompleted", state.TimerState)
	}
}

func Test残り時間は経過とともに減少する(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test")
	}

	cfg := config.Default()
	cfg.WorkDuration = 5 * time.Second
	tmr := New(cfg)

	tmr.Start(SessionWork)
	time.Sleep(100 * time.Millisecond)

	initialState := tmr.State()
	initialRemaining := initialState.CurrentSession.Remaining

	time.Sleep(2 * time.Second)

	laterState := tmr.State()
	laterRemaining := laterState.CurrentSession.Remaining

	if laterRemaining >= initialRemaining {
		t.Errorf("remaining did not decrease: initial=%v, later=%v", initialRemaining, laterRemaining)
	}

	tmr.Stop()
}

func TestBreakセッションはCompletedWorkをインクリメントしない(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test")
	}

	cfg := config.Default()
	cfg.ShortBreakDuration = 1 * time.Second
	tmr := New(cfg)

	tmr.Start(SessionShortBreak)
	time.Sleep(2 * time.Second)

	state := tmr.State()
	if state.TimerState != StateCompleted {
		t.Errorf("state = %v, want StateCompleted", state.TimerState)
	}
	if state.CompletedWork != 0 {
		t.Errorf("CompletedWork = %d, want 0 (break should not increment)", state.CompletedWork)
	}
}

func TestLongBreakセッションが正しく動作する(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test")
	}

	cfg := config.Default()
	cfg.LongBreakDuration = 1 * time.Second
	tmr := New(cfg)

	tmr.Start(SessionLongBreak)
	time.Sleep(2 * time.Second)

	state := tmr.State()
	if state.TimerState != StateCompleted {
		t.Errorf("state = %v, want StateCompleted", state.TimerState)
	}
	if state.CurrentSession.Type != SessionLongBreak {
		t.Errorf("session type = %v, want SessionLongBreak", state.CurrentSession.Type)
	}
}
