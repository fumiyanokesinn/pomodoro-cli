package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"pomodoro-cli/internal/config"
	"pomodoro-cli/internal/timer"
)

// ----------------------------------------------------------------------------
// 基本出力関数
// ----------------------------------------------------------------------------

// printLine はrawモード対応の出力（stdout + \r\n）
func printLine(s string) {
	fmt.Print(s + "\r\n")
}

// ----------------------------------------------------------------------------
// rawモード前に使用（stderr + \n）
// flag.PrintDefaults()と合わせるためstderrに出力
// ----------------------------------------------------------------------------

// ShowUsage はヘルプメッセージを表示する
func ShowUsage() {
	fmt.Fprintln(os.Stderr, "Usage: pomodoro [options] [command]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  start              Start pomodoro timer (default)")
	fmt.Fprintln(os.Stderr, "  config             Show current configuration")
	fmt.Fprintln(os.Stderr, "  init               Create default config file")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Options:")
	fmt.Fprintln(os.Stderr, "  -w, --work         Work duration (e.g., -w 25m)")
	fmt.Fprintln(os.Stderr, "  -s, --short-break  Short break duration (e.g., -s 5m)")
	fmt.Fprintln(os.Stderr, "  -l, --long-break   Long break duration (e.g., -l 15m)")
	fmt.Fprintln(os.Stderr, "  -n, --sessions     Sessions until long break (e.g., -n 4)")
	fmt.Fprintln(os.Stderr, "      --no-sound     Disable notification sound")
	fmt.Fprintln(os.Stderr, "      --no-notify    Disable system notifications")
	fmt.Fprintln(os.Stderr, "      --no-auto-break  Disable auto-start breaks")
	fmt.Fprintln(os.Stderr, "      --no-auto-work   Disable auto-start work")
	fmt.Fprintln(os.Stderr, "  -v, --version      Show version")
	fmt.Fprintln(os.Stderr, "  -h, --help         Show help")
}

// ShowUnknownCommand は不明なコマンドのエラーを表示する
func ShowUnknownCommand(command string) {
	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
}

// ShowError はエラーメッセージを表示する
func ShowError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

// ----------------------------------------------------------------------------
// rawモード前に使用（stdout + \n）
// ----------------------------------------------------------------------------

// ShowVersion はバージョンを表示する
func ShowVersion(version string) {
	fmt.Printf("pomodoro version %s\n", version)
}

// ShowConfig は設定を表示する
func ShowConfig(cfg *config.Config) {
	fmt.Println()
	fmt.Println("  ┌─────────────────────────────────────────────┐")
	fmt.Println("  │         CURRENT CONFIGURATION               │")
	fmt.Println("  ├─────────────────────────────────────────────┤")
	fmt.Println("  │  Timing                                     │")
	fmt.Printf("  │    Work duration:      %-20v│\n", cfg.WorkDuration)
	fmt.Printf("  │    Short break:        %-20v│\n", cfg.ShortBreakDuration)
	fmt.Printf("  │    Long break:         %-20v│\n", cfg.LongBreakDuration)
	fmt.Printf("  │    Sessions until long: %-19d│\n", cfg.SessionsUntilLong)
	fmt.Println("  ├─────────────────────────────────────────────┤")
	fmt.Println("  │  Behavior                                   │")
	fmt.Printf("  │    Auto-start breaks:  %-20v│\n", boolToYesNo(cfg.AutoStartBreaks))
	fmt.Printf("  │    Auto-start work:    %-20v│\n", boolToYesNo(cfg.AutoStartWork))
	fmt.Println("  ├─────────────────────────────────────────────┤")
	fmt.Println("  │  Notifications                              │")
	fmt.Printf("  │    Sound enabled:      %-20v│\n", boolToYesNo(cfg.SoundEnabled))
	fmt.Printf("  │    Notify enabled:     %-20v│\n", boolToYesNo(cfg.NotifyEnabled))
	fmt.Println("  └─────────────────────────────────────────────┘")
}

