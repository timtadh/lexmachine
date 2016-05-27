package frontend

import "testing"
import "github.com/timtadh/data-structures/test"

import (
	"github.com/timtadh/lexmachine/inst"
	"github.com/timtadh/lexmachine/machines"
)

func TestParse(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("ab(a|c|d)?we*\\\\\\[\\..[s-f]+|qyx"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Alternation (Concat (Character a), (Character b), (? (Alternation (Character a), (Alternation (Character c), (Character d)))), (Character w), (* (Character e)), (Character \\), (Character [), (Character .), (Range 0 255), (+ (Range 115 102))), (Concat (Character q), (Character y), (Character x))))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
}

func t_match(program inst.InstSlice, text string, t *test.T) {
	expected := []machines.Match{machines.Match{len(program) - 1, 0, 1, 1, 1, len(text), []byte(text)}}
	if expected[0].EndColumn == 0 {
		expected[0].EndColumn = 1
	}
	i := 0
	scan := machines.LexerEngine(program, []byte(text))
	for tc, m, err, scan := scan(0); scan != nil; tc, m, err, scan = scan(tc) {
		t.Log("match", m)
		if err != nil {
			t.Error("error", err)
		} else if !m.Equals(&expected[i]) {
			t.Error(m, expected[i])
		}
		i++
	}
	t.Assert(i == len(expected), "unconsumed matches %v", expected[i:])

	dfa := machines.ToDFA(program)
	t.Log(dfa)
	l := len(text)
	if l == 0 {
		l += 1
	}
	expected = []machines.Match{machines.Match{len(dfa)-1, 0, 1, 1, 1, l, []byte(text)}}
	i = 0
	scan = machines.DFALexerEngine(dfa, []byte(text))
	for tc, m, err, scan := scan(0); scan != nil; tc, m, err, scan = scan(tc) {
		t.Log("match", m)
		if err != nil {
			t.Error("error", err)
		} else if !m.Equals(&expected[i]) {
			t.Error(m, expected[i])
		}
		i++
	}
	t.Assert(i == len(expected), "unconsumed matches %v", expected[i:])
}

func TestParseConcatAlts(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("A|((C|D|E)(F|G)(H|I)B)"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Alternation (Character A), (Concat (Alternation (Character C), (Alternation (Character D), (Character E))), (Alternation (Character F), (Character G)), (Alternation (Character H), (Character I)), (Character B))))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
	program, err := Generate(ast)
	if err != nil {
		t.Error(err)
	}
	t.Log(program)
	t_match(program, "A", t)
	t_match(program, "CFHB", t)
	t_match(program, "CFIB", t)
	t_match(program, "CGHB", t)
	t_match(program, "CGIB", t)
	t_match(program, "DFHB", t)
	t_match(program, "DFIB", t)
	t_match(program, "DGHB", t)
	t_match(program, "DGIB", t)
	t_match(program, "EFHB", t)
	t_match(program, "EFIB", t)
	t_match(program, "EGHB", t)
	t_match(program, "EGIB", t)
}

func TestParseConcatAltMaybes(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("((A?)?|(B|C))(D|E?)"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (Alternation (? (? (Character A))), (Alternation (Character B), (Character C))), (Alternation (Character D), (? (Character E)))))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
	program, err := Generate(ast)
	if err != nil {
		t.Error(err)
	}
	t.Log(program)
	t_match(program, "", t)
	t_match(program, "E", t)
	t_match(program, "D", t)
	t_match(program, "A", t)
	t_match(program, "AE", t)
	t_match(program, "AD", t)
	t_match(program, "B", t)
	t_match(program, "BE", t)
	t_match(program, "BD", t)
	t_match(program, "C", t)
	t_match(program, "CE", t)
	t_match(program, "CD", t)
}

func TestParseConcatAltPlus(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("(A|(B|C))+(D|E?)"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (+ (Alternation (Character A), (Alternation (Character B), (Character C)))), (Alternation (Character D), (? (Character E)))))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
	program, err := Generate(ast)
	if err != nil {
		t.Error(err)
	}
	t.Log(program)
	t_match(program, "A", t)
	t_match(program, "AAA", t)
	t_match(program, "AAABBCC", t)
	t_match(program, "AAABBCC", t)
	t_match(program, "AAABBCCD", t)
}

func TestParseConcatAltStar(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("(A|[C-G])*(X|Y?)"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (* (Alternation (Character A), (Range 67 71))), (Alternation (Character X), (? (Character Y)))))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
	program, err := Generate(ast)
	if err != nil {
		t.Error(err)
	}
	t.Log(program)
	t_match(program, "", t)
	t_match(program, "X", t)
	t_match(program, "Y", t)
	t_match(program, "A", t)
	t_match(program, "AAA", t)
	t_match(program, "AAACC", t)
	t_match(program, "AAACC", t)
	t_match(program, "AAACCFFF", t)
	t_match(program, "CAACCGEDFX", t)
}

func TestIdent(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("([a-z]|[A-Z])([a-z]|[A-Z]|[0-9]|_)*"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (Alternation (Range 97 122), (Range 65 90)), (* (Alternation (Range 97 122), (Alternation (Range 65 90), (Alternation (Range 48 57), (Character _)))))))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
	program, err := Generate(ast)
	if err != nil {
		t.Error(err)
	}
	t.Log(program)
	t_match(program, "X", t)
	t_match(program, "asdfY0923", t)
	t_match(program, "A", t)
	t_match(program, "AAA", t)
	t_match(program, "AAACC", t)
}

func TestLineComment(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("//[^\n]*"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (Character /), (Character /), (* (Alternation (Range 0 9), (Range 11 255)))))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
	program, err := Generate(ast)
	if err != nil {
		t.Error(err)
	}
	t.Log(program)
	t_match(program, "// adfawefawe awe", t)
}
