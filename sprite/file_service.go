package sprite

import (
	"io"
	"os"
	"path/filepath"
)

// FileService implement Service interface to load/save image to/from file
// system.
type fileService struct {
	srcPaths []string
	outPath  string
}

// Create a Service work with file system.
//
//  srcPaths: Base path used to resolve image files referenced in css. Normally
//  it is the directory where src .css file is. Can be multiple path, if the
//  .img not exist in 1st path, search it in next path.
//  outPath: Base path used to resolve generated sprite image files. Normally
//  it is the diretory where out .css file is.
func NewFileService(srcPaths []string, outPath string) Service {
	return &fileService{srcPaths, outPath}
}

func (f *fileService) OpenImage(path string) (r io.Reader, err error) {
	for _, srcPath := range f.srcPaths {
		var p = filepath.Join(srcPath, path)
		if r, err = os.Open(p); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		return r, nil
	}
	return
}

func (f *fileService) CreateSpriteImage(path string) (io.Writer, error) {
	var p = filepath.Join(f.outPath, path)
	return os.Create(p)
}
