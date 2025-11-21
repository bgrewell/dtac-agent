package webhost

import (
	"context"
	"fmt"
	"net"
	"os"
)

// DefaultWebhostHost manages the lifecycle of a module and the control plane (gRPC) when DTAC owns it.
type DefaultWebhostHost struct {
	Module WebhostModule
}

func NewHost(m WebhostModule) (*DefaultWebhostHost, error) {
	if m == nil {
		return nil, fmt.Errorf("module cannot be nil")
	}
	return &DefaultWebhostHost{Module: m}, nil
}

func (h *DefaultWebhostHost) Serve() error {
	if os.Getenv("DTAC_WEBHOSTS") != "" {
		return h.serveManaged()
	}
	return h.serveStandalone()
}

func (h *DefaultWebhostHost) serveStandalone() error {
	cfg := &InitConfig{
		ModuleName: h.Module.Name(),
		Listen: ListenConfig{
			Addr: "127.0.0.1",
			Port: 45000,
		},
	}
	if err := h.Module.OnInit(cfg); err != nil {
		return err
	}
	ctx := context.Background()
	return h.Module.Run(ctx)
}

func (h *DefaultWebhostHost) serveManaged() error {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	port := l.Addr().(*net.TCPAddr).Port

	fmt.Printf("CONNECT{{%s:%s:%s:%s:%s:%d:%s:[%s]}}\n",
		h.Module.Name(),
		"",
		"webhost_grpc",
		"tcp",
		"127.0.0.1",
		port,
		"webhost_api_1.0",
		"")

	_ = l.Close()
	return fmt.Errorf("managed serve not yet implemented in skeleton")
}
