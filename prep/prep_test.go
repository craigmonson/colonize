package prep_test

import (
	"io/ioutil"
	"os"

	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/log_mock"
	. "github.com/craigmonson/colonize/prep"
	"github.com/craigmonson/colonize/util_mock"

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
		_ = os.Remove(conf.CombinedVarsFilePath)
		_ = os.Remove(conf.CombinedTfFilePath)
		_ = os.Remove(conf.CombinedDerivedVarsFilePath)
		_ = os.Remove(conf.CombinedDerivedValsFilePath)
		_ = os.Remove(conf.CombinedRemoteFilePath)
	})

	Describe("Run", func() {
		var mLog *log_mock.MockLog
		var err error

		BeforeEach(func() {
			mLog = &log_mock.MockLog{}
			util_mock.MCmd = &util_mock.MockCmd{}
			util_mock.MockTheCommand()

			err = Run(conf, mLog)
		})

		AfterEach(func() {
			util_mock.ResetTheCommand()
		})

		It("should not raise an error", func() {
			Ω(err).ToNot(HaveOccurred())
		})

		It("should log the steps", func() {
			expected := `
Removing .terraform directory...
Building combined terraform variable assignment files...
Building combined variable files...
Building combined terraform files...
Building combined derived files...
Building remote config script...
Fetching terraform modules...`
			Ω(mLog.Output).To(Equal(expected))
		})

		It("should have terraform get modules", func() {
			Ω(util_mock.MCmd.Cmd).To(Equal("\nterraform get -update"))
		})

		It("should run the exec.Run command", func() {
			Ω(util_mock.MCmd.CallCount).To(Equal(1))
		})

	})

	Describe("BuildCombinedValsFile", func() {
		BeforeEach(func() {
			err = BuildCombinedValsFile(conf)
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
			expected := "provider \"aws\" {}\nvariable \"vpc_tf_test\" {}\n# test env file\n# test base file\n"
			Ω(string(contents)).To(Equal(expected))
		})
	})

	Describe("BuildCombinedDerivedFiles", func() {
		BeforeEach(func() {
			// reads in vals file in order to do substitutions
			err = BuildCombinedValsFile(conf)
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

	Describe("BuildRemoteFile", func() {
		BeforeEach(func() {
			// reads in vals and derived files in order to do substitutions
			err = BuildCombinedValsFile(conf)
			err = BuildCombinedDerivedFiles(conf)
			err = BuildRemoteFile(conf)
		})

		It("should create the combined remote shell script", func() {
			Ω(conf.CombinedRemoteFilePath).To(BeARegularFile())
		})

		It("should have the right contents", func() {
			contents, _ := ioutil.ReadFile(conf.CombinedRemoteFilePath)
			expected := `terraform remote config \
-backend=s3 \
-backend-config="bucket=dev" \
-backend-config="region=dev_root_var" \
-backend-config="key=foo-dev-bar-vpc"
`
			Ω(string(contents)).To(Equal(expected))
		})
	})
})
