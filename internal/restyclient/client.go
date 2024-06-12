package restyclient

import (
	"context"
	"crypto/md5"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/robotjoosen/go-tplink-deco-client/internal/crypto"
	"github.com/robotjoosen/go-tplink-deco-client/internal/model"
)

type RestyClient struct {
	client *resty.Client
}

func NewClient(ip string) *RestyClient {
	return &RestyClient{
		client: resty.
			NewWithClient(http.DefaultClient).
			SetBaseURL("http://" + ip + "/cgi-bin/luci/"),
	}
}

func (c *RestyClient) Login(ctx context.Context, username, password string) (model.LoginResponse, error) {
	var res model.LoginResponse

	passwordKey, err := c.getPasswordKey(ctx)
	if err != nil {
		return res, err
	}

	sessionKey, sequence, err := c.getSessionKey(ctx)
	if err != nil {
		return res, err
	}

	decoMiddleware := NewDecoMiddleware(
		crypto.GenerateAESKey(),
		fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s", username, password)))),
		sessionKey,
		sequence,
	)
	c.client.OnBeforeRequest(decoMiddleware.ParseRequest)
	c.client.OnAfterResponse(decoMiddleware.ParseResponse)

	encryptedPassword, err := crypto.EncryptRsa(password, passwordKey)
	if err != nil {
		return res, err
	}

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "login").
		SetBody(model.OperationRequest{
			Operation: "login",
			Params:    map[string]interface{}{"password": encryptedPassword},
		}).
		Post("/login")
	if err != nil {
		return res, err
	}

	if !rsp.IsSuccess() {
		return res, errors.New("failed to login")
	}

	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return res, err
	}

	decoMiddleware.SetStok(res.Result.Stok)

	return res, nil
}

func (c *RestyClient) GetDevices(ctx context.Context) (model.DeviceListResponse, error) {
	var res model.DeviceListResponse

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "device_list").
		SetBody(model.OperationRequest{Operation: "read"}).
		Post("/admin/device")
	if err != nil {
		return res, err
	}

	if rsp.String() == "" {
		return res, errors.New("empty response")
	}

	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *RestyClient) getPasswordKey(ctx context.Context) (*rsa.PublicKey, error) {
	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "keys").
		SetBody(model.OperationRequest{Operation: "read"}).
		Post(";stok=/login")
	if err != nil {
		return nil, err
	}

	if !rsp.IsSuccess() {
		slog.Error("password fail", slog.String("err", string(rsp.Body())))

		return nil, errors.New("failed to get password key")
	}

	var res model.LoginKeyResponse
	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return nil, err
	}

	key, err := crypto.GenerateRsaKey(res.Result.Password[0], res.Result.Password[1])
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (c *RestyClient) getSessionKey(ctx context.Context) (*rsa.PublicKey, uint, error) {
	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "auth").
		SetBody(model.OperationRequest{Operation: "read"}).
		Post(";stok=/login")
	if err != nil {
		return nil, 0, err
	}

	if !rsp.IsSuccess() {
		return nil, 0, errors.New("failed to get password key")
	}

	var res model.SessionKeyResponse
	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return nil, 0, err
	}

	key, err := crypto.GenerateRsaKey(res.Result.Key[0], res.Result.Key[1])
	if err != nil {
		return nil, 0, err
	}

	return key, res.Result.Seq, nil
}
