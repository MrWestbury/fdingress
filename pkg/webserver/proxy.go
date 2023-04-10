package webserver

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func (server *WebServer) proxy(c *gin.Context) {
	allowed := server.checkAuthorisation(c)
	if !allowed {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	ingresses := server.scanner.Ingresses()
	var upstreamUrl string
	for _, ingress := range ingresses {
		if c.Request.Host != ingress.Host {
			continue
		}
		match := false
		switch ingress.PathType {
		case "Exact":
			if c.Request.URL.Path == ingress.Path {
				match = true
			}
		case "Prefix":
			if strings.HasPrefix(c.Request.URL.Path, ingress.Path) {
				match = true
			}
		case "ImplementationSpecific":
			re := regexp.MustCompile(ingress.Path)
			if re.MatchString(c.Request.URL.Path) {
				match = true
			}
		}
		if !match {
			continue
		}

		upstreamUrl = fmt.Sprintf("http://%s.%s:%d", ingress.ServiceName, ingress.Namespace, ingress.ServicePort)
		break
	}

	if upstreamUrl == "" {
		c.Status(http.StatusNotFound)
		return
	}

	remote, err := url.Parse(upstreamUrl)
	if err != nil {
		log.Printf("failed to parse upstreamUrl: %s\n", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = c.Request.URL.Path
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func (server *WebServer) checkAuthorisation(c *gin.Context) bool {

	requestFdId := c.GetHeader("x-azure-fdid")
	if requestFdId != "" {
		frontdoorIds := server.scanner.FrontdoorIds()
		for _, fdid := range frontdoorIds {
			if fdid == requestFdId {
				return true
			}
		}
	}

	remoteAddr := c.Request.RemoteAddr
	allowedIps := server.scanner.IpAddresses()
	for _, ipAddress := range allowedIps {
		remoteIp := strings.Split(remoteAddr, ":")[0]
		if ipAddress == remoteIp {
			return true
		}
	}

	return false
}
