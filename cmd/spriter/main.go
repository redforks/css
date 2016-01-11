package main

import (
	"flag"
	"os"
	"path/filepath"
)

func main() {
	srcCssFile := flag.String("i", "", "Input css file")
	dstCssFile := flag.String("o", "", "Output css file, can be the same as input css file")
	imgBaseDiretory := flag.String("base", "", "Base directory to resolve image files referenced in input css file. Default to input css file directory")

	flag.Parse()

	if *srcCssFile == "" || *dstCssFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *imgBaseDiretory == "" {
		*imgBaseDiretory = filepath.Dir(*srcCssFile)
	}
}
