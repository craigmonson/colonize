package log_test

import (
	"io/ioutil"
	"os"

	. "github.com/craigmonson/colonize/log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Log", func() {
	var log Log
	var origStdout, write, read *os.File

	var getStdoutStr = func() string {
		write.Close()
		out, _ := ioutil.ReadAll(read)
		return string(out)
	}

	BeforeEach(func() {
		log = Log{}
		origStdout = os.Stdout
		read, write, _ = os.Pipe()
		os.Stdout = write
	})
	AfterEach(func() {
		os.Stdout = origStdout
	})

	Describe("Log", func() {
		It("should log something", func() {
			log.Log("something")
			out := getStdoutStr()
			Ω(out).To(Equal("something\n"))
		})
	})
})
