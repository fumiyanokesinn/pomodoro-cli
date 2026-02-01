package timer

import (
	"context"
	"sync"
	"time"

	"pomodoro-cli/internal/config"
)

// Timer はポモドーロタイマーを管理する
type Timer struct {
	config *config.Config
	state  *PomodoroState
	mu     sync.RWMutex
	cancel context.CancelFunc
}

// New は新しいTimerを作成する
func New(cfg *config.Config) *Timer {
	return &Timer{
		config: cfg,
		state: &PomodoroState{
			TimerState: StateIdle,
		},
	}
}

// Start は指定された種類のセッションを開始する
func (t *Timer) Start(sessionType SessionType) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 既存のタイマーを停止
	if t.cancel != nil {
		t.cancel()
	}

	duration := t.getDuration(sessionType)
	t.state.CurrentSession = &Session{
		Type:      sessionType,
		Duration:  duration,
		Remaining: duration,
		StartedAt: time.Now(),
	}
	t.state.TimerState = StateRunning

	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	go t.run(ctx)
}

// Pause は現在のセッションを一時停止する
func (t *Timer) Pause() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state.TimerState == StateRunning {
		t.state.TimerState = StatePaused
		t.state.CurrentSession.PausedAt = time.Now()
	}
}

// Resume は一時停止中のセッションを再開する
func (t *Timer) Resume() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state.TimerState == StatePaused {
		t.state.TimerState = StateRunning
	}
}

// Stop は現在のセッションを停止する
func (t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.cancel != nil {
		t.cancel()
		t.cancel = nil
	}
	t.state.TimerState = StateIdle
}

// State は現在の状態のコピーを返す（スレッドセーフ）
func (t *Timer) State() *PomodoroState {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// コピーを返してレースコンディションを防ぐ
	stateCopy := &PomodoroState{
		CompletedWork: t.state.CompletedWork,
		TimerState:    t.state.TimerState,
	}
	if t.state.CurrentSession != nil {
		sessionCopy := *t.state.CurrentSession
		stateCopy.CurrentSession = &sessionCopy
	}
	return stateCopy
}

// run はタイマーのカウントダウンを実行するgoroutine
func (t *Timer) run(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.mu.Lock()
			if t.state.TimerState == StateRunning {
				t.state.CurrentSession.Remaining -= time.Second
				if t.state.CurrentSession.Remaining <= 0 {
					t.complete()
					t.mu.Unlock()
					return
				}
			}
			t.mu.Unlock()
		}
	}
}

// complete はセッション完了時の処理を行う
func (t *Timer) complete() {
	if t.state.CurrentSession.Type == SessionWork {
		t.state.CompletedWork++
	}
	t.state.TimerState = StateCompleted
}

// getDuration はセッション種類に応じた時間を返す
func (t *Timer) getDuration(sessionType SessionType) time.Duration {
	switch sessionType {
	case SessionWork:
		return t.config.WorkDuration
	case SessionShortBreak:
		return t.config.ShortBreakDuration
	case SessionLongBreak:
		return t.config.LongBreakDuration
	default:
		return t.config.WorkDuration
	}
}
