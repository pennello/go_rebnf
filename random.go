// chris 072815

package rebnf

import (
	"errors"
	"io"
	"unicode"

	mathrand "math/rand"
	"unicode/utf8"

	fixrand "chrispennello.com/go/util/fix/math/rand"
	"golang.org/x/exp/ebnf"
)

const (
	maxRepetitions = 100
	maxRecursionDepth = 100
)

// ErrNoStart is returned by Random when the specified start production
// cannot be found in the grammar.
var ErrNoStart = errors.New("start production not found")

// Random generates random productions of the given grammar starting at
// the given start production, and writes them into the destination
// io.Writer.
func Random(dst io.Writer, grammar ebnf.Grammar, start string) error {
	prod, ok := grammar[start]
	if !ok {
		return ErrNoStart
	}
	return random(dst, grammar, prod.Expr, 0)
}

// IsCapital returns a boolean indicating whether or not the first rune
// of the given string is upper case.
func IsCapital(s string) bool {
	ch, _ := utf8.DecodeRuneInString(s)
	return !unicode.IsUpper(ch)
}

// IsTerminal returns a boolean that indicates whether the given
// Expression is a terminal one.  Ranges and Tokens are unconditionally
// considered to be terminal, and Names are terminal iff they're
// capitalized.  Productions are not considered because Alternatives
// contain Names, and you have to loo up the production by name in the
// grammar--it's just not a use case handled by this library, but could
// be added easily if needed.
func IsTerminal(expr ebnf.Expression) bool {
	switch expr.(type) {
	case *ebnf.Name:
		name := expr.(*ebnf.Name)
		return !IsCapital(name.String)
	case *ebnf.Range:
		return true
	case *ebnf.Token:
		return true
	default:
		return false
	}
}

// findTerminals is like filter(IsTerminal, exprs).  It ranges over all
// of the given expressions and produces a new slice of expressions
// containing only those for which IsTerminal returns true.
func findTerminals(exprs []ebnf.Expression) []ebnf.Expression {
	r := make([]ebnf.Expression, 0, len(exprs))
	for _, expr := range exprs {
		if IsTerminal(expr) {
			r = append(r, expr)
		}
	}
	return r
}

func random(dst io.Writer, grammar ebnf.Grammar, expr ebnf.Expression, depth int) error {
	switch expr.(type) {
	case ebnf.Alternative:
		// Choose a random alternative.
		alt := expr.(ebnf.Alternative)
		var exprs []ebnf.Expression
		// If maximum recursion depth has been exceeded, attempt
		// to select from only terminal expressions.
		if depth > maxRecursionDepth {
			exprs = findTerminals(alt)
			if len(exprs) == 0 {
				// No luck, we have no choice but to
				// explore one of the non-terminals in
				// this alternative.
				exprs = alt
			}
		} else {
			exprs = alt
		}
		err := random(dst, grammar, exprs[mathrand.Intn(len(exprs))], depth + 1)
		if err != nil {
			return err
		}

	case *ebnf.Group:
		// Evalute the group.
		gr := expr.(*ebnf.Group)
		err := random(dst, grammar, gr.Body, depth + 1)
		if err != nil {
			return err
		}

	case *ebnf.Name:
		// The name refers to a production; look it up and
		// continue the recursion.
		name := expr.(*ebnf.Name)
		err := random(dst, grammar, grammar[name.String], depth + 1)
		if err != nil {
			return err
		}

	case *ebnf.Option:
		// Randomly include the option.
		opt := expr.(*ebnf.Option)
		// If recursion depth has been exceeded, and option is
		// non-termainl, unconditionally omit.
		if depth > maxRecursionDepth && !IsTerminal(opt.Body) {
			// Omit.
		} else if fixrand.Bool() {
			// Otherwise, proceed with usual random
			// inclusion of option.
			err := random(dst, grammar, opt.Body, depth + 1)
			if err != nil {
				return err
			}
		}

	case *ebnf.Production:
		// Produce the production.
		prod := expr.(*ebnf.Production)
		err := random(dst, grammar, prod.Expr, depth + 1)
		if err != nil {
			return err
		}

	case *ebnf.Range:
		rng := expr.(*ebnf.Range)
		ch, err := fixrand.ChooseString(rng.Begin.String, rng.End.String)
		if err != nil {
			return err
		}
		if _, err := io.WriteString(dst, ch); err != nil {
			return err
		}

	case *ebnf.Repetition:
		// Choose a random number of repetitions.
		rep := expr.(*ebnf.Repetition)
		// If the recursion depth has been exceeded, and the
		// repetition is non-terminal, unconditionally omit it.
		if depth > maxRecursionDepth && !IsTerminal(rep.Body) {
			// Omit.
		} else {
			// Otherwise, do normal inclusion of a random
			// number of repetitions.
			for i := 0; i < mathrand.Intn(maxRepetitions+1); i++ {
				err := random(dst, grammar, rep.Body, depth + 1)
				if err != nil {
					return err
				}
			}
		}

	case ebnf.Sequence:
		// Recurse on each of the expressions.
		seq := expr.(ebnf.Sequence)
		for _, e := range seq {
			err := random(dst, grammar, e, depth + 1)
			if err != nil {
				return err
			}
		}

	case *ebnf.Token:
		// Emit the token.
		tok := expr.(*ebnf.Token)
		if _, err := io.WriteString(dst, tok.String); err != nil {
			return err
		}
	}
	return nil
}
