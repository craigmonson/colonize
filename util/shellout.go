package util

import (
	"os"
	"os/exec"
)

type Cmder interface {
	Run() error
}

type Cmd struct {
	Cmder
}

var NewCmd = func(cmd string, arg ...string) Cmder {
	newCmd := exec.Command(cmd, arg...)
	newCmd.Stdout = os.Stdout
	newCmd.Stderr = os.Stderr
	return newCmd
}

func RunCmd(cmd string, arg ...string) error {
	newCmd := NewCmd(cmd, arg...)
	return newCmd.Run()
}
