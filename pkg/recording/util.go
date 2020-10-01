package recording

import (
	"io"
	"wrbb-stream-recorder/pkg/spinitron"
)

// I Have to make this file because I essentially need to do what io.Copy
// is doing but I need it to stop once the show ends

// I Know there is probably a much more clean and eligant way to do this BUT
// I couldn't figure it out and anything online I looked up didnt help so here we are!


// copyBuffer is the actual implementation of Copy and CopyBuffer.
func copyShow(dst io.Writer, src io.Reader, show spinitron.Show) (written int64, err error) {
	var buf []byte
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for show.IsLive() {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
