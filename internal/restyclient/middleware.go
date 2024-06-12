package restyclient

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/robotjoosen/go-tplink-deco-client/internal/crypto"
)

type BodyAware interface {
	Body() []byte
	SetBody(body []byte) interface{}
}

type DecoMiddleware struct {
	stok     string
	sequence uint
	hash     string
	aes      *crypto.AESKey
	rsa      *rsa.PublicKey
}

func NewDecoMiddleware(aes *crypto.AESKey, hash string, rsa *rsa.PublicKey, sequence uint) *DecoMiddleware {
	return &DecoMiddleware{
		aes:      aes,
		hash:     hash,
		rsa:      rsa,
		sequence: sequence,
	}
}

func (c *DecoMiddleware) SetStok(stok string) {
	c.stok = stok
}

func (c *DecoMiddleware) ParseRequest(_ *resty.Client, request *resty.Request) error {
	body, err := c.encrypt(request.Body)
	if err != nil {
		return err
	}

	request.URL = "/;stok=" + c.stok + request.URL
	request.SetBody(body)
	request.SetHeader("Content-Type", "application/json")

	return nil
}

func (c *DecoMiddleware) ParseResponse(_ *resty.Client, response *resty.Response) error {
	body, err := c.decrypt(response.Body())
	if err != nil {
		return err
	}

	response.SetBody(body)

	return nil
}

func (c *DecoMiddleware) decrypt(data []byte) ([]byte, error) {
	var result struct {
		Data string `json:"data"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	decoded, err := crypto.AES256Decrypt(result.Data, *c.aes)
	if err != nil {
		return nil, err
	}

	return []byte(decoded), nil
}

func (c *DecoMiddleware) encrypt(data interface{}) (string, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	encryptedData, err := crypto.AES256Encrypt(string(body), *c.aes)
	if err != nil {
		return "", err
	}

	var sign string
	length := int(c.sequence) + len(encryptedData)

	switch c.stok {
	case "":
		sign = fmt.Sprintf("k=%s&i=%s&h=%s&s=%v", c.aes.Key, c.aes.Iv, c.hash, length)
	default:
		sign = fmt.Sprintf("h=%s&s=%v", c.hash, length)
	}

	switch {
	case len(sign) > 53:
		first, _ := crypto.EncryptRsa(sign[:53], c.rsa)
		second, _ := crypto.EncryptRsa(sign[53:], c.rsa)
		sign = fmt.Sprintf("%s%s", first, second)
	default:
		sign, _ = crypto.EncryptRsa(sign, c.rsa)
	}

	return fmt.Sprintf(
		"sign=%s&data=%s",
		url.QueryEscape(sign),
		url.QueryEscape(encryptedData),
	), nil
}
