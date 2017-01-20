package branch_test

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "testing"
)

func TestBranch(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "Generate Branch Suite")
}
