package recording

import (
	"fmt"
	"io"
	"time"
	"wrbb-stream-recorder/internal/util"
)

// I Have to make this file because I essentially need to do what io.Copy
// is doing but I need it to stop once the show ends

// I Know there is probably a much more clean and elegant way to do this BUT
// I couldn't figure it out and anything online I looked up didnt help so here we are!

// copyShow is the actual implementation of Copy and CopyBuffer.
func copyShow(dst io.Writer, src io.Reader, duration time.Duration, name string) {
	var buf []byte
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
	timer := time.NewTimer(duration)
	go writeToFile(dst, src, buf, timer, name)
	return
}

// writeToFile continues to write to a file until the timer's channel receives a signal.
// name is used for debugging and to remove the show from the current recording list
func writeToFile(dst io.Writer, src io.Reader,  buf []byte, timer *time.Timer, name string) {
	var err error
	for {
		select {
		case <-timer.C:
			// Show is over, stop recording
			currentRecording.mu.Lock()
			delete(currentRecording.shows, name)
			currentRecording.mu.Unlock()
			util.InfoLogger.Printf("Finished recording %s\n", name)
			return
		default:
			nr, er := src.Read(buf)
			if nr > 0 {
				nw, ew := dst.Write(buf[0:nr])
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
		if err != nil {
			util.ErrorLog(fmt.Sprintf("Error writing to file %s: %s", name, err.Error()))
			return
		}
	}
}