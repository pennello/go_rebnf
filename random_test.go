// chris 080415

package rebnf

import (
	"testing"

	"golang.org/x/exp/ebnf"
)

func testIsCapital(t *testing.T, x string, expect bool) {
	out := IsCapital(x)
	if out != expect {
		t.Errorf("IsCapital(%q) != %v (got %v)", x, out, expect)
	}
}

func TestIsCapital(t *testing.T) {
	testIsCapital(t, "Hello", true)
	testIsCapital(t, "there", false)
	testIsCapital(t, "", false)
	testIsCapital(t, "Σ1", true)
}

func testIsTerminal(t *testing.T, expr ebnf.Expression, expect bool) {
	out := IsTerminal(expr)
	if out != expect {
		t.Errorf("IsTerminal(%#v) != %v (got %v)", expr, out, expect)
	}
}

func TestIsTerminal(t *testing.T) {
	var expr ebnf.Expression
	var t1, t2 ebnf.Token

	expr = &ebnf.Name{String: "Production"}
	testIsTerminal(t, expr, false)

	expr = &ebnf.Name{String: "token"}
	testIsTerminal(t, expr, true)

	expr = &ebnf.Token{String: "blah"}
	testIsTerminal(t, expr, true)

	t1 = ebnf.Token{String: "a"}
	t2 = ebnf.Token{String: "z"}
	expr = &ebnf.Range{Begin: &t1, End: &t2}
	testIsTerminal(t, expr, true)

	expr = &ebnf.Repetition{Body: &ebnf.Token{String: "foo"}}
	testIsTerminal(t, expr, false)
}

func sliceExprEquals(a, b []ebnf.Expression) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func testFindTerminals(t *testing.T, x []ebnf.Expression, expect []ebnf.Expression) {
	out := findTerminals(x)
	if !sliceExprEquals(out, expect) {
		t.Errorf("findTerminals(%#v) != %#v (got %#v)", x, expect, out)
	}
}

func TestFindTerminals(t *testing.T) {
	var term ebnf.Expression
	var x, expect []ebnf.Expression

	term = &ebnf.Name{String: "some_terminal"}

	x = []ebnf.Expression{
		term,
		&ebnf.Name{String: "NonTerminal"},
		&ebnf.Option{Body: &ebnf.Name{String: "whatever"}},
	}
	expect = []ebnf.Expression{term}
	testFindTerminals(t, x, expect)
}
