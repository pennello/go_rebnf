// chris 072315

package main

import (
	"errors"
	"fmt"

	"math/rand"
)

var ErrBadRange = errors.New("bad range")

// ChooseInt chooses a random integer in [begin, end].
func ChooseInt(begin, end int) int {
	if begin > end {
		panic(fmt.Sprintf("%d < %d", begin, end))
	}
	return begin + rand.Intn(end - begin + 1)
}

// ChooseRune chooses a random rune that is lexically includively
// between begin and end.  Specifically, the output x satisfies begin <=
// x <= end.
func ChooseRune(begin, end rune) rune {
	return rune(ChooseInt(int(begin), int(end)))
}

// ChooseString chooses a random string that is lexically inclusively
// between begin and end.  Specifically, the output x satisfies begin <=
// x <= end.  If the begin and end strings don't contain the same number
// of runes, the empty string is returned with an error.
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
