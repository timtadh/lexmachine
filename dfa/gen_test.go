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
	ast, err := frontend.Parse([]byte("(a|b)?[xyz]+c"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ast.String())
	dfa := Generate(ast)
	fmt.Println(dfa)
}
