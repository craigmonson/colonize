package leaf_test

import (
  . "github.com/craigmonson/colonize/generate/leaf"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "io/ioutil"
  "os"
  "github.com/craigmonson/colonize/log"
)

var test_dir string
var cwd string
var build_order *os.File
var err error

// Before we run tests, we create a temp dir to work out of
var _ = BeforeSuite(func() {

  cwd,err = os.Getwd()
  if err != nil {
    panic(err.Error())
  }

  test_dir, err = ioutil.TempDir("", "leaf_test")
  if err != nil {
    panic(err.Error())
  }

  err = os.Chdir(test_dir)
  if err != nil {
    panic(err.Error())
  }

  build_order, err = os.Create("build_order.txt")
  if err != nil {
    panic(err.Error())
  }

  Run(nil, log.Log{}, RunArgs{
    Name: "test",
    BuildOrder: build_order,
  })

})

// After all tests, cleanup
var _ = AfterSuite(func() {
  build_order.Close()
  os.RemoveAll(test_dir)
  os.Chdir(cwd)
})


var _ = Describe("generate/leaf/leaf", func() {

  Describe("Run", func() {

    It("should have created a directory for the leaf", func() {
      _, err := os.Stat("test")
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("should create a main.tf file inside the leaf", func() {
      _, err := os.Stat("test/main.tf")
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("should have appended the leaf name to build order", func() {
      contents,err := ioutil.ReadFile("build_order.txt")
      Ω(err).ShouldNot(HaveOccurred())
      Ω(string(contents)).Should(Equal("test\n"))
    })
  })
})
