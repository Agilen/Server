package mycrypto_test

import (
	"fmt"
	"testing"

	"github.com/Agilen/Server/mycrypto"
	"github.com/stretchr/testify/assert"
)

func TestKeyExchange(t *testing.T) {
	alice, err := mycrypto.NewCryptoContext()
	if !assert.NoError(t, err, " New cc shouldn't have an error") {
		panic(err)
	}
	bob, err := mycrypto.NewCryptoContext()
	if !assert.NoError(t, err, " New cc shouldn't have an error") {
		panic(err)
	}

	ares, err := alice.ECDH(bob.PublicInfo[:])
	if !assert.NoError(t, err, "ECDH alice should not throw error with valid input") {
		panic(err)
	}

	bres, err := bob.ECDH(alice.PublicInfo[:])
	if !assert.NoError(t, err, "ECDH bob should not throw error with valid input") {
		panic(err)
	}

	if !assert.Equal(t, bres, ares, "Wrap Mismatch: Actual wrapped ciphertext should equal expected for test case '%s'") {
		panic(fmt.Errorf("not equal"))
	}

	encData, err := alice.Encrypt([]byte("hello"), bob.PublicInfo[:])
	if !assert.NoError(t, err, "Encrypt alice should not throw error with valid input") {
		panic(err)
	}

	decData, err := bob.Decrypt(encData)
	if !assert.NoError(t, err, "Encrypt alice should not throw error with valid input") {
		panic(err)
	}

	if !assert.Equal(t, []byte("hello"), decData, "ENCDEC error '%s'") {
		panic(fmt.Errorf("not equal"))
	}
}
