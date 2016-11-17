package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestColonize(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Colonize Suite")
}
