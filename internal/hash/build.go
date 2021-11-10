package hash

import "crypto/sha256"

func BuildHash(tokens ...string) []byte {
	h := sha256.New()
	for _, t := range tokens {
		_, err := h.Write([]byte(t))
		if err != nil {
			panic(err)
		}
	}
	return h.Sum(nil)
}
