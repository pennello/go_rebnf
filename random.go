package main

import (
	"errors"
	"io"

	"math/rand"

	"golang.org/x/exp/ebnf"
)

var ErrNoStart = errors.New("start production not found")

const maxRepetitions = 100

func Random(dst io.Writer, grammar ebnf.Grammar, start string) error {
	prod, ok := grammar[start]
	if !ok {
		return ErrNoStart
	}
	return random(dst, grammar, prod.Expr)
}

func random(dst io.Writer, grammar ebnf.Grammar, expr ebnf.Expression) error {
	switch expr.(type) {
	case ebnf.Alternative:
		// Choose a random alternative.
		alt := expr.(ebnf.Alternative)
		err := random(dst, grammar, alt[rand.Intn(len(alt))])
		if err != nil {
			return err
		}

	case *ebnf.Group:
		// Evalute the group.
		gr := expr.(*ebnf.Group)
		if err := random(dst, grammar, gr.Body); err != nil {
			return err
		}

	case *ebnf.Name:
		// The name refers to a production; look it up and
		// continue the recursion.
		name := expr.(*ebnf.Name)
		if err := random(dst, grammar, grammar[name.String]); err != nil {
			return err
		}

	case *ebnf.Option:
		// Randomly include the option.
		opt := expr.(*ebnf.Option)
		if 0 == rand.Intn(2) { // 50/50 chance.
			if err := random(dst, grammar, opt.Body); err != nil {
				return err
			}
		}

	case *ebnf.Production:
		// Produce the production.
		prod := expr.(*ebnf.Production)
		if err := random(dst, grammar, prod.Expr); err != nil {
			return err
		}

	case *ebnf.Range:
		rng := expr.(*ebnf.Range)
		ch, err := ChooseString(rng.Begin.String, rng.End.String)
		if err != nil {
			return err
		}
		if _, err := io.WriteString(dst, ch); err != nil {
			return err
		}

	case *ebnf.Repetition:
		// Choose a random number of repetitions.
		rep := expr.(*ebnf.Repetition)
		for i := 0; i < rand.Intn(maxRepetitions + 1); i++ {
			if err := random(dst, grammar, rep.Body); err != nil {
				return err
			}
		}

	case ebnf.Sequence:
		// Recurse on each of the expressions.
		seq := expr.(ebnf.Sequence)
		for _, e := range seq {
			if err := random(dst, grammar, e); err != nil {
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
