package machines

import "testing"
import "strings"
import "github.com/timtadh/lexmachine/inst"


func TestToDFA(t *testing.T) {
	text := []byte("ababcbcbb")
	//. (a|b)*cba?(c|b)bb
	program := make(inst.InstSlice, 17)

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

	t.Log(program)
	dfa := ToDFA(program)

	t.Log(dfa)
	t.Log(string(text))
	t.Log(len(text))
	mtext := []byte("ababcbcbb")
	t.Log(program)
	t.Log(len(program)-1)
	expected := []Match{
		Match{len(program)-1, 0, 1, 1, 1, len(mtext), mtext},
	}
	i := 0
	for tc, m, err, scan := DFALexerEngine(dfa, text)(0); scan != nil; tc, m, err, scan = scan(tc) {
		t.Log(tc, m)
		if err != nil {
			t.Error(err)
		} else if !m.Equals(&expected[i]) {
			t.Error(m, expected[i])
		}
		i++
	}
	if i + 1 < len(expected) {
		t.Error("unconsumed matches", expected[i:])
	}
}


func TestDFANoMatch(t *testing.T) {
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
	dfa := ToDFA(program)
	t.Log(program)
	t.Log(dfa)

	for tc, m, err, scan := DFALexerEngine(dfa, text)(0); scan != nil; tc, m, err, scan = scan(tc) {
		if err == nil || !strings.HasPrefix(err.Error(), "Unconsumed text") {
			t.Error("no error!", m, err)
		}
	}
}

func TestToDFAMulti(t *testing.T) {
	text := []byte("c")
	//. (a|b)*cba?(c|b)bb
	program := make(inst.InstSlice, 20)

	program[0]  = inst.New(inst.SPLIT, 1, 6)
	program[1]  = inst.New(inst.SPLIT, 2, 4)
	program[2]  = inst.New(inst.CHAR, 97, 97)
	program[3]  = inst.New(inst.JMP, 5, 0)
	program[4]  = inst.New(inst.CHAR, 98, 98)
	program[5]  = inst.New(inst.MATCH, 0, 0)
	program[6]  = inst.New(inst.SPLIT, 7, 9)
	program[7]  = inst.New(inst.CHAR, 99, 99)
	program[8]  = inst.New(inst.JMP, 10, 0)
	program[9]  = inst.New(inst.CHAR, 100, 100)
	program[10] = inst.New(inst.MATCH, 0, 0)

	t.Log(program)
	dfa := ToDFA(program)

	t.Log(dfa)
	t.Log(string(text))
	t.Log(len(text))
	mtext := []byte("c")
	expected := []Match{
		Match{len(dfa)-1, 0, 1, 1, 1, len(mtext), mtext},
	}
	i := 0
	for tc, m, err, scan := DFALexerEngine(dfa, text)(0); scan != nil; tc, m, err, scan = scan(tc) {
		t.Log(tc, m)
		if err != nil {
			t.Error(err)
		} else if !m.Equals(&expected[i]) {
			t.Error(m, expected[i])
		}
		i++
	}
	if i + 1 < len(expected) {
		t.Error("unconsumed matches", expected[i:])
	}
}
