package prep_test

import (
	"io/ioutil"
	"os"

	"github.com/craigmonson/colonize/config"
	. "github.com/craigmonson/colonize/prep"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Prep", func() {
	var conf *config.ColonizeConfig
	var err error

	BeforeEach(func() {
		conf, err = config.LoadConfigInTree("../test/vpc", "dev")
	})
	AfterEach(func() {
		// remove files.
		_ = os.Remove(conf.CombinedValsFilePath)
	})

	PDescribe("Run", func() {})

	Describe("BuildCombinedValuesFile", func() {
		BeforeEach(func() {
			err = BuildCombinedValuesFile(conf)
		})

		It("should create the combined values file", func() {
			Ω(conf.CombinedValsFilePath).To(BeARegularFile())
		})

		It("should have the right contents", func() {
			contents, _ := ioutil.ReadFile(conf.CombinedValsFilePath)
			expected := "root_var = \"dev_root_var\"\nvpc_var = \"dev_vpc_var\"\n"
			Ω(string(contents)).To(Equal(expected))
		})
	})
})
