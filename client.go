package tplink_deco_client

import (
	"context"
	"encoding/json"
	"errors"
	"net"

	"github.com/robotjoosen/go-tplink-deco-client/internal/httpclient"
	"github.com/robotjoosen/go-tplink-deco-client/internal/model"
	exportModel "github.com/robotjoosen/go-tplink-deco-client/pkg/model"
)

const username = "admin"

type ClientAware interface {
	Login(ctx context.Context, username, password string) (model.LoginResponse, error)
	GetDevices(ctx context.Context) (model.DeviceListResponse, error)
	GetClients(ctx context.Context) (model.ClientListResponse, error)
	GetPerformance(ctx context.Context) (model.PerformanceResponse, error)
	RebootDevice(ctx context.Context, macAddrs []string) error
	Custom(ctx context.Context, path string, form string, body []byte) (interface{}, error)
	GetWANIPv4Info(ctx context.Context) (model.WANIPv4Response, error)
	GetWiFiSettings(ctx context.Context) (model.WiFiResponse, error)
	SetWiFiSettings(ctx context.Context, settings map[string]interface{}) error
}

type Client struct {
	client        ClientAware
	authenticated bool
}

func New(ip string) *Client {
	return &Client{
		client: httpclient.NewClient(ip),
		//client: restyclient.NewClient(ip),
	}
}

func (c *Client) Authenticate(ctx context.Context, password string) (*Client, error) {
	if _, err := c.client.Login(ctx, username, password); err != nil {
		return nil, err
	}

	c.authenticated = true

	return c, nil
}

func (c *Client) GetDevices(ctx context.Context) (exportModel.Devices, error) {
	if !c.authenticated {
		return exportModel.Devices{}, errors.New("not authenticated")
	}

	res, err := c.client.GetDevices(ctx)
	if err != nil {
		return nil, err
	}

	output := make(exportModel.Devices, 0, len(res.Result.DeviceList))
	for _, item := range res.Result.DeviceList {
		macAddress, err := net.ParseMAC(item.Mac)
		if err != nil {
			return nil, err
		}

		output = append(output, exportModel.Device{
			ID:         item.DeviceID,
			Name:       item.Nickname,
			IPAddress:  net.ParseIP(item.DeviceIP),
			MACAddress: macAddress,
		})
	}

	return output, nil
}

func (c *Client) GetClients(ctx context.Context) (exportModel.Clients, error) {
	if !c.authenticated {
		return exportModel.Clients{}, errors.New("not authenticated")
	}

	res, err := c.client.GetClients(ctx)
	if err != nil {
		return nil, err
	}

	output := make(exportModel.Clients, 0, len(res.Result.ClientList))
	for _, item := range res.Result.ClientList {
		macAddress, err := net.ParseMAC(item.Mac)
		if err != nil {
			return nil, err
		}

		output = append(output, exportModel.Client{
			Name:           item.Name,
			OwnerId:        item.OwnerID,
			Online:         item.Online,
			WireType:       item.WireType,
			ClientType:     item.ClientType,
			ConnectionType: item.ConnectionType,
			AccessHost:     item.AccessHost,
			SpaceID:        item.SpaceID,
			ClientMesh:     item.ClientMesh,
			EnablePriority: item.EnablePriority,
			RemainTime:     item.RemainTime,
			UpSpeed:        item.UpSpeed,
			DownSpeed:      item.DownSpeed,
			IPAddress:      net.ParseIP(item.IP),
			MACAddress:     macAddress,
		})
	}

	return output, nil
}

func (c *Client) GetPerformance(ctx context.Context) (exportModel.Performance, error) {
	if !c.authenticated {
		return exportModel.Performance{}, errors.New("not authenticated")
	}

	res, err := c.client.GetPerformance(ctx)
	if err != nil {
		return exportModel.Performance{}, err
	}

	return exportModel.Performance{
		CPU: res.Result.CPU,
		MEM: res.Result.MEM,
	}, nil
}

