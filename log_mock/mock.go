package log_mock

import (
	"github.com/craigmonson/colonize/log"
)

type MockLog struct {
	log.Logger

	Output string
}

func (l *MockLog) Log(s string) {
	l.Output = l.Output + "\n" + s
}

func (l *MockLog) Print(s string) {
	l.Output = l.Output + s
}
