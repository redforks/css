package sprite

import (
	"io"
	"os"
	"path/filepath"
)

// FileService implement Service interface to load/save image to/from file
// system.
type fileService struct {
	srcPath, outPath string
}

// Create a Service work with file system.
//
//  srcPath: Base path used to resolve image files referenced in css. Normally
//  it is the directory where src .css file is.
//  outPath: Base path used to resolve generated sprite image files. Normally
//  it is the diretory where out .css file is.
func NewFileService(srcPath, outPath string) Service {
	return &fileService{srcPath, outPath}
}

func (f *fileService) OpenImage(path string) (io.Reader, error) {
	var p = filepath.Join(f.srcPath, path)
	return os.Open(p)
}

func (f *fileService) CreateSpriteImage(path string) (io.Writer, error) {
	var p = filepath.Join(f.outPath, path)
	return os.Create(p)
}