// boolToYesNo はboolをYes/Noに変換する
func boolToYesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

// ShowConfigCreated は設定ファイル作成成功メッセージを表示する
func ShowConfigCreated(path string) {
	fmt.Println()
	fmt.Println("  ╔═════════════════════════════════════════════╗")
	fmt.Println("  ║  ✓ Configuration saved successfully!        ║")
	fmt.Println("  ╚═════════════════════════════════════════════╝")
	fmt.Printf("  File: %s\n", path)
	fmt.Println()
}

// ShowInitHeader はinit開始時のヘッダーを表示する
func ShowInitHeader() {
	fmt.Println()
	fmt.Println("  ╔═════════════════════════════════════════════╗")
	fmt.Println("  ║       POMODORO CONFIGURATION SETUP          ║")
	fmt.Println("  ╠═════════════════════════════════════════════╣")
	fmt.Println("  ║  Press Enter to keep current values.        ║")
	fmt.Println("  ╚═════════════════════════════════════════════╝")
	fmt.Println()
}

// ----------------------------------------------------------------------------
// 対話的プロンプト関数（init用）
// ----------------------------------------------------------------------------

var scanner = bufio.NewScanner(os.Stdin)

// GetScanner は現在のscannerを返す（テスト用）
func GetScanner() *bufio.Scanner {
	return scanner
}

// SetScanner はscannerを設定する（テスト用）
func SetScanner(s *bufio.Scanner) {
	scanner = s
}

// PromptDuration は時間を入力させる
func PromptDuration(label string, currentVal time.Duration, example string) time.Duration {
	fmt.Printf("%s [current: %s] (e.g. %s): ", label, FormatDuration(currentVal), example)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return currentVal
	}
	d, err := time.ParseDuration(input)
	if err != nil {
		fmt.Println("  Invalid format, using current value")
		return currentVal
	}
	return d
}

// PromptInt は整数を入力させる
func PromptInt(label string, currentVal int, example int) int {
	fmt.Printf("%s [current: %d] (e.g. %d): ", label, currentVal, example)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return currentVal
	}
	n, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("  Invalid number, using current value")
		return currentVal
	}
	return n
}

// PromptBool はyes/noを入力させる
func PromptBool(label string, currentVal bool, example bool) bool {
	currentStr := "n"
	if currentVal {
		currentStr = "y"
	}
	exampleStr := "n"
	if example {
		exampleStr = "y"
	}
	fmt.Printf("%s [current: %s] (e.g. %s): ", label, currentStr, exampleStr)
	scanner.Scan()
	input := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if input == "" {
		return currentVal
	}
	return input == "y" || input == "yes"
}

// ----------------------------------------------------------------------------
// rawモード中に使用（stdout + \r\n）
// ----------------------------------------------------------------------------

// RenderTimer はタイマーの状態を表示する
func RenderTimer(session *timer.Session, state timer.TimerState) {
	if session == nil {
		return
	}

	remaining := session.Remaining
	minutes := int(remaining.Minutes())
	seconds := int(remaining.Seconds()) % 60

	progress := 1.0 - (float64(remaining) / float64(session.Duration))
	bar := progressBar(progress, 30)

	stateStr := "▶"
	if state == timer.StatePaused {
		stateStr = "⏸"
	}

	fmt.Printf("\r%s %s [%s] %02d:%02d", stateStr, session.Type.String(), bar, minutes, seconds)
}

