// chris 072115

package main

import (
	"bytes"
	"path/filepath"

	"log"
)

// Markers around EBNF sections in .html files.
var (
	open  = []byte(`<pre class="ebnf">`)
	close = []byte(`</pre>`)
)

// ExtractEBNF extracts the EBNF text from a file with the open and
// close markers defined above.  E.g., the Go Programming Language
// Specification HTML page.
//
// Grammar productions are grouped in boxes demarcated by the following
// HTML elements.
//
//	<pre class="ebnf">
//	</pre>
//
// ExtractEBNF is itself extracted from golang.org/x/exp/ebnflint.
// Unfortunately, it's not exported there, so we duplicate it here and
// export it.
func ExtractEBNF(src []byte) []byte {
	var buf bytes.Buffer

	for {
		// i = beginning of EBNF text
		i := bytes.Index(src, open)
		if i < 0 {
			break // no EBNF found - we are done
		}
		i += len(open)

		// write as many newlines as found in the excluded text
		// to maintain correct line numbers in error messages
		for _, ch := range src[0:i] {
			if ch == '\n' {
				buf.WriteByte('\n')
			}
		}

		// j = end of EBNF text (or end of source)
		j := bytes.Index(src[i:], close) // close marker
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
	if filepath.Ext(filename) == ".html"  || bytes.Index(src, open) >= 0 {
		src = ExtractEBNF(src)

		log.Print(string(src))
	}
	return src
}
