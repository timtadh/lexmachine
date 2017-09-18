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
	parsed := "(Match (Alternation (Concat (Character a), (Character b), (? (Alternation (Character a), (Alternation (Character c), (Character d)))), (Character w), (* (Character e)), (Character \\), (Character [), (Character .), (Range 0 255), (+ (Range 102 115))), (Concat (Character q), (Character y), (Character x))))"
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
			t.Errorf("got %q expected %q", m, expected[i])
		}
		i++
	}
	t.Assert(i == len(expected), "unconsumed matches %v", expected[i:])
}

func t_nomatch(program inst.InstSlice, text string, t *test.T) {
	scan := machines.LexerEngine(program, []byte(text))
	for tc, m, err, scan := scan(0); scan != nil; tc, m, err, scan = scan(tc) {
		if err == nil {
			t.Errorf("expected no match got %q, for %q", m, text)
		} else {
			break
		}
	}
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

func TestDigitClass(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("\\d+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (+ (Range 48 57)))"
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
	t_match(program, "0123456789", t)
	t_nomatch(program, "a234", t)
}

func TestNotDigitClass(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("\\D+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (+ (Alternation (Range 0 47), (Range 58 255))))"
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
	t_match(program, "wacky wizards", t)
	t_nomatch(program, "234", t)
}

func TestSpaceClass(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("\\s+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (+ (Alternation (Range 9 10), (Alternation (Range 12 13), (Range 32 32)))))"
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
	t_match(program, " \t\f\r", t)
	t_nomatch(program, "\vasdf", t)
}

func TestNoSpaceClass(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("\\S+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (+ (Alternation (Range 0 8), (Alternation (Range 11 11), (Alternation (Range 14 31), (Range 33 255))))))"
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
	t_match(program, "\vasdf", t)
	t_nomatch(program, " \t\f\r", t)
}

func TestWordClass(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("\\w+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (+ (Alternation (Range 48 57), (Alternation (Range 65 90), (Alternation (Range 95 95), (Range 97 122))))))"
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
	t_match(program, "asdf_asdf", t)
	t_nomatch(program, " asdf", t)
	t_nomatch(program, "@#$@#$", t)
}

func TestNoWordClass(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("\\W+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (+ (Alternation (Range 0 47), (Alternation (Range 58 64), (Alternation (Range 91 94), (Alternation (Range 96 96), (Range 123 255)))))))"
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
	t_match(program, " @#$", t)
	t_nomatch(program, "asdf_asdf", t)
}

func TestMultiRangeClasses(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("([a-zA-Z])([a-zA-Z0-9_])*"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (Alternation (Range 65 90), (Range 97 122)), (* (Alternation (Range 48 57), (Alternation (Range 65 90), (Alternation (Range 95 95), (Range 97 122)))))))"
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

func TestMultiRangeClasses2(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("([\\._/:a-zA-Z]+):\"(.+)\""))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (+ (Alternation (Range 46 47), (Alternation (Range 58 58), (Alternation (Range 65 90), (Alternation (Range 95 95), (Range 97 122)))))), (Character :), (Character \"), (+ (Range 0 255)), (Character \")))"
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
	t_match(program, ".X:\"a\"", t)
	t_nomatch(program, ".X:a\"", t)
}

func TestInvertRangeClasses1(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("[^abcd]+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (+ (Alternation (Range 0 96), (Range 101 255))))"
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
	t_match(program, "oiwe", t)
	t_match(program, "ef", t)
	t_match(program, "fin", t)
	t_nomatch(program, "a", t)
	t_nomatch(program, "b", t)
	t_nomatch(program, "c", t)
	t_nomatch(program, "d", t)
}

func TestInvertRangeClasses2(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("[^a-d]+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (+ (Alternation (Range 0 96), (Range 101 255))))"
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
	t_match(program, "oiwe", t)
	t_match(program, "ef", t)
	t_match(program, "fin", t)
	t_nomatch(program, "a", t)
	t_nomatch(program, "b", t)
	t_nomatch(program, "c", t)
	t_nomatch(program, "d", t)
}

func TestInvertRangeClasses3(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("[^a-dxyz]+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (+ (Alternation (Range 0 96), (Alternation (Range 101 119), (Range 123 255)))))"
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
	t_match(program, "oiwe", t)
	t_match(program, "ef", t)
	t_match(program, "fin", t)
	t_nomatch(program, "a", t)
	t_nomatch(program, "b", t)
	t_nomatch(program, "c", t)
	t_nomatch(program, "d", t)
	t_nomatch(program, "x", t)
	t_nomatch(program, "y", t)
	t_nomatch(program, "z", t)
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
