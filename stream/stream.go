package stream

// Stream represents a stream of bytes. Its interface is analogous to
// bufio.Scanner. Here is an example for how to read all the bytes in a stream
// (and print them one by one):
//
//     s := BufferedStream(reader)
//     for s.Advance(1) {
//         fmt.Println(s.Byte())
//     }
//     if s.Err() != nil {
//         return s.Err()
//     }
//
type Stream interface {

	// Byte returns the current byte in the stream. This method will panic if
	// Advance has not been called before this method or Advance has returned
	// false.
	Byte() byte

	// Position returns the position of the current byte: text counter, line,
	// and column. This method will panic if Advance has not been called before
	// this method or Advance has returned false.
	Position() (tc, line, column int)

	// Peek returns byte at the current cursor + the lookahead in the stream if
	// one exists. If lookahead == 0, Peek will panic, if lookahead == 1, it
	// returns the next byte, and so on. Peek does not advance the cursor. If
	// there are no further bytes in the stream (or lookahead causes a read
	// past the end of the stream) Peek returns has == false. You may call this
	// method before Advance.
	Peek(lookahead int) (char byte, has bool)

	// Advance moves the cursor i bytes forward in the stream. If there is a
	// byte to read it returns true. If it reaches the end of the stream (EOS)
	// it returns false. Advance with i > than number of bytes remaining moves
	// the cursor to the end of stream (may be less than i) and returns false
	// (as you cannot read past the end of the stream). Advance must be called
	// with positive movement otherwise it will panic.
	Advance(i int) bool

	// Started returns true if at least 1 call to Advance has been made.
	Started() bool

	// EOS returns true if the stream has reached the end of the stream.
	EOS() bool

	// Err returns an error if there was an error reading from the underlying
	// source of the bytes. Panics if called before Advance returns false.
	// Err() will never return io.EOF (it will be nil in this case -- following
	// the behavior of ioutil.ReadAll)
	Err() error
}
