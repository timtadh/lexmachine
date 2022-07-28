package stream

import (
	"fmt"
	"io"
	"sync"
)

type bufferedStream struct {
	lock    sync.Mutex
	r       io.Reader
	tc      int
	line    int
	column  int
	started bool
	eos     bool
	buf     []Character
	err     error
}

// BufferedStream makes a Stream which is backed by an expandable buffer.
func BufferedStream(r io.Reader) Stream {
	b := &bufferedStream{
		r:      r,
		tc:     -1,
		line:   1,
		column: 0,
	}
	return b
}

// Character returns the character at the cursor
func (b *bufferedStream) Character() Character {
	b.lock.Lock()
	defer b.lock.Unlock()
	if !b.started {
		panic(fmt.Errorf("Call to Byte() before first call to Advance"))
	} else if b.eos {
		panic(fmt.Errorf("Call to Byte() after first call to Advance returned false"))
	}
	return b.buf[0]
}

// Peek gets the character at lookahead i
func (b *bufferedStream) Peek(i int) (char Character, has bool) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if !b.started {
		panic(fmt.Errorf("Call to Peek() before first call to Advance"))
	} else if b.eos {
		panic(fmt.Errorf("Call to Peek() after first call to Advance returned false"))
	}
	if i < 0 {
		panic(fmt.Errorf("Peek() must be called with lookahead >= 0 got %d", i))
	}
	if len(b.buf) >= i+1 {
		return b.buf[i], true
	}
	if !b.read(i) {
		return Character{}, false
	}
	return b.buf[i], true
}

// Started indicates if Advance has been called at least once.
func (b *bufferedStream) Started() bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.started
}

// EOS indicates whether the stream has reached End Of Stream
func (b *bufferedStream) EOS() bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.eos
}

// Err returns the error from the underlying io.Reader if io.Read() returned
// a non-EOF error.
func (b *bufferedStream) Err() error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if !b.started {
		panic(fmt.Errorf("Call to Err() before first call to Advance"))
	} else if !b.eos {
		panic(fmt.Errorf("Call to Err() before call to Advance returned false"))
	}
	return b.err
}

// Advance moves the cursor forward by i
func (b *bufferedStream) Advance(i int) bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.advance(i)
}

// advance moves the cursor forward by i
func (b *bufferedStream) advance(i int) bool {
	if i == 0 {
		return true
	}
	if i < 0 {
		panic(fmt.Errorf("Advance() must be called with move >= 0 got %d", i))
	}
	// the "cursor" technically starts at -1, this does that adjustment
	if !b.started {
		b.started = true
		i--
		// ensures a read happens even if i==0 when the buf is empty
		if len(b.buf) <= 0 && !b.read(1) {
			b.eos = true
			return false
		}
	}
	i = i - b.trimBuffer(i)
	if len(b.buf) <= i {
		if !b.read(i) {
			b.eos = true
			return false
		}
	}
	if i > 0 {
		i = i - b.trimBuffer(i)
		if i != 0 {
			panic(fmt.Errorf("i != 0 (i = %d)", i))
		}
	}
	return true
}

// trims the buffer by up i bytes and returns the number of bytes trimmed.
func (b *bufferedStream) trimBuffer(i int) int {
	if len(b.buf) > i {
		// we already recorded the position
		// of b.buf[0]. we need to track all the chars
		// we are dropping by the skip
		copy(b.buf[:len(b.buf)-i], b.buf[i:])
		b.buf = b.buf[:len(b.buf)-i]
		return i
	} else {
		trimmed := len(b.buf)
		b.buf = b.buf[:0]
		return trimmed
	}
	return 0
}

// updates the position information for the given character.  only call once
// per character in the stream.
func (b *bufferedStream) trackPos(char byte) {
	b.tc++
	if char == '\n' {
		b.line++
		b.column = 0
	} else {
		b.column++
	}
}

// reads at least i bytes from the underlying reader into the buffer.
func (b *bufferedStream) read(i int) bool {
	if b.eos {
		return false
	}
	buf := make([]byte, 4096)
	for {
		n, err := b.r.Read(buf)
		if err != nil {
			if err != io.EOF {
				// only set err if it is an unexpected error.
				b.err = err
			}
			return false
		}
		for _, c := range buf[:n] {
			b.trackPos(c)
			b.buf = append(b.buf, Character{
				Byte:   c,
				TC:     b.tc,
				Line:   b.line,
				Column: b.column,
			})
		}
		if len(b.buf) >= i+1 {
			break
		}
	}
	return true
}
