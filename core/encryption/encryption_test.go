package encryption

import (
	"bytes"
	"testing"
)

func TestEncryption(t *testing.T) {

	en := NewEncryption()
	en.SetAccessKey("TestingKey")

	data, err := en.Encrypt([]byte("test"))
	if err != nil {
		t.Error(err)
	}

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
