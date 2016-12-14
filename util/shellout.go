package util

import (
	"os/exec"
)

type Cmder interface {
	Run() error
}

type Cmd struct {
	Cmder
}

var NewCmd = func(cmd string, arg ...string) Cmder {
	return exec.Command(cmd, arg...)
}

func RunCmd(cmd string, arg ...string) error {
	newCmd := NewCmd(cmd, arg...)
	return newCmd.Run()
}
