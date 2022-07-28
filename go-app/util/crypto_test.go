package util

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// using a 16 bytes key
const mySecret = "SECRET_KEY_0016i"

var tests = []struct {
	plain     string
	encripted string
}{
	{"This is a string to encrypt", ""},
	{"mail@test.com", ""},
	{"this_is_another@email.com.mx", ""},
	{"This is not an email but will serve to test a long string", ""},
}

// TestEncrypt calls util.Encrypt with a string an expects to return an
// encrypted string
func TestEncrypt(t *testing.T) {
	for _, tt := range tests {
		enc, err := Encrypt(tt.plain, mySecret)
		if err != nil {
			t.Fatalf("Error encrypting: %s", err)
		}
		t.Logf("Plain: %s\n", tt.plain)
		t.Logf("Encripted: %s\n", enc)
	}
}

func TestDecript(t *testing.T) {
	for _, tt := range tests {
		// add expiry time to token for 24 hours
		expiry := time.Now().Add(time.Hour * 24)
		// simulate token value with expiry
		token := fmt.Sprintf("%s|%s", tt.plain, expiry.Format(time.RFC3339))

		// encript it
		enc, err := Encrypt(token, mySecret)
		if err != nil {
			t.Fatalf("Error encrypting: %s", err)
		}
		t.Logf("Plain: %s\n", token)
		t.Logf("Encripted: %s\n", enc)

		// decript it
		dec, err := Decrypt(enc, mySecret)
		if err != nil {
			t.Fatalf("Error decrypting: %s", err)
		}
		if dec != token {
			t.Errorf("Decripted value (%s) is not equal to token value (%s)\n", dec, token)
		}
		t.Logf("Decripted: %s\n", dec)

		// extract plain value
		split := strings.Split(dec, "|")
		// first split should be the plain value
		if split[0] != tt.plain {
			t.Errorf("Decripted split value (%s) is not equal to plain value (%s)\n", split[0], tt.plain)
		}
		t.Logf("Split value: %v", split)
	}
}
