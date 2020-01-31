package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
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
			return hashToEmoji(id)
		}
	}
	return string(id)
}

func hashToEmoji(id []byte) string {
	shorty := hex.EncodeToString(id[:7])
	dec, _ := strconv.ParseInt(shorty, 16, 64)
	n := int(dec) % len(emoji)
	return fmt.Sprintf("%s %s", emoji[n], shorty)
}
