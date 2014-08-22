package helpers

import (
	"crypto/rand"
)

// Random string generator, thanks to:
// http://devpy.wordpress.com/2013/10/24/create-random-string-in-golang/

func RandString(str_size int) string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, str_size)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
