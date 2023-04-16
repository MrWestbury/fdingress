package main

import (
	"log"
	"os"

	"github.com/mrwestbury/frontdoor-ingress/pkg"
	"github.com/mrwestbury/frontdoor-ingress/pkg/k8s"
	"github.com/mrwestbury/frontdoor-ingress/pkg/webserver"
)

func main() {
	envClassName := os.Getenv("FDINGRESS_CLASSNAME")
	if envClassName == "" {
		envClassName = "fdingress"
	}
	scanner := k8s.NewScanner(envClassName)
	scanner.Start()
	web := webserver.NewWebServer(scanner)

	log.Printf("Starting fdingress controller. Version %s\n", pkg.Version)
	if err := web.Run(); err != nil {
		log.Printf("Web server failed: %s\n", err)
	}
}
