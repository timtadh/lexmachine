package stream

import "fmt"

// Stream represents a stream of bytes. Its interface is analogous to
// bufio.Scanner. Here is an example for how to read all the bytes in a stream
// (and print them one by one):
//
//     s := BufferedStream(reader)
//     for s.Advance(1) {
//         fmt.Println(s.Character().Byte)
//     }
//     if s.Err() != nil {
//         return s.Err()
//     }
//
type Stream interface {

	// Character returns the current byte in the stream. This method will panic
	// if Advance has not been called before this method or Advance has
	// returned false.
	Character() Character

	// Peek returns byte at the current cursor + the lookahead in the stream if
	// one exists. If lookahead == 0, it returns the same character Character()
	// returns. If lookahead == 1, it returns the next byte, and so on. Peek
	// does not advance the cursor. If there are no further bytes in the stream
	// (or lookahead causes a read past the end of the stream) Peek returns has
	// == false. If you call Peek() before Advance() has been called it will
	// panic.
	Peek(lookahead int) (char Character, has bool)

	// Advance moves the cursor i bytes forward in the stream. If there is a
	// byte to read it returns true. If it reaches the end of the stream (EOS)
	// it returns false. Advance with i > than number of bytes remaining moves
	// the cursor to the end of stream (may be less than i) and returns false
	// (as you cannot read past the end of the stream). Advance must be called
	// with movement >= 0 otherwise it will panic. If Advance is called with
	// i == 0 it does nothing (including setting the stream to started).
	Advance(i int) bool

	// Started returns true the stream has been started (eg. a call to Advance
	// has been made with a positive movement).
	Started() bool

	// EOS returns true if the stream has reached the end of the stream.
	EOS() bool

	// Err returns an error if there was an error reading from the underlying
	// source of the bytes. Panics if called before Advance returns false.
	// Err() will never return io.EOF (it will be nil in this case -- following
	// the behavior of ioutil.ReadAll)
	Err() error
}

// Character represents one byte in a stream with position information.
type Character struct {
	Byte   byte
	TC     int
	Line   int
	Column int
}

// String humanizes the character
func (c Character) String() string {
	return fmt.Sprintf("<%q tc:%d @ %d:%d>", c.Byte, c.TC, c.Line, c.Column)
}
