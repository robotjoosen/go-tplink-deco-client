package restyclient

import (
	"context"
	"crypto/md5"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/robotjoosen/go-tplink-deco-client/internal/crypto"
	"github.com/robotjoosen/go-tplink-deco-client/internal/model"
)

type RestyClient struct {
	client         *resty.Client
	decoMiddleware *DecoMiddleware
	authenticated  bool
}

func NewClient(ip string) *RestyClient {
	decoMiddleware := NewDecoMiddleware()

	client := resty.
		NewWithClient(http.DefaultClient).
		SetBaseURL("http://" + ip + "/cgi-bin/luci/")
	client.OnBeforeRequest(decoMiddleware.ParseRequest)
	client.OnAfterResponse(decoMiddleware.ParseResponse)

	return &RestyClient{
		client:         client,
		decoMiddleware: decoMiddleware,
	}
}

func (c *RestyClient) Login(ctx context.Context, username, password string) (model.LoginResponse, error) {
	var res model.LoginResponse

	if c.authenticated {
		return res, errors.New("already authenticated")
	}

	passwordKey, err := c.getPasswordKey(ctx)
	if err != nil {
		return res, err
	}

	sessionKey, sequence, err := c.getSessionKey(ctx)
	if err != nil {
		return res, err
	}

	c.decoMiddleware.Initialize(
		crypto.GenerateAESKey(),
		fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s", username, password)))),
		sessionKey,
		sequence,
	)

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

	c.decoMiddleware.SetStok(res.Result.Stok)
	c.authenticated = true

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

func (c *RestyClient) GetClients(ctx context.Context) (model.ClientListResponse, error) {
	var res model.ClientListResponse

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "client_list").
		SetBody(model.OperationRequest{
			Operation: "read",
			Params:    map[string]interface{}{"device_mac": "default"},
		}).
		Post("/admin/client")
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return res, err
	}

	for i := range res.Result.ClientList {
		name, err := base64.StdEncoding.DecodeString(res.Result.ClientList[i].Name)
		if err == nil {
			res.Result.ClientList[i].Name = string(name)
		}
	}

	return res, nil
}

func (c *RestyClient) GetPerformance(ctx context.Context) (model.PerformanceResponse, error) {
	var res model.PerformanceResponse

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "performance").
		SetBody(model.OperationRequest{Operation: "read"}).
		Post("/admin/network")
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *RestyClient) RebootDevice(ctx context.Context, macAddrs []string) error {
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

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "system").
		SetBody(rebootReq).
		Post("/admin/device")
	if err != nil {
		return err
	}

	var res model.ErrorResponse
	return json.Unmarshal(rsp.Body(), &res)
}

func (c *RestyClient) Custom(ctx context.Context, path string, form string, body []byte) (interface{}, error) {
	var result interface{}

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", form).
		SetBody(body).
		Post(path)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(rsp.Body(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *RestyClient) GetWANIPv4Info(ctx context.Context) (model.WANIPv4Response, error) {
	var res model.WANIPv4Response

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "wan_ipv").
		SetBody(model.OperationRequest{Operation: "read"}).
		Post("/admin/network")
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *RestyClient) GetWiFiSettings(ctx context.Context) (model.WiFiResponse, error) {
	var res model.WiFiResponse

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "wlan").
		SetBody(model.OperationRequest{Operation: "read"}).
		Post("/admin/wireless")
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *RestyClient) SetWiFiSettings(ctx context.Context, settings map[string]interface{}) error {
	writeReq := model.OperationRequest{
		Operation: "write",
		Params:    settings,
	}

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "wlan").
		SetBody(writeReq).
		Post("/admin/wireless")
	if err != nil {
		return err
	}

	var res model.ErrorResponse
	return json.Unmarshal(rsp.Body(), &res)
}

func (c *RestyClient) GetInternetStatus(ctx context.Context) (model.InternetStatusResponse, error) {
	var res model.InternetStatusResponse

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "internet").
		SetBody(model.OperationRequest{Operation: "read"}).
		Post("/admin/network")
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *RestyClient) GetIPv6Settings(ctx context.Context) (model.IPv6Response, error) {
	var res model.IPv6Response

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "ipv").
		SetBody(model.OperationRequest{Operation: "read"}).
		Post("/admin/network")
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *RestyClient) GetLANIPSettings(ctx context.Context) (model.LANIPResponse, error) {
	var res model.LANIPResponse

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "lan_ip").
		SetBody(model.OperationRequest{Operation: "read"}).
		Post("/admin/network")
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *RestyClient) GetDHCPDialSettings(ctx context.Context) (model.DHCPDialResponse, error) {
	var res model.DHCPDialResponse

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "dhcp_dial").
		SetBody(model.OperationRequest{Operation: "read"}).
		Post("/admin/network")
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *RestyClient) GetLEDPower(ctx context.Context) (model.LEDPowerResponse, error) {
	var res model.LEDPowerResponse

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "power").
		SetBody(model.OperationRequest{Operation: "read"}).
		Post("/admin/wireless")
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *RestyClient) Logout(ctx context.Context) error {
	logoutReq := model.OperationRequest{Operation: "logout"}

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "logout").
		SetBody(logoutReq).
		Post("/admin/system")
	if err != nil {
		return err
	}

	var res model.ErrorResponse
	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return err
	}

	c.decoMiddleware.SetStok("")
	c.authenticated = false

	return nil
}

func (c *RestyClient) GetSystemLog(ctx context.Context) (model.SysLogResponse, error) {
	var res model.SysLogResponse

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "log").
		SetBody(model.OperationRequest{Operation: "read"}).
		Post("/admin/syslog")
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *RestyClient) SaveLog(ctx context.Context) (model.LogExportResponse, error) {
	var res model.LogExportResponse

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "save_log").
		SetBody(model.OperationRequest{Operation: "save_log"}).
		Post("/admin/log_export")
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *RestyClient) GetIGMPSettings(ctx context.Context) (model.IGMPSettingResponse, error) {
	var res model.IGMPSettingResponse

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "igmp_setting").
		SetBody(model.OperationRequest{Operation: "read"}).
		Post("/admin/network")
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *RestyClient) GetFastXmitSettings(ctx context.Context) (model.FastXmitSettingResponse, error) {
	var res model.FastXmitSettingResponse

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "fast_xmit_setting").
		SetBody(model.OperationRequest{Operation: "read"}).
		Post("/admin/network")
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *RestyClient) GetQoSSettings(ctx context.Context) (model.FlowControlResponse, error) {
	var res model.FlowControlResponse

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "flow_control").
		SetBody(model.OperationRequest{Operation: "read"}).
		Post("/admin/network")
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(rsp.Body(), &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *RestyClient) GetFlowControlLANWAN(ctx context.Context) (model.FlowControlLANWANResponse, error) {
	var res model.FlowControlLANWANResponse

	rsp, err := c.client.
		NewRequest().
		SetContext(ctx).
		SetQueryParam("form", "flow_control_lan_wan").
		SetBody(model.OperationRequest{Operation: "read"}).
		Post("/admin/network")
	if err != nil {
		return res, err
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
		return nil, 0, errors.New("failed to get session key")
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
