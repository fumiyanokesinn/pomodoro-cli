package ui

import (
	"os/exec"
	"runtime"

	"pomodoro-cli/internal/timer"
)

// NotifySessionComplete はセッション完了通知を送信する
func NotifySessionComplete(sessionType timer.SessionType) error {
	return notify("Pomodoro", sessionType.String()+" completed")
}

// notify はシステム通知を送信する
func notify(title, message string) error {
	switch runtime.GOOS {
	case "darwin":
		script := `display notification "` + message + `" with title "` + title + `"`
		return exec.Command("osascript", "-e", script).Run()
	case "linux":
		return exec.Command("notify-send", title, message).Run()
	default:
		print("\a")
		return nil
	}
}

// PlaySound は通知音を再生する
func PlaySound() error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("afplay", "/System/Library/Sounds/Glass.aiff").Run()
	default:
		print("\a")
		return nil
	}
}
