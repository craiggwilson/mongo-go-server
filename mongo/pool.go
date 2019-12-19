package mongo

import (
	"bufio"
	"io"
	"sync"
)

var (
	bufioReaderPool    sync.Pool
	bufioWriter2KPool  sync.Pool
	bufioWriter4KPool  sync.Pool
	bufioWriter8KPool  sync.Pool
	bufioWriter16KPool sync.Pool
)

func bufioWriterPool(size int) *sync.Pool {
	switch size {
	case 2 << 10:
		return &bufioWriter2KPool
	case 4 << 10:
		return &bufioWriter4KPool
	case 8 << 10:
		return &bufioWriter8KPool
	case 16 << 10:
		return &bufioWriter16KPool
	default:
		return nil
	}
}

func newBufioReader(r io.Reader) *bufio.Reader {
	if v := bufioReaderPool.Get(); v != nil {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}

	return bufio.NewReader(r)
}

func putBufioReader(br *bufio.Reader) {
	br.Reset(nil)
	bufioReaderPool.Put(br)
}

func newBufioWriterSize(w io.Writer, size int) *bufio.Writer {
	pool := bufioWriterPool(size)
	if pool != nil {
		if v := pool.Get(); v != nil {
			bw := v.(*bufio.Writer)
			bw.Reset(w)
			return bw
		}
	}
	return bufio.NewWriterSize(w, size)
}

func putBufioWriter(bw *bufio.Writer) {
	bw.Reset(nil)
	if pool := bufioWriterPool(bw.Available()); pool != nil {
		pool.Put(bw)
	}
}
