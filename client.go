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
	GetClients(ctx context.Context) (model.ClientListResponse, error)
	GetPerformance(ctx context.Context) (model.PerformanceResponse, error)
	RebootDevice(ctx context.Context, macAddrs []string) error
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
