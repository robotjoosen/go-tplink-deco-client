package aes256

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"math/rand"
	"time"
)

const (
	Key string = "key"
	IV  string = "iv"
)

type Crypto struct {
	key []byte
	iv  []byte
}

func New(key, iv []byte) Crypto {
	if key == nil {
		key = make([]byte, 16)
		binary.LittleEndian.PutUint64(key, uint64(time.Now().UnixMilli()+1e9*rand.Int63n(1e16)))
	}

	if iv == nil {
		iv = make([]byte, 16)
		binary.LittleEndian.PutUint64(iv, uint64(time.Now().UnixMilli()+1e9*rand.Int63n(1e16)))
	}

	return Crypto{key: key, iv: iv}
}

func (c Crypto) GetKeys() map[string]interface{} {
	return map[string]interface{}{
		Key: c.key,
		IV:  c.iv,
	}
}

func (c Crypto) Decrypt(content []byte) ([]byte, error) {
	if len(content) == 0 {
		return nil, errors.New("empty")
	}

	//crypt, err := base64.StdEncoding.DecodeString(content)
	//if err != nil {
	//	return nil, err
	//}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	ecb := cipher.NewCBCDecrypter(block, c.iv)

	decrypted := make([]byte, len(content))
	ecb.CryptBlocks(decrypted, content)

	return pkcs5Trimming(decrypted), nil
}

func (c Crypto) Encrypt(content []byte) ([]byte, error) {
	if len(content) == 0 {
		return nil, errors.New("empty")
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	ecb := cipher.NewCBCEncrypter(block, c.iv)
	content = pkcs5Padding(content, block.BlockSize())

	encrypted := make([]byte, len(content))
	ecb.CryptBlocks(encrypted, content)

	//encoded := make([]byte, base64.StdEncoding.EncodedLen(len(encrypted)))
	//base64.StdEncoding.Encode(encoded, encrypted)

	return encrypted, nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

func pkcs5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]

	return encrypt[:len(encrypt)-int(padding)]
}
