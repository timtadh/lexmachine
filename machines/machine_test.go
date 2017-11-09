package machines

import "testing"
import "github.com/timtadh/lexmachine/inst"

func TestLexerMatch(t *testing.T) {
	text := []byte("ababcbcbb")
	//. (a|b)*cba?(c|b)bb
	program := make(inst.Slice, 20)

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
	mtext := []byte("ababcbcbb")
	expected := []Match{
		{16, 0, 1, 1, 1, len(mtext), mtext},
	}
	i := 0
	for tc, m, err, scan := LexerEngine(program, text)(0); scan != nil; tc, m, err, scan = scan(tc) {
		t.Log(tc, m)
		if err != nil {
			t.Error(err)
		} else if !m.Equals(&expected[i]) {
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
	program := make(inst.Slice, 20)

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

	for tc, _, err, scan := LexerEngine(program, text)(0); scan != nil; tc, _, err, scan = scan(tc) {
		if err == nil {
			t.Fatal("no error!", err)
		}
		_, unconsumed := err.(*UnconsumedInput)
		if !unconsumed {
			t.Fatalf("unexpected error type (expected *UnconsumedInput) got %v", err)
		}
	}
}

func TestLexerThreeStrings(t *testing.T) {
	var text = []byte("struct  *")
	program := make(inst.Slice, 30)

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
		{8, 0, 1, 1, 1, 6, []byte("struct")},
		{13, 6, 1, 7, 1, 8, []byte("  ")},
		{15, 8, 1, 9, 1, 9, []byte("*")},
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

func TestLexerRestart(t *testing.T) {
	var text = []byte("struct\n  *")
	program := make(inst.Slice, 30)

	program[0] = inst.New(inst.SPLIT, 2, 1)  // go to 1 or 2/3
	program[1] = inst.New(inst.SPLIT, 9, 20) // go to 2 or 3
	program[2] = inst.New(inst.CHAR, 's', 's')
	program[3] = inst.New(inst.CHAR, 't', 't')
	program[4] = inst.New(inst.CHAR, 'r', 'r')
	program[5] = inst.New(inst.CHAR, 'u', 'u')
	program[6] = inst.New(inst.CHAR, 'c', 'c')
	program[7] = inst.New(inst.CHAR, 't', 't')
	program[8] = inst.New(inst.MATCH, 0, 0)
	program[9] = inst.New(inst.SPLIT, 10, 12)
	program[10] = inst.New(inst.CHAR, ' ', ' ')
	program[11] = inst.New(inst.JMP, 13, 0)
	program[12] = inst.New(inst.CHAR, '\n', '\n')
	program[13] = inst.New(inst.SPLIT, 14, 19)
	program[14] = inst.New(inst.SPLIT, 15, 17)
	program[15] = inst.New(inst.CHAR, ' ', ' ')
	program[16] = inst.New(inst.JMP, 18, 0)
	program[17] = inst.New(inst.CHAR, '\n', '\n')
	program[18] = inst.New(inst.JMP, 13, 0)
	program[19] = inst.New(inst.MATCH, 0, 0)
	program[20] = inst.New(inst.CHAR, '*', '*')
	program[21] = inst.New(inst.MATCH, 0, 0)

	t.Log(string(text))
	t.Log(len(text))
	t.Log(program)
	expected := []Match{
		{8, 0, 1, 1, 1, 6, []byte("struct")},
		{19, 6, 2, 0, 2, 2, []byte("\n  ")},
		{21, 9, 2, 3, 2, 3, []byte("*")},
	}

	check := func(m *Match, i int, err error) {
		t.Log(m)
		if err != nil {
			t.Error(err)
		} else if !m.Equals(&expected[i]) {
			t.Error(m, expected[i])
		}
	}

	i := 0
	tc, m, err, scan := LexerEngine(program, text)(0)
	check(m, i, err)
	i++

	tc, m, err, scan = scan(tc)
	check(m, i, err)
	i++

	tc, m, err, scan = scan(tc)
	check(m, i, err)
	i -= 2

	tc, m, err, scan = scan(tc - 10) // backtrack
	check(m, i, err)
	i++

	tc, m, err, scan = scan(tc)
	check(m, i, err)
	i++

	tc, m, err, scan = scan(tc)
	check(m, i, err)
	i--

	tc, m, err, scan = scan(tc - 4)
	check(m, i, err)
	i++

	tc, m, err, scan = scan(tc)
	check(m, i, err)
	i++

	_, _, _, scan = scan(tc)
	if scan != nil {
		t.Error("scan should have ended")
	}
	if i != len(expected) {
		t.Error("unconsumed matches", expected[i-1:])
	}
}
