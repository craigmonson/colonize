package util

import (
	"os/exec"
)

type Cmder interface {
	Run() error
        CombinedOutput() ([]byte, error)
}

type Cmd struct {
	Cmder
}

var NewCmd = func(cmd string, arg ...string) Cmder {
	return exec.Command(cmd, arg...)
}

func RunCmd(cmd string, arg ...string) (error,string) {
	newCmd := NewCmd(cmd, arg...)
        output, err := newCmd.CombinedOutput()

        return err,string(output)
}
