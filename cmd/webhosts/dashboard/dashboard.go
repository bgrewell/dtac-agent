package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"

	"github.com/bgrewell/dtac-agent/pkg/webhost"
)

//go:embed static/*
var staticFS embed.FS

type DashboardModule struct {
	cfg *webhost.InitConfig
}

func NewDashboardModule() *DashboardModule { return &DashboardModule{} }
func (d *DashboardModule) Name() string    { return "dashboard" }
func (d *DashboardModule) OnInit(cfg *webhost.InitConfig) error {
	d.cfg = cfg
	return nil
}
func (d *DashboardModule) Run(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", d.cfg.Listen.Addr, d.cfg.Listen.Port)
	mux := http.NewServeMux()
	fs := http.FS(staticFS)
	mux.Handle("/", http.FileServer(fs))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	srv := &http.Server{Addr: addr, Handler: mux}
	go func() {
		<-ctx.Done()
		srv.Close()
	}()
	return srv.ListenAndServe()
}
