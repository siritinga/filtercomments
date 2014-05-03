package filtercomments

import (
	"bufio"
	"io"
)

type Reader struct {
	rd       *bufio.Reader
	commchar byte
	buf      []byte
	err      error
}

func NewReader(r io.Reader, c byte) *Reader {
	w := new(Reader)
	w.rd = bufio.NewReader(r)
	w.commchar = c
	w.buf = []byte{}
	w.err = nil

	return w
}

func (b *Reader) Read(p []byte) (n int, err error) {
	// If there has been an error...
	if b.err != nil {
		return 0, b.err
	}

	// Send saved text from previously.
	if len(b.buf) > 0 {
		if len(p) >= len(b.buf) {
			//The entire buf can be send
			n := copy(p, b.buf)
			b.buf = []byte{}
			return n, nil
		} else {
			// only part of buf can be send
			n := copy(p, b.buf)
			b.buf = b.buf[n:]
			return n, nil
		}
	}

	// Get new data
	for {
		// TODO: Prefix is ignored
		b.buf, _, err = (*b.rd).ReadLine()
		if err != nil {
			b.err = err
			return 0, b.err
		}
		// Check for comment or empty line
		if len(b.buf) < 2 {
			continue
		}
		if b.buf[0] == b.commchar {
			continue
		}
		break
	}
	b.buf = append(b.buf, '\n')
	// Now that buf is full, return 0 so this method is called again.
	// TODO: this is not very elegant...
	return 0, nil
}
