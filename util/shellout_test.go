package util_test

import (
	//"fmt"
	"os/exec"

	. "github.com/craigmonson/colonize/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockCmd struct {
	Cmder

	callCount int
}

func (c *mockCmd) Run() error {
	c.callCount++
	return nil
}

// NOTE: this is actually testing the 'test' directory in the root of the
// project.
var _ = Describe("util/shellout", func() {
	Describe("NewCmd", func() {
		var cmd Cmder

		BeforeEach(func() {
			cmd = NewCmd("ls", "-l")
		})

		It("should return a command type", func() {
			Ω(cmd).To(BeAssignableToTypeOf(&exec.Cmd{}))
		})

		It("should set the path", func() {
			cmd := cmd.(*exec.Cmd)
			Ω(cmd.Path).To(Equal("/bin/ls"))
		})

		It("should set the args", func() {
			cmd := cmd.(*exec.Cmd)
			Ω(cmd.Args).To(Equal([]string{"ls", "-l"}))
		})
	})

	Describe("RunCmd", func() {
		var origFunc func(string, ...string) Cmder
		var err error
		var mocked *mockCmd

		BeforeEach(func() {
			mocked = &mockCmd{}

			origFunc = NewCmd
			NewCmd = func(cmd string, arg ...string) Cmder {
				return mocked
			}

			err = RunCmd("ls", "-l")

		})
		AfterEach(func() {
			NewCmd = origFunc
		})

		It("should not cause an error", func() {
			Ω(err).ToNot(HaveOccurred())
		})

		It("should have called Run", func() {
			Ω(mocked.callCount).To(Equal(1))
		})
	})
})
