// chris 072115

package rebnf

import (
	"bytes"
	"fmt"
	"html"

	"path/filepath"
)

// Markers around EBNF sections in HTML files.
var (
	ebnfopen  = []byte(`<pre class="ebnf">`)
	ebnfclose = []byte(`</pre>`)
)

func preserveNewlines(buf *bytes.Buffer, excl []byte) {
	// write as many newlines as found in the excluded text
	// to maintain correct line numbers in error messages
	for _, ch := range excl {
		if ch == '\n' {
			buf.WriteByte('\n')
		}
	}
}

// ExtractEBNF extracts the EBNF text from a file in which grammar
// productions are grouped in boxes demarcated by the following HTML
// elements.
//
//	<pre class="ebnf">
//	</pre>
//
// An example of such a file is the Go Programming Language
// Specification HTML page.
//
// ExtractEBNF was itself initially extracted from
// golang.org/x/exp/ebnflint.  Unfortunately, it's not exported there,
// so we duplicate it here and export it.
func ExtractEBNF(src []byte) []byte {
	buf := new(bytes.Buffer)

	for {
		// i = beginning of EBNF text
		i := bytes.Index(src, ebnfopen)
		if i < 0 {
			break // no EBNF found - we are done
		}
		i += len(ebnfopen)

		preserveNewlines(buf, src[:i])

		// j = end of EBNF text (or end of source)
		j := bytes.Index(src[i:], ebnfclose) // close marker
		if j < 0 {
			j = len(src) - i
		}
		j += i

		// copy EBNF text
		buf.Write(src[i:j])

		// advance
		src = src[j:]
	}

	return buf.Bytes()
}

// StripTag strips out the named tag from a slice of bytes.  The tag is
// expected to open with a string of the following form.
//
//	<tagname ...>
//
// Where ... can be anything.  If any newlines appear in the ..., they
// will be preserved.
//
// The tag is expected to close with a string of the following form.
//
//	</tagname>
//
func StripTag(tagname string, src []byte) []byte {
	buf := new(bytes.Buffer)

	var (
		openbegin = []byte(fmt.Sprintf("<%s ", tagname))
		openend   = []byte(`>`)
		close     = []byte(fmt.Sprintf("</%s>", tagname))
	)

	for {
		// i = beginning of the opening of the tag
		i := bytes.Index(src, openbegin)
		if i < 0 {
			// no tag found - copy rest, and we are done
			buf.Write(src)
			break
		}

		// j = end of the opening of the tag
		j := bytes.Index(src[i:], openend)
		if j < 0 {
			// tag doesn't close - we are done
			break
		}
		j += i // Since we sliced src after the beginning of the open.
		j += len(openend)

		preserveNewlines(buf, src[i:j])

		// elide opening tag
		buf.Write(src[:i])
		src = src[j:]

		// k = end of tag (or end of source)
		k := bytes.Index(src, close) // close marker
		if k < 0 {
			// Great, we want everything then, and we're
			// done.
			buf.Write(src)
			break
		}

		// Copy up to the close.
		buf.Write(src[:k])

		k += len(close)

		// advance
		src = src[k:]
	}

	return buf.Bytes()
}

// CheckRead checks the input filename to see if it ends in ".html" or,
// if not, then checks if an opening token is present.  If either is
// true, then the source byte array is treated as an HTML document and
// EBNF text is extracted according to the following opening and closing
// tokens.
//
//	<pre class="ebnf">
//	</pre>
//
// In this case, the extracted byte slice is returned instead of the
// original.  Otherwise, the original byte slice is returned.
//
// The logic in CheckRead is extracted from golang.org/x/exp/ebnflint.
// Unfortunately, it's not exported there, so we duplicate it here and
// export it.
func CheckRead(filename string, src []byte) []byte {
	if filepath.Ext(filename) == ".html" || bytes.Index(src, ebnfopen) >= 0 {
		src = ExtractEBNF(src)

		// The Go Programming Language Specification HTML page
		// wraps production names with links.
		src = StripTag("a", src)

		// We also want to unescape HTML escape sequences.
		src = []byte(html.UnescapeString(string(src)))
	}
	return src
}
