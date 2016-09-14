package xcf

import (
	"io"

	"github.com/MJKWoolnough/byteio"
)

type writer struct {
	*byteio.StickyWriter
	io.WriterAt
}

func newWriter(w io.WriterAt) writer {
	var (
		wr io.Writer
		ok bool
	)
	if wr, ok = w.(io.Writer); !ok {
		wr = &writerAtWriter{WriterAt: w}
	}

	return writer{
		&byteio.StickyWriter{Writer: byteio.BigEndianWriter{Writer: wr}},
		w,
	}
}

type writerAtWriter struct {
	io.WriterAt
	pos int64
}

func (w *writerAtWriter) Write(p []byte) (int, error) {
	n, err := w.WriteAt(p, w.pos)
	w.pos += int64(n)
	return n, err
}