package Utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Rand struct {
}

func (Rand) RandCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var buffer strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&buffer, "%d", numeric[rand.Intn(r)])
	}
	return buffer.String()
}
