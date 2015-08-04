// chris 080415

package rebnf

import (
	"testing"
)

func testIsCapital(t *testing.T, x string, expect bool) {
	if IsCapital(x) != expect {
		t.Fail()
	}
}

func TestIsCapital(t *testing.T) {
	testIsCapital(t, "Hello", true)
	testIsCapital(t, "there", false)
	testIsCapital(t, "", false)
	testIsCapital(t, "Σ1", true)
}
