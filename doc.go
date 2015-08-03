// chris 073015

// Package rebnf implements random productions of arbitrary EBNF
// grammars.
//
// It uses golang.org/x/exp/ebnf to parse the grammars.  It also offers
// updated support for parsing EBNF grammars from HTML pages, a feature
// in exp/ebnf, but, at the time of writing, out of date with respect to
// the formatting of the Go Programming Language Specification HTML
// page.
//
// Usage
//
// Call Parse to extract a grammar from an input file or reader.
//
// Create a context to specify several options to control the random
// productions.
//
// First, recall EBNF repetitions.
//
//	Repetition = "{" Expression "}" .
//
// These are expressions that can be repeated 0 or more times.  The
// approach taken with these is to pick a random number of times to
// repeat such an expression, between 0 and the specified maximum number
// of repetitions.  This is one of the arguments you must provide when
// creating a context.  A reasonable value is 100.
//
// The grammar is randomly explored recursively.  Thus, another way you
// may control the random productions is via a recursion depth limit.
// Recall that exp/ebnf EBNF grammars distinguish between non-terminal
// and terminal productions by using capitalized or uncapitalized names.
// When the recursion depth limit is exceeded, the algorithm will make
// an effort to favor terminal symbols over non-terminal ones.  Note
// that this does NOT protect, you, however from pathological grammars
// such as "S = S", or other grammars that necessitate infinite
// productions.  This recursion depth limit is another one of the
// arguments you must specify when creating a context.  A reasonable
// value is 30.
//
// Several other context parameters include a whitespace padding feature
// for non-terminals and debug output.  See the Ctx documentation for
// more details.
//
// Pass the grammar and the name of the start production into Random, a
// function defined on contexts.  The output instance of the language
// will be written to the given destination writer.
//
// A command-line implementation is made available at cmd/rebnf.
package rebnf
