package toolkit

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var (
	randCh     = make(chan *rand.Rand, runtime.NumCPU())
	randChOnce sync.Once
)

const (
	chars    = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsLen = len(chars)
)

func initRandCh() {
	for i := 0; i < runtime.NumCPU(); i++ {
		randCh <- rand.New(rand.NewSource(time.Now().UnixNano()))
	}
}

func randomString(l uint) string {
	randChOnce.Do(initRandCh)

	r := <-randCh
	s := make([]byte, l)
	for i := 0; i < int(l); i++ {
		s[i] = chars[r.Intn(charsLen)]
	}
	randCh <- r
	return string(s)
}

func randomInt(min, max int) int {
	randChOnce.Do(initRandCh)
	r := <-randCh
	i := r.Intn(max-min) + min
	randCh <- r
	return i
}
