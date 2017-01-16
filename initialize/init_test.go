package initialize_test

import (
  "reflect"

  "github.com/craigmonson/colonize/config"
  "github.com/craigmonson/colonize/log_mock"
  . "github.com/craigmonson/colonize/initialize"

  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("Init", func() {

  Describe("Run", func() {
    var conf config.Config
    var mLog *log_mock.MockLog
    var err error

    BeforeEach(func() {
      mLog = &log_mock.MockLog{}
      err = Run(&conf, mLog, true)
    })

    It("should not raise an error", func() {
      Ω(err).ToNot(HaveOccurred())
    })

    It("should setup config data w/ defaults", func() {

      defaults := reflect.ValueOf(&config.ConfigFileDefaults).Elem()
      actuals := reflect.ValueOf(&conf.ConfigFile).Elem()

      for i := 0; i< defaults.NumField(); i++ {
        Ω(actuals.Field(i).String()).To(Equal(defaults.Field(i).String()))
      }

    })
  })
})
