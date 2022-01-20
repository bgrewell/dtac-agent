package main

import (
	"context"
	"github.com/BGrewell/system-api/plugin/core"
	api "github.com/BGrewell/system-api/plugin/go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"time"
)

type HelloServer struct {
	api.UnimplementedPluginServer
	port    int
	running bool
}

func HelloHandler(c *gin.Context) {
	// Update Routes
	start := time.Now()
	core.WriteResponseJSON(c, time.Since(start), "Hello Plugin!")
}

func main() {

	r := gin.Default()
	r.GET("/hello", HelloHandler)
	s := &HelloServer{
		port:    12345,
		running: false,
	}
	err := s.Serve(r)
	if err != nil {
		panic(err)
	}
}

func (s *HelloServer) Serve(r *gin.Engine) error {
	log.WithFields(log.Fields{
		"proto": "tcp",
		"port":  s.port,
	}).Info("setting up plugin server")

	listener, err := net.Listen("tcp", string(s.port))
	if err != nil {
		return err
	}
	log.Debug("plugin server listener created")

	gs := grpc.NewServer()
	log.Debug("setup new grpc server")

	api.RegisterPluginServer(gs, s)
	log.Debug("registered plugin control server")

	go func() {
		err = gs.Serve(listener)
		log.WithFields(log.Fields{
			"error": err,
		}).Debug("error serving plugin control server")
	}()
	log.Debug("finished starting plugin control server")

	go func() {
		for s.running {
			time.Sleep(100 * time.Millisecond)
		}
		gs.Stop()
		log.Debug("stopped plugin control server")
	}()

	return nil
}

func (s *HelloServer) Stop() {
	s.running = false
	time.Sleep(200 * time.Millisecond)
}

func (*HelloServer) Register(ctx context.Context, req *api.RegisterRequest) (resp *api.RegisterResponse, err error) {
	req.
}
func (*HelloServer) Heartbeat(ctx context.Context, req *api.HeartbeatRequest) (resp *api.HeartbeatResponse, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}
func (*HelloServer) Execute(ctx context.Context, req *api.ExecutionRequest) (resp *api.ExecutionResponse, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method Execute not implemented")
}
