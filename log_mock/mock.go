package log_mock

import (
	"github.com/craigmonson/colonize/log"
	"github.com/fatih/color"
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

func (l *MockLog) LogPretty(s string, p ...color.Attribute) {
	l.Log(s)
}

func (l *MockLog) PrintPretty(s string, p ...color.Attribute) {
	l.Print(s)
}
