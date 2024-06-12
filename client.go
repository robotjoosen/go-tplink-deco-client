package tplink_deco_client

import (
	"context"
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

	output := make([]exportModel.Device, 0, len(res.Result.DeviceList))
	for _, device := range res.Result.DeviceList {
		macAddress, err := net.ParseMAC(device.Mac)
		if err != nil {
			return nil, err
		}

		output = append(output, exportModel.Device{
			ID:         device.DeviceId,
			Name:       device.Nickname,
			IPAddress:  net.ParseIP(device.DeviceIp),
			MACAddress: macAddress,
		})
	}

	return output, nil
}
