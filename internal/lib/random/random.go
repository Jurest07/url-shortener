package random

import (
	"math/rand/v2"
	"time"
)

func NewRandomString(length int) string {
	t := time.Now().UnixNano()
	source := rand.NewPCG(uint64(t), rand.Uint64())
	rnd := rand.New(source)
	resBuf := make([]rune, length)
	for i := range resBuf {
		r := rnd.IntN(25) + 97
		resBuf[i] = rune(r)
	}
	return string(resBuf)


}