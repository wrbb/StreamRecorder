package recording

import (
	"fmt"
	"io"
	"time"
)

const BufferSize = 1024


// This file contains copyShow which is a modified version of Copy/CopyBuffer.
// The function is almost the same but now features a select statement in
// and infinite loop. The select statement waits for the duration provided to
// send a signal on its channel to indicate the show is over and the Copy
// should finish

// copyShow is the actual implementation of Copy and CopyBuffer.
func copyShow(dst io.Writer, src io.ReadCloser, duration time.Duration) error {
	buf := make([]byte, BufferSize)
	timer := time.NewTimer(duration)
	return writeToFile(dst, src, buf, timer)
}

// writeToFile continues to write to a file until the timer's channel receives a signal.
// name is used for debugging and to remove the show from the current recording list
func writeToFile(dst io.Writer, src io.Reader, buf []byte, timer *time.Timer) error{
	var err error
	for {
		select {
		case <-timer.C:
			// Show is over, stop recording
			return fmt.Errorf("uh on")
		default:
			nr, er := src.Read(buf)
			if nr > 0 {
				nw, ew := dst.Write(buf[0:nr])
				if ew != nil {
					return ew
				}
				if nr != nw {
					return io.ErrShortWrite
				}
			}
			if er != nil && er != io.EOF {
				return err
			}
		}
	}
}
