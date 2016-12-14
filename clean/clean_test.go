package clean_test

import (
	. "github.com/craigmonson/colonize/clean"
	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/log_mock"
	um "github.com/craigmonson/colonize/util_mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clean", func() {
	var conf *config.ColonizeConfig
	var mLog *log_mock.MockLog
	var err error

	BeforeEach(func() {
		conf, err = config.LoadConfigInTree("../test/vpc", "dev")
		mLog = &log_mock.MockLog{}
		um.MCmd = &um.MockCmd{}
		um.MockTheCommand()
	})
	AfterEach(func() {
		um.ResetTheCommand()
	})

	Describe("Run", func() {
		Context("Given the proper inputs", func() {
			BeforeEach(func() {
				Run(conf, mLog)
			})
			It("should run the remote config", func() {
				Ω(um.MCmd.Cmd).To(MatchRegexp(conf.CombinedRemoteFilePath))
			})

			It("should remove the files", func() {
				Ω(um.MCmd.Cmd).To(MatchRegexp("rm -f "))
			})

			It("should log some stuff", func() {
				Ω(mLog.Output).To(MatchRegexp("Cleaning up"))
				Ω(mLog.Output).To(MatchRegexp("rm -f "))
			})
		})
	})
})
