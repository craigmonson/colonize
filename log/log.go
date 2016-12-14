package log

import (
	"fmt"
)

type Logger interface {
	Log(string)
}

type Log struct {
	Logger
}

func (l Log) Log(s string) {
	fmt.Println(s)
}
