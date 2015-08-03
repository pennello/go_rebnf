// chris 072115

// rebnf outputs random productions of arbitrary EBNF grammars.
//
// The program uses golang.org/x/exp/ebnf to parse the grammars and
// randomly explores the grammar, sending its output to standard out.
// It also offers updated support for parsing EBNF grammars from HTML
// pages, a feature in exp/ebnf, but, at the time of writing, out of
// date with respect to the formatting of the Go Programming Language
// Specification HTML page.
//
// There is an optional command-line argument to provide a file
// containing an EBNF grammar or HTML page with a grammar embedded
// within.  If the command-line argument is not provided, then the
// input will be read from standard in.
//
// There are several options to control the random productions.  You may
// specify the name of the start production.  You may also specify the
// seed used to seed the random number generator.  By default, rebnf
// will use the current system time, in UTC and in nanoseconds, as the
// seed.  In either case, the seed used will be printed to standard
// error before productions commence, preceeded by the word "seed".
//
// There are two more options that you may use to control the random
// productions.  First, recall EBNF repetitions.
//
//	Repetition = "{" Expression "}" .
//
// These are expressions that can be repeated 0 or more times.  The
// approach taken with these is to pick a random number of times to
// repeat such an expression, between 0 and the specified maximum number
// of repetitions.
//
// The grammar is randomly explored recursively.  Thus, another way you
// may control the random productions is via a recursion depth limit.
// Recall that exp/ebnf EBNF grammars distinguish between non-terminal
// and terminal productions by using capitalized or uncapitalized names.
// When the recursion depth limit is exceeded, the algorithm will make
// an effort to favor terminal symbols over non-terminal ones.  Note
// that this does NOT protect, you, however from pathological grammars
// such as "S = S", or other grammars that necessitate infinite
// productions.
//
// There is also an option to insert padding characters on either side
// of non-terminal symbols.  The default is to insert a space, but you
// may specify any string of padding characters for the algorithm to
// randomly choose from.
//
// Usage is as follows.
//
//	rebnf [-h] [options] [grammarfile]
//	  -maxdepth=30:   maximum recursion depth
//	  -maxreps=100:   maximum number of repetitions
//	  -padding=" ":   non-terminal padding characters
//	  -seed=-1:       random seed
//	  -start="Start": name of start production
//	  -debug=false:   grammar traversal debug output
//
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"math/rand"

	"chrispennello.com/go/rebnf"
	"golang.org/x/exp/ebnf"
)

var args struct {
	prog     string
	seed     int64
	start    string
	maxreps  int
	maxdepth int
	padding  string
	debug    bool
}

func init() {
	log.SetFlags(0)

	seed     := flag.Int64("seed", -1, "random seed")
	start    := flag.String("start", "Start", "name of start production")
	maxreps  := flag.Int("maxreps", 100, "maximum number of repetitions")
	maxdepth := flag.Int("maxdepth", 30, "maximum recursion depth")
	padding  := flag.String("padding", " ", "non-terminal padding characters")
        debug    := flag.Bool("debug", false, "grammar traversal debug output")
	flag.Parse()

	args.prog     = os.Args[0]
	args.seed     = *seed
	args.start    = *start
	args.maxreps  = *maxreps
	args.maxdepth = *maxdepth
	args.padding  = *padding
	args.debug    = *debug

	if args.seed == -1 {
		args.seed = time.Now().UTC().UnixNano()
	}
	rand.Seed(args.seed)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [flags] [filename]\n", args.prog)
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	var (
		name string
		r    io.Reader
	)
	switch flag.NArg() {
	case 0:
		name, r = "-", os.Stdin
	case 1:
		name = flag.Arg(0)
	default:
		usage()
	}

	grammar, err := rebnf.Parse(name, r)
	if err != nil {
		log.Fatal(err)
	}

	// Verify grammar before attempting any productions.
	err = ebnf.Verify(grammar, args.start)
	if err != nil {
		log.Fatal(err)
	}

	ctx := rebnf.NewCtx(args.maxreps, args.maxdepth, args.padding, args.debug)
	log.Printf("seed %d", args.seed)
	err = ctx.Random(os.Stdout, grammar, args.start)
	if err != nil {
		log.Fatal(err)
	}
}
