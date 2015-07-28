package main

import (
	"bytes"
	"io"
	"os"

	"io/ioutil"

	"golang.org/x/exp/ebnf"
)

func Parse(name, start string, r io.Reader) (ebnf.Grammar, error) {
	if r == nil {
		f, err := os.Open(name)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	}

	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	src = CheckRead(name, src)

	grammar, err := ebnf.Parse(name, bytes.NewBuffer(src))
	if err != nil {
		return nil, err
	}

	return grammar, err
}
