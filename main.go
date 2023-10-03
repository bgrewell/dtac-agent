package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/BGrewell/go-conversions"
	"github.com/gin-gonic/gin"
	. "github.com/intel-innersource/frameworks.automation.dtac.agent/common"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/configuration"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/handlers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/httprouting"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/middleware"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/plugin"
	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	date    = time.Now().Format("2006-01-02 15:04:05")
	rev     = "DEBUG"
	branch  = "DEBUG"
	version = "DEBUG"
)

type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		log.Info("Running interactively")
	} else {
		log.Info("Running as a service")
	}
	p.exit = make(chan struct{})

	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	log.Info("Stopping...")
	close(p.exit)
	return nil
}

func (p *program) run() {

	gin.SetMode(gin.ReleaseMode)

	// Default Router
	r := gin.Default()

	// Add Middleware (registration is further down after configuration is loaded)
	r.Use(middleware.LockoutMiddleware())

	// Check for custom config file location
	cfgfile := ""
	customCfgFile := os.Getenv("DTAC_CFG_LOCATION")
	if customCfgFile != "" {
		cfgfile = customCfgFile
	}

	// Load configuration
	err := configuration.Load(cfgfile)
	if err != nil {
		log.Errorf("failed to load configuration file: %v", err)
	}
	c := configuration.Config

	// General Routes
	httprouting.AddGeneralHandlers(r)

	// OS Specific Routes
	httprouting.AddOSSpecificHandlers(r)

	// Custom Routes
	httprouting.AddCustomHandlers(c, r)

	if c.Lockout.Enabled {
		middleware.RegisterLockoutHandler(r, c.Lockout.AutoUnlockTime)
	}

	// Check for updates
	//go runUpdateChecker(c)

	// Initialize internal modules
	// module.Initialize(c.Modules, r) TODO: Modules need to be ported to plugins

	// Initialize plugins
	if c.Plugins.Enabled {
		err = plugin.Initialize(c.Plugins.PluginDir, c.Plugins.Entries, r)
		if err != nil {
			log.Errorf("failed to load plugins: %s\n", err)
		}
	}

	// Setup custom 404 handler
	r.NoRoute(func(c *gin.Context) {
		WriteNotFoundResponseJSON(c)
	})

	// Before starting update the handlers Routes var
	handlers.Routes = r.Routes()

	// Setup the http(s) server
	var srvFunc func() error
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.Listener.Port),
		Handler: r,
	}

	// Setup cert manager if https is enabled
	if c.Listener.Https.Enabled {

		if c.Listener.Https.Type == TLS_TYPE_SELF_SIGNED {
			// Create default files if not specified and save to config
			if c.Listener.Https.CertFile == "" || c.Listener.Https.KeyFile == "" {
				c.Listener.Https.CertFile = "/etc/dtac/certs/tls.crt"
				c.Listener.Https.KeyFile = "/etc/dtac/certs/tls.key"
			}

			// Ensure the directories exist and are secure
			if err := os.MkdirAll(filepath.Dir(c.Listener.Https.CertFile), 0700); err != nil {
				log.Fatalf("Failed to create certificate directory: %v", err)
			}
			if err := os.MkdirAll(filepath.Dir(c.Listener.Https.KeyFile), 0700); err != nil {
				log.Fatalf("Failed to create certificate key directory: %v", err)
			}
			if _, err := os.Stat(c.Listener.Https.CertFile); os.IsNotExist(err) {
				if _, err := os.Stat(c.Listener.Https.KeyFile); os.IsNotExist(err) {
					log.Printf("[INFO] Generating self-signed certificate (%s) and key (%s)\n",
						c.Listener.Https.CertFile, c.Listener.Https.KeyFile)
					if err := GenerateSelfSignedCertKey(
						c.Listener.Https.CertFile,
						c.Listener.Https.KeyFile,
						365,
						c.Listener.Https.Domains); err != nil {
						log.Fatalf("Failed to generate self-signed certs: %v", err)
					}
				} else if err != nil {
					log.Fatalf("Failed to check key file: %v", err)
				}
			} else if err != nil {
				log.Fatalf("Failed to check cert file: %v", err)
			}
		}
	}

	// Setup serving function
	if c.Listener.Https.Enabled {
		wrapper := func() error {
			return srv.ListenAndServeTLS(c.Listener.Https.CertFile, c.Listener.Https.KeyFile)
		}
		srvFunc = wrapper
	} else {
		srvFunc = srv.ListenAndServe
	}

	// Run in a goroutine so that it won't block the graceful shutdown handling
	go func() {
		proto := "http"
		if c.Listener.Https.Enabled {
			proto = "https"
		}
		log.Printf("DTAC-Agent server is running on %s://localhost:%d\n", proto, c.Listener.Port)
		fmt.Printf("DTAC-Agent server is running on %s://localhost:%d\n", proto, c.Listener.Port)
		if err := srvFunc(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v\n", err)
		}
	}()

	<-p.exit

	// Exit has been requested give the service 5 seconds to finish its work
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("forcing server to shutdown: %v", err)
	}

	log.Info("server has exited")
}

