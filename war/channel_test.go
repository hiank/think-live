package war_test

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestChannelBuffer(t *testing.T) {

	ch := make(chan int, 3)
	ch <- 0
	ch <- 1
	ch <- 2
	assert.Equal(t, len(ch), 3)
	assert.Equal(t, cap(ch), 3)
}

func TestSelect(t *testing.T) {

	ch := make(chan int, 1)
	for i := 0; i < 10; i++ {
		select {
		case x := <-ch:
			t.Log(x)
		case ch <- i:
		}
	}
}
