package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/redforks/css/sprite"
	"github.com/redforks/errors"
	"github.com/redforks/errors/cmdline"
)

type basePathSlice []string

func (bps *basePathSlice) String() string {
	return fmt.Sprintf("%s", *bps)
}

func (bps *basePathSlice) Set(value string) error {
	*bps = append(*bps, value)
	return nil
}

func main() {
	cmdline.Go(func() error {
		var bps basePathSlice
		srcCssFile := flag.String("i", "", "Input css file")
		dstCssFile := flag.String("o", "", "Output css file, can be the same as input css file")
		flag.Var(&bps, "base", "Base directory to resolve image files. Default to input css file directory. Can be specified multiple times, it is useful if the input css file is created by tools such as scss from multiple .css files in different directories.")

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

		if len(bps) == 0 {
			bps = basePathSlice{filepath.Dir(*srcCssFile)}
		}

		if css, err = ioutil.ReadFile(*srcCssFile); err != nil {
			return errors.NewInput(err)
		}

		spriter := sprite.New(string(css), sprite.NewFileService(([]string)(bps), filepath.Dir(*dstCssFile)))
		if out, err = spriter.Gen(); err != nil {
			return err
		}

		if err = ioutil.WriteFile(*dstCssFile, ([]byte)(out), 0); err != nil {
			return errors.NewRuntime(err)
		}

		return nil
	})
}
