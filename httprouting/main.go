package httprouting

import (
	"github.com/BGrewell/system-api/handlers"
	"github.com/gin-gonic/gin"
)

func AddGeneralHandlers(r *gin.Engine) {
	// GET Routes - Retrieve information
	r.GET("/", handlers.HomeHandler)                                              // postman | implemented
	r.GET("/os", handlers.GetOSHandler)                                           // postman | implemented
	r.GET("/cpu", handlers.GetCPUHandler)                                         // postman | implemented
	r.GET("/cpu/usage", handlers.GetCPUUsageHandler)								// postman | implemented
	r.GET("/memory", handlers.GetMemoryHandler)                                   // postman | implemented
	r.GET("/disk", handlers.GetDiskHandler)                                       // postman | implemented
	r.GET("/docker", handlers.GetDockerHandler)                                   // postman | implemented
	r.GET("/process", handlers.GetProcessesHandler)                               // postman | implemented
	r.GET("/process/:pid", handlers.GetProcessHandler)                            // postman | implemented
	r.GET("/docker/containers", handlers.GetDockerContainersHandler)              // postman
	r.GET("/docker/container/:id", handlers.GetDockerContainerHandler)            // postman
	r.GET("/docker/images", handlers.GetDockerImagesHandler)                      // postman
	r.GET("/docker/image/:id", handlers.GetDockerImageHandler)                    // postman
	r.GET("/docker/networks", handlers.GetDockerNetworksHandler)                  // postman
	r.GET("/docker/network/:id", handlers.GetDockerNetworkHandler)                // postman
	r.GET("/host", handlers.GetHostHandler)                                       // postman | implemented
	r.GET("/network", handlers.GetNetworkHandler)                                 // postman | implemented
	r.GET("/endpoints", handlers.GetEndpointsHandler)                             // postman | implemented
	r.GET("/internal/logs", handlers.GetLogsHandler)                              // postman | implemented
	r.GET("/internal/logs/stream", handlers.GetLogsStreamHandler)                 // postman
	r.GET("/ping", handlers.GetPingHandler)                                       // postman | implemented
	r.GET("/ping/udp/:target", handlers.SendTimedUdpPingHandler)                  // postman | implemented
	r.GET("/ping/udp/results/:id", handlers.GetUdpPingWorkerHandler)              // postman | implemented
	r.GET("/ping/tcp/results/:id", handlers.GetTcpPingWorkerHandler)              // postman | implemented
	r.GET("/ping/tcp/:target", handlers.SendTimedTcpPingHandler)                  // postman | implemented
	r.GET("/ping/reflectors", handlers.GetReflectors)                             // postman | implemented
	r.GET("/network/interfaces", handlers.GetInterfacesHandler)                   // postman | implemented
	r.GET("/network/interfaces/names", handlers.GetInterfaceNamesHandler)         // postman | implemented
	r.GET("/network/interface/:name", handlers.GetInterfaceByNameHandler)         // postman | implemented
	r.GET("/network/routes", handlers.GetRoutesHandler)                           // postman | implemented
	r.GET("/internal/tests/secret", handlers.SecretTestHandler)                   // postman | implemented
	r.GET("/iperf/client/results/:id", handlers.GetIperfClientTestResultsHandler) // postman | implemented
	r.GET("/iperf/client/live/:id", handlers.GetIperfClientTestLiveHandler)       // postman | implemented
	r.GET("/iperf/server/results/:id", handlers.GetIperfServerTestResultsHandler) // postman | implemented
	r.GET("/internal/tests/unimplemented", handlers.UnimplementedHandler)         // postman | implemented

	// PUT Routes - Update information
	r.PUT("/network/route", handlers.UpdateRouteHandler)                // postman
	r.PUT("/system/reboot", handlers.SystemRebootHandler)               // postman
	r.PUT("/system/shutdown", handlers.SystemShutdownHandler)           // postman
	r.PUT("/systemapi/restart", handlers.SystemApiRestartHandler)       // postman
	r.PUT("/systemapi/restart/:time", handlers.SystemApiRestartHandler) // postman

	// Delete Routes - Remove information
	r.DELETE("/network/route", handlers.DeleteRouteHandler)                   // postman
	r.DELETE("/ping/udp/:id", handlers.DeleteUdpPingWorkerHandler)            // postman
	r.DELETE("/ping/tcp/:id", handlers.DeleteTcpPingWorkerHandler)            // postman
	r.DELETE("/ping/reset", handlers.DeleteResetAllPingWorkersHandler)        // postman
	r.DELETE("/iperf/client/stop/:id", handlers.DeleteIperfClientTestHandler) // postman
	r.DELETE("/iperf/server/stop/:id", handlers.DeleteIperfServerTestHandler) // postman
	r.DELETE("/iperf/reset", handlers.DeleteIperfResetHandler)                // postman

	// POST Routes - Create information
	r.POST("/login", handlers.LoginHandler)                                    // postman
	r.POST("/ping/udp/:target", handlers.CreateUdpPingWorkerHandler)           // postman
	r.POST("/ping/tcp/:target", handlers.CreateTcpPingWorkerHandler)           // postman
	r.POST("/network/route", handlers.CreateRouteHandler)                      // postman
	r.POST("/iperf/client/start/:host", handlers.CreateIperfClientTestHandler) // postman
	r.POST("/iperf/server/start", handlers.CreateIperfServerTestHandler)       // postman
	r.POST("/iperf/server/start/:bind", handlers.CreateIperfServerTestHandler) // postman

	// ===== Unified Endpoints - Provide consistent input/output across OS's =====
	// Unified Firewall
	r.POST("/u/network/firewall/rule", handlers.CreateFirewallRuleUniversalHandler)       // postman
	r.GET("/u/network/firewall/rule/:id", handlers.GetFirewallRuleUniversalHandler)       // postman
	r.GET("/u/network/firewall/rules", handlers.GetFirewallRulesUniversalHandler)         // postman
	r.PUT("/u/network/firewall/rule/:id", handlers.UpdateFirewallRuleUniversalHandler)    // postman
	r.DELETE("/u/network/firewall/rule/:id", handlers.DeleteFirewallRuleUniversalHandler) // postman

	// Unified Qos
	r.POST("/u/network/qos/rule", handlers.CreateQosRuleUniversalHandler)       // postman
	r.GET("/u/network/qos/rule/:id", handlers.GetQosRuleUniversalHandler)       // postman
	r.GET("/u/network/qos/rules", handlers.GetQosRulesUniversalHandler)         // postman
	r.PUT("/u/network/qos/rule/:id", handlers.UpdateQosRuleUniversalHandler)    // postman
	r.DELETE("/u/network/qos/rule/:id", handlers.DeleteQosRuleUniversalHandler) // postman

	// Unified Route Info
	r.POST("/u/network/route/rule", handlers.CreateRouteRuleUniversalHandler)       // postman
	r.GET("/u/network/route/rule/:id", handlers.GetRouteRuleUniversalHandler)       // postman
	r.GET("/u/network/route/rules", handlers.GetRouteRulesUniversalHandler)         // postman
	r.PUT("/u/network/route/rule/:id", handlers.UpdateRouteRuleUniversalHandler)    // postman
	r.DELETE("/u/network/route/rule/:id", handlers.DeleteRouteRuleUniversalHandler) // postman

	// Unified Interface Info
	r.GET("/u/network/interfaces", handlers.GetInterfacesUniversalHandler)   // postman
	r.GET("/u/network/interface/:id", handlers.GetInterfaceUniversalHandler) // postman

	// TODO: Unified CPU Info

	// TODO: Unified Memory Info

	// TODO: Unified Disk Info
}
