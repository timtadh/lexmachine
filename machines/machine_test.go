package machines

import "testing"
import "strings"
import "github.com/timtadh/lexmachine/inst"

func TestLexerMatch(t *testing.T) {
	text := []byte("ababcbcbb")
	//. (a|b)*cba?(c|b)bb
	program := make(inst.InstSlice, 20)

	program[0] = inst.New(inst.SPLIT, 1, 6)
	program[1] = inst.New(inst.SPLIT, 2, 4)
	program[2] = inst.New(inst.CHAR, 'a', 'a')
	program[3] = inst.New(inst.JMP, 5, 0)
	program[4] = inst.New(inst.CHAR, 'b', 'b')
	program[5] = inst.New(inst.JMP, 0, 0)
	program[6] = inst.New(inst.CHAR, 'c', 'c')
	program[7] = inst.New(inst.CHAR, 'b', 'b')
	program[8] = inst.New(inst.SPLIT, 9, 10)
	program[9] = inst.New(inst.CHAR, 'a', 'a')
	program[10] = inst.New(inst.SPLIT, 11, 13)
	program[11] = inst.New(inst.CHAR, 'c', 'c')
	program[12] = inst.New(inst.JMP, 14, 0)
	program[13] = inst.New(inst.CHAR, 'b', 'b')
	program[14] = inst.New(inst.CHAR, 'b', 'b')
	program[15] = inst.New(inst.CHAR, 'b', 'b')
	program[16] = inst.New(inst.MATCH, 0, 0)

	t.Log(string(text))
	t.Log(len(text))
	t.Log(program)
	expected := []Match{
		Match{16, []byte("ababcbcbb")},
	}
	i := 0
	for tc, m, err, scan := LexerEngine(program, text)(0); scan != nil; tc, m, err, scan = scan(tc) {
		t.Log(m)
		if err != nil {
			t.Error(err)
		}
		if !m.Equals(&expected[i]) {
			t.Error(m, expected[i])
		}
		i++
	}
	if i != len(expected) {
		t.Error("unconsumed matches", expected[i:])
	}
}

func TestLexerNoMatch(t *testing.T) {
	text := []byte("ababcbcb")
	//. (a|b)*cba?(c|b)bb
	program := make(inst.InstSlice, 20)

	program[0] = inst.New(inst.SPLIT, 1, 6)
	program[1] = inst.New(inst.SPLIT, 2, 4)
	program[2] = inst.New(inst.CHAR, 'a', 'a')
	program[3] = inst.New(inst.JMP, 5, 0)
	program[4] = inst.New(inst.CHAR, 'b', 'b')
	program[5] = inst.New(inst.JMP, 0, 0)
	program[6] = inst.New(inst.CHAR, 'c', 'c')
	program[7] = inst.New(inst.CHAR, 'b', 'b')
	program[8] = inst.New(inst.SPLIT, 9, 10)
	program[9] = inst.New(inst.CHAR, 'a', 'a')
	program[10] = inst.New(inst.SPLIT, 11, 13)
	program[11] = inst.New(inst.CHAR, 'c', 'c')
	program[12] = inst.New(inst.JMP, 14, 0)
	program[13] = inst.New(inst.CHAR, 'b', 'b')
	program[14] = inst.New(inst.CHAR, 'b', 'b')
	program[15] = inst.New(inst.CHAR, 'b', 'b')
	program[16] = inst.New(inst.MATCH, 0, 0)

	t.Log("(a|b)*cba?(c|b)bb")
	t.Log(string(text))
	t.Log(program)

	for tc, m, err, scan := LexerEngine(program, text)(0); scan != nil; tc, m, err, scan = scan(tc) {
		if err == nil || !strings.HasPrefix(err.Error(), "Unconsumed text") {
			t.Error("no error!", m, err)
		}
	}
}

func TestLexerThreeStrings(t *testing.T) {
	var text []byte = []byte{'s', 't', 'r', 'u', 'c', 't', ' ', ' ', '*'}
	program := make(inst.InstSlice, 30)

	program[0] = inst.New(inst.SPLIT, 2, 1)  // go to 1 or 2/3
	program[1] = inst.New(inst.SPLIT, 9, 14) // go to 2 or 3
	program[2] = inst.New(inst.CHAR, 's', 's')
	program[3] = inst.New(inst.CHAR, 't', 't')
	program[4] = inst.New(inst.CHAR, 'r', 'r')
	program[5] = inst.New(inst.CHAR, 'u', 'u')
	program[6] = inst.New(inst.CHAR, 'c', 'c')
	program[7] = inst.New(inst.CHAR, 't', 't')
	program[8] = inst.New(inst.MATCH, 0, 0)
	program[9] = inst.New(inst.SPLIT, 10, 12)
	program[10] = inst.New(inst.CHAR, ' ', ' ')
	program[11] = inst.New(inst.JMP, 9, 0)
	program[12] = inst.New(inst.CHAR, ' ', ' ')
	program[13] = inst.New(inst.MATCH, 0, 0)
	program[14] = inst.New(inst.CHAR, '*', '*')
	program[15] = inst.New(inst.MATCH, 0, 0)

	t.Log(string(text))
	t.Log(len(text))
	t.Log(program)
	expected := []Match{
		Match{8, []byte("struct")},
		Match{13, []byte("  ")},
		Match{15, []byte("*")},
	}

	i := 0
	for tc, m, err, scan := LexerEngine(program, text)(0); scan != nil; tc, m, err, scan = scan(tc) {
		t.Log(m)
		if err != nil {
			t.Error(err)
		} else if !m.Equals(&expected[i]) {
			t.Error(m, expected[i])
		}
		i++
	}
	if i != len(expected) {
		t.Error("unconsumed matches", expected[i-1:])
	}
}

