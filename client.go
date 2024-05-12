package tplink_deco_client

import (
	"context"
	"crypto/md5"
	rsa2 "crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/robotjoosen/go-tplink-deco-client/internal/aes256"
	dto "github.com/robotjoosen/go-tplink-deco-client/internal/model"
	"github.com/robotjoosen/go-tplink-deco-client/internal/rsa"
	"github.com/robotjoosen/go-tplink-deco-client/pkg/model"
)

type Cipher interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
	GetKeys() map[string]interface{}
}

type Client struct {
	signatureRSA Cipher
	passwordRSA  Cipher
	aes          Cipher
	stok         string
	client       *resty.Client
}

func New(baseURL string) *Client {
	return &Client{
		client: resty.
			NewWithClient(http.DefaultClient).
			SetBaseURL(baseURL),
	}
}

// link: https://gist.github.com/rosmo/29200c1aedb991ce55942c4ae8b54edd
func (c *Client) Login(password string) error {
	c.aes = aes256.New(nil, nil) // key and iv are auto generated

	// Fetch 1024-bit RSA public key for password encryption
	passwordPubkey, err := c.getPasswordEncryptionPublicKey()
	if err != nil {
		slog.Error("getting password pub key failed", "err", err.Error())
		return err
	}
	c.passwordRSA = rsa.New(passwordPubkey, -1)

	// Fetch 512-bit RSA public key used for signature encryption
	signatureSeq, signaturePubKey, err := c.getSignaturePublicKey()
	if err != nil {
		slog.Error("getting signature public key failed", "err", err.Error())
		return err
	}
	c.signatureRSA = rsa.New(signaturePubKey, signatureSeq)

	encryptedAuthData, err := c.generateAuthenticationData(password)
	if err != nil {
		slog.Error("auth body failed", "err", err.Error())
		return err
	}

	encryptedSignature, err := c.generateSignature(password, signatureSeq, len(encryptedAuthData))
	if err != nil {
		slog.Error("signature failed", "err", err.Error())
		return err
	}

	urlParams := url.Values{}
	urlParams.Add("form", "login")
	rsp, err := c.client.NewRequest().
		SetHeaders(map[string]string{
			"Host":             "192.168.2.1",
			"Origin":           c.client.BaseURL,
			"Referer":          c.client.BaseURL + "/webpages/index.html",
			"Content-Type":     "application/json",
			"X-Requested-With": "XMLHttpRequest",
		}).
		SetBody("sign=" + string(encryptedSignature) + "&data=" + url.QueryEscape(string(encryptedAuthData))).
		Post("/cgi-bin/luci/;stok=/login?" + urlParams.Encode())
	if err != nil {
		return err
	}

	if !rsp.IsSuccess() {
		fmt.Println(rsp.StatusCode(), rsp.Request.Body)

		return errors.New("failed to get authentication")
	}

	content, _ := c.aes.Decrypt(rsp.Body())
	fmt.Println(content)

	return err
}

func (c *Client) GetClients() (model.Clients, error) {
	data := `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`

	content, err := c.aes.Decrypt([]byte(data))
	if err != nil {
		return nil, err
	}

	//if valid := json.Valid(content); !valid {
	//	return nil, errors.New("invalid JSON")
	//}

	res := dto.ClientListResponse{}
	if err := json.Unmarshal(content, &res); err != nil {
		fmt.Println(string(content))

		return nil, err
	}

	slog.Info("got", slog.Any("data", res.Result))

	return make(model.Clients, 0), nil
}

func (c *Client) GetDevices() (model.Devices, error) {
	data := `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`

	content, err := c.aes.Decrypt([]byte(data))
	if err != nil {
		return nil, err
	}

	//if valid := json.Valid(content); !valid {
	//	return nil, errors.New("invalid JSON")
	//}

	res := dto.DeviceListResponse{}
	if err := json.Unmarshal(content, &res); err != nil {
		fmt.Println(string(content))

		return nil, err
	}

	slog.Info("got", slog.Any("data", res.Result))

	return make(model.Devices, 0), nil
}

// Fetch 1024-bit RSA public key used for encrypting password
func (c *Client) generateAuthenticationData(password string) ([]byte, error) {
	encryptedPassword, err := c.passwordRSA.Encrypt([]byte(password))
	if err != nil {
		slog.Error("password encryption failed", "err", err.Error())
		return nil, err
	}

	//`{"params":{"password":%s},"operation":"login"}`,
	encryptedBody, err := c.aes.Encrypt([]byte(fmt.Sprintf(`password=%s&operation=login`, encryptedPassword)))
	if err != nil {
		slog.Error("auth body encryption failed", "err", err.Error())
		return nil, err
	}

	b64EncBody := make([]byte, base64.StdEncoding.EncodedLen(len(encryptedBody)))
	base64.StdEncoding.Encode(b64EncBody, encryptedBody)

	return b64EncBody, nil
}

