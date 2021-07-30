package test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fp Suite")
}
