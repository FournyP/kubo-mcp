package common_tests

import "crypto/rand"

// RandomContent returns unique bytes, so its CID is never shared with another test.
func RandomContent(size int) []byte {
	content := make([]byte, size)
	_, _ = rand.Read(content)
	return content
}
