package index

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {

	ch := make(chan int)

	go func() {
		ch <- 1
	}()
	v, b := <-ch
	fmt.Println(v, b)
	close(ch)

	v, b = <-ch
	fmt.Println(v, b)
}
