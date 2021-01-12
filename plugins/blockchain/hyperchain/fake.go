package hyperchain

import (
	"github.com/meshplus/gosdk/common"
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
	chars    = "abcdef0123456789"
	charsLen = len(chars)
	signLen  = 130
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

func fakeSign() string {
	return common.StringToHex(randomString(signLen))
}
