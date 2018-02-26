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
	parsed := "(Match (Concat (Alternation (Concat (Character a), (Character b), (? (Alternation (Character a), (Alternation (Character c), (Character d)))), (Character w), (* (Character e)), (Character \\), (Character [), (Character .), (Range 0 255), (+ (Range 102 115))), (Concat (Character q), (Character y), (Character x))), (EOS)))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
}

func tMatch(program inst.Slice, text string, t *test.T) {
	expected := []machines.Match{{len(program) - 1, 0, 1, 1, 1, len(text), []byte(text)}}
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

func tNoMatch(program inst.Slice, text string, t *test.T) {
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
	parsed := "(Match (Concat (Alternation (Character A), (Concat (Alternation (Character C), (Alternation (Character D), (Character E))), (Alternation (Character F), (Character G)), (Alternation (Character H), (Character I)), (Character B))), (EOS)))"
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
	tMatch(program, "A", t)
	tMatch(program, "CFHB", t)
	tMatch(program, "CFIB", t)
	tMatch(program, "CGHB", t)
	tMatch(program, "CGIB", t)
	tMatch(program, "DFHB", t)
	tMatch(program, "DFIB", t)
	tMatch(program, "DGHB", t)
	tMatch(program, "DGIB", t)
	tMatch(program, "EFHB", t)
	tMatch(program, "EFIB", t)
	tMatch(program, "EGHB", t)
	tMatch(program, "EGIB", t)
}

func TestParseConcatAltMaybes(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("((A?)?|(B|C))(D|E?)"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (Concat (Alternation (? (? (Character A))), (Alternation (Character B), (Character C))), (Alternation (Character D), (? (Character E)))), (EOS)))"
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
	tNoMatch(program, "", t) // will get empty string error
	tMatch(program, "E", t)
	tMatch(program, "D", t)
	tMatch(program, "A", t)
	tMatch(program, "AE", t)
	tMatch(program, "AD", t)
	tMatch(program, "B", t)
	tMatch(program, "BE", t)
	tMatch(program, "BD", t)
	tMatch(program, "C", t)
	tMatch(program, "CE", t)
	tMatch(program, "CD", t)
}

func TestParseConcatAltPlus(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("(A|(B|C))+(D|E?)"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (Concat (+ (Alternation (Character A), (Alternation (Character B), (Character C)))), (Alternation (Character D), (? (Character E)))), (EOS)))"
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
	tMatch(program, "A", t)
	tMatch(program, "AAA", t)
	tMatch(program, "AAABBCC", t)
	tMatch(program, "AAABBCC", t)
	tMatch(program, "AAABBCCD", t)
}

func TestParseAltOps(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a|b"))
	t.AssertNil(err)
	parsed := "(Match (Concat (Alternation (Character a), (Character b)), (EOS)))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}

	ast, err = Parse([]byte("a?|b"))
	t.AssertNil(err)
	parsed = "(Match (Concat (Alternation (? (Character a)), (Character b)), (EOS)))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}

	ast, err = Parse([]byte("a|b?"))
	t.AssertNil(err)
	parsed = "(Match (Concat (Alternation (Character a), (? (Character b))), (EOS)))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}

	ast, err = Parse([]byte("a?|b?"))
	t.AssertNil(err)
	parsed = "(Match (Concat (Alternation (? (Character a)), (? (Character b))), (EOS)))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
}

func TestChainedOps(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("A?+*B*?+C+*?(x+?)**"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (Concat (* (+ (? (Character A)))), (+ (? (* (Character B)))), (? (* (+ (Character C)))), (* (* (? (+ (Character x)))))), (EOS)))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
}

func TestParseConcatAltStar(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("(A|[C-G])*(X|Y?)"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (Concat (* (Alternation (Character A), (Range 67 71))), (Alternation (Character X), (? (Character Y)))), (EOS)))"
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
	tNoMatch(program, "", t) // will get empty string error
	tMatch(program, "X", t)
	tMatch(program, "Y", t)
	tMatch(program, "A", t)
	tMatch(program, "AAA", t)
	tMatch(program, "AAACC", t)
	tMatch(program, "AAACC", t)
	tMatch(program, "AAACCFFF", t)
	tMatch(program, "CAACCGEDFX", t)
}

func TestIdent(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("([a-z]|[A-Z])([a-z]|[A-Z]|[0-9]|_)*"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (Concat (Alternation (Range 97 122), (Range 65 90)), (* (Alternation (Range 97 122), (Alternation (Range 65 90), (Alternation (Range 48 57), (Character _)))))), (EOS)))"
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
	tMatch(program, "X", t)
	tMatch(program, "asdfY0923", t)
	tMatch(program, "A", t)
	tMatch(program, "AAA", t)
	tMatch(program, "AAACC", t)
}

