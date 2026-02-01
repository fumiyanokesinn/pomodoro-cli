package timer

import "time"

// SessionType はポモドーロセッションの種類を表す
type SessionType int

const (
	SessionWork SessionType = iota
	SessionShortBreak
	SessionLongBreak
)

// String はセッション種類の名前を返す
func (s SessionType) String() string {
	switch s {
	case SessionWork:
		return "Work"
	case SessionShortBreak:
		return "Short Break"
	case SessionLongBreak:
		return "Long Break"
	default:
		return "Unknown"
	}
}

// TimerState はタイマーの状態を表す
type TimerState int

const (
	StateIdle TimerState = iota
	StateRunning
	StatePaused
	StateCompleted
)

// Session は1つのポモドーロセッションを表す
type Session struct {
	Type      SessionType
	Duration  time.Duration
	Remaining time.Duration
	StartedAt time.Time
	PausedAt  time.Time
}

// PomodoroState はポモドーロ全体の進行状態を追跡する
type PomodoroState struct {
	CurrentSession *Session
	CompletedWork  int // 完了した作業セッション数
	TimerState     TimerState
}

// NextSessionType は次に来るべきセッション種類を決定する
func (p *PomodoroState) NextSessionType(sessionsUntilLong int) SessionType {
	if p.CurrentSession == nil || p.CurrentSession.Type != SessionWork {
		return SessionWork
	}
	// 作業後は休憩を決定
	if (p.CompletedWork+1)%sessionsUntilLong == 0 {
		return SessionLongBreak
	}
	return SessionShortBreak
}
