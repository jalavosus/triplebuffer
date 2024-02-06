package triplebuffer_test

import (
	"fmt"
	"time"

	"github.com/jalavosus/triplebuffer"
)

// Example of instantiating an empty buffer and
// writing to it every 750 milliseconds. Thanks to the timer in produce(),
// reads will rarely be brand new pending reads, but sometimes they can be if the
// stars align.
func ExampleNewBuffer() {
	const maxProduceInt int = 42

	produce := func(buf *triplebuffer.Buffer[int]) {
		ticker := time.NewTicker(750 * time.Millisecond)

		for i := 0; i <= maxProduceInt; i++ {
			buf.Write(&i)
			buf.Commit()
			<-ticker.C
		}
	}

	buf := triplebuffer.NewBuffer[int]()
	go produce(buf)

	var prevRead int

	for {
		val, pending := buf.Read()
		if val != nil && *val != prevRead {
			dVal := *val
			fmt.Printf("val: %d - pending %t\n", dVal, pending)

			if dVal >= maxProduceInt {
				break
			}

			prevRead = dVal
		}
	}
}
