package log

import (
	"fmt"

	"github.com/fatih/color"
)

type Logger interface {
	Log(string)
	Print(string)
	LogPretty(string, ...color.Attribute)
}

type Log struct {
	Logger
}

func (l Log) Log(s string) {
	fmt.Println(s)
}

func (l Log) Print(s string) {
	fmt.Print(s)
}

func (l Log) LogPretty(s string, p ...color.Attribute) {
	color.Set(p...)
	l.Log(s)
	color.Unset()
}
