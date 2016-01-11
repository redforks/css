package main

import (
	"flag"
	"io/ioutil"
	"path/filepath"

	"github.com/redforks/css/sprite"
	"github.com/redforks/errors"
	"github.com/redforks/errors/cmdline"
)

func main() {
	cmdline.Go(func() error {
		srcCssFile := flag.String("i", "", "Input css file")
		dstCssFile := flag.String("o", "", "Output css file, can be the same as input css file")
		imgBaseDiretory := flag.String("base", "", "Base directory to resolve image files referenced in input css file. Default to input css file directory")

		flag.Parse()

		var (
			css []byte
			out string
			err error
		)

		if *srcCssFile == "" || *dstCssFile == "" {
			flag.Usage()
			return cmdline.NewExitError(2)
		}

		if *imgBaseDiretory == "" {
			*imgBaseDiretory = filepath.Dir(*srcCssFile)
		}

		if css, err = ioutil.ReadFile(*srcCssFile); err != nil {
			return errors.NewInput(err)
		}

		spriter := sprite.New(string(css), sprite.NewFileService(*imgBaseDiretory, filepath.Dir(*dstCssFile)))
		if out, err = spriter.Gen(); err != nil {
			return err
		}

		if err = ioutil.WriteFile(*dstCssFile, ([]byte)(out), 0); err != nil {
			return errors.NewRuntime(err)
		}

		return nil
	})
}
