package main

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sosedoff/pgweb/pkg/cli"
	"github.com/sosedoff/pgweb/pkg/command"
)

func main() {

	opts, err := command.ParseOptions(os.Args)
	if err != nil {
		log.Fatalf("Failed to parse options: %v", err)
	}
	target, err := url.Parse(fmt.Sprintf("http://%s:%d", opts.Host, opts.Port))
	if err != nil {
		log.Fatalf("Invalid target URL: %v", err)
	}

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Initialize Gin router
	router := gin.Default()

	// Forward all requests
	router.Any("/*path", func(c *gin.Context) {
		c.Request.URL.Scheme = target.Scheme
		c.Request.URL.Host = target.Host
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// Load TLS certificate & key
	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")

	httpsHost := os.Getenv("HTTPS_HOST")
	httpsPort := os.Getenv("HTTPS_PORT")

	// Start pgweb CLI
	go cli.Run()

	// Start HTTPS server
	log.Printf("Starting TLS server on :%s...\n", httpsPort)
	if err := router.RunTLS(fmt.Sprintf("%s:%s", httpsHost, httpsPort), certFile, keyFile); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
