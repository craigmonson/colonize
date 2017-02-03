package branch_test

import (
  . "github.com/craigmonson/colonize/generate/branch"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "io/ioutil"
  "os"
  "path"

  "github.com/craigmonson/colonize/config"
  "github.com/craigmonson/colonize/log_mock"
  "github.com/craigmonson/colonize/util"
)

var test_dir string
var cwd string
var build_order *os.File
var err error

// Before we run tests, we create a temp dir and req files to work with
var _ = BeforeSuite(func() {

  cwd,err = os.Getwd()
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

  os.Mkdir("env", 0777)
  util.Touch(test_dir,"build_order.txt")
  util.Touch(path.Join(test_dir,"env"),"dev.tfvars")
  util.Touch(path.Join(test_dir,"env"),"prod.tfvars")

  Run(&config.Config{
      RootPath: test_dir,
      TmplPath: test_dir,

      ConfigFile: config.ConfigFile{
        Environments_Dir: "env",
        Branch_Order_File: "build_order.txt",
      },
    }, &log_mock.MockLog{}, RunArgs{
    Name:  "testbranch",
    Leafs: []string{"leaf1","leaf2"},
  })

})

// After all tests, cleanup
var _ = AfterSuite(func() {
  os.Chdir(cwd)
  build_order.Close()
  os.RemoveAll(test_dir)
})

var _ = Describe("generate/branch/branch", func() {

  It("should have created a directory for the branch", func() {
    _, err := os.Stat("../testbranch")
    Ω(err).ShouldNot(HaveOccurred())
  })

  It("should have appended the branch name to the parent build order", func() {
    contents,err := ioutil.ReadFile("../build_order.txt")
    Ω(err).ShouldNot(HaveOccurred())
    Ω(string(contents)).Should(Equal("testbranch\n"))
  })

  It("should have created an env directory in the branch", func() {
    _, err := os.Stat("env")
    Ω(err).ShouldNot(HaveOccurred())
  })

  It("should have created environment tfvars inside the env directory", func() {
    _, err := os.Stat("env/dev.tfvars")
    Ω(err).ShouldNot(HaveOccurred())
    _, err = os.Stat("env/prod.tfvars")
    Ω(err).ShouldNot(HaveOccurred())
  })

  It("should have created leaf directories (leaf1,leaf2)", func() {
    _, err := os.Stat("leaf1")
    Ω(err).ShouldNot(HaveOccurred())
    _, err = os.Stat("leaf2")
    Ω(err).ShouldNot(HaveOccurred())
  })

  It("should have listed leafs in build_order.txt", func() {
    contents,err := ioutil.ReadFile("build_order.txt")
    Ω(err).ShouldNot(HaveOccurred())
    Ω(string(contents)).Should(Equal("leaf1\nleaf2\n"))
  })

})
