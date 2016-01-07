package writer_test

import (
	"bytes"
	"fmt"

	"github.com/gorilla/css/scanner"
	"github.com/redforks/css/writer"
)

func ExampleWriter() {
	css := `
		// comment
		.foo {
			color: white;
		}
		`
	s, buf := scanner.New(css), bytes.Buffer{}
	w := writer.New(&buf)
	for to := s.Next(); to.Type != scanner.TokenEOF; to = s.Next() {
		w.Write(to)
	}
	if err := w.Close(); err != nil {
		fmt.Print(err)
	} else {
		fmt.Print((string)(buf.Bytes()))
	}
	// Output:
	//
	//		// comment
	//		.foo {
	//			color: white;
	//		}
	//
}