// ShowWelcome はウェルカムメッセージを表示する
func ShowWelcome(work, shortBreak, longBreak time.Duration) {
	printLine("")
	printLine("  ╔══════════════════════════════════════════════════════════════════════════╗")
	printLine("  ║                                                                          ║")
	printLine("  ║   ██████╗  ██████╗ ███╗   ███╗ ██████╗ ██████╗  ██████╗ ██████╗  ██████╗ ║")
	printLine("  ║   ██╔══██╗██╔═══██╗████╗ ████║██╔═══██╗██╔══██╗██╔═══██╗██╔══██╗██╔═══██╗║")
	printLine("  ║   ██████╔╝██║   ██║██╔████╔██║██║   ██║██║  ██║██║   ██║██████╔╝██║   ██║║")
	printLine("  ║   ██╔═══╝ ██║   ██║██║╚██╔╝██║██║   ██║██║  ██║██║   ██║██╔══██╗██║   ██║║")
	printLine("  ║   ██║     ╚██████╔╝██║ ╚═╝ ██║╚██████╔╝██████╔╝╚██████╔╝██║  ██║╚██████╔╝║")
	printLine("  ║   ╚═╝      ╚═════╝ ╚═╝     ╚═╝ ╚═════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ║")
	printLine("  ║                                                                          ║")
	printLine("  ╚══════════════════════════════════════════════════════════════════════════╝")
	printLine("")
	printLine("  ┌────────────────────────────────────────────────────────────────────────┐")
	printLine(fmt.Sprintf("  │  Work: %-10v   Short Break: %-10v   Long Break: %-10v │", work, shortBreak, longBreak))
	printLine("  └────────────────────────────────────────────────────────────────────────┘")
	printLine("")
	printLine("  ┌─ Keyboard Shortcuts ───────────────────────────────────────────────────┐")
	printLine("  │  [Space] Pause/Resume    [s] Skip    [r] Reset    [q] Quit             │")
	printLine("  └────────────────────────────────────────────────────────────────────────┘")
	printLine("")
}

// ShowSessionComplete はセッション完了メッセージを表示する
func ShowSessionComplete(sessionType timer.SessionType) {
	printLine("")
	switch sessionType {
	case timer.SessionWork:
		printLine("  ╔════════════════════════════════════════════╗")
		printLine("  ║  ✓ Work session complete!                  ║")
		printLine("  ║    Time for a well-deserved break.         ║")
		printLine("  ╚════════════════════════════════════════════╝")
	case timer.SessionShortBreak, timer.SessionLongBreak:
		printLine("  ╔════════════════════════════════════════════╗")
		printLine("  ║  ✓ Break over!                             ║")
		printLine("  ║    Time to get back to work.               ║")
		printLine("  ╚════════════════════════════════════════════╝")
	}
}

// ShowStartSession はセッション開始メッセージを表示する
func ShowStartSession(sessionType timer.SessionType) {
	printLine("")
	printLine(fmt.Sprintf("  >>> Starting %s...", sessionType.String()))
}

// ShowPaused は一時停止メッセージを表示する
func ShowPaused() {
	printLine("")
	printLine("  || Paused")
}

// ShowResumed は再開メッセージを表示する
func ShowResumed() {
	printLine("")
	printLine("  >> Resumed")
}

// ShowSkipped はスキップメッセージを表示する
func ShowSkipped(nextType timer.SessionType) {
	printLine("")
	printLine(fmt.Sprintf("  >> Skipped. Starting %s...", nextType.String()))
}

// ShowReset はリセットメッセージを表示する
func ShowReset() {
	printLine("")
	printLine("  <> Reset")
}

// ShowExit は終了メッセージを表示する
func ShowExit() {
	printLine("")
	printLine("  ┌──────────────────────────────┐")
	printLine("  │  Thanks for using Pomodoro!  │")
	printLine("  │  See you next time!          │")
	printLine("  └──────────────────────────────┘")
}

// ----------------------------------------------------------------------------
// ヘルパー関数
// ----------------------------------------------------------------------------

// progressBar はプログレスバーを生成する
func progressBar(progress float64, width int) string {
	filled := int(progress * float64(width))
	empty := width - filled
	return strings.Repeat("█", filled) + strings.Repeat("░", empty)
}

// FormatDuration は時間を人間が読みやすい形式に変換する
func FormatDuration(d time.Duration) string {
	m := int(d.Minutes())
	return fmt.Sprintf("%dm", m)
}
