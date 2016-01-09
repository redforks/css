package sprite

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"strings"

	"github.com/gorilla/css/scanner"
	"github.com/redforks/css/writer"
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
}

// Create Spriter.
//
//  css: css file content
func New(css string, service Service) *Spriter {
	return &Spriter{
		css: css,
		sv:  service,
	}
}

// SetSpritePath set the relative path of generated sprite image files. Default
// is "", which means store sprite in the same css directory. Other possible
// value can be "../images", "images/".
func (s *Spriter) SetSpritePath(path string) {
}

// Do the generation, return translated css file content. Generated sprite
// image files are saved using Service interface.
func (s *Spriter) Gen() (css string, err error) {
	var tks []*scanner.Token
	tks, err = scan(s.css)
	if err != nil {
		return
	}

	needTranslate := false
	groups := make(map[string][]*stamp)
	for _, tk := range tks {
		switch tk.Type {
		case scanner.TokenIdent:
			needTranslate = tk.Value == "background"
		case scanner.TokenURI:
			if needTranslate {
				var (
					st *stamp
					g  string
				)
				if st, g, err = parseStamp(s.sv, tk); err != nil {
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
		for i, st := range sts {
			st.tk.Value = "url(" + g + ".png) no-repeat"
			if i != 0 {
				st.tk.Value += fmt.Sprintf(" -%dpx 0", p.X)
			}
			b := st.bounds()
			draw.Draw(sprite, b.Add(p), st.img, b.Min, draw.Src)
			p.X += st.dx()
		}

		var f io.Writer
		f, err = s.sv.CreateSpriteImage(g + ".png")
		if err != nil {
			return
		}
		defer closeClosable(f)
		if err = png.Encode(f, sprite); err != nil {
			return
		}
	}

	return writer.Dumps(tks)
}

func getSpriteSize(imgs []*stamp) image.Point {
	p := image.Point{}
	for _, img := range imgs {
		p.X += img.dx()
		if p.Y < img.dy() {
			p.Y = img.dy()
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

func scan(css string) ([]*scanner.Token, error) {
	s := scanner.New(css)
	tks := []*scanner.Token{}
	for {
		tk := s.Next()
		switch tk.Type {
		case scanner.TokenEOF:
			return tks, nil
		case scanner.TokenError:
			return nil, errors.New(tk.Value)
		default:
			tks = append(tks, tk)
		}
	}
}

func extractUriFile(uri string) (file string, err error) {
	return uri[4 : len(uri)-1], nil
}

// extract file name, expect [group].[name].png. Group name is empty string if
// not expected format, or extension not supported
func extractGroup(file string) (group string) {
	words := strings.Split(file, ".")
	if len(words) != 3 || words[2] != "png" {
		return
	}

	return words[0]
}

// Represent a image inside sprite
type stamp struct {
	tk       *scanner.Token
	filename string // Filename of the image
	img      image.Image
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
func parseStamp(sv Service, tk *scanner.Token) (stmp *stamp, groupName string, err error) {
	var fn string
	if fn, err = extractUriFile(tk.Value); err != nil {
		return
	}

	groupName = extractGroup(fn)
	if groupName == "" {
		return
	}

	var f io.Reader
	if f, err = sv.OpenImage(fn); err != nil {
		return
	}
	defer closeClosable(f)

	var img image.Image
	if img, _, err = image.Decode(f); err != nil {
		return
	}

	stmp = &stamp{
		tk, fn, img,
	}
	return
}
