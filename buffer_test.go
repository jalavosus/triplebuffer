package triplebuffer

/**
This file exports a few methods on Buffer for testing purposes.
These methods are not publicly accessible otherwise.
*/

// Back returns the data currently in the back buffer,
// or nil.
func (b *Buffer[T]) Back() *T {
	return b.back
}

// Middle returns the data currently in the middle buffer, as
// well as its pending status.
// Pending will be false if the data is nil.
func (b *Buffer[T]) Middle() (*T, bool) {
	m := b.middle.Load()
	if m != nil {
		return m.data, m.pending
	}

	return nil, false
}

// Front returns the data currently in the front buffer,
// or nil.
func (b *Buffer[T]) Front() *T {
	return b.front
}
