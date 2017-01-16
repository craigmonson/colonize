package leaf_test

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "testing"
)

func TestLeaf(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "Generate Leaf Suite")
}
