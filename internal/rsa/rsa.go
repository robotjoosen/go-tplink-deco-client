package rsa

import (
	"crypto/rand"
	"crypto/rsa"
)

const (
	Exponent string = "exponent"
	ModulusN string = "modulus_n"
	Sequence string = "seq"
)

type Crypto struct {
	publicKey rsa.PublicKey
	sequence  int
}

func New(pubKey rsa.PublicKey, seq int) Crypto {
	return Crypto{
		publicKey: pubKey,
		sequence:  seq,
	}
}

func (c Crypto) GetKeys() map[string]interface{} {
	return map[string]interface{}{
		ModulusN: c.publicKey.N,
		Exponent: c.publicKey.E,
		Sequence: c.sequence,
	}
}

func (c Crypto) Encrypt(input []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(
		rand.Reader,
		&c.publicKey,
		input,
	)
}

func (c Crypto) Decrypt(_ []byte) ([]byte, error) {
	panic("not implemented")
}
