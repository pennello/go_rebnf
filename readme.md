TODO
----
 - Add CR license.
 - Reach out to Go mailing list; mention weird issues with random
   productions of Go grammar.

`rebnf` is a tool to generate random productions of arbitrary EBNF
grammars.

Documentation
-------------
 - [GoDoc Documentation](https://godoc.org/chrispennello.com/go/rebnf)

Tools
-----
The `tools` directory contains a Python tool for fetching Unicode
character ranges from `fileformat.info` and consolidating them into EBNF
alternative range expressions, suitable for consumption by `rebnf`.

Test Languages
--------------
The `testlang` directory contains grammars for testing and
experimentation.  In particular, `testlang/go` contains a snapshot of
the Go Programming Language Specification and an extracted grammar,
including one in which the Unicode productions have been defined.

Future Work
-----------
 - Non-uniform production weighting.
 - Mode to prefer non-weird characters (e.g., put it in "ASCII mode").
