package util_mock

import (
	"strings"

	"github.com/craigmonson/colonize/util"
)

type MockCmd struct {
	util.Cmder

	CallCount int
	Cmd       string
}

func (c *MockCmd) Run() error {
	c.CallCount++
	return nil
}

var MCmd = &MockCmd{}
var OrigCmd func(string, ...string) util.Cmder

var mockedNewCmd = func(cmd string, args ...string) util.Cmder {
	MCmd.Cmd = MCmd.Cmd + "\n" + cmd + " " + strings.Join(args, " ")

	return MCmd
}

func MockTheCommand() {
	OrigCmd = util.NewCmd
	util.NewCmd = mockedNewCmd
}

func ResetTheCommand() {
	util.NewCmd = OrigCmd
}
