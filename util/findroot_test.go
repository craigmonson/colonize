package util_test

import (
	"github.com/craigmonson/colonize/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// NOTE: this is actually testing the 'test' directory
var _ = Describe("Util/Rootfind", func() {
	var err error
	path := "../test/"

	Describe("func FindRoot", func() {
		Context("given a path that contains a "+util.RootFile, func() {
			var testPath string

			BeforeEach(func() {
				testPath, err = util.FindRoot(path)
			})

			It("should return the full path to the config file", func() {
				Ω(testPath).To(Equal(path + util.RootFile))
			})

			It("should not return an error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

		Context("given a path that doesn't have a "+util.RootFile, func() {
			var testPath string

			BeforeEach(func() {
				testPath, err = util.FindRoot("/etc")
			})

			It("should return empty string.", func() {
				Ω(testPath).To(Equal(""))
			})

			It("should return an error", func() {
				Ω(err).Should(MatchError(util.RootFile + " not found in the directory tree."))
			})
		})
	})
})
