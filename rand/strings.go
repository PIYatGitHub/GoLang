package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytes = 32

//Bytes will take in an integer n - the bytes count and will
//return the bytes or an err if it was one.
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//String will generate a string drom a byte slice (as a BASE6 encoded string!!!)
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

//RememberToken generatest the RememberToken for a predifined # of bytes
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
