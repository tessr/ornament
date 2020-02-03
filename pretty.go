package main

import (
	"fmt"

	"github.com/tessr/hmoj"
)

func nodeEncoder(id []byte, depth int, isLeaf bool) string {
	prefix := fmt.Sprintf("-%d ", depth)
	if isLeaf {
		prefix = fmt.Sprintf("\t*%d ", depth)
	}
	if len(id) == 0 {
		return fmt.Sprintf("%s<nil>", prefix)
	}
	return fmt.Sprintf("%s%s", prefix, encodeID(id))
}

// casts to a string if it is printable ascii, emoji-hex-encodes otherwise
func encodeID(id []byte) string {
	for _, b := range id {
		if b < 0x20 || b >= 0x80 {
			return hmoj.Encode(id)
		}
	}
	return string(id)
}
