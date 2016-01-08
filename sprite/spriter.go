package sprite

import (
	"errors"
	"image"
	_ "image/png"
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
		return "", err
	}

	needTranslate := false
	groups := make(map[string][]*scanner.Token)
	for _, tk := range tks {
		switch tk.Type {
		case scanner.TokenIdent:
			needTranslate = tk.Value == "background"
		case scanner.TokenURI:
			if needTranslate {
				var imgFile string
				if imgFile, err = extractUriFile(tk.Value); err != nil {
					return "", err
				}
				g := extractGroup(imgFile)
				if g != "" {
					groups[g] = append(groups[g], tk)
				}
			}
		}
	}

	for g, tks := range groups {
		var sprite = image.NewRGBA(image.Rect(0, 0, 0, 0))
		for _, tk := range tks {
			tk.Value = "url(" + g + ".png) no-repeat"

			f, err := extractUriFile(tk.Value)
			if err != nil {
				return "", err
			}

			r, err := s.sv.OpenImage(f)
			defer closeClosable(r)

			img, _, err := image.Decode(r)
			if err != nil {
				return "", err
			}

			sprite = appendImage(sprite, img)
		}
		// f, err := s.sv.CreateSpriteImage(g + ".png")
		// if err != nil {
		// 	return nil, err
		// }
		// defer closeClosable(f)
	}

	return writer.Dumps(tks)
}

// Create a new image append img2 after img1, img1, img2 not affect.
func appendImage(img1 *image.RGBA, img2 image.Image) *image.RGBA {
	// r:=image.NewRGBA(image.Rect(0, 0, img1.))
	return nil
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