func TestDigitClass(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("\\d+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (+ (Range 48 57)), (EOS)))"
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
	tMatch(program, "0123456789", t)
	tNoMatch(program, "a234", t)
}

func TestNotDigitClass(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("\\D+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (+ (Alternation (Range 0 47), (Range 58 255))), (EOS)))"
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
	tMatch(program, "wacky wizards", t)
	tNoMatch(program, "234", t)
}

func TestSpaceClass(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("\\s+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (+ (Alternation (Range 9 10), (Alternation (Range 12 13), (Range 32 32)))), (EOS)))"
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
	tMatch(program, " \t\f\r", t)
	tNoMatch(program, "\vasdf", t)
}

func TestNoSpaceClass(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("\\S+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (+ (Alternation (Range 0 8), (Alternation (Range 11 11), (Alternation (Range 14 31), (Range 33 255))))), (EOS)))"
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
	tMatch(program, "\vasdf", t)
	tNoMatch(program, " \t\f\r", t)
}

func TestWordClass(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("\\w+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (+ (Alternation (Range 48 57), (Alternation (Range 65 90), (Alternation (Range 95 95), (Range 97 122))))), (EOS)))"
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
	tMatch(program, "asdf_asdf", t)
	tNoMatch(program, " asdf", t)
	tNoMatch(program, "@#$@#$", t)
}

func TestNoWordClass(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("\\W+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (+ (Alternation (Range 0 47), (Alternation (Range 58 64), (Alternation (Range 91 94), (Alternation (Range 96 96), (Range 123 255)))))), (EOS)))"
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
	tMatch(program, " @#$", t)
	tNoMatch(program, "asdf_asdf", t)
}

func TestMultiRangeClasses(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("([a-zA-Z])([a-zA-Z0-9_])*"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (Concat (Alternation (Range 65 90), (Range 97 122)), (* (Alternation (Range 48 57), (Alternation (Range 65 90), (Alternation (Range 95 95), (Range 97 122)))))), (EOS)))"
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
	tMatch(program, "X", t)
	tMatch(program, "asdfY0923", t)
	tMatch(program, "A", t)
	tMatch(program, "AAA", t)
	tMatch(program, "AAACC", t)
}

func TestMultiRangeClasses2(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("([\\._/:a-zA-Z]+):\"(.+)\""))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (Concat (+ (Alternation (Range 46 47), (Alternation (Range 58 58), (Alternation (Range 65 90), (Alternation (Range 95 95), (Range 97 122)))))), (Character :), (Character \"), (+ (Range 0 255)), (Character \")), (EOS)))"
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
	tMatch(program, ".X:\"a\"", t)
	tNoMatch(program, ".X:a\"", t)
}

func TestInvertRangeClasses1(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("[^abcd]+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (+ (Alternation (Range 0 96), (Range 101 255))), (EOS)))"
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
	tMatch(program, "X", t)
	tMatch(program, "oiwe", t)
	tMatch(program, "ef", t)
	tMatch(program, "fin", t)
	tNoMatch(program, "a", t)
	tNoMatch(program, "b", t)
	tNoMatch(program, "c", t)
	tNoMatch(program, "d", t)
}

func TestInvertRangeClasses2(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("[^a-d]+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (+ (Alternation (Range 0 96), (Range 101 255))), (EOS)))"
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
	tMatch(program, "X", t)
	tMatch(program, "oiwe", t)
	tMatch(program, "ef", t)
	tMatch(program, "fin", t)
	tNoMatch(program, "a", t)
	tNoMatch(program, "b", t)
	tNoMatch(program, "c", t)
	tNoMatch(program, "d", t)
}

func TestInvertRangeClasses3(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("[^a-dxyz]+"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (+ (Alternation (Range 0 96), (Alternation (Range 101 119), (Range 123 255)))), (EOS)))"
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
	tMatch(program, "X", t)
	tMatch(program, "oiwe", t)
	tMatch(program, "ef", t)
	tMatch(program, "fin", t)
	tNoMatch(program, "a", t)
	tNoMatch(program, "b", t)
	tNoMatch(program, "c", t)
	tNoMatch(program, "d", t)
	tNoMatch(program, "x", t)
	tNoMatch(program, "y", t)
	tNoMatch(program, "z", t)
}

func TestLineComment(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("//[^\n]*"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := "(Match (Concat (Concat (Character /), (Character /), (* (Alternation (Range 0 9), (Range 11 255)))), (EOS)))"
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
	tMatch(program, "// adfawefawe awe", t)
}
