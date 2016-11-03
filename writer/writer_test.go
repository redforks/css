package writer

import (
	"bytes"

	"github.com/golang/mock/gomock"
	bdd "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/redforks/css-1/scanner"
	"github.com/redforks/testing/iotest"
)

var _ = bdd.Describe("Writer", func() {
	var (
		buf bytes.Buffer
		w   *Writer
	)

	assertClose := func(content string) {
		Ω(w.Close()).Should(Succeed())
		Ω(buf.Bytes()).Should(BeEquivalentTo(content))
	}

	bdd.BeforeEach(func() {
		buf = bytes.Buffer{}
		w = New(&buf)
	})

	bdd.It("Empty", func() {
		assertClose("")
	})

	bdd.It("Ident", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenIdent,
			Value: "foo",
		})
		assertClose("foo")
	})

	bdd.It("At Keyword", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenAtKeyword,
			Value: "@foo",
		})
		assertClose("@foo")
	})

	bdd.It("String", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenString,
			Value: `"foo"`,
		})
		assertClose(`"foo"`)
	})

	bdd.It("Hash", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenHash,
			Value: "#name",
		})
		assertClose("#name")
	})

	bdd.It("Number", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenNumber,
			Value: "42",
		})
		assertClose("42")
	})

	bdd.It("Percentage", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenPercentage,
			Value: "42%",
		})
		assertClose("42%")
	})

	bdd.It("Dimension", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenDimension,
			Value: "42px",
		})
		assertClose("42px")
	})

	bdd.It("URI", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenURI,
			Value: "url('http://www.google.com/')",
		})
		assertClose("url('http://www.google.com/')")
	})

	bdd.It("UnicodeRange", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenUnicodeRange,
			Value: "U+0042",
		})
		assertClose("U+0042")
	})

	bdd.It("CDO", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenCDO,
			Value: "<!--",
		})
		assertClose("<!--")
	})

	bdd.It("CDC", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenCDC,
			Value: "-->",
		})
		assertClose("-->")
	})

	bdd.It("S", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenS,
			Value: "   \n   \t   \n",
		})
		assertClose("   \n   \t   \n")
	})

	bdd.It("Comment", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenComment,
			Value: "/* foo */",
		})
		assertClose("/* foo */")
	})

	bdd.It("Function", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenFunction,
			Value: "bar(",
		})
		assertClose("bar(")
	})

	bdd.It("Includes", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenIncludes,
			Value: "~=",
		})
		assertClose("~=")
	})

	bdd.It("DashMatch", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenIncludes,
			Value: "|=",
		})
		assertClose("|=")
	})

	bdd.It("PrefixMatch", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenPrefixMatch,
			Value: "^=",
		})
		assertClose("^=")
	})

	bdd.It("SuffixMatch", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenSuffixMatch,
			Value: "$=",
		})
		assertClose("$=")
	})

	bdd.It("SubstringMatch", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenSubstringMatch,
			Value: "*=",
		})
		assertClose("*=")
	})

	bdd.It("Char", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenChar,
			Value: "{",
		})
		assertClose("{")
	})

	bdd.It("BOM", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenBOM,
			Value: "\uFEFF",
		})
		assertClose("\uFEFF")
	})

	bdd.It("Parse and Write", func() {
		css := `
		// comment
		.foo {
			color: white;
		}
		`
		s := scanner.New(css)
		for to := s.Next(); to.Type != scanner.TokenEOF; to = s.Next() {
			w.Write(to)
		}
		assertClose(css)
	})

	bdd.It("Close closable writer", func() {
		controller := gomock.NewController(t())
		defer controller.Finish()
		bufMock := NewMockWriteCloser(controller)
		bufMock.EXPECT().Close()

		w := New(bufMock)
		Ω(w.Close()).Should(Succeed())
	})

	bdd.It("Inner writer error", func() {
		w := New(iotest.ErrorWriter(5))
		w.Write(&scanner.Token{
			Type:  scanner.TokenIdent,
			Value: "foo",
		})
		w.Write(&scanner.Token{
			Type:  scanner.TokenIdent,
			Value: "bar",
		})
		w.Write(&scanner.Token{
			Type:  scanner.TokenIdent,
			Value: "foobar",
		})
		Ω(w.Close()).Should(MatchError(iotest.ErrWriter))
	})

	bdd.It("Dumps", func() {
		css := `
		// comment
		.foo {
			color: white;
		}
		`
		s := scanner.New(css)
		var tokens []*scanner.Token
		for to := s.Next(); to.Type != scanner.TokenEOF; to = s.Next() {
			tokens = append(tokens, to)
		}
		Ω(Dumps(tokens)).Should(Equal(css))
	})

})
