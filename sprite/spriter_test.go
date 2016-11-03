//go:generate go-bindata -pkg sprite -o bindata_test.go testdata/

package sprite

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("sprite", func() {

	It("Empty", func() {
		ts := newTestService(nil)

		s := New("", ts)
		Ω(s.Gen()).Should(Equal(""))
		Ω(ts.sprites).Should(BeEmpty())
	})

	It("Two Icons", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png":       "t1.png",
			"image/g1.t2.png": "t2.png",
		})
		s := New(`
	.foo { background: url(g1.t1.png); }
	.bar { background: url(image/g1.t2.png); }
		`, ts)
		Ω(s.Gen()).Should(Equal(`
	.foo { background: url(l01cVKU8.png) no-repeat; }
	.bar { background: url(l01cVKU8.png) no-repeat -16px 0; }
		`))

		ts.assertSprite("l01cVKU8.png", 32, 16)
	})

	It("Group has one file", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png": "t1.png",
			"g1.t2.png": "t2.png",
		})
		s := New(`
	.foo { background: url(g1.t1.png); }
		`, ts)
		Ω(s.Gen()).Should(Equal(`
	.foo { background: url(cobz_bF6.png) no-repeat; }
		`))
		ts.assertSprite("cobz_bF6.png", 16, 16)
	})

	It("Reference two identity file", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png": "t1.png",
			"g1.t2.png": "t2.png",
		})
		s := New(`
	.foo { background: url(g1.t1.png); }
	.bar { background: url(g1.t2.png); }
	.foobar { background: url(g1.t1.png); }
	.foo-bar { background: url(g1.t2.png); }
		`, ts)
		Ω(s.Gen()).Should(Equal(`
	.foo { background: url(l01cVKU8.png) no-repeat; }
	.bar { background: url(l01cVKU8.png) no-repeat -16px 0; }
	.foobar { background: url(l01cVKU8.png) no-repeat; }
	.foo-bar { background: url(l01cVKU8.png) no-repeat -16px 0; }
		`))
		ts.assertSprite("l01cVKU8.png", 32, 16)
	})

	It("Only Two identity file", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png": "t1.png",
			"g1.t2.png": "t2.png",
		})
		s := New(`
	.foo { background: url(g1.t1.png); }
	.foobar { background: url(g1.t1.png); }
		`, ts)
		Ω(s.Gen()).Should(Equal(`
	.foo { background: url(cobz_bF6.png) no-repeat; }
	.foobar { background: url(cobz_bF6.png) no-repeat; }
		`))
		ts.assertSprite("cobz_bF6.png", 16, 16)
	})

	XIt("background-image")

	It("Two Groups", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png": "t1.png",
			"g1.t2.png": "t2.png",
			"g2.t1.png": "t1.png",
			"g2.t2.png": "t2.png",
		})
		s := New(`
	.foo { background: url(g1.t1.png); }
	.bar { background: url(g1.t2.png); }
	.foobar { background: url(g2.t2.png); }
	.foo-bar { background: url(g2.t1.png); }
		`, ts)
		Ω(s.Gen()).Should(Equal(`
	.foo { background: url(l01cVKU8.png) no-repeat; }
	.bar { background: url(l01cVKU8.png) no-repeat -16px 0; }
	.foobar { background: url(piouFODI.png) no-repeat; }
	.foo-bar { background: url(piouFODI.png) no-repeat -16px 0; }
		`))
		ts.assertSprite("l01cVKU8.png", 32, 16)
		ts.assertSprite("piouFODI.png", 32, 16)
	})

	It("Ignore images", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png": "t1.png",
			"g1.t2.png": "t2.png",
		})
		s := New(`
	.foo { background: url(g1.t1.png); }
	.foobar { background: url(g1.t2.png); }
	.bar { background: url(bar.png); }
		`, ts)
		Ω(s.Gen()).Should(Equal(`
	.foo { background: url(l01cVKU8.png) no-repeat; }
	.foobar { background: url(l01cVKU8.png) no-repeat -16px 0; }
	.bar { background: url(bar.png); }
		`))
		ts.assertSprite("l01cVKU8.png", 32, 16)
	})

	It("Icons not the same size", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png": "24.png",
			"g1.t2.png": "t2.png",
		})
		s := New(`
	.foo { background: url(g1.t1.png); }
	.foobar { background: url(g1.t2.png); }
		`, ts)
		Ω(s.Gen()).Should(Equal(`
	.foo { background: url(TRw0KSHq.png) no-repeat; }
	.foobar { background: url(TRw0KSHq.png) no-repeat -24px 0; }
		`))
		ts.assertSprite("TRw0KSHq.png", 40, 24)
	})

	It("url('img')", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png": "t1.png",
			"g1.t2.png": "t2.png",
		})
		s := New(`
	.foo { background: url('g1.t1.png'); }
	.bar { background: url("g1.t2.png"); }
		`, ts)
		Ω(s.Gen()).Should(Equal(`
	.foo { background: url(l01cVKU8.png) no-repeat; }
	.bar { background: url(l01cVKU8.png) no-repeat -16px 0; }
		`))
		ts.assertSprite("l01cVKU8.png", 32, 16)
	})

	It("url() not after background", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png": "t1.png",
		})
		s := New(`
	.foo { bkg: url(g1.t1.png); }
		`, ts)
		Ω(s.Gen()).Should(Equal(`
	.foo { bkg: url(g1.t1.png); }
		`))
	})

	XIt("background has more info than url()")

})

// Implement Service interface for testing
type testService struct {
	images  map[string][]byte        // source images path -> image
	sprites map[string]*bytes.Buffer // created sprites
}

// images: filename -> resource name
func newTestService(images map[string]string) *testService {
	imgs := make(map[string][]byte)
	for filename, resName := range images {
		content, err := Asset(path.Join("testdata", resName))
		Ω(err).Should(Succeed())
		imgs[filename] = content
	}
	return &testService{
		images:  imgs,
		sprites: make(map[string]*bytes.Buffer),
	}
}

func (s *testService) OpenImage(path string) (io.Reader, error) {
	r := s.images[path]
	if r == nil {
		return nil, fmt.Errorf("file %s not exist", path)
	}
	return bytes.NewReader(r), nil
}

func (s *testService) CreateSpriteImage(path string) (io.Writer, error) {
	r := &bytes.Buffer{}
	s.sprites[path] = r
	return r, nil
}

func (s *testService) assertSprite(path string, width, height int) {
	buf := s.sprites[path]
	Ω(buf).ShouldNot(BeNil())
	config, err := png.DecodeConfig(bytes.NewReader(buf.Bytes()))
	Ω(err).Should(Succeed())
	Ω(config.Width).Should(Equal(width))
	Ω(config.Height).Should(Equal(height))
}
