// chris 072115

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
)

var args struct {
	prog     string
	start    string
	maxreps  int
	maxdepth int
}

func init() {
	log.SetFlags(0)
	rand.Seed(time.Now().UTC().UnixNano())

	start    := flag.String("start", "Start", "name of start production")
	maxreps  := flag.Int("maxreps", 100, "maximum number of repetitions")
	maxdepth := flag.Int("maxdepth", 100, "maximum recursion depth")

	flag.Parse()

	args.prog     = os.Args[0]
	args.start    = *start
	args.maxreps  = *maxreps
	args.maxdepth = *maxdepth
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

	grammar, err := rebnf.Parse(name, args.start, r)
	if err != nil {
		log.Fatal(err)
	}

	ctx := rebnf.NewCtx(args.maxreps, args.maxdepth)
	err = ctx.Random(os.Stdout, grammar, args.start)
	if err != nil {
		log.Fatal(err)
	}
}
