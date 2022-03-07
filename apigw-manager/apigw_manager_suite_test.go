package manager_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestApigwManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ApigwManager Suite")
}
