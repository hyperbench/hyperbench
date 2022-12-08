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

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randomStringWithSeed(r *rand.Rand, l uint) string {
	b := make([]byte, l)
	for i, cache, remain := len(b)-1, r.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = r.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(chars) {
			b[i] = chars[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)

}

func randomInt(min, max int) int {
	randChOnce.Do(initRandCh)
	r := <-randCh
	i := r.Intn(max-min) + min
	randCh <- r
	return i
}
