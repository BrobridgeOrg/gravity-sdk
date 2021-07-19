package encryption

import (
	"bytes"
	"testing"
)

func TestEncryption(t *testing.T) {

	en := NewEncryption()
	en.Key = "TestingKEY01234567890123456789AB"

	data := en.Encrypt([]byte("test"))

	if bytes.Compare(data, []byte("test")) == 0 {
		t.Fail()
	}

	raw, err := en.Decrypt(data)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(raw, []byte("test")) != 0 {
		t.Fail()
	}
}
