package start

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pomodoro-cli/internal/config"
	"pomodoro-cli/internal/timer"
	"pomodoro-cli/internal/ui"
)

// Run はタイマーを実行する
func Run(cfg *config.Config) error {
	if err := ui.InitInput(); err != nil {
		return fmt.Errorf("failed to initialize input: %w", err)
	}
	defer ui.RestoreInput()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	t := timer.New(cfg)
	var prevState timer.TimerState

	ui.ShowWelcome(cfg.WorkDuration, cfg.ShortBreakDuration, cfg.LongBreakDuration)
	t.Start(timer.SessionWork)
	ui.ShowStartSession(timer.SessionWork)

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			ui.ShowExit()
			return nil
		case key := <-ui.KeyChan():
			if handleKeyInput(t, cfg, key) {
				return nil
			}
			state := t.State()
			ui.RenderTimer(state.CurrentSession, state.TimerState)
		case <-ticker.C:
			state := t.State()
			ui.RenderTimer(state.CurrentSession, state.TimerState)
			handleSessionComplete(t, cfg, state, &prevState)
		}
	}
}

// handleKeyInput はキー入力を処理する（終了時 true を返す）
func handleKeyInput(t *timer.Timer, cfg *config.Config, key ui.KeyEvent) bool {
	state := t.State()
	switch key {
	case ui.KeySpace:
		switch state.TimerState {
		case timer.StateRunning:
			t.Pause()
			ui.ShowPaused()
		case timer.StatePaused:
			t.Resume()
			ui.ShowResumed()
		}
	case ui.KeyQ:
		ui.ShowExit()
		return true
	case ui.KeyS:
		t.Stop()
		nextType := t.State().NextSessionType(cfg.SessionsUntilLong)
		t.Start(nextType)
		ui.ShowSkipped(nextType)
	case ui.KeyR:
		if state.CurrentSession != nil {
			sessionType := state.CurrentSession.Type
			t.Stop()
			t.Start(sessionType)
			ui.ShowReset()
		}
	}
	return false
}

// handleSessionComplete はセッション完了時の処理を行う
func handleSessionComplete(t *timer.Timer, cfg *config.Config, state *timer.PomodoroState, prevState *timer.TimerState) {
	if state.TimerState != timer.StateCompleted || *prevState == timer.StateCompleted {
		*prevState = state.TimerState
		return
	}

	ui.ShowSessionComplete(state.CurrentSession.Type)
	if cfg.NotifyEnabled {
		if err := ui.NotifySessionComplete(state.CurrentSession.Type); err != nil {
			ui.ShowError("Notification failed: " + err.Error())
		}
	}
	if cfg.SoundEnabled {
		if err := ui.PlaySound(); err != nil {
			ui.ShowError("Sound playback failed: " + err.Error())
		}
	}

	nextType := state.NextSessionType(cfg.SessionsUntilLong)
	if ShouldAutoStart(cfg, nextType) {
		t.Start(nextType)
		ui.ShowStartSession(nextType)
	}
	*prevState = state.TimerState
}

// ShouldAutoStart は自動開始すべきかを判定する
func ShouldAutoStart(cfg *config.Config, nextType timer.SessionType) bool {
	if nextType == timer.SessionWork {
		return cfg.AutoStartWork
	}
	return cfg.AutoStartBreaks
}
