package main

import (
	"fmt"
	"time"

	"github.com/jalavosus/triplebuffer"
)

const (
	maxProduceInt int = 25
)

func produce(buf *triplebuffer.Buffer[int]) {
	ticker := time.NewTicker(750 * time.Millisecond)

	for i := 0; i <= maxProduceInt; i++ {
		buf.Write(&i)
		buf.Commit()
		<-ticker.C
	}
}

func main() {
	buf := triplebuffer.NewBuffer[int]()
	go produce(buf)

	var prevRead int
	// ticker := time.NewTicker(748 * time.Millisecond)

	for {
		// <-ticker.C
		val, ok := buf.Read()
		if val != nil && *val != prevRead {
			dVal := *val
			fmt.Printf("val: %d - ok %t\n", dVal, ok)

			if dVal >= maxProduceInt {
				break
			}

			prevRead = dVal
		}
	}
}
