package main

import (
	"context"
	"flag"
	"github.com/BGrewell/go-update"
	"github.com/BGrewell/go-update/stores/github"
	"github.com/BGrewell/system-api/configuration"
	"github.com/BGrewell/system-api/handlers"
	"github.com/BGrewell/system-api/httprouting"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"runtime"
	"time"
)

var (
	date = time.Now().Format("2006-01-02 15:04:05")
	rev = "DEBUG"
	branch = "DEBUG"
	version = "DEBUG"
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

func checkForUpdates() {

	token := "7888124ef00163d6bddc618bcc627b1a4c0d0564"
	binaryName, _ := os.Executable()

	m := &update.Manager{
		Command: binaryName,
		Store: &github.Store{
			Owner:   "BGrewell",
			Repo:    "system-api",
			Version: "",
			Token: &token,
		},
	}

	releases, err := m.LatestReleases()
	if err != nil {
		log.Infof("error getting releases: %s\n", err)
		return
	}

	if len(releases) == 0 {
		log.Info("no updates available")
		return
	}

	latest := releases[0]

	if latest.Newer(version) {
		archive := latest.FindTarball(runtime.GOOS, runtime.GOARCH)
		if archive == nil {
			log.Info("unable to find binary for this system")
			return
		}

		tarball, err := archive.DownloadSecure(token)
		if err != nil {
			log.Infof("failed to download update: %s\n", err)
			return
		}

		log.Printf("tarball: %s", tarball)
		if err := m.Install(tarball); err != nil {
			log.Infof("failed to install update: %s\n", err)
			return
		}

		log.Infof("updated to version %s\n", latest.Version)
	} else {
		log.Info("local version is newer than online version")
	}

}

func main() {

	log.Printf("Date: %s\n", date)
	log.Printf("Rev: %s\n", rev)
	log.Printf("Branch: %s\n", branch)
	log.Printf("Version: %s\n", version)
	log.Println("checking for updates")
	checkForUpdates()

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
