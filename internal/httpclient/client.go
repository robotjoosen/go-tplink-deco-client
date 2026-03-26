// Package httpclient is a rewrite of the MrMarble/deco package
package httpclient

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/robotjoosen/go-tplink-deco-client/internal/crypto"
	"github.com/robotjoosen/go-tplink-deco-client/internal/model"
)

const defaultPath = ";stok=/login"

type HTTPClient struct {
	client   *http.Client
	baseUrl  url.URL
	aes      *crypto.AESKey
	rsa      *rsa.PublicKey
	hash     string
	stok     string
	sequence uint
}

func NewClient(ip string) *HTTPClient {
	jar, _ := cookiejar.New(nil)

	return &HTTPClient{
		client: &http.Client{Timeout: 10 * time.Second, Jar: jar},
		baseUrl: url.URL{
			Host:   ip,
			Scheme: "http",
			Path:   "/cgi-bin/luci/",
		},
	}
}

func (c *HTTPClient) Login(ctx context.Context, username, password string) (model.LoginResponse, error) {
	var res model.LoginResponse

	c.aes = crypto.GenerateAESKey()
	c.hash = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s", username, password))))

	passwordKey, err := c.getPasswordKey(ctx)
	if err != nil {
		return res, err
	}

	sessionKey, sequence, err := c.getSessionKey(ctx)
	if err != nil {
		return res, err
	}

	c.rsa = sessionKey
	c.sequence = sequence

	encryptedPassword, err := crypto.EncryptRsa(password, passwordKey)
	if err != nil {
		return res, err
	}

	loginReq := model.OperationRequest{
		Operation: "login",
		Params:    map[string]interface{}{"password": encryptedPassword},
	}

	loginJSON, err := json.Marshal(loginReq)
	if err != nil {
		return res, err
	}

	args := url.Values{}
	args.Set("form", "login")

	err = c.encryptPost(ctx, defaultPath, args, loginJSON, true, &res)
	if err != nil {
		return res, err
	}

	c.stok = res.Result.Stok

	return res, nil
}

func (c *HTTPClient) GetDevices(ctx context.Context) (model.DeviceListResponse, error) {
	var res model.DeviceListResponse

	readBody, err := json.Marshal(model.OperationRequest{Operation: "read"})
	if err != nil {
		return res, err
	}

	args := url.Values{}
	args.Add("form", "device_list")

	err = c.encryptPost(ctx, fmt.Sprintf(";stok=%s/admin/device", c.stok), args, readBody, false, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *HTTPClient) GetClients(ctx context.Context) (model.ClientListResponse, error) {
	var res model.ClientListResponse

	readBody, err := json.Marshal(model.OperationRequest{
		Operation: "read",
		Params: map[string]interface{}{
			"device_mac": "default",
		},
	})
	if err != nil {
		return res, err
	}

	args := url.Values{}
	args.Add("form", "client_list")

	err = c.encryptPost(ctx, fmt.Sprintf(";stok=%s/admin/client", c.stok), args, readBody, false, &res)
	if err != nil {
		return res, err
	}

	// decode client names
	for i := range res.Result.ClientList {
		name, err := base64.StdEncoding.DecodeString(res.Result.ClientList[i].Name)
		if err == nil {
			res.Result.ClientList[i].Name = string(name)
		}
	}

	return res, nil
}

func (c *HTTPClient) GetPerformance(ctx context.Context) (model.PerformanceResponse, error) {
	var res model.PerformanceResponse

	readBody, err := json.Marshal(model.OperationRequest{Operation: "read"})
	if err != nil {
		return res, err
	}

	args := url.Values{}
	args.Add("form", "performance")

	err = c.encryptPost(ctx, fmt.Sprintf(";stok=%s/admin/network", c.stok), args, readBody, false, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *HTTPClient) RebootDevice(ctx context.Context, macAddrs []string) error {
	var macList []map[string]string
	for _, mac := range macAddrs {
		macList = append(macList, map[string]string{
			"mac": strings.ToUpper(mac),
		})
	}

	rebootReq := model.OperationRequest{
		Operation: "reboot",
		Params: map[string]interface{}{
			"mac_list": macList,
		},
	}

	jsonBody, err := json.Marshal(rebootReq)
	if err != nil {
		return err
	}

	args := url.Values{}
	args.Add("form", "system")

	var res model.RebootResponse
	return c.encryptPost(ctx, fmt.Sprintf(";stok=%s/admin/device", c.stok), args, jsonBody, false, &res)
}

func (c *HTTPClient) getPasswordKey(ctx context.Context) (*rsa.PublicKey, error) {
	readBody, err := json.Marshal(model.OperationRequest{Operation: "read"})
	if err != nil {
		return nil, err
	}

	args := url.Values{}
	args.Add("form", "keys")

	var res model.LoginKeyResponse
	if err := c.post(ctx, defaultPath, args, readBody, &res); err != nil {
		return nil, err
	}

	key, err := crypto.GenerateRsaKey(res.Result.Password[0], res.Result.Password[1])
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (c *HTTPClient) getSessionKey(ctx context.Context) (*rsa.PublicKey, uint, error) {
	readBody, err := json.Marshal(model.OperationRequest{Operation: "read"})
	if err != nil {
		return nil, 0, err
	}

	args := url.Values{}
	args.Add("form", "auth")

	var res model.SessionKeyResponse
	if err := c.post(ctx, defaultPath, args, readBody, &res); err != nil {
		return nil, 0, err
	}

	key, err := crypto.GenerateRsaKey(res.Result.Key[0], res.Result.Key[1])
	if err != nil {
		return nil, 0, err
	}

	return key, res.Result.Seq, nil
}

func (c *HTTPClient) post(ctx context.Context, path string, params url.Values, body []byte, result interface{}) error {
	endpoint := c.baseUrl.ResolveReference(&url.URL{Path: path}).String()

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.URL.RawQuery = params.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(&result)
}

func (c *HTTPClient) encryptPost(
	ctx context.Context,
	path string,
	params url.Values,
	body []byte,
	isLogin bool,
	result interface{},
) error {
	encryptedData, err := crypto.AES256Encrypt(string(body), *c.aes)
	if err != nil {
		return err
	}

	var sign string
	length := int(c.sequence) + len(encryptedData)

	switch isLogin {
	case true:
		sign = fmt.Sprintf("k=%s&i=%s&h=%s&s=%v", c.aes.Key, c.aes.Iv, c.hash, length)
	case false:
		sign = fmt.Sprintf("h=%s&s=%v", c.hash, length)
	}

	if len(sign) > 53 {
		first, _ := crypto.EncryptRsa(sign[:53], c.rsa)
		second, _ := crypto.EncryptRsa(sign[53:], c.rsa)
		sign = fmt.Sprintf("%s%s", first, second)
	} else {
		sign, _ = crypto.EncryptRsa(sign, c.rsa)
	}

	postData := fmt.Sprintf("sign=%s&data=%s", url.QueryEscape(sign), url.QueryEscape(encryptedData))

	var res struct {
		Data string `json:"data"`
	}
	err = c.post(ctx, path, params, []byte(postData), &res)
	if err != nil {
		return err
	}

	decoded, err := crypto.AES256Decrypt(res.Data, *c.aes)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(decoded), &result)
}
