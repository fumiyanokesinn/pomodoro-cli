package ui

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// KeyEvent はキーボードイベントを表す
type KeyEvent int

const (
	KeyNone KeyEvent = iota
	KeySpace
	KeyQ
	KeyS
	KeyR
	KeyUnknown
)

// 入力用のグローバル状態（パッケージ内で管理）
var (
	oldTermState *term.State
	keyChan      chan KeyEvent
)

// InitInput はターミナルをrawモードに設定し、キー入力の監視を開始する
func InitInput() error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	oldTermState = oldState
	keyChan = make(chan KeyEvent, 1)

	go readLoop()
	return nil
}

// RestoreInput はターミナルを元の状態に戻す
func RestoreInput() {
	if oldTermState != nil {
		if err := term.Restore(int(os.Stdin.Fd()), oldTermState); err != nil {
			fmt.Fprintf(os.Stderr, "failed to restore terminal: %v\n", err)
		}
	}
}

// ReadKey はキー入力を読み取る（ノンブロッキング）
func ReadKey() KeyEvent {
	select {
	case key := <-keyChan:
		return key
	default:
		return KeyNone
	}
}

// KeyChan はキーイベントのチャンネルを返す（select文で使用）
func KeyChan() <-chan KeyEvent {
	return keyChan
}

// readLoop はバックグラウンドでキー入力を読み取る
func readLoop() {
	buf := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			continue
		}

		var key KeyEvent
		switch buf[0] {
		case ' ':
			key = KeySpace
		case 'q', 'Q', 3: // 3 = Ctrl+C
			key = KeyQ
		case 's', 'S':
			key = KeyS
		case 'r', 'R':
			key = KeyR
		default:
			key = KeyUnknown
		}

		select {
		case keyChan <- key:
		default:
		}
	}
}
