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
			"g1.t1.png": "t1.png",
			"g1.t2.png": "t2.png",
		})
		s := New(`
	.foo { background: url(g1.t1.png); }
	.bar { background: url(g1.t2.png); }
		`, ts)
		out, err := s.Gen()
		assert.NoError(t(), err)
		assert.Equal(t(), `
	.foo { background: url(g1.png) no-repeat; }
	.bar { background: url(g1.png) no-repeat; }
		`, out)
		ts.assertSprite("g1.png", 32, 16)
	})

	bdd.XIt("Group has one file")

	bdd.XIt("background-image")

	bdd.XIt("Two Groups")

	bdd.XIt("Ignore images")

	bdd.XIt("SetSpritePath")

	bdd.XIt("Icons not the same size")

	bdd.XIt("url('img')")

	bdd.XIt("url() not after background")

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