func (c *Client) RebootDevice(ctx context.Context, macAddrs ...string) error {
	if !c.authenticated {
		return errors.New("not authenticated")
	}

	return c.client.RebootDevice(ctx, macAddrs)
}

func (c *Client) Custom(ctx context.Context, path string, form string, operation string, params map[string]interface{}) (interface{}, error) {
	if !c.authenticated {
		return nil, errors.New("not authenticated")
	}

	req := model.OperationRequest{
		Operation: operation,
		Params:    params,
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return c.client.Custom(ctx, path, form, jsonBody)
}

func (c *Client) GetWANIPv4Info(ctx context.Context) (exportModel.WANInfo, error) {
	if !c.authenticated {
		return exportModel.WANInfo{}, errors.New("not authenticated")
	}

	res, err := c.client.GetWANIPv4Info(ctx)
	if err != nil {
		return exportModel.WANInfo{}, err
	}

	wanMAC, err := net.ParseMAC(res.WAN.IPInfo.MAC)
	if err != nil {
		return exportModel.WANInfo{}, err
	}

	lanMAC, err := net.ParseMAC(res.LAN.IPInfo.MAC)
	if err != nil {
		return exportModel.WANInfo{}, err
	}

	_, wanSubnet, err := net.ParseCIDR(res.WAN.IPInfo.Mask)
	if err != nil {
		return exportModel.WANInfo{}, err
	}

	_, lanSubnet, err := net.ParseCIDR(res.LAN.IPInfo.Mask)
	if err != nil {
		return exportModel.WANInfo{}, err
	}

	return exportModel.WANInfo{
		WANMAC:     wanMAC,
		WANIP:      net.ParseIP(res.WAN.IPInfo.IP),
		WANGateway: net.ParseIP(res.WAN.IPInfo.Gateway),
		WANDNS1:    net.ParseIP(res.WAN.IPInfo.DNS1),
		WANDNS2:    net.ParseIP(res.WAN.IPInfo.DNS2),
		WANSubnet:  wanSubnet.Mask,
		WANType:    res.WAN.DialType,
		LinkStatus: res.WAN.LinkStatus,
		LANMAC:     lanMAC,
		LANIP:      net.ParseIP(res.LAN.IPInfo.IP),
		LANSubnet:  lanSubnet.Mask,
	}, nil
}

func (c *Client) GetWiFiSettings(ctx context.Context) (exportModel.WiFiSettings, error) {
	if !c.authenticated {
		return exportModel.WiFiSettings{}, errors.New("not authenticated")
	}

	res, err := c.client.GetWiFiSettings(ctx)
	if err != nil {
		return exportModel.WiFiSettings{}, err
	}

	return exportModel.WiFiSettings{
		Band24: exportModel.WiFiBandSettings{
			Host:  toNetworkSettings(res.Band24.Host),
			Guest: toNetworkSettings(res.Band24.Guest),
			IoT:   toNetworkSettings(res.Band24.IoT),
		},
		Band5: exportModel.WiFiBandSettings{
			Host:  toNetworkSettings(res.Band5.Host),
			Guest: toNetworkSettings(res.Band5.Guest),
			IoT:   toNetworkSettings(res.Band5.IoT),
		},
		Band6: exportModel.WiFiBandSettings{
			Host:  toNetworkSettings(res.Band6.Host),
			Guest: toNetworkSettings(res.Band6.Guest),
			IoT:   toNetworkSettings(res.Band6.IoT),
		},
	}, nil
}

func (c *Client) SetWiFiSettings(ctx context.Context, settings map[string]interface{}) error {
	if !c.authenticated {
		return errors.New("not authenticated")
	}

	return c.client.SetWiFiSettings(ctx, settings)
}

func toNetworkSettings(n model.WiFiNetwork) exportModel.WiFiNetworkSettings {
	return exportModel.WiFiNetworkSettings{
		Enabled:  n.Enable,
		SSID:     n.SSID,
		Password: n.Key,
	}
}
