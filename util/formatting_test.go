package util_test

import (
	. "github.com/craigmonson/colonize/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("util/formatting", func() {

	Describe("PadLeft", func() {
		initial := "abc"
		actual := PadLeft(initial, "X", 6)
		expected := "XXXabc"

		It("should a string of correct length", func() {
			Ω(len(actual)).To(Equal(6))
		})

		It("should produce the expected result", func() {
			Ω(actual).To(Equal(expected))
		})

	})

	Describe("PadRight", func() {
		initial := "abc"
		actual := PadRight(initial, "X", 6)
		expected := "abcXXX"

		It("should a string of correct length", func() {
			Ω(len(actual)).To(Equal(6))
		})

		It("should produce the expected result", func() {
			Ω(actual).To(Equal(expected))
		})

	})
})
