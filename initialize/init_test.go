package initialize_test

import (
	. "github.com/craigmonson/colonize/initialize"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"reflect"

	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/log_mock"
)

var test_dir string
var cwd string
var err error
var cfg config.Config

var _ = BeforeSuite(func() {

	cwd, err = os.Getwd()
	if err != nil {
		panic(err.Error())
	}

	test_dir, err = ioutil.TempDir("", "branch_test")
	if err != nil {
		panic(err.Error())
	}

	err = os.Chdir(test_dir)
	if err != nil {
		panic(err.Error())
	}

	cfg = config.Config{
		RootPath:   test_dir,
		TmplPath:   test_dir,
		ConfigFile: config.ConfigFile{},
	}

	Run(&cfg, &log_mock.MockLog{}, RunArgs{
		AcceptDefaults:   true,
		InitEnvironments: "dev,prod",
	})
})

var _ = AfterSuite(func() {
	os.Chdir(cwd)
	os.RemoveAll(test_dir)
})

var _ = Describe("init", func() {

	It("should not raise an error", func() {
		Ω(err).ToNot(HaveOccurred())
	})

	It("should setup config data w/ defaults", func() {
		defaults := reflect.ValueOf(&config.ConfigFileDefaults).Elem()
		actuals := reflect.ValueOf(&cfg.ConfigFile).Elem()

		for i := 0; i < defaults.NumField(); i++ {
			Ω(actuals.Field(i).String()).To(Equal(defaults.Field(i).String()))
		}

	})

	It("should have created a .colonize.yaml", func() {
		_, err := os.Stat(".colonize.yaml")
		Ω(err).ShouldNot(HaveOccurred())
	})

	It("should have created a branch order file", func() {
		_, err := os.Stat("build_order.txt")
		Ω(err).ShouldNot(HaveOccurred())
	})

	It("should have created a directory for the env", func() {
		_, err := os.Stat("env")
		Ω(err).ShouldNot(HaveOccurred())
	})

	It("should have created environment tfvars inside the env directory", func() {
		_, err := os.Stat("env/dev.tfvars")
		Ω(err).ShouldNot(HaveOccurred())
		_, err = os.Stat("env/prod.tfvars")
		Ω(err).ShouldNot(HaveOccurred())
	})

})