func runUpdateChecker(c *configuration.Configuration) {
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

	//binaryName, _ := os.Executable()
	//
	//m := &update.Manager{
	//	Command: binaryName,
	//	Store: &github.Store{
	//		Owner:   "BGrewell",
	//		Repo:    "dtac-agent",
	//		Version: "",
	//		Token:   token,
	//	},
	//}
	//
	//releases, err := m.LatestReleases()
	//if err != nil {
	//	log.Infof("error getting releases: %s\n", err)
	//	return
	//}
	//
	//if len(releases) == 0 {
	//	log.Info("no updates available")
	//	return
	//}
	//
	//latest := releases[0]
	//
	//if latest.Newer(version) {
	//	archive := latest.FindTarball(runtime.GOOS, runtime.GOARCH)
	//	if archive == nil {
	//		log.Info("unable to find binary for this system")
	//		return false, errors.New("unable to find binary for this system")
	//	}
	//
	//	tarball, err := archive.DownloadSecure(*token)
	//	if err != nil {
	//		log.Infof("failed to download update: %s\n", err)
	//		return false, err
	//	}
	//
	//	log.Printf("tarball: %s", tarball)
	//	if err := m.Install(tarball); err != nil {
	//		log.Infof("failed to install update: %s\n", err)
	//		return false, err
	//	}
	//
	//	log.Infof("updated to version %s\n", latest.Version)
	//	return true, nil
	//} else {
	//	log.Info("local version is the latest version")
	//}
	return false, nil

}

func isRoot() bool {
	if runtime.GOOS == "windows" {
		// Check for admin privileges on Windows
		_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
		if err != nil {
			return false
		}
		return true
	} else {
		// Check for root privileges on UNIX
		return syscall.Geteuid() == 0
	}
}

func main() {

	if !isRoot() {
		fmt.Println("[!] This program must be run as root/admin")
		return
	}

	filename := "/var/log/dtac-agentd/dtac-agentd.log"
	if runtime.GOOS == "windows" {
		filename = "C:\\Logs\\dtac-agentd.log"
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
	fmt.Printf("=======================================================\n")
	fmt.Printf("Build Info:\n")
	fmt.Printf("\t|- Version: %s\n", version)
	fmt.Printf("\t|- Rev: %s\n", rev)
	fmt.Printf("\t|- Branch: %s\n", branch)
	fmt.Printf("\t|- Compiled Date: %s\n", date)
	fmt.Printf("=======================================================\n")

	flagInstall := flag.Bool("install", false, "Install service")
	flagUninstall := flag.Bool("uninstall", false, "Uninstall service")
	flagStart := flag.Bool("start", false, "Start service")
	flagStop := flag.Bool("stop", false, "Stop service")
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
		Name:         "dtac-agent",
		DisplayName:  "DTAC-Agent Service",
		Description:  "DTAC-Agent provides access to many system details and controls  via REST endpoints",
		Dependencies: dependencies,
		Option:       options,
	}

	p := &program{}
	s, err := service.New(p, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// errs := make(chan error, 5)
	// logger, err = s.Logger(errs)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// handle any errors that happen
	// go func() {
	// 	for {
	// 		err := <-errs
	// 		if err != nil {
	// 			log.Print(err)
	// 		}
	// 	}
	// }()

	if *flagInstall {
		if BinaryIsCorrect() {
			err = s.Install()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Service installed successfully")
			return
		} else {
			fmt.Println("Preparing to install...")
			handlers.DB.Close()
			PrepForInstall()
			return
		}
	}

	if *flagUninstall {
		if BinaryIsCorrect() {
			err = s.Uninstall()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Service uninstalled successfully")
			return
		} else {
			fmt.Printf("Current binary is not installed. Try running binary at %s\n", configuration.BINARY_NAME)
			return
		}

	}

	if *flagStart {
		if BinaryIsCorrect() {
			err = s.Start()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Service started")
			return
		} else {
			fmt.Printf("Current binary is not installed. Try running binary at %s\n", configuration.BINARY_NAME)
			return
		}

	}

	if *flagStop {
		if BinaryIsCorrect() {
			err = s.Stop()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Service stopped")
			return
		} else {
			fmt.Printf("Current binary is not installed. Try running binary at %s\n", configuration.BINARY_NAME)
			return
		}
	}

	// By default, let's run the service logic.
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
