package plugin

//TODO: THIS IS NO LONGER USED, SYSTEM-API IS THE CLIENT IN THE PLUGIN ARCHITECTURE THIS CAN BE REMOVED NOW

import (
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
	once     sync.Once
	instance *server
)

type server struct {
	api.UnimplementedPluginServer
	port          string
	running       bool
	engine        *gin.Engine
	plugins       map[string]*Shim
	routeHandlers map[string]string
	lock          *sync.Mutex
}

func NewServer(port int) *server {
	once.Do(func() {
		instance = &server{
			port:          fmt.Sprintf(":%d", port),
			lock:          &sync.Mutex{},
			plugins:       make(map[string]*Shim),
			routeHandlers: make(map[string]string),
		}
	})

	return instance
}

func (s *server) Serve(r *gin.Engine) error {
	log.WithFields(log.Fields{
		"proto": "tcp",
		"port":  s.port,
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
