package cmd_test

import (
	"errors"
	"os"

	. "github.com/craigmonson/colonize/cmd"
	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/log"
	"github.com/craigmonson/colonize/log_mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var cwd, _ = os.Getwd()

type mockRunner struct {
	callCount int
	origin    string
}

var r = mockRunner{}

func myRun(c *config.Config, l log.Logger, args interface{}) error {
	r.callCount++
	r.origin = c.OriginPath
	return nil
}

func myErrRun(c *config.Config, l log.Logger, args interface{}) error {
	return errors.New("Failed")
}

var _ = Describe("Runner", func() {
	Describe("Run", func() {
		BeforeEach(func() {
			r.callCount = 0
		})
		Context("in normal order", func() {
			It("should make the correct calls", func() {
				c, err := config.LoadConfigInTree(cwd+"/../test", "dev")
				err = Run("TEST", myRun, c, &log_mock.MockLog{}, false, nil)
				Ω(r.callCount).To(Equal(3))
				Ω(err).ToNot(HaveOccurred())
				Ω(r.origin).To(MatchRegexp("../test/microservices"))
			})
		})

		Context("in reverse order", func() {
			It("should make the correct calls", func() {
				c, err := config.LoadConfigInTree(cwd+"/../test", "dev")
				err = Run("TEST", myRun, c, &log_mock.MockLog{}, true, nil)
				Ω(r.callCount).To(Equal(3))
				Ω(err).ToNot(HaveOccurred())
				Ω(r.origin).To(MatchRegexp("../test/vpc"))
			})
		})

		Context("when run returns an error", func() {
			It("should return an error", func() {
				c, err := config.LoadConfigInTree(cwd+"/../test", "dev")
				err = Run("TEST", myErrRun, c, &log_mock.MockLog{}, true, nil)
				Ω(err).To(HaveOccurred())
			})
		})
	})

})
