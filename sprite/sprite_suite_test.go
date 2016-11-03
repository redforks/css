package sprite

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var t = GinkgoT

func TestSprite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sprite Suite")
}
