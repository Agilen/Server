package mycrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/crypto/curve25519"
)

type CryptoContext struct {
	PrivateKey [32]byte
	PublicKey  [32]byte
	PublicInfo [32]byte
}

type CMS struct {
	PublicInfo []byte
	WrapedKey  []byte
	EncData    []byte
}

func NewCryptoContext() (*CryptoContext, error) {
	cc := new(CryptoContext)

	if err := cc.GenerateKeys(); err != nil {
		return nil, err
	}

	return cc, nil
}

func Encrypt(key, data []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create new cipher: %v", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(data))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("could not encrypt: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], data)

	return cipherText, nil
}

func Decrypt(key, data []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create new cipher: %v", err)
	}

	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("invalid ciphertext block size")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return data, nil
}

//-----------------------------------------------------

func (cc *CryptoContext) GenerateKeys() error {

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("cant generate private key: %v", err)
	}
	copy(cc.PrivateKey[:], privateKey.D.Bytes())
	copy(cc.PublicKey[:], privateKey.PublicKey.X.Bytes())
	var dst [32]byte
	curve25519.ScalarBaseMult(&dst, &cc.PrivateKey)
	copy(cc.PublicInfo[:], dst[:])

	return nil
}

func (cc *CryptoContext) ECDH(peersPublicKey []byte) ([]byte, error) {
	if peersPublicKey == nil {
		return nil, fmt.Errorf("peersPublicKey == nil")
	}

	if bytes.Equal(cc.PublicKey[:], peersPublicKey) {
		return nil, fmt.Errorf("PublicKey == peersPublicKey")
	}

	var sharedSecret [32]byte
	var ppk [32]byte
	copy(ppk[:], peersPublicKey)
	curve25519.ScalarMult(&sharedSecret, &cc.PrivateKey, &ppk)

	return sharedSecret[:], nil
}

func (cc *CryptoContext) Encrypt(data []byte, publicKey []byte) ([]byte, error) {
	var key [32]byte
	_, err := rand.Read(key[:])
	if err != nil {
		return nil, fmt.Errorf("cant generate key: %v", err)
	}

	encData, err := Encrypt(key[:], data)
	if err != nil {
		return nil, fmt.Errorf("error in encrypt: %v", err)
	}

	sharedSecret, err := cc.ECDH(publicKey)
	if err != nil {
		return nil, fmt.Errorf("cant make shared secret: %v", err)
	}

	chipher, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return nil, fmt.Errorf("cant create new chipher: %v", err)
	}

	wrapKey, err := Wrap(chipher, key[:])
	if err != nil {
		return nil, fmt.Errorf("error in key wrap: %v", err)
	}

	cms := CMS{
		PublicInfo: cc.PublicInfo[:],
		WrapedKey:  wrapKey,
		EncData:    encData,
	}

	finalData, err := json.Marshal(cms)
	if err != nil {
		return nil, fmt.Errorf("error in marshal: %v", err)
	}

	return finalData, nil
}

func (cc *CryptoContext) Decrypt(data []byte) ([]byte, error) {
	var cms CMS
	err := json.Unmarshal(data, &cms)
	if err != nil {
		return nil, fmt.Errorf("cant unmarshal cms: %v", err)
	}

	secret, err := cc.ECDH(cms.PublicInfo)
	if err != nil {
		return nil, fmt.Errorf("cant make secret: %v", err)
	}

	cihper, err := aes.NewCipher(secret)
	if err != nil {
		return nil, fmt.Errorf("cant create new chipher: %v", err)
	}

	key, err := Unwrap(cihper, cms.WrapedKey)
	if err != nil {
		return nil, fmt.Errorf("cant unwrap key:%v", err)
	}

	finalData, err := Decrypt(key, cms.EncData)
	if err != nil {
		return nil, fmt.Errorf("cant decrypt data: %v", err)
	}

	return finalData, nil
}
