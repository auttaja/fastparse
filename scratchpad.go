package fastparse

import (
	"errors"
	"sync"
)

type scratchpad struct {
	// Defines the length of the content in the pad.
	len int

	// Defines the slice which is used in the pad.
	slice []byte

	// Defines the length of the slice.
	sliceLen int

	// Defines the pool (if this was pre-allocated).
	pool *sync.Pool
}

// Resets the pad.
func (s *scratchpad) reset() {
	s.len = 0
}

// Readds the pad to the pool if it was re-allocated. This should be called after usage.
func (s *scratchpad) readd() {
	if s.pool != nil {
		s.pool.Put(s)
	}
}

// AddByte puts a byte into the pad.
func (s *scratchpad) AddByte(b byte) error {
	if s.len == s.sliceLen {
		return errors.New("pad is full")
	}
	s.slice[s.len] = b
	s.len++
	return nil
}

// String is used to return a string representation of what is in the pad.
func (s *scratchpad) String() string {
	return string(s.slice[:s.len])
}
