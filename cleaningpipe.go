package cleaningpipe

import (
	"bytes"
	"io"
)

/**
cleans a Reader's contents while reading
will consume but NOT CLOSE the wrapped reader.
*/

type CleaningFunc func([]byte) []byte

type CleaningPipe struct {
	in      io.Reader
	cleaner func([]byte) []byte
}

func NewCleaningPipe(c CleaningFunc, in io.Reader) CleaningPipe {
	return CleaningPipe{
		in:      in,
		cleaner: c,
	}
}

func (r CleaningPipe) Read(p []byte) (n int, err error) {
	if r.in == nil {
		return 0, io.EOF
	}
	n, err = r.in.Read(p)

	if n <= 0 {
		return n, err
	}

	tmp := r.cleaner(p)

	if &tmp != &p {
		copy(p, tmp)
	}

	if len(tmp) < len(p) {
		//something was deleted
		if firstZero := bytes.IndexByte(p, '\x00'); firstZero >= 0 {
			//case 1: p was not full and we deleted bytes --> there are trailing zeroes in p and tmp
			return firstZero, err
		}
		//case 2: p was full and we deleted bytes --> n = length of tmp as tmp must be full
		return len(tmp), err
	}
	//case 3: nothing was deleted
	return n, err
}
