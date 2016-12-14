package apply_test

import (
	. "github.com/craigmonson/colonize/apply"
	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/log_mock"
	um "github.com/craigmonson/colonize/util_mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Plan", func() {
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
				Run(conf, mLog, false, false)
			})
			It("should run the remote config", func() {
				Ω(um.MCmd.Cmd).To(MatchRegexp(conf.CombinedRemoteFilePath))
			})

			It("should run the proper apply command", func() {
				Ω(um.MCmd.Cmd).To(MatchRegexp(
					"terraform apply -parallelism 1 " +
						"-var-file " + conf.CombinedValsFilePath +
						" -var-file " + conf.CombinedDerivedValsFilePath,
				))
			})

			It("should log some stuff", func() {
				Ω(mLog.Output).To(MatchRegexp("Running remote setup"))
				Ω(mLog.Output).To(MatchRegexp("Executing terraform apply"))
			})
		})

		Context("when skipRemote is true", func() {
			Context("when remoteAfterApply is false", func() {
				BeforeEach(func() {
					Run(conf, mLog, true, false)
				})

				It("should NOT run the remote config", func() {
					Ω(um.MCmd.Cmd).ToNot(MatchRegexp(conf.CombinedRemoteFilePath))
				})

				It("should log that it's skipping the remote setup", func() {
					Ω(mLog.Output).ToNot(MatchRegexp("Running remote setup"))
					Ω(mLog.Output).To(MatchRegexp("Skipping remote setup"))
					Ω(mLog.Output).To(MatchRegexp("Executing terraform apply"))
					Ω(mLog.Output).To(MatchRegexp("REMOTE NOT SYNC'D!"))
				})

				It("should run terraform apply", func() {
					Ω(um.MCmd.Cmd).To(MatchRegexp("terraform apply -parallelism"))
				})
			})

			Context("when remoteAfterApply is true", func() {
				BeforeEach(func() {
					Run(conf, mLog, true, true)
				})

				It("should run the remote config", func() {
					Ω(um.MCmd.Cmd).To(MatchRegexp(conf.CombinedRemoteFilePath))
				})

				It("should log that it's skipping the remote setup", func() {
					Ω(mLog.Output).ToNot(MatchRegexp("Running remote setup"))
					Ω(mLog.Output).To(MatchRegexp("Skipping remote setup"))
					Ω(mLog.Output).To(MatchRegexp("Executing terraform apply"))
				})

				It("should run terraform apply", func() {
					Ω(um.MCmd.Cmd).To(MatchRegexp("terraform apply -parallelism"))
				})
			})
		})
	})
})
