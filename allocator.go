package signal

import "golang.org/x/exp/constraints"

// Allocator provides allocation of various signal buffers.
type Allocator struct {
	Channels int
	Length   int
	Capacity int
}

func Allocate[T constraints.Float](a Allocator) *F[T] {
	return &F[T]{
		buffer:   make([]T, a.Channels*a.Length, a.Channels*a.Capacity),
		channels: channels(a.Channels),
	}
}
