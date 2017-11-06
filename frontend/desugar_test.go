package frontend

import "testing"
import "github.com/timtadh/data-structures/test"

func TestDesugarRanges_any(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("."))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ast.String())
	parsed := "(Match (Concat (Range 0 255), (EOS)))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
	// asserts this doesn't infinte loop
	DesugarRanges(ast)
}

func TestDesugarRanges(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("(([a-z]+[A-Z])*[0-9])?wizard"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ast.String())
	parsed := "(Match (Concat (Concat (? (Concat (* (Concat (+ (Range 97 122)), (Range 65 90))), (Range 48 57))), (Character w), (Character i), (Character z), (Character a), (Character r), (Character d)), (EOS)))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
	ast = DesugarRanges(ast)
	desugared := "(Match (Concat (Concat (? (Concat (* (Concat (+ (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Character a), (Character b)), (Character c)), (Character d)), (Character e)), (Character f)), (Character g)), (Character h)), (Character i)), (Character j)), (Character k)), (Character l)), (Character m)), (Character n)), (Character o)), (Character p)), (Character q)), (Character r)), (Character s)), (Character t)), (Character u)), (Character v)), (Character w)), (Character x)), (Character y)), (Character z))), (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Character A), (Character B)), (Character C)), (Character D)), (Character E)), (Character F)), (Character G)), (Character H)), (Character I)), (Character J)), (Character K)), (Character L)), (Character M)), (Character N)), (Character O)), (Character P)), (Character Q)), (Character R)), (Character S)), (Character T)), (Character U)), (Character V)), (Character W)), (Character X)), (Character Y)), (Character Z)))), (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Alternation (Character 0), (Character 1)), (Character 2)), (Character 3)), (Character 4)), (Character 5)), (Character 6)), (Character 7)), (Character 8)), (Character 9)))), (Character w), (Character i), (Character z), (Character a), (Character r), (Character d)), (EOS)))"
	if ast.String() != desugared {
		t.Log(ast.String())
		t.Log(desugared)
		t.Error("Did not desugar correctly")
	}
}
