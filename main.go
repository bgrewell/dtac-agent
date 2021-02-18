package main

import (
	"context"
	"flag"
	"github.com/BGrewell/system-api/configuration"
	"github.com/BGrewell/system-api/handlers"
	"github.com/BGrewell/system-api/httprouting"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	logger service.Logger
)

type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running interactively")
	} else {
		logger.Info("Running as a service")
	}
	p.exit = make(chan struct{})

	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	logger.Info("Stopping...")
	close(p.exit)
	return nil
}

func (p *program) run() {
	// Default Router
	r := gin.Default()

	// General Routes
	httprouting.AddGeneralHandlers(r)

	// OS Specific Routes
	httprouting.AddOSSpecificHandlers(r)

	// Custom Configuration Routes
	c, err := configuration.Load("support/config/config.yaml")
	if err != nil {
		logger.Errorf("failed to load configuration file: %v", err)
	} else {
		httprouting.AddCustomHandlers(c, r)
	}


	// Before starting update the handlers Routes var
	handlers.Routes = r.Routes()

	log.Println("system-api server is running http://localhost:8080")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	// Run in a goroutine so that it won't block the graceful shutdown handling
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server: %v\n", err)
		}
	}()

	<-p.exit

	// Exit has been requested give the service 5 seconds to finish its work
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("forcing server to shutdown: %v", err)
	}

	logger.Info("server has exited")
}

func main() {

	svcFlag := flag.String("service", "", "control the service")
	flag.Parse()

	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SKIGKILL"
	svcConfig := &service.Config{
		Name:        "system-api.service",
		DisplayName: "System-API Service",
		Description: "System-API provides access to many system details via REST endpoints",
		Dependencies: []string{
			"Requires=network.target",
			"After=network-online.target syslog.target",
		},
		Option: options,
	}

	p := &program{}
	s, err := service.New(p, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	// handle any errors that happen
	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
