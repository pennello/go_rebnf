// chris 073015

// extractebnf extracts and prints EBNF grammars.
//
// In particular, it extracts grammars from HTML documents, as specified
// by golang.org/x/exp/ebnflint, which includes the Go Programming
// Language Specification HTML page.
//
// Options are:
//
//	-strip=false: strip superfluous newlines
//
package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"path"

	"io/ioutil"

	"github.com/crunchyroll/rebnf"
)

var args struct {
	progName string
	strip    bool
}

func init() {
	log.SetFlags(0)
	strip := flag.Bool("strip", false, "strip superfluous newlines")
	flag.Parse()
	args.progName = path.Base(os.Args[0])
	args.strip = *strip
}

func usage() {
	log.Printf("usage: %s [grammarfile]\n", args.progName)
	os.Exit(2)
}

// stripNewlines consolidates contiguous newlines and returns a new
// slice of bytes with duplicated newlines omitted.
func stripNewlines(x []byte) []byte {
	buf := new(bytes.Buffer)
	last := byte('\n')
	for _, b := range x {
		if last != '\n' || b != '\n' {
			buf.WriteByte(b)
		}
		last = b
	}
	return buf.Bytes()
}

func main() {
	var (
		filename string
		r        io.Reader
	)

	if flag.NArg() == 0 {
		filename, r = "-", os.Stdin
	} else if flag.NArg() == 1 {
		filename = flag.Args()[0]
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		r = file
	} else {
		usage()
	}

	src, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	src = rebnf.CheckExtract(filename, src)

	if args.strip {
		src = stripNewlines(src)
	}

	_, err = io.Copy(os.Stdout, bytes.NewBuffer(src))
	if err != nil {
		log.Fatal(err)
	}
}
