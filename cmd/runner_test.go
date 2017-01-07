package cmd_test

import (
	"errors"

	. "github.com/craigmonson/colonize/cmd"
	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
				c, err := config.LoadConfigInTree("../test", "dev")
				err = Run(myRun, c, Log, false, nil)
				Ω(r.callCount).To(Equal(2))
				Ω(err).ToNot(HaveOccurred())
				Ω(r.origin).To(Equal("../test/microservices"))
			})
		})

		Context("in reverse order", func() {
			It("should make the correct calls", func() {
				c, err := config.LoadConfigInTree("../test", "dev")
				err = Run(myRun, c, Log, true, nil)
				Ω(r.callCount).To(Equal(2))
				Ω(err).ToNot(HaveOccurred())
				Ω(r.origin).To(Equal("../test/vpc"))
			})
		})

		Context("when run returns an error", func() {
			It("should return an error", func() {
				c, err := config.LoadConfigInTree("../test", "dev")
				err = Run(myErrRun, c, Log, true, nil)
				Ω(err).To(HaveOccurred())
			})
		})
	})

})
