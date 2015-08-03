// chris 072815

package rebnf

import (
	"errors"
	"fmt"
	"io"
	"log"
	"unicode"

	mathrand "math/rand"
	"unicode/utf8"

	fixrand "chrispennello.com/go/util/fix/math/rand"
	"golang.org/x/exp/ebnf"
)

// ErrNoStart is returned by Random when the specified start production
// cannot be found in the grammar.
var ErrNoStart = errors.New("start production not found")

// Random generates random productions of the given grammar starting at
// the given start production, and writes them into the destination
// io.Writer.
func (c *Ctx) Random(dst io.Writer, grammar ebnf.Grammar, start string) error {
	prod, ok := grammar[start]
	if !ok {
		return ErrNoStart
	}
	return c.random(dst, grammar, prod.Expr, 0)
}

// IsCapital returns a boolean indicating whether or not the first rune
// of the given string is upper case.
func IsCapital(s string) bool {
	ch, _ := utf8.DecodeRuneInString(s)
	return unicode.IsUpper(ch)
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

// pad picks a random rune from the padding string specified by the
// context and writes it to the destination writer.
func (c *Ctx) pad(dst io.Writer) error {
	runes := []rune(c.padding)
	if len(runes) == 0 {
		return nil
	}
	r := runes[mathrand.Intn(len(runes))]
	_, err := io.WriteString(dst, string([]rune{r}))
	return err
}

func (c *Ctx) log(format string, a ...interface{}) {
	if !c.debug {
		return
	}
	// Logging errors ignored.
	log.Printf(format, a...)
}

// random is the inner, recursive implementation of Random.  It handles
// each of the ebnf.Expression implementations, outputting productions
// randomly to the destination writer.  It implements a recursion depth
// counter, and once the counter exceeds the limit, it favors producing
// terminals over non-terminals.  Note that this does not guarantee
// termination, however.  For example, the pathological grammar "S = S"
// will still loop forever.
func (c *Ctx) random(dst io.Writer, grammar ebnf.Grammar, expr ebnf.Expression, depth int) error {
	c.log("recursion depth %d\n", depth)
	c.log("%#v\n\n", expr)

	switch expr.(type) {
	// Choose a random alternative.
	case ebnf.Alternative:
		alt := expr.(ebnf.Alternative)
		var exprs []ebnf.Expression
		// If maximum recursion depth has been exceeded, attempt
		// to select from only terminal expressions.
		if depth > c.maxdepth {
			exprs = findTerminals(alt)
			c.log("alternative: found %d terminals", len(exprs))
			if len(exprs) == 0 {
				// No luck, we have no choice but to
				// explore one of the non-terminals in
				// this alternative.
				c.log("alternative: no terminals\n")
				exprs = alt
			}
		} else {
			exprs = alt
		}
		err := c.random(dst, grammar, exprs[mathrand.Intn(len(exprs))], depth+1)
		if err != nil {
			return err
		}

	// Evalute the group.
	case *ebnf.Group:
		gr := expr.(*ebnf.Group)
		err := c.random(dst, grammar, gr.Body, depth+1)
		if err != nil {
			return err
		}

	// The name refers to a production; look it up and continue the
	// recursion.
	case *ebnf.Name:
		name := expr.(*ebnf.Name)
		// Pad non-terminals.
		pad := !IsTerminal(expr)
		if pad {
			c.pad(dst)
		}
		err := c.random(dst, grammar, grammar[name.String], depth+1)
		if err != nil {
			return err
		}
		if pad {
			c.pad(dst)
		}

	// Randomly include the option.
	case *ebnf.Option:
		opt := expr.(*ebnf.Option)
		// If recursion depth has been exceeded, and option is
		// non-termainl, unconditionally omit.
		if depth > c.maxdepth && !IsTerminal(opt.Body) {
			// Omit.
			c.log("option: non-terminal omitted due to having exceeded recursion depth limit\n")
		} else if fixrand.Bool() {
			// Otherwise, proceed with usual random
			// inclusion of option.
			err := c.random(dst, grammar, opt.Body, depth+1)
			if err != nil {
				return err
			}
		}

	// Produce the production.
	case *ebnf.Production:
		prod := expr.(*ebnf.Production)
		err := c.random(dst, grammar, prod.Expr, depth+1)
		if err != nil {
			return err
		}

	// Generate a random string in the given range.
	case *ebnf.Range:
		rng := expr.(*ebnf.Range)
		ch, err := fixrand.ChooseString(rng.Begin.String, rng.End.String)
		if err != nil {
			return err
		}
		if _, err := io.WriteString(dst, ch); err != nil {
			return err
		}

	// Choose a random number of repetitions.
	case *ebnf.Repetition:
		rep := expr.(*ebnf.Repetition)
		// If the recursion depth has been exceeded, and the
		// repetition is non-terminal, unconditionally omit it.
		if depth > c.maxdepth && !IsTerminal(rep.Body) {
			// Omit.
			c.log("repetition: non-terminal omitted due to having exceeded recursion depth limit\n")
		} else {
			// Otherwise, do normal inclusion of a random
			// number of repetitions.
			reps := mathrand.Intn(c.maxreps + 1)
			c.log("repetition: chose %d repetitions\n", reps)
			for i := 0; i < reps; i++ {
				err := c.random(dst, grammar, rep.Body, depth+1)
				if err != nil {
					return err
				}
			}
		}

	// Recurse on each of the expressions.
	case ebnf.Sequence:
		seq := expr.(ebnf.Sequence)
		for _, e := range seq {
			err := c.random(dst, grammar, e, depth+1)
			if err != nil {
				return err
			}
		}

	// Emit the token.
	case *ebnf.Token:
		tok := expr.(*ebnf.Token)
		if _, err := io.WriteString(dst, tok.String); err != nil {
			return err
		}

	default:
		// This indicates a bug in the code.
		panic(fmt.Sprintf("bad expression %#v", expr))
	}

	return nil
}
