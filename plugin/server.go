package plugin

import (
	"context"
	"fmt"
	api "github.com/BGrewell/system-api/plugin/go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

var (
	once sync.Once
	instance *server
)

type server struct {
	api.UnimplementedPluginServer
	port string
	running bool
	engine *gin.Engine
	plugins map[string]*Shim
}

func NewServer(port int) *server {
	once.Do(func() {
		instance = &server{port: fmt.Sprintf(":%d", port)}
	})

	return instance
}

func (s *server) Serve(r *gin.Engine) error {
	log.WithFields(log.Fields{
		"proto": "tcp",
		"port": s.port,
	}).Info("setting up plugin server")

	s.engine = r

	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}
	log.Debug("plugin server listener created")

	gs := grpc.NewServer()
	log.Debug("setup new grpc server")

	api.RegisterPluginServer(gs, instance)
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

func (s *server) Stop() {
	s.running = false
	time.Sleep(200 * time.Millisecond)
}


func (s *server) ControlChannel(channelServer api.Plugin_ControlChannelServer) error {
	// TODO: Take all of the routes and the base and combine them and register them with the web server. Any
	//		 calls to those routes will be passed back through the command channel to the plugin along with any
	//		 body, headers and url so that the plugin can do what it likes and return the result

	if _, ok := s.plugins[request.PluginName]; ok {
		panic("plugin exists, we need to return a real error")
	}
	panic("implement me")
}