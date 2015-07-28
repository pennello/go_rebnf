// chris 072315

package main

import (
	"errors"
	"fmt"

	"math/rand"
)

// ErrBadRange is returned by ChooseString when the input strings do not
// contain the same number of runes.
var ErrBadRange = errors.New("bad range")

// ChooseInt returns, as an int, a random integer in [begin, end] from
// the default Source.  It panics if begin > end.
func ChooseInt(begin, end int) int {
	if begin > end {
		panic(fmt.Sprintf("invalid arguments to ChooseInt: begin %d > end %d", begin, end))
	}
	return begin + rand.Intn(end-begin+1)
}

// ChooseRune returns a random rune that is lexically and inclusively
// between begin and end.  Specifically, the output x satisfies
// begin <= x <= end.
func ChooseRune(begin, end rune) rune {
	return rune(ChooseInt(int(begin), int(end)))
}

// ChooseString returns a random string that is lexically and
// inclusively between begin and end.  Specifically, the output x
// satisfies begin <= x <= end.  If the begin and end strings don't
// contain the same number of runes, the empty string is returned with
// ErrBadRange.
func ChooseString(begin, end string) (string, error) {
	br := []rune(begin)
	er := []rune(end)
	if len(br) != len(er) {
		return "", ErrBadRange
	}
	ret := make([]rune, len(br))
	for i := range ret {
		ret[i] = ChooseRune(br[i], er[i])
	}
	return string(ret), nil
}
