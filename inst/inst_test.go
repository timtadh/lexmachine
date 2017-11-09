package inst

import "testing"

func TestPrint(t *testing.T) {
	i := New(CHAR, uint32('a'), 0)
	j := New(MATCH, 0, 0)
	k := New(JMP, 14, 0)
	l := New(SPLIT, 15, 17)
	t.Log(i)
	t.Log(j)
	t.Log(k)
	t.Log(l)
	s := make(Slice, 4)
	s[0] = i
	s[1] = j
	s[2] = k
	s[3] = l
	t.Log(s)
}
