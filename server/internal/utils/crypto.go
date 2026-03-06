package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

// ref: https://blog.logrocket.com/learn-golang-encryption-decryption/
// ref: https://pkg.go.dev/golang.org/x/crypto/pbkdf2

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

// set keySize value 16, 24 or 32
const keySize = 16

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// Encrypt method encripts the text using mySecret key
func Encrypt(text string, mySecret string) (string, error) {

	// derived key
	dk := pbkdf2.Key([]byte(mySecret), bytes, 4096, keySize, sha1.New)

	block, err := aes.NewCipher(dk)
	if err != nil {
		return "", err
	}
	plainText := []byte(text)

	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text, mySecret string) (string, error) {

	// derived key
	dk := pbkdf2.Key([]byte(mySecret), bytes, 4096, keySize, sha1.New)

	block, err := aes.NewCipher(dk)
	if err != nil {
		return "", err
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
