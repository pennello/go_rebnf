TODO
----
 - Write out actual ranges for Go grammar character productions.
 - Make readme publicly consumable.
 - Add CR license.
 - Document test languages in repo.
 - Figure out why some JSON productions are invalid (according to
   `python -m json.tool`).

Tools
-----
The `tools` directory contains a Python tool for fetching Unicode
character ranges from `fileformat.info` and consolidating them into EBNF
alternative range expressions, suitable for consumption by `rebnf`.

Future Work
-----------
 - Non-uniform production weighting.
 - Mode to prefer non-weird characters.
