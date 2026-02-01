package timer

import "testing"

func TestSessionTypeString(t *testing.T) {
	tests := []struct {
		input    SessionType
		expected string
	}{
		{SessionWork, "Work"},
		{SessionShortBreak, "Short Break"},
		{SessionLongBreak, "Long Break"},
		{SessionType(99), "Unknown"},
	}
	for _, tt := range tests {
		if got := tt.input.String(); got != tt.expected {
			t.Errorf("SessionType(%d).String() = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestNextSessionType(t *testing.T) {
	tests := []struct {
		name              string
		currentType       SessionType
		completedWork     int
		sessionsUntilLong int
		expected          SessionType
	}{
		{"nil session returns Work", SessionWork, 0, 4, SessionWork},
		{"after work, session 1 of 4", SessionWork, 0, 4, SessionShortBreak},
		{"after work, session 3 of 4", SessionWork, 2, 4, SessionShortBreak},
		{"after work, session 4 of 4 (long break)", SessionWork, 3, 4, SessionLongBreak},
		{"after short break", SessionShortBreak, 1, 4, SessionWork},
		{"after long break", SessionLongBreak, 4, 4, SessionWork},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &PomodoroState{
				CurrentSession: &Session{Type: tt.currentType},
				CompletedWork:  tt.completedWork,
			}
			if tt.name == "nil session returns Work" {
				state.CurrentSession = nil
			}
			got := state.NextSessionType(tt.sessionsUntilLong)
			if got != tt.expected {
				t.Errorf("NextSessionType() = %v, want %v", got, tt.expected)
			}
		})
	}
}
