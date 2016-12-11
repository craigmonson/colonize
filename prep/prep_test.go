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
		//_ = os.Remove(conf.CombinedValsFilePath)
		//_ = os.Remove(conf.CombinedVarsFilePath)
		_ = os.Remove(conf.CombinedTfFilePath)
		//_ = os.Remove(conf.CombinedDerivedFilePath)
	})

	PDescribe("Run", func() {})

	Describe("BuildCombinedValuesFile", func() {
		BeforeEach(func() {
			err = BuildCombinedValuesFile(conf)
		})

		It("should create the combined file", func() {
			Ω(conf.CombinedValsFilePath).To(BeARegularFile())
		})

		It("should have the right contents (derived too)", func() {
			contents, _ := ioutil.ReadFile(conf.CombinedValsFilePath)
			expected := `environment = "dev"
origin_path = "../test/vpc"
tmpl_name = "vpc"
tmpl_path_dashed = "vpc"
tmpl_path_underscored = "vpc"
root_path = "../test"
root_var = "dev_root_var"
vpc_var = "dev_vpc_var"
`
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

		It("should have the right contents (derived too)", func() {
			contents, _ := ioutil.ReadFile(conf.CombinedVarsFilePath)
			expected := `variable "environment" {}
variable "origin_path" {}
variable "tmpl_name" {}
variable "tmpl_path_dashed" {}
variable "tmpl_path_underscored" {}
variable "root_path" {}
variable "root_var" {}
variable "vpc_var" {}
`
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

	Describe("BuildCombinedDerivedFiles", func() {
		BeforeEach(func() {
			err = BuildCombinedDerivedFiles(conf)
		})

		It("should create the combined val file", func() {
			Ω(conf.CombinedDerivedValsFilePath).To(BeARegularFile())
		})

		It("should create the combined var file", func() {
			Ω(conf.CombinedDerivedVarsFilePath).To(BeARegularFile())
		})

		It("should have the right val contents", func() {
			contents, _ := ioutil.ReadFile(conf.CombinedDerivedValsFilePath)
			expected := "test_derived = \"foo-dev-bar-vpc\"\n"
			Ω(string(contents)).To(Equal(expected))
		})

		It("should have the right var contents", func() {
			contents, _ := ioutil.ReadFile(conf.CombinedDerivedVarsFilePath)
			expected := "variable \"test_derived\" {}\n"
			Ω(string(contents)).To(Equal(expected))
		})
	})
})
