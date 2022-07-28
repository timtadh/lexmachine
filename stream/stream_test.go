package stream

import (
	"bytes"
	"testing"
)

func TestReadFullStream(t *testing.T) {
	text := "hello world"
	var buf bytes.Buffer
	s := BufferedStream(bytes.NewBufferString(text))
	for s.Advance(1) {
		if err := buf.WriteByte(s.Character().Byte); err != nil {
			if err != nil {
				t.Fatalf("err writing %v", err)
			}
		}
	}
	if s.Err() != nil {
		t.Fatalf("stream err %v", s.Err())
	}
	if buf.String() != text {
		t.Fatalf("expect %q got %q", text, buf.String())
	}
}

func TestReadEveryOther(t *testing.T) {
	text := "hello world"
	expected := "el ol"
	var buf bytes.Buffer
	s := BufferedStream(bytes.NewBufferString(text))
	for s.Advance(2) {
		if err := buf.WriteByte(s.Character().Byte); err != nil {
			if err != nil {
				t.Fatalf("err writing %v", err)
			}
		}
	}
	if s.Err() != nil {
		t.Fatalf("stream err %v", s.Err())
	}
	if buf.String() != expected {
		t.Fatalf("expect %q got %q", expected, buf.String())
	}
}

func TestReadEvery3(t *testing.T) {
	text := "hello world"
	expected := "l r"
	var buf bytes.Buffer
	s := BufferedStream(bytes.NewBufferString(text))
	for s.Advance(3) {
		if err := buf.WriteByte(s.Character().Byte); err != nil {
			if err != nil {
				t.Fatalf("err writing %v", err)
			}
		}
	}
	if s.Err() != nil {
		t.Fatalf("stream err %v", s.Err())
	}
	if buf.String() != expected {
		t.Fatalf("expect %q got %q", expected, buf.String())
	}
}

func TestPeekTillW(t *testing.T) {
	text := "hello world"
	expected := "world"
	var buf bytes.Buffer
	s := BufferedStream(bytes.NewBufferString(text))
	if !s.Started() {
		s.Advance(1)
	}
	for i := 0; ; i++ {
		b, has := s.Peek(i)
		if !has {
			break
		}
		if b.Byte == 'w' {
			s.Advance(i)
			break
		}
	}
	if s.Character().Byte != 'w' {
		t.Fatalf("expected w got %v", s.Character().Byte)
	}
	for {
		if err := buf.WriteByte(s.Character().Byte); err != nil {
			if err != nil {
				t.Fatalf("err writing %v", err)
			}
		}
		if !s.Advance(1) {
			break
		}
	}
	if s.Err() != nil {
		t.Fatalf("stream err %v", s.Err())
	}
	if buf.String() != expected {
		t.Fatalf("expect %q got %q", expected, buf.String())
	}
}

func TestPeekTillWThenL(t *testing.T) {
	text := "hello world"
	expected := "ld"
	var buf bytes.Buffer
	s := BufferedStream(bytes.NewBufferString(text))
	if !s.Started() {
		s.Advance(1)
	}
	for i := 0; ; i++ {
		b, has := s.Peek(i)
		if !has {
			break
		}
		if b.Byte == 'w' {
			s.Advance(i)
			break
		}
	}
	if s.Character().Byte != 'w' {
		t.Fatalf("expected w got %v", s.Character().Byte)
	}
	for i := 1; ; i++ {
		b, has := s.Peek(i)
		if !has {
			break
		}
		if b.Byte == 'l' {
			s.Advance(i)
			break
		}
	}
	if s.Character().Byte != 'l' {
		t.Fatalf("expected l got %v", s.Character().Byte)
	}
	for {
		if err := buf.WriteByte(s.Character().Byte); err != nil {
			if err != nil {
				t.Fatalf("err writing %v", err)
			}
		}
		if !s.Advance(1) {
			break
		}
	}
	if s.Err() != nil {
		t.Fatalf("stream err %v", s.Err())
	}
	if buf.String() != expected {
		t.Fatalf("expect %q got %q", expected, buf.String())
	}
}

func TestPeekTillWThenLThenEnd(t *testing.T) {
	text := "hello world"
	expected := ""
	var buf bytes.Buffer
	s := BufferedStream(bytes.NewBufferString(text))
	if !s.Started() {
		s.Advance(1)
	}
	for i := 0; ; i++ {
		b, has := s.Peek(i)
		if !has {
			break
		}
		if b.Byte == 'w' {
			s.Advance(i)
			break
		}
	}
	if s.Character().Byte != 'w' {
		t.Fatalf("expected w got %v", s.Character().Byte)
	}
	for i := 1; ; i++ {
		b, has := s.Peek(i)
		if !has {
			break
		}
		if b.Byte == 'l' {
			s.Advance(i)
			break
		}
	}
	if s.Character().Byte != 'l' {
		t.Fatalf("expected l got %v", s.Character().Byte)
	}
	for i := 1; ; i++ {
		_, has := s.Peek(i)
		if !has {
			s.Advance(i)
			break
		}
	}
	if !s.EOS() {
		t.Fatalf("expected EOS")
	}
	if s.Err() != nil {
		t.Fatalf("stream err %v", s.Err())
	}
	if buf.String() != expected {
		t.Fatalf("expect %q got %q", expected, buf.String())
	}
}

