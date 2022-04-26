package safeweb_lib_helper

import (
	"errors"

	"github.com/Luzifer/go-openssl"
)

// Decrypt from CryptoJS AES
func Decrypt(encrypted string) (string, error) {
	secret := "safeweb"

	o := openssl.New()

	dec, err := o.DecryptBytes(secret, []byte(encrypted))
	if err != nil {
		return "", errors.New("Decrypt error: %s")
	}

	return string(dec), nil

}
