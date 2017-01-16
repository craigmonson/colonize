package generate_test

import (
  . "github.com/craigmonson/colonize/generate"

  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)


var _ = Describe("generate/validate", func() {

  args := []string{"abc","def","geh"}

  Describe("ValidateArgsLength", func() {
    It("should return an error if len(args) < min", func() {
      err := ValidateArgsLength("Test", args, 5, 6)
      Ω(err).Should(HaveOccurred())
    })

    It("should return an error if len(args) > max", func() {
      err := ValidateArgsLength("Test", args, 1, 2)
      Ω(err).Should(HaveOccurred())
    })

    It("should not error if min < len(args) < max", func() {
      err := ValidateArgsLength("Test", args, 1, 5)
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("should not error if max == -1 and len(args) > min", func() {
      err := ValidateArgsLength("Test", args, 1, -1)
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("should error if len(args) != (min == max)", func() {
      err := ValidateArgsLength("Test", args, 1, 1)
      Ω(err).Should(HaveOccurred())
    })

    It("should not error if len(args) == (min == max)", func() {
      err := ValidateArgsLength("Test", args, 3, 3)
      Ω(err).ShouldNot(HaveOccurred())
    })
  })

  Describe("ValidateNameAvailable", func() {

    It("should error if file exists", func() {
      err := ValidateNameAvailable("Test", "validate_test.go")
      Ω(err).Should(HaveOccurred())
    })

    It("should not error if does not exist", func() {
      err := ValidateNameAvailable("Test", "foo")
      Ω(err).ShouldNot(HaveOccurred())
    })
  })
})
