package writer

import (
	"io"

	"github.com/gorilla/css/scanner"
)

// Writer to write gorilla css token stream to a io.Writer.
type Writer struct {
	w io.Writer
	e error
}

// Create a new Writer instance, must call .Close() after write out all tokens.
func New(w io.Writer) *Writer {
	return &Writer{w, nil}
}

// Write token out. Any error occurred will report on calling .Close() method.
// .Write() itself always succeed.
func (w *Writer) Write(token *scanner.Token) {
	if w.e == nil {
		_, w.e = w.w.Write(([]byte)(token.Value))
	}
}

// Close() close the writer if it support io.Closer interface.
//
// Must call .Close() even the writer do not need Close(). Because
// Writer.Write() method do not report error, the error returned when
// calling .Close().
func (w *Writer) Close() error {
	closer, ok := w.w.(io.Closer)
	if ok {
		return closer.Close()
	}
	return w.e
}
