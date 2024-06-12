package tplink_deco_client_test

import (
	"context"
	"log/slog"

	decoClient "github.com/robotjoosen/go-tplink-deco-client"
)

func ExampleClient() {
	ctx := context.Background()

	client, err := decoClient.
		New("192.168.2.1").
		Authenticate(ctx, "xbt7v)p/zP5)$hx")
	if err != nil {
		slog.Error("auth", slog.String("error", err.Error()))

		return
	}

	devices, err := client.GetDevices(ctx)
	if err != nil {
		slog.Error("devices", slog.String("error", err.Error()))

		return
	}

	clients, err := client.GetClients(ctx)
	if err != nil {
		slog.Error("clients", slog.String("error", err.Error()))

		return
	}

	slog.Info("client", slog.Any("devices", devices), slog.Any("clients", clients))

	// output:
}
