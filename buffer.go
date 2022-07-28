package lexmachine

import (
	"fmt"

	"github.com/timtadh/lexmachine/stream"
)

// Buffer is a abstracts to implementations of "text". The first is a []byte with a
type Buffer interface {
	Byte(i int) byte
	HasByte(i int) bool
	TC() int
	SetTC(i int)
}

type SliceBuffer struct {
	Text        []byte
	TextCounter int
}

func sliceBuffer(text []byte, tc int) *SliceBuffer {
	return &SliceBuffer{
		Text:        text,
		TextCounter: tc,
	}
}

func (s *SliceBuffer) Byte(i int) byte {
	return s.Text[i]
}

func (s *SliceBuffer) HasByte(i int) bool {
	return i >= 0 && i < len(s.Text)
}

func (s *SliceBuffer) TC() int {
	return s.TextCounter
}

func (s *SliceBuffer) SetTC(tc int) {
	s.TextCounter = tc
}

func (s *SliceBuffer) finalize() int {
	return s.TextCounter
}

type StreamBuffer struct {
	Text      stream.Stream
	Lookahead int
}

func streamBuffer(text stream.Stream) *StreamBuffer {
	return &StreamBuffer{
		Text:      text,
		Lookahead: 0,
	}
}

func (s *StreamBuffer) Byte(i int) byte {
	c, has := s.Text.Peek(i)
	if !has {
		panic(fmt.Errorf("read past the end of the buffer"))
	}
	return c.Byte
}

func (s *StreamBuffer) HasByte(i int) bool {
	_, has := s.Text.Peek(i)
	return has
}

func (s *StreamBuffer) TC() int {
	return s.Lookahead
}

func (s *StreamBuffer) SetTC(tc int) {
	s.Lookahead = tc
}

func (s *StreamBuffer) finalize() {
	s.Text.Advance(s.Lookahead)
}
