package stream

import (
	"bufio"
	"fmt"
	"io"
	"sync"
)

type bufferedStream struct {
	lock    sync.Mutex
	r       *bufio.Reader
	tc      int
	line    int
	column  int
	started bool
	eos     bool
	buf     []byte
	err     error
}

func BufferedStream(r io.Reader) Stream {
	b := &bufferedStream{
		r:      bufio.NewReader(r),
		tc:     -1,
		line:   1,
		column: 0,
	}
	return b
}

func (b *bufferedStream) Byte() byte {
	b.lock.Lock()
	defer b.lock.Unlock()
	if !b.started {
		panic(fmt.Errorf("Call to Byte() before first call to Advance"))
	} else if b.eos {
		panic(fmt.Errorf("Call to Byte() after first call to Advance returned false"))
	}
	return b.buf[0]
}

func (b *bufferedStream) Position() (tc, line, column int) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if !b.started {
		panic(fmt.Errorf("Call to Position() before first call to Advance"))
	} else if b.eos {
		panic(fmt.Errorf("Call to Position() after first call to Advance returned false"))
	}
	return b.tc, b.line, b.column
}

func (b *bufferedStream) Peek(i int) (char byte, has bool) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.eos {
		panic(fmt.Errorf("Call to Byte() after first call to Advance returned false"))
	}
	if i <= 0 {
		panic(fmt.Errorf("Peek() must be called with positive lookahead got %d", i))
	}
	// the "cursor" technically starts at -1, this does that adjustment
	if !b.started {
		i--
	}
	if len(b.buf) >= i+1 {
		return b.buf[i], true
	}
	if !b.read(i) {
		return 0, false
	}
	return b.buf[i], true
}

func (b *bufferedStream) Started() bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.eos
}

func (b *bufferedStream) EOS() bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.eos
}

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

func (b *bufferedStream) Advance(i int) bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.advance(i)
}

func (b *bufferedStream) advance(i int) bool {
	if i <= 0 {
		panic(fmt.Errorf("Advance() must be called with positive move got %d", i))
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
	b.trackPos(b.buf[0])
	return true
}

// trims the buffer by up i bytes and returns the number of
// bytes trimmed.
func (b *bufferedStream) trimBuffer(i int) int {
	for j := 1; j < i && j < len(b.buf); j++ {
		b.trackPos(b.buf[j])
	}
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

// updates the position information for the given character.
// only call once per character in the stream.
func (b *bufferedStream) trackPos(char byte) {
	b.tc++
	if char == '\n' {
		b.line++
		b.column = 0
	} else {
		b.column++
	}
}

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
		b.buf = append(b.buf, buf[:n]...)
		if len(b.buf) >= i+1 {
			break
		}
	}
	return true
}
