//go:generate go-bindata -pkg sprite -o bindata_test.go testdata/

package sprite

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"path"

	bdd "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

var _ = bdd.Describe("sprite", func() {

	bdd.It("Empty", func() {
		ts := newTestService(nil)

		s := New("", ts)
		out, err := s.Gen()
		assert.Equal(t(), "", out)
		assert.NoError(t(), err)
		assert.Empty(t(), ts.sprites)
	})

	bdd.It("Two Icons", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png":       "t1.png",
			"image/g1.t2.png": "t2.png",
		})
		s := New(`
	.foo { background: url(g1.t1.png); }
	.bar { background: url(image/g1.t2.png); }
		`, ts)
		out, err := s.Gen()
		assert.NoError(t(), err)
		assert.Equal(t(), `.foo {
  background: url(g1.png) no-repeat;
}
.bar {
  background: url(g1.png) no-repeat -16px 0;
}`, out)
		ts.assertSprite("g1.png", 32, 16)
	})

	bdd.It("Group has one file", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png": "t1.png",
			"g1.t2.png": "t2.png",
		})
		s := New(`
	.foo { background: url(g1.t1.png); }
		`, ts)
		out, err := s.Gen()
		assert.NoError(t(), err)
		assert.Equal(t(), `.foo {
  background: url(g1.png) no-repeat;
}`, out)
		ts.assertSprite("g1.png", 16, 16)
	})

	bdd.It("Reference two identity file", func() {
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
		out, err := s.Gen()
		assert.NoError(t(), err)
		assert.Equal(t(), `.foo {
  background: url(g1.png) no-repeat;
}
.bar {
  background: url(g1.png) no-repeat -16px 0;
}
.foobar {
  background: url(g1.png) no-repeat;
}
.foo-bar {
  background: url(g1.png) no-repeat -16px 0;
}`, out)
		ts.assertSprite("g1.png", 32, 16)
	})

	bdd.It("Only Two identity file", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png": "t1.png",
			"g1.t2.png": "t2.png",
		})
		s := New(`
	.foo { background: url(g1.t1.png); }
	.foobar { background: url(g1.t1.png); }
		`, ts)
		out, err := s.Gen()
		assert.NoError(t(), err)
		assert.Equal(t(), `.foo {
  background: url(g1.png) no-repeat;
}
.foobar {
  background: url(g1.png) no-repeat;
}`, out)
		ts.assertSprite("g1.png", 16, 16)
	})

	bdd.XIt("background-image")

	bdd.It("Two Groups", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png": "t1.png",
			"g1.t2.png": "t2.png",
			"g2.t1.png": "t1.png",
			"g2.t2.png": "t2.png",
		})
		s := New(`
	.foo { background: url(g1.t1.png); }
	.bar { background: url(g1.t2.png); }
	.foobar { background: url(g2.t1.png); }
	.foo-bar { background: url(g2.t2.png); }
		`, ts)
		out, err := s.Gen()
		assert.NoError(t(), err)
		assert.Equal(t(), `.foo {
  background: url(g1.png) no-repeat;
}
.bar {
  background: url(g1.png) no-repeat -16px 0;
}
.foobar {
  background: url(g2.png) no-repeat;
}
.foo-bar {
  background: url(g2.png) no-repeat -16px 0;
}`, out)
		ts.assertSprite("g1.png", 32, 16)
		ts.assertSprite("g2.png", 32, 16)
	})

	bdd.It("Ignore images", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png": "t1.png",
			"g1.t2.png": "t2.png",
		})
		s := New(`
	.foo { background: url(g1.t1.png); }
	.foobar { background: url(g1.t2.png); }
	.bar { background: url(bar.png); }
		`, ts)
		out, err := s.Gen()
		assert.NoError(t(), err)
		assert.Equal(t(), `.foo {
  background: url(g1.png) no-repeat;
}
.foobar {
  background: url(g1.png) no-repeat -16px 0;
}
.bar {
  background: url(bar.png);
}`, out)
		ts.assertSprite("g1.png", 32, 16)
	})

	bdd.It("Icons not the same size", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png": "24.png",
			"g1.t2.png": "t2.png",
		})
		s := New(`
	.foo { background: url(g1.t1.png); }
	.foobar { background: url(g1.t2.png); }
		`, ts)
		out, err := s.Gen()
		assert.NoError(t(), err)
		assert.Equal(t(), `.foo {
  background: url(g1.png) no-repeat;
}
.foobar {
  background: url(g1.png) no-repeat -24px 0;
}`, out)
		ts.assertSprite("g1.png", 40, 24)
	})

	bdd.It("url('img')", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png": "t1.png",
			"g1.t2.png": "t2.png",
		})
		s := New(`
	.foo { background: url('g1.t1.png'); }
	.bar { background: url("g1.t2.png"); }
		`, ts)
		out, err := s.Gen()
		assert.NoError(t(), err)
		assert.Equal(t(), `.foo {
  background: url(g1.png) no-repeat;
}
.bar {
  background: url(g1.png) no-repeat -16px 0;
}`, out)
		ts.assertSprite("g1.png", 32, 16)
	})

	bdd.It("url() not after background", func() {
		ts := newTestService(map[string]string{
			"g1.t1.png": "t1.png",
		})
		s := New(`
	.foo { bkg: url(g1.t1.png); }
		`, ts)
		out, err := s.Gen()
		assert.NoError(t(), err)
		assert.Equal(t(), `.foo {
  bkg: url(g1.t1.png);
}`, out)
	})

	bdd.XIt("background has more info than url()")

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
		assert.NoError(t(), err)
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
	assert.NotNil(t(), buf)
	config, err := png.DecodeConfig(bytes.NewReader(buf.Bytes()))
	assert.NoError(t(), err)
	assert.Equal(t(), width, config.Width)
	assert.Equal(t(), height, config.Height)
}
