package main

import (
	"log"

	"github.com/mrwestbury/frontdoor-ingress/pkg/k8s"
	"github.com/mrwestbury/frontdoor-ingress/pkg/webserver"
)

func main() {
	scanner := k8s.NewScanner()
	scanner.Start()
	web := webserver.NewWebServer(scanner)

	if err := web.Run(); err != nil {
		log.Printf("Web server failed: %s\n", err)
	}
}
