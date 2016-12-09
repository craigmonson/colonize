package prep_test

import (
	"io/ioutil"
	//"os"

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
		//_ = os.Remove(conf.CombinedValsFilePath)
		//_ = os.Remove(conf.CombinedVarsFilePath)
		//_ = os.Remove(conf.CombinedTfFilePath)
	})

	PDescribe("Run", func() {})

	Describe("BuildCombinedValuesFile", func() {
		BeforeEach(func() {
			err = BuildCombinedValuesFile(conf)
		})

		It("should create the combined file", func() {
			Ω(conf.CombinedValsFilePath).To(BeARegularFile())
		})

		It("should have the right contents", func() {
			contents, _ := ioutil.ReadFile(conf.CombinedValsFilePath)
			expected := "root_var = \"dev_root_var\"\nvpc_var = \"dev_vpc_var\"\n"
			Ω(string(contents)).To(Equal(expected))
		})
	})

	Describe("BuildCombinedVarsFile", func() {
		BeforeEach(func() {
			err = BuildCombinedVarsFile(conf)
		})

		It("should create the combined file", func() {
			Ω(conf.CombinedVarsFilePath).To(BeARegularFile())
		})

		It("should have the right contents", func() {
			contents, _ := ioutil.ReadFile(conf.CombinedVarsFilePath)
			expected := "variable \"root_var\" {}\nvariable \"vpc_var\" {}\n"
			Ω(string(contents)).To(Equal(expected))
		})
	})

	Describe("BuildCombinedTfFile", func() {
		BeforeEach(func() {
			err = BuildCombinedTfFile(conf)
		})

		It("should create the combined file", func() {
			Ω(conf.CombinedTfFilePath).To(BeARegularFile())
		})

		It("should have the right contents", func() {
			contents, _ := ioutil.ReadFile(conf.CombinedTfFilePath)
			expected := "provider \"aws\" {}\nvariable \"vpc_tf_test\" {}\n"
			Ω(string(contents)).To(Equal(expected))
		})
	})

	Describe("BuildCombinedDerivedFile", func() {
		BeforeEach(func() {
			err = BuildCombinedDerivedFile(conf)
		})

		It("should create the combined file", func() {
			Ω(conf.CombinedDerivedFilePath).To(BeARegularFile())
		})

		It("should have the right contents", func() {
			contents, _ := ioutil.ReadFile(conf.CombinedDerivedFilePath)
			expected := "provider \"aws\" {}\nvariable \"vpc_tf_test\" {}\n"
			Ω(string(contents)).To(Equal(expected))
		})
	})
})
