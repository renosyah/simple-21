package util

import (
	"math/rand"
	"time"
)

type Boolgen struct {
	src       rand.Source
	cache     int64
	remaining int
}

func (b *Boolgen) Bool() bool {
	if b.remaining == 0 {
		b.cache, b.remaining = b.src.Int63(), 63
	}

	result := b.cache&0x01 == 1
	b.cache >>= 1
	b.remaining--

	return result
}

func NewBoolgen() *Boolgen {
	return &Boolgen{src: rand.NewSource(time.Now().UnixNano())}
}
