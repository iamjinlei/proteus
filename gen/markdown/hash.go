package markdown

import (
	"crypto/sha256"
	"fmt"
)

func hash20(b []byte) string {
	h := sha256.New()
	data := h.Sum(b)
	n := len(data)
	for i := 0; i < n/2; i++ {
		data[i] = data[i] ^ data[n-1-i]
	}
	n /= 2
	for i := 0; i < n/2; i++ {
		data[i] = data[i] ^ data[n-1-i]
	}
	data = data[:n/2]
	return fmt.Sprintf("%x", data)
}
