package destroy_test

import (
	"github.com/craigmonson/colonize/config"
	. "github.com/craigmonson/colonize/destroy"
	"github.com/craigmonson/colonize/log_mock"
	um "github.com/craigmonson/colonize/util_mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Plan", func() {
	var conf *config.Config
	var mLog *log_mock.MockLog
	var err error
	var args = RunArgs{SkipRemote: false}

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
				Run(conf, mLog, args)
			})
			It("should run the remote config", func() {
				Ω(um.MCmd.Cmd).To(MatchRegexp(conf.CombinedRemoteFilePath))
			})

			It("should run the proper destroy command", func() {
				Ω(um.MCmd.Cmd).To(MatchRegexp(
					"terraform destroy -force " +
						"-var-file " + conf.CombinedValsFilePath +
						" -var-file " + conf.CombinedDerivedValsFilePath,
				))
			})

			It("should log some stuff", func() {
				Ω(mLog.Output).To(MatchRegexp("Running remote setup"))
				Ω(mLog.Output).To(MatchRegexp("Executing terraform destroy"))
			})
		})

		Context("when skipRemote is true", func() {
			BeforeEach(func() {
				args.SkipRemote = true
				Run(conf, mLog, args)
			})

			It("should NOT run the remote config", func() {
				Ω(um.MCmd.Cmd).ToNot(MatchRegexp(conf.CombinedRemoteFilePath))
			})

			It("should log that it's skipping the remote setup", func() {
				Ω(mLog.Output).ToNot(MatchRegexp("Running remote setup"))
				Ω(mLog.Output).To(MatchRegexp("Skipping remote setup"))
				Ω(mLog.Output).To(MatchRegexp("Executing terraform destroy"))
			})

			It("should run terraform destroy", func() {
				Ω(um.MCmd.Cmd).To(MatchRegexp("terraform destroy -force "))
			})
		})
	})
})
