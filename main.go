package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/BGrewell/go-conversions"
	"github.com/BGrewell/go-update"
	"github.com/BGrewell/go-update/stores/github"
	. "github.com/BGrewell/system-api/common"
	"github.com/BGrewell/system-api/configuration"
	"github.com/BGrewell/system-api/handlers"
	"github.com/BGrewell/system-api/httprouting"
	"github.com/BGrewell/system-api/middleware"
	"github.com/BGrewell/system-api/plugin/core"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
	"runtime"
	"time"
)

var (
	date    = time.Now().Format("2006-01-02 15:04:05")
	rev     = "DEBUG"
	branch  = "DEBUG"
	version = "DEBUG"
	logger  service.Logger
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

	// Add Middleware (registration is further down after configuration is loaded)
	r.Use(middleware.LockoutMiddleware())

	// General Routes
	httprouting.AddGeneralHandlers(r)

	// OS Specific Routes
	httprouting.AddOSSpecificHandlers(r)

	// Load Configuration and Custom Routes
	cfgfile := "/etc/system-api/config.yaml"
	if runtime.GOOS == "windows" {
		cfgfile = "c:\\Program Files\\Intel\\System-Api\\config.yaml"
	}

	// Check for custom config file location
	customCfgFile := os.Getenv("SYSTEMAPI_CFG_LOCATION")
	if customCfgFile != "" {
		cfgfile = customCfgFile
	}

	c, err := configuration.Load(cfgfile)
	if err != nil {
		logger.Errorf("failed to load configuration file: %v", err)
	} else {
		httprouting.AddCustomHandlers(c, r)
	}

	middleware.RegisterLockoutHandler(r, c.LockoutTime)

	// Check for updates
	//go runUpdateChecker(c)
	// TODO: Need to build the mapping for plugin to client for calls from REST API - this is just to prevent errors during dev
	cplugs := make(map[string]*core.Client)

	// TODO: Deploy any plugins
	for _, p := range c.Plugins.ActivePlugins {
		for name, cfg := range p {
			pluginClient, err := core.NewClient(name, cfg)
			if err != nil {
				log.Errorf("failed to load plugin: %v", err)
				continue
			}
			cplugs[name] = pluginClient
			log.Infof("loaded plugin: %v", name)
		}
	}

	// Setup custom 404 handler
	r.NoRoute(func(c *gin.Context) {
		WriteNotFoundResponseJSON(c)
	})

	// Before starting update the handlers Routes var
	handlers.Routes = r.Routes()

	log.Println("system-api server is running http://localhost:8080")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.ListenPort),
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

func runUpdateChecker(c *configuration.Config) {
	//todo: make run periodic checks
	sleepTime, err := conversions.ConvertStringTimeToNanoseconds(c.Updater.Interval)
	if err != nil {
		log.Infof("failed to convert update interval: %s\n", err)
	}
	errorTime, err := conversions.ConvertStringTimeToNanoseconds(c.Updater.ErrorFallback)
	if err != nil {
		log.Infof("failed to convert update fallback interval: %s\n", err)
	}
	for {
		log.Info("checking for updates...")
		updated, err := checkForUpdates(&c.Updater.Token)
		//todo: make restart on update
		if updated && c.Updater.RestartOnUpdate {
			log.Println("application updated. need to restart")
		} else if updated && !c.Updater.RestartOnUpdate {
			log.Println("application updated but auto-restart is off. updater is now disabled")
			return
		}

		if c.Updater.Mode != "auto" {
			return
		}

		t := sleepTime
		if err != nil {
			t = errorTime
		}

		time.Sleep(time.Duration(t) * time.Nanosecond)
	}

}

func checkForUpdates(token *string) (applied bool, err error) {

	binaryName, _ := os.Executable()

	m := &update.Manager{
		Command: binaryName,
		Store: &github.Store{
			Owner:   "BGrewell",
			Repo:    "system-api",
			Version: "",
			Token:   token,
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
			return false, errors.New("unable to find binary for this system")
		}

		tarball, err := archive.DownloadSecure(*token)
		if err != nil {
			log.Infof("failed to download update: %s\n", err)
			return false, err
		}

		log.Printf("tarball: %s", tarball)
		if err := m.Install(tarball); err != nil {
			log.Infof("failed to install update: %s\n", err)
			return false, err
		}

		log.Infof("updated to version %s\n", latest.Version)
		return true, nil
	} else {
		log.Info("local version is the latest version")
	}
	return false, nil
}

func main() {

	filename := "/var/log/system-apid/system-apid.log"
	if runtime.GOOS == "windows" {
		filename = "C:\\Logs\\system-apid.log"
	}

	log.SetOutput(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	})
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	log.SetReportCaller(true)
	log.ParseLevel("debug")
	log.Printf("Date: %s", date)
	log.Printf("Rev: %s", rev)
	log.Printf("Branch: %s", branch)
	log.Printf("Version: %s", version)

	svcFlag := flag.String("service", "", "control the service")
	flag.Parse()

	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SKIGKILL"
	var dependencies []string
	if runtime.GOOS != "windows" {
		dependencies = []string{
			"Requires=network.target",
			"After=network-online.target syslog.target",
		}
	}
	svcConfig := &service.Config{
		Name:         "system-api.service",
		DisplayName:  "System-API Service",
		Description:  "System-API provides access to many system details via REST endpoints",
		Dependencies: dependencies,
		Option:       options,
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