func TestPeekThenReadFullStream(t *testing.T) {
	text := "hello world"
	var peek bytes.Buffer
	var read bytes.Buffer
	s := BufferedStream(bytes.NewBufferString(text))
	if !s.Started() {
		s.Advance(1)
	}
	for i := 0; ; i++ {
		b, has := s.Peek(i)
		if !has {
			break
		}
		if err := peek.WriteByte(b.Byte); err != nil {
			if err != nil {
				t.Fatalf("err writing %v", err)
			}
		}
	}
	for !s.EOS() {
		if err := read.WriteByte(s.Character().Byte); err != nil {
			if err != nil {
				t.Fatalf("err writing %v", err)
			}
		}
		s.Advance(1)
	}
	if s.Err() != nil {
		t.Fatalf("stream err %v", s.Err())
	}
	if peek.String() != text {
		t.Fatalf("expect %q got %q", text, peek.String())
	}
	if read.String() != text {
		t.Fatalf("expect %q got %q", text, read.String())
	}
}

func TestLineColumns(t *testing.T) {
	text := `b
	this
	is
	wizard
`
	var expected = []struct {
		tc, line, column int
		char             byte
	}{
		{0, 1, 1, 'b'},
		{1, 2, 0, '\n'},
		{2, 2, 1, '\t'},
		{3, 2, 2, 't'},
		{4, 2, 3, 'h'},
		{5, 2, 4, 'i'},
		{6, 2, 5, 's'},
		{7, 3, 0, '\n'},
		{8, 3, 1, '\t'},
		{9, 3, 2, 'i'},
		{10, 3, 3, 's'},
		{11, 4, 0, '\n'},
		{12, 4, 1, '\t'},
		{13, 4, 2, 'w'},
		{14, 4, 3, 'i'},
		{15, 4, 4, 'z'},
		{16, 4, 5, 'a'},
		{17, 4, 6, 'r'},
		{18, 4, 7, 'd'},
		{19, 5, 0, '\n'},
	}
	s := BufferedStream(bytes.NewBufferString(text))
	// pre-peek everything just to futz with the interior state
	if !s.Started() {
		s.Advance(1)
	}
	for i := 0; ; i++ {
		_, has := s.Peek(i)
		if !has {
			break
		}
	}
	for i := 0; !s.EOS(); i++ {
		char := s.Character()
		if char.Byte != expected[i].char {
			t.Fatalf("got %v expected %v", char.Byte, expected[i].char)
		}
		if char.TC != expected[i].tc {
			t.Fatalf("got %v expected %v", char.TC, expected[i].tc)
		}
		if char.Line != expected[i].line {
			t.Fatalf("got %v expected %v", char.Line, expected[i].line)
		}
		if char.Column != expected[i].column {
			t.Fatalf("got %v expected %v", char.Column, expected[i].column)
		}
		s.Advance(1)
	}
}

func TestEveryOtherLineColumns(t *testing.T) {
	text := `b
	this
	is
	wizard
`
	var expected = []struct {
		tc, line, column int
		char             byte
	}{
		{1, 2, 0, '\n'},
		{3, 2, 2, 't'},
		{5, 2, 4, 'i'},
		{7, 3, 0, '\n'},
		{9, 3, 2, 'i'},
		{11, 4, 0, '\n'},
		{13, 4, 2, 'w'},
		{15, 4, 4, 'z'},
		{17, 4, 6, 'r'},
		{19, 5, 0, '\n'},
	}
	s := BufferedStream(bytes.NewBufferString(text))
	for i := 0; s.Advance(2); i++ {
		c := s.Character()
		if c.Byte != expected[i].char {
			t.Fatalf("got %v expected %v", c.Byte, expected[i].char)
		}
		if c.TC != expected[i].tc {
			t.Fatalf("got %v expected %v", c.TC, expected[i].tc)
		}
		if c.Line != expected[i].line {
			t.Fatalf("got %v expected %v", c.Line, expected[i].line)
		}
		if c.Column != expected[i].column {
			t.Fatalf("got %v expected %v", c.Column, expected[i].column)
		}
	}
}