func (c *Client) generateSignature(password string, sequence, bodyLength int) ([]byte, error) {
	aesKeys := c.aes.GetKeys()
	signature := []byte(fmt.Sprintf(`k=%s&i=%s&h=%s&s=%d`,
		aesKeys[aes256.Key],
		aesKeys[aes256.IV],
		getMD5Hash("admin"+password),
		sequence+bodyLength,
	))

	var err error
	encryptedSignature := make([]byte, 0)
	if len(signature) > 53 {
		encryptedSignature, err = c.signatureRSA.Encrypt(signature[0:53])
		if err != nil {
			return nil, err
		}

		res, err := c.signatureRSA.Encrypt(signature[53:])
		if err != nil {
			return nil, err
		}

		encryptedSignature = append(encryptedSignature, res...)
	} else {
		encryptedSignature, err = c.signatureRSA.Encrypt(signature)
		if err != nil {
			return nil, err
		}
	}

	b64EncSignature := make([]byte, base64.StdEncoding.EncodedLen(len(encryptedSignature)))
	base64.StdEncoding.Encode(b64EncSignature, encryptedSignature)

	return b64EncSignature, nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))

	return hex.EncodeToString(hash[:])
}

func (c *Client) getPasswordEncryptionPublicKey() (rsa2.PublicKey, error) {
	urlParams := url.Values{}
	urlParams.Add("form", "keys")
	rsp, err := c.client.NewRequest().
		SetContext(context.Background()).
		SetHeader("Content-Type", "application/json").
		SetBody(dto.OperationRequest{Operation: "read"}).
		Post("/cgi-bin/luci/;stok=/login?" + urlParams.Encode())
	if err != nil {
		return rsa2.PublicKey{}, err
	}

	if !rsp.IsSuccess() {
		err = errors.New("failed to get RSA Key")

		return rsa2.PublicKey{}, err
	}
	body := rsp.Body()
	//var err error
	//body := []byte(`{ "result":{ "username":"", "password":[ "D1E79FF135D14E342D76185C23024E6DEAD4D6EC2C317A526C811E83538EA4E5ED8E1B0EEE5CE26E3C1B6A5F1FE11FA804F28B7E8821CA90AFA5B2F300DF99FDA27C9D2131E031EA11463C47944C05005EF4C1CE932D7F4A87C7563581D9F27F0C305023FCE94997EC7D790696E784357ED803A610EBB71B12A8BE5936429BFD", "010001" ] }, "error_code":0 }`)

	res := dto.LoginKeyResponse{}
	if err = json.Unmarshal(body, &res); err != nil {
		return rsa2.PublicKey{}, err
	}

	if len(res.Result.Password) < 2 {
		return rsa2.PublicKey{}, errors.New("missing passwordRSA values")
	}

	return generateRSAPublicKey(res.Result.Password[0], res.Result.Password[1])
}

func (c *Client) getSignaturePublicKey() (int, rsa2.PublicKey, error) {
	urlParams := url.Values{}
	urlParams.Add("form", "auth")
	rsp, err := c.client.NewRequest().
		SetContext(context.Background()).
		SetHeader("Content-Type", "application/json").
		SetBody(dto.OperationRequest{Operation: "read"}).
		Post("/cgi-bin/luci/;stok=/login?" + urlParams.Encode())
	if err != nil {
		return 0, rsa2.PublicKey{}, err
	}

	if !rsp.IsSuccess() {
		return 0, rsa2.PublicKey{}, errors.New("failed to get RSA Key")
	}
	body := rsp.Body()
	//var err error
	//body := []byte(`{ "result":{ "key":[ "963046394FC5C5B06CD2ED7C4C837CA533621FC93BC9B85C42FF6FCF615E6BE5A3473928CE0EEC0791CFF319830056437CA59322A88C7F48E13EB2DE312D0B4B", "010001" ], "seq":718248345 }, "error_code":0 }`)

	res := dto.LoginAuthResponse{}
	if err = json.Unmarshal(body, &res); err != nil {
		return 0, rsa2.PublicKey{}, err
	}

	pubKey, err := generateRSAPublicKey(res.Result.Key[0], res.Result.Key[1])

	return res.Result.Seq, pubKey, err
}

func generateRSAPublicKey(n, e string) (rsa2.PublicKey, error) {
	modulus, err := decodeKey(strings.ToLower(n))
	if err != nil {
		return rsa2.PublicKey{}, err
	}

	exponent, err := decodeBitString(e)
	if err != nil {
		return rsa2.PublicKey{}, err
	}

	return rsa2.PublicKey{N: modulus, E: exponent}, nil
}

func decodeKey(input string) (*big.Int, error) {
	bytes, err := hex.DecodeString(input)
	if err != nil {
		return nil, err
	}

	output := big.NewInt(0)
	output.SetBytes(bytes)

	return output, nil
}

func decodeBitString(input string) (int, error) {
	decodedExponent, err := hex.DecodeString(input)
	if err != nil {
		return 0, err
	}

	var number uint32
	for _, bit := range decodedExponent {
		number = (number << 1) | uint32(bit)
	}

	return int(number), nil
}
