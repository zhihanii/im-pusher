package timingwheel

import (
	"log"
	"testing"
	"time"
)

func TestTimingWheel(t *testing.T) {
	tw := NewTimingWheel(time.Second, 10)
	tw.Start()
	done := make(chan struct{})
	log.Printf("start")
	tw.Add(500*time.Millisecond, func() {
		log.Println("yes")
		done <- struct{}{}
	})
	select {
	case <-done:
	}
}
