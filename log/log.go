package log

import (
	"fmt"
)

type Logger interface {
	Log(string)
}

type Log struct{}

func (l *Log) Log(s string) {
	fmt.Println(s)
}
