package writer

import (
	"bytes"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/redforks/css-1/scanner"
	"github.com/redforks/testing/iotest"
)

var _ = Describe("Writer", func() {
	var (
		buf bytes.Buffer
		w   *Writer
	)

	assertClose := func(content string) {
		Ω(w.Close()).Should(Succeed())
		Ω(buf.Bytes()).Should(BeEquivalentTo(content))
	}

	BeforeEach(func() {
		buf = bytes.Buffer{}
		w = New(&buf)
	})

	It("Empty", func() {
		assertClose("")
	})

	It("Ident", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenIdent,
			Value: "foo",
		})
		assertClose("foo")
	})

	It("At Keyword", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenAtKeyword,
			Value: "@foo",
		})
		assertClose("@foo")
	})

	It("String", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenString,
			Value: `"foo"`,
		})
		assertClose(`"foo"`)
	})

	It("Hash", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenHash,
			Value: "#name",
		})
		assertClose("#name")
	})

	It("Number", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenNumber,
			Value: "42",
		})
		assertClose("42")
	})

	It("Percentage", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenPercentage,
			Value: "42%",
		})
		assertClose("42%")
	})

	It("Dimension", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenDimension,
			Value: "42px",
		})
		assertClose("42px")
	})

	It("URI", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenURI,
			Value: "url('http://www.google.com/')",
		})
		assertClose("url('http://www.google.com/')")
	})

	It("UnicodeRange", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenUnicodeRange,
			Value: "U+0042",
		})
		assertClose("U+0042")
	})

	It("CDO", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenCDO,
			Value: "<!--",
		})
		assertClose("<!--")
	})

	It("CDC", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenCDC,
			Value: "-->",
		})
		assertClose("-->")
	})

	It("S", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenS,
			Value: "   \n   \t   \n",
		})
		assertClose("   \n   \t   \n")
	})

	It("Comment", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenComment,
			Value: "/* foo */",
		})
		assertClose("/* foo */")
	})

	It("Function", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenFunction,
			Value: "bar(",
		})
		assertClose("bar(")
	})

	It("Includes", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenIncludes,
			Value: "~=",
		})
		assertClose("~=")
	})

	It("DashMatch", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenIncludes,
			Value: "|=",
		})
		assertClose("|=")
	})

	It("PrefixMatch", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenPrefixMatch,
			Value: "^=",
		})
		assertClose("^=")
	})

	It("SuffixMatch", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenSuffixMatch,
			Value: "$=",
		})
		assertClose("$=")
	})

	It("SubstringMatch", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenSubstringMatch,
			Value: "*=",
		})
		assertClose("*=")
	})

	It("Char", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenChar,
			Value: "{",
		})
		assertClose("{")
	})

	It("BOM", func() {
		w.Write(&scanner.Token{
			Type:  scanner.TokenBOM,
			Value: "\uFEFF",
		})
		assertClose("\uFEFF")
	})

	It("Parse and Write", func() {
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

	It("Close closable writer", func() {
		controller := gomock.NewController(t())
		defer controller.Finish()
		bufMock := NewMockWriteCloser(controller)
		bufMock.EXPECT().Close()

		w := New(bufMock)
		Ω(w.Close()).Should(Succeed())
	})

	It("Inner writer error", func() {
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

	It("Dumps", func() {
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
