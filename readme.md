TODO
----
 - Move exported functions to importable (non-main) package.
 - GoDoc methods, etc.
 - On error, output chosen seed.  Accept seed as an input so you can
   reproduce errors.
 - Bound output.  Recursion depth limit?  Need to understand how to get
   to only terminal symbols.  Graph analysis of grammar?
 - Document how library does not avoid infinite productions (i.e.,
   `testlang/inf`).
 - Context for `maxRepetitions` to allow customization.
