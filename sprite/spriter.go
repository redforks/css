package sprite

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/redforks/douceur/css"
	"github.com/redforks/douceur/parser"
	"github.com/redforks/errors"
)

// Service provide/isolate I/O interface of Spriter.
type Service interface {
	// OpenImage opens image file as io.Reader, path is the relative path defined
	// in css. Spriter will close the io.Reader if it also implements io.Closer.
	OpenImage(path string) (io.Reader, error)

	// Create sprite image file return as io.Writer. Spriter will close the
	// io.Writer if it also implements io.Closer. path is relative to output css
	// file.
	CreateSpriteImage(path string) (io.Writer, error)
}

// Generate css sprite image by scan .css file, generate updated .css file.
//
// Image must use relative path, absolute path or other web site report as warning.
//
// Only png file supported, ignore other image file format.
//
// Image file name need to be in [Group].[Name].png format, images with the
// same group name will generate a sprite image [Group].png. Images without
// group name leave it untouched.
type Spriter struct {
	css string
	sv  Service

	loadedImages map[string]*stamp
}

// Create Spriter.
//
//  css: css file content
func New(css string, service Service) *Spriter {
	return &Spriter{
		css:          css,
		sv:           service,
		loadedImages: make(map[string]*stamp),
	}
}

// Do the generation, return translated css file content. Generated sprite
// image files are saved using Service interface.
func (s *Spriter) Gen() (out string, err error) {
	var sheet *css.Stylesheet
	sheet, err = scan(s.css)
	if err != nil {
		return
	}
	rules := sheet.Rules

	groups := make(map[string][]*cssImage)
	for _, rule := range rules {
		for _, attr := range rule.Declarations {
			if attr.Property == "background" {
				var (
					st *cssImage
					g  string
				)
				if st, g, err = s.parseCssImage(attr); err != nil {
					return
				}

				if st != nil {
					groups[g] = append(groups[g], st)
				}
			}
		}
	}

	for g, sts := range groups {
		size := getSpriteSize(sts)
		var sprite = image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: size})
		p := image.Pt(0, 0)
		for _, st := range sts {
			st.tk.Value = "url(" + g + ".png) no-repeat"
			if st.img.sp.X != 0 {
				st.tk.Value += fmt.Sprintf(" %dpx 0", st.img.sp.X)
			}
			b := st.img.bounds()
			draw.Draw(sprite, b.Add(p), st.img.img, b.Min, draw.Src)
			p.X += st.img.dx()
		}

		var f io.Writer
		f, err = s.sv.CreateSpriteImage(g + ".png")
		if err != nil {
			return
		}
		defer closeClosable(f)
		if err = png.Encode(f, sprite); err != nil {
			err = errors.NewRuntime(err)
			return
		}
	}

	return sheet.String(), nil
}

func getSpriteSize(imgs []*cssImage) image.Point {
	p := image.Point{}
	for _, img := range imgs {
		st := img.img
		if st.sp.X == -1 {
			st.sp = image.Pt(-p.X, 0)
			p.X += st.dx()
			if p.Y < st.dy() {
				p.Y = st.dy()
			}
		}
	}
	return p
}

// Call .Close() if object implements io.Closer.
func closeClosable(o interface{}) {
	closable, ok := o.(io.Closer)
	if ok {
		closable.Close()
	}
}

func scan(input string) (*css.Stylesheet, error) {
	sheet, err := parser.Parse(input)
	if err != nil {
		return nil, errors.NewInput(err)
	}
	return sheet, nil
}

func extractUriFile(uri string) (file string, err error) {
	s := uri[4 : len(uri)-1]
	switch s[0] {
	case '"', '\'':
		s = s[1 : len(s)-1]
	}
	return s, nil
}

// extract file name, expect [group].[name].png. Group name is empty string if
// not expected format, or extension not supported
func extractGroup(path string) (group string) {
	words := strings.Split(filepath.Base(path), ".")
	if len(words) != 3 || words[2] != "png" {
		return
	}

	return words[0]
}

// Represent a css image style
type cssImage struct {
	tk  *css.Declaration
	img *stamp
}

// Represent a image inside sprite
type stamp struct {
	filename string // Filename of the image
	img      image.Image
	sp       image.Point // Start position in sprite
}

func (st *stamp) bounds() image.Rectangle {
	return st.img.Bounds()
}

func (st *stamp) dx() int {
	return st.img.Bounds().Dx()
}

func (st *stamp) dy() int {
	return st.img.Bounds().Dy()
}

// Parse stamp from a image url css token. stamp is nil if the url need
// ignored: not png, not expected filename format.
func (s *Spriter) parseCssImage(attr *css.Declaration) (cssImg *cssImage, groupName string, err error) {
	var fn string
	if fn, err = extractUriFile(attr.Value); err != nil {
		return
	}

	groupName = extractGroup(fn)
	if groupName == "" {
		return
	}

	var st *stamp
	if st, err = s.parseImage(fn); err != nil {
		return
	}

	cssImg = &cssImage{
		attr,
		st,
	}
	return
}

func (s *Spriter) parseImage(imgFile string) (*stamp, error) {
	if img, ok := s.loadedImages[imgFile]; ok {
		return img, nil
	}

	if f, err := s.sv.OpenImage(imgFile); err != nil {
		return nil, errors.NewInput(err)
	} else {
		defer closeClosable(f)

		img, _, err := image.Decode(f)
		if err != nil {
			return nil, errors.NewInput(err)
		}
		st := &stamp{
			imgFile,
			img,
			image.Pt(-1, -1),
		}
		s.loadedImages[imgFile] = st
		return st, nil
	}
}
