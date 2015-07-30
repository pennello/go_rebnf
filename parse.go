// chris 072815

package rebnf

import (
	"bytes"
	"io"
	"os"

	"io/ioutil"

	"golang.org/x/exp/ebnf"
)

// Parse parses an EBNF grammar from the file with the given filename
// and an optional io.Reader.  If the reader is not provided, the code
// will attempt to open the file specified by name.  In either case, the
// filename will be used in error output and debug messages.
//
// The logic in Parse is extracted from golang.org/x/exp/ebnflint.
// Unfortunately, it's not exported there, so we duplicate it here and
// export it.  It is modified a bit, though, to be more generic.
func Parse(filename string, r io.Reader) (ebnf.Grammar, error) {
	if r == nil {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	}

	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	src = CheckExtract(filename, src)

	grammar, err := ebnf.Parse(filename, bytes.NewBuffer(src))
	if err != nil {
		return nil, err
	}

	return grammar, err
}
