package dfa

import (
	"fmt"
	"testing"

	"github.com/timtadh/data-structures/test"
	"github.com/timtadh/lexmachine/frontend"
)

func TestGen(x *testing.T) {
	t := (*test.T)(x)
	// ast, err := frontend.Parse([]byte("(([a-z]+[A-Z])*[0-9])?wizard"))
	ast, err := frontend.Parse([]byte("(a|b)+[xyz]*"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ast.String())
	dfa := Generate(ast)
	fmt.Println(dfa)
}

func TestGenCharacter(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("a"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ast.String())
	dfa := Generate(ast)
	fmt.Println(dfa)
}

func TestGenRange(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("[a-c]"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ast.String())
	dfa := Generate(ast)
	fmt.Println(dfa)
}

func TestGenAltMatch(x *testing.T) {
	t := (*test.T)(x)
	ast := frontend.NewAltMatch(
		frontend.NewMatch(
			frontend.NewRange('a', 'd'),
		),
		frontend.NewMatch(
			frontend.NewAlternation(
				frontend.NewCharacter('b'),
				frontend.NewCharacter('e'),
			),
		),
	)
	t.Log(ast.String())
	dfa := Generate(ast)
	fmt.Println(dfa)
}
