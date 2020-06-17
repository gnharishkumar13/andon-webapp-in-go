package hash

import (
	"crypto/sha512"
	"encoding/base64"

	"github.com/user/andon-webapp-in-go/src/config"
)

// Hash returns a base64-encoded, salted hash of the provided strings
func Hash(data ...string) string {
	h := sha512.New()
	for _, d := range data {
		h.Write([]byte(d))
	}
	h.Write([]byte(config.Get().Salt))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
