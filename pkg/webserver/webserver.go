package webserver

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrwestbury/frontdoor-ingress/pkg"
	"github.com/mrwestbury/frontdoor-ingress/pkg/k8s"
)

type WebServer struct {
	engine  *gin.Engine
	scanner *k8s.Scanner
}

func NewWebServer(scanner *k8s.Scanner) *WebServer {
	server := &WebServer{
		engine:  gin.Default(),
		scanner: scanner,
	}

	server.engine.GET("/health", server.health)
	server.engine.GET("/ingresses", server.ingresses)

	server.engine.NoRoute(server.proxy)

	return server
}

func (server *WebServer) Run() error {
	return server.engine.Run(":8090")
}

func (server *WebServer) health(c *gin.Context) {
	log.Printf("Healthcheck on %s\n", c.Request.Host)
	healthCheck := map[string]interface{}{
		"ingresses": len(server.scanner.Ingresses()),
	}
	c.IndentedJSON(http.StatusOK, healthCheck)
}

func (server *WebServer) ingresses(c *gin.Context) {
	result := map[string]interface{}{
		"app_version": pkg.Version,
		"ingresses":   server.scanner.Ingresses(),
		"frontdoors":  server.scanner.FrontdoorIds(),
		"ipAddresses": server.scanner.IpAddresses(),
	}

	c.IndentedJSON(http.StatusOK, result)
}
