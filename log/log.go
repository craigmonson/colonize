package log

import (
	"fmt"
)

type Logger interface {
	Log(string)
	Print(string)
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
