package triplebuffer_test

import (
	"math/big"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/jalavosus/triplebuffer"
)

func toPointer[T any](t *testing.T, val T) *T {
	t.Helper()
	return &val
}

func TestNewBuffer(t *testing.T) {
	b := triplebuffer.NewBuffer[int]()
	assert.NotNil(t, b)
	assert.Nil(t, b.Front(), "expected front to be nil")

	middle, pending := b.Middle()
	assert.Nil(t, middle, "expected middle to be nil")
	assert.False(t, pending, "expected pending to be false")

	assert.Nil(t, b.Back(), "expected back to be nil")
}

func TestNewBufferWithFront(t *testing.T) {
	t.Run("scalar_type", func(t *testing.T) {
		front := toPointer(t, int(42))

		b := triplebuffer.NewBufferWithFront(front)

		t.Run("correct_value", func(t *testing.T) {
			assert.NotNil(t, b.Front())

			// pointer test
			// ...not sure if this is overkill
			// or even the right way to do it but eh?
			ptr1 := uintptr(unsafe.Pointer(front))
			ptr2 := uintptr(unsafe.Pointer(b.Front()))
			assert.Equal(t, ptr1, ptr2)

			// value test
			assert.Equal(t, *front, *b.Front())
		})

		t.Run("value_updates", func(t *testing.T) {
			*front = 43

			// value test
			assert.Equal(t, *front, *b.Front())
		})
	})

	t.Run("custom_type", func(t *testing.T) {
		type testType struct {
			stringData string
			bigIntData *big.Int
		}

		front := &testType{
			stringData: "hello, world!",
			bigIntData: big.NewInt(42),
		}

		b := triplebuffer.NewBufferWithFront(front)

		t.Run("correct_value", func(t *testing.T) {
			assert.NotNil(t, b.Front())

			// pointer test
			// ...not sure if this is overkill
			// or even the right way to do it but eh?
			ptr1 := uintptr(unsafe.Pointer(front))
			ptr2 := uintptr(unsafe.Pointer(b.Front()))
			assert.Equal(t, ptr1, ptr2)

			// value test
			assert.Equal(t, front, b.Front())
			assert.Equal(t, front.stringData, b.Front().stringData)
			assert.Equal(t, front.bigIntData.Int64(), b.Front().bigIntData.Int64())
		})

		t.Run("value_updates", func(t *testing.T) {
			*front = testType{
				stringData: "KOOLAIDMAN",
				bigIntData: nil,
			}

			// value test
			assert.Equal(t, front, b.Front())
			assert.Nil(t, b.Front().bigIntData)
			assert.Equal(t, front.stringData, b.Front().stringData)
		})
	})
}

func TestNewPopulatedBuffer(t *testing.T) {
	t.Run("scalar_type", func(t *testing.T) {
		front := toPointer(t, int(42))
		back := toPointer(t, int(420))

		b := triplebuffer.NewPopulatedBuffer(back, nil, front)

		t.Run("correct_value", func(t *testing.T) {
			assert.NotNil(t, b.Back())
			assert.NotNil(t, b.Front())
			gotMid, gotPending := b.Middle()
			assert.Nil(t, gotMid)
			assert.False(t, gotPending)

			// pointer test
			// ...not sure if this is overkill
			// or even the right way to do it but eh?
			ptr1 := uintptr(unsafe.Pointer(front))
			ptr2 := uintptr(unsafe.Pointer(b.Front()))
			assert.Equal(t, ptr1, ptr2)

			ptr1 = uintptr(unsafe.Pointer(back))
			ptr2 = uintptr(unsafe.Pointer(b.Back()))
			assert.Equal(t, ptr1, ptr2)

			// value test
			assert.Equal(t, *front, *b.Front())
		})

		t.Run("value_updates", func(t *testing.T) {
			*front = 43
			*back = 107

			// value test
			assert.Equal(t, *front, *b.Front())
			assert.Equal(t, *back, *b.Back())
		})
	})

	t.Run("custom_type", func(t *testing.T) {
		type testType struct {
			stringData string
			bigIntData *big.Int
		}

		back := &testType{
			stringData: "hello, world!",
			bigIntData: big.NewInt(42),
		}

		b := triplebuffer.NewPopulatedBuffer(back, nil, nil)

		t.Run("correct_value", func(t *testing.T) {
			assert.NotNil(t, b.Back())
			assert.Nil(t, b.Front())

			// pointer test
			// ...not sure if this is overkill
			// or even the right way to do it but eh?
			ptr1 := uintptr(unsafe.Pointer(back))
			ptr2 := uintptr(unsafe.Pointer(b.Back()))
			assert.Equal(t, ptr1, ptr2)

			// value test
			assert.Equal(t, back, b.Back())
			assert.Equal(t, back.stringData, b.Back().stringData)
			assert.Equal(t, back.bigIntData.Int64(), b.Back().bigIntData.Int64())
		})

		t.Run("value_updates", func(t *testing.T) {
			*back = testType{
				stringData: "KOOLAIDMAN",
				bigIntData: nil,
			}

			// value test
			assert.Equal(t, back, b.Back())
			assert.Nil(t, b.Back().bigIntData)
			assert.Equal(t, back.stringData, b.Back().stringData)
		})
	})
}

func TestBuffer_ProducerAPI(t *testing.T) {
	b := triplebuffer.NewBuffer[int]()

	data := toPointer(t, int(31337))

	t.Run("write", func(t *testing.T) {
		b.Write(data)
		assert.NotNil(t, b.Back())

		gotMid, gotPending := b.Middle()
		assert.Nil(t, gotMid)
		assert.False(t, gotPending)
	})

	t.Run("commit", func(t *testing.T) {
		b.Commit()
		assert.Nil(t, b.Back())

		gotMid, gotPending := b.Middle()
		assert.NotNil(t, gotMid)
		assert.True(t, gotPending)

		t.Run("overwrites", func(t *testing.T) {
			newData := toPointer(t, int(1337))

			b.Write(newData)
			assert.Equal(t, *data, *gotMid)

			b.Commit()
			gotMid, gotPending = b.Middle()

			ptr1 := uintptr(unsafe.Pointer(newData))
			ptr2 := uintptr(unsafe.Pointer(gotMid))
			assert.Equal(t, ptr1, ptr2)

			assert.Equal(t, *newData, *gotMid)
		})
	})
}

func TestBuffer_ConsumerAPI(t *testing.T) {
	b := triplebuffer.NewBuffer[int]()

	data := toPointer(t, int(31337))

	b.Write(data)
	b.Commit()

	t.Run("dirty_read", func(t *testing.T) {
		gotData, gotPending := b.Read()
		assert.NotNil(t, gotData)
		assert.True(t, gotPending)

		assert.Equal(t, *data, *gotData)
	})

	t.Run("stale_read", func(t *testing.T) {
		gotData, gotPending := b.Read()
		assert.NotNil(t, gotData)
		assert.False(t, gotPending)

		assert.Equal(t, *data, *gotData)
	})

	t.Run("produce_new", func(t *testing.T) {
		newData := toPointer(t, int(1337))
		b.Write(newData)

		gotData, gotPending := b.Read()
		assert.NotNil(t, gotData)
		assert.False(t, gotPending)

		assert.NotEqual(t, *newData, *gotData)

		b.Commit()

		gotData, gotPending = b.Read()
		assert.NotNil(t, gotData)
		assert.True(t, gotPending)

		assert.Equal(t, *newData, *gotData)
	})
}
