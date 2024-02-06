package triplebuffer

import (
	"sync/atomic"
)

// Buffer is a triple buffer! Enjoy!
// A Buffer has three (private) fields, to wit:
//   - back: the Back buffer, which is whatever a producer is working on.
//     This is a pointer, so it is modifiable after setting without any extra work.
//     Initially set by Write, then promoted via atomic swap to the Middle buffer
//     by calling Commit after the producer has finished whatever it needs to do with the data.
//   - middle: the Middle buffer, also known as the very special atomic buffer.
//     Data in the Middle buffer can be in one of two states: "pending", or "stale".
//     Pending data is that which has been committed (via Commit), and is considered new data.
//     When Read is next called, the new data will be promoted to the Front buffer,
//     with the previous Front buffer data being atomically swapped into the Middle buffer and marked as "stale".
//     Stale data has been consumed already.
//   - front: the Front buffer. Whatever is in here will be returned by Read if there is no
//     pending data in the Middle buffer. If there is pending data in Middle, that data gets promoted
//     and returned by Read.
//
// Note that Commit can (and, likely at some point *will*) overwrite whatever is in the Middle buffer,
// even if that data is still pending a read.
type Buffer[T any] struct {
	back *T
	// using atomics is fun and exciting! ...
	middle atomic.Pointer[bufferItem[T]]
	front  *T
}

type bufferItem[T any] struct {
	data    *T
	pending bool
}

// NewBuffer creates a new empty Buffer.
func NewBuffer[T any]() *Buffer[T] {
	return new(Buffer[T])
}

// NewBufferWithFront returns a Buffer with the Front buffer populated,
// (ie it's ready to CONSUME!)
func NewBufferWithFront[T any](front *T) *Buffer[T] {
	b := new(Buffer[T])
	b.front = front

	return b
}

// NewPopulatedBuffer returns a fully populated Buffer.
// Good luck.
// Note: the Middle buffer item (`middle`) will be considered Pending
// (ie. fresh off the line from a producer, and subject to swappage when Buffer.Read is called).
func NewPopulatedBuffer[T any](back, middle, front *T) *Buffer[T] {
	b := new(Buffer[T])

	b.back = back
	b.middle.Store(&bufferItem[T]{
		data:    middle,
		pending: middle != nil,
	})
	b.front = front

	return b
}

// Read is part of the consumer API.
// It returns the value currently in the Front buffer.
// If the value currently in the Middle buffer is pending,
// (ie. it has just been written to by a producer)
// then the Front and Middle buffers will be swapped and the *new* Front buffer returned.
// Otherwise, the stale Front buffer value will be returned.
// If there is no value in the Middle buffer, no swap occurs.
// The returned boolean corresponds to whether the value is
// "stale": false, or "pending": true.
func (b *Buffer[T]) Read() (*T, bool) {
	var dirtyRead bool

	if prev := b.middle.Load(); prev != nil && prev.pending {
		dirtyRead = true

		newFront := b.middle.Swap(&bufferItem[T]{
			data:    b.front,
			pending: false,
		})

		b.front = newFront.data
	}

	return b.front, dirtyRead
}

// Write is part of the producer API,
// It is passed a pointer to the data currently being "produced", whatever that may be.
// Once the production is finished, then Buffer.Commit must be called to commit
// the new Back buffer to the Middle buffer.
func (b *Buffer[T]) Write(data *T) {
	b.back = data
}

// Commit is called once whatever data in the Back buffer has finished being created.
// On commit, the following happens:
// Back buffer is promoted to the Middle buffer.
// Middle buffer is swapped into the Back buffer, which is more or less useless.
func (b *Buffer[T]) Commit() {
	if b.back == nil {
		panic("cannot commit nil back buffer data")
	}

	prev := b.middle.Swap(&bufferItem[T]{
		data:    b.back,
		pending: true,
	})

	if prev != nil {
		b.back = prev.data
	} else {
		b.back = nil
	}

	return
}
