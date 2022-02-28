package bkapi_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBkapi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bkapi Suite")
}
