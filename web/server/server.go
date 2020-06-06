// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package server contains a web server that can serve an instance of
// the dataaccess.Repository interface via HTTP and HTTPs.
package server

import (
	"github.com/elWyatt/allmark/common/config"
	"github.com/elWyatt/allmark/common/logger"
	"github.com/elWyatt/allmark/dataaccess"
	"github.com/elWyatt/allmark/services/converter/markdowntohtml"
	"github.com/elWyatt/allmark/services/converter/markdowntohtml/imageprovider"
	"github.com/elWyatt/allmark/services/parser"
	"github.com/elWyatt/allmark/services/thumbnail"
	"github.com/elWyatt/allmark/web/handlers"
	"github.com/elWyatt/allmark/web/header"
	"github.com/elWyatt/allmark/web/orchestrator"
	"github.com/elWyatt/allmark/web/view/templates"
	"github.com/elWyatt/allmark/web/webpaths"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/skratchdot/open-golang/open"
	"net"
	"net/http"
	"strings"
)

// New creates a new Server instance for the given repository.
func New(logger logger.Logger, config config.Config, repository dataaccess.Repository, parser parser.Parser, thumbnailIndex *thumbnail.Index) (*Server, error) {

	patherFactory := webpaths.NewFactory(logger, repository)
	webPathProvider := webpaths.NewWebPathProvider(patherFactory, handlers.BasePath, handlers.TagPathPrefix)

	// image provider
	imageProvider := imageprovider.NewImageProvider(webPathProvider.AbsolutePather("/"), thumbnailIndex)

	// converter
	converter := markdowntohtml.New(logger, imageProvider)

	orchestratorFactory := orchestrator.NewFactory(logger, config, repository, parser, converter, webPathProvider)
	reindexInterval := config.Indexing.IntervalInSeconds
	headerWriterFactory := header.NewHeaderWriterFactory(reindexInterval)
	templateProvider := templates.NewProvider(config.TemplatesFolder())
	requestHandlers := handlers.GetBaseHandlers(logger, config, templateProvider, *orchestratorFactory, headerWriterFactory)

	return &Server{
		logger: logger,
		config: config,

		headerWriterFactory: headerWriterFactory,
		requestHandlers:     requestHandlers,
	}, nil

}

// Server represents a web server instance for a given repository.
type Server struct {
	logger logger.Logger
	config config.Config

	headerWriterFactory header.WriterFactory

	requestHandlers handlers.HandlerList
}

// Start starts the current web server.
func (server *Server) Start() chan error {

	result := make(chan error)

	standardRequestRouter := server.getStandardRequestRouter()

	// bindings
	httpEndpoint, httpEnabled := server.httpEndpoint()
	httpsEndpoint, httpsEnabled := server.httpsEndpoint()

	// abort if no tcp bindings are configured
	if len(httpEndpoint.Bindings()) == 0 && len(httpsEndpoint.Bindings()) == 0 {
		result <- fmt.Errorf("No TCP bindings configured")
		return result
	}

	uniqueURLs := make(map[string]string)

	// http
	if httpEnabled {

		for _, tcpBinding := range httpEndpoint.Bindings() {

			tcpBinding.AssignFreePort()

			tcpAddr := tcpBinding.GetTCPAddress()
			address := tcpAddr.String()

			// start listening
			go func() {
				server.logger.Info("HTTP Endpoint: %s", address)

				if httpEndpoint.ForceHTTPS() {

					// Redirect HTTP → HTTPS
					redirectTarget := httpsEndpoint.DefaultURL()
					httpsRedirectRouter := server.getRedirectRouter(redirectTarget, standardRequestRouter)

					if err := http.ListenAndServe(address, httpsRedirectRouter); err != nil {
						result <- fmt.Errorf("Server failed with error: %v", err)
					} else {
						result <- nil
					}

				} else {

					// Standard HTTP Request Router
					if err := http.ListenAndServe(address, standardRequestRouter); err != nil {
						result <- fmt.Errorf("Server failed with error: %v", err)
					} else {
						result <- nil
					}

				}

			}()

			// store the URL for later opening
			if httpsEnabled == false {
				endpointURL := httpEndpoint.DefaultURL()
				uniqueURLs[endpointURL] = endpointURL
			}

		}
	}

	// https
	if httpsEnabled {

		for _, tcpBinding := range httpsEndpoint.Bindings() {

			tcpBinding.AssignFreePort()

			tcpAddr := tcpBinding.GetTCPAddress()
			address := tcpAddr.String()

			// start listening
			go func() {
				server.logger.Info("HTTPS Endpoint: %s", address)

				// Standard HTTPS Request Router
				if err := http.ListenAndServeTLS(address, httpsEndpoint.CertFilePath(), httpsEndpoint.KeyFilePath(), standardRequestRouter); err != nil {
					result <- fmt.Errorf("Server failed with error: %v", err)
				} else {
					result <- nil
				}

			}()

			// store the URL for later opening
			endpointURL := httpsEndpoint.DefaultURL()
			uniqueURLs[endpointURL] = endpointURL
		}

	}

	// docx conversion endpoint (unencrypted, no authentication)
	if server.config.Conversion.DOCX.IsEnabled() {

		// start listening
		go func() {
			conversionEndpointTCPAddress := server.config.Conversion.EndpointBinding().GetTCPAddress()
			conversionEndpointAddress := conversionEndpointTCPAddress.String()

			server.logger.Info("Docx Conversion Endpoint: %s", conversionEndpointAddress)

			// Standard HTTPS Request Router
			if err := http.ListenAndServe(conversionEndpointAddress, server.getLocalRequestRouter()); err != nil {
				result <- fmt.Errorf("Docx Conversion endpoint failed with error: %v", err)
			} else {
				result <- nil
			}

		}()

	}

	// open HTTP URL(s) in a browser
	for _, url := range uniqueURLs {
		server.logger.Info("Open URL: %s", url)
		go open.Run(url)
	}

	return result
}

// getRedirectRouter returns a router which redirects all requests to the url with the given base.
func (server *Server) getRedirectRouter(baseURITarget string, baseHandler http.Handler) *mux.Router {
	redirectRouter := mux.NewRouter()

	for _, requestHandler := range handlers.GetRedirectHandlers(server.logger, baseURITarget, baseHandler) {
		requestRoute := requestHandler.Route
		requestHandler := requestHandler.Handler

		redirectRouter.Handle(requestRoute, requestHandler)
	}

	return redirectRouter
}

// Get an instance of the standard request router for all repository related routes.
func (server *Server) getStandardRequestRouter() *mux.Router {

	// register requst routers
	requestRouter := mux.NewRouter()

	for _, requestHandler := range server.requestHandlers {
		requestRoute := requestHandler.Route
		requestHandler := requestHandler.Handler

		// add logging
		requestHandler = handlers.LogRequests(requestHandler)

		// add compression
		requestHandler = handlers.CompressResponses(requestHandler)

		// add authentication
		if _, httpsEnabled := server.httpsEndpoint(); httpsEnabled && server.config.AuthenticationIsEnabled() {
			secretProvider := server.config.GetAuthenticationUserStore()
			if secretProvider == nil {
				panic("Authentication is enabled but the supplied secret provider is nil.")
			}

			requestHandler = handlers.RequireDigestAuthentication(server.logger, requestHandler, secretProvider)
		}

		requestRouter.Handle(requestRoute, requestHandler)
	}

	return requestRouter
}

// getLocalRequestRouter returns a local request router without compression and without authentication.
func (server *Server) getLocalRequestRouter() *mux.Router {

	// register requst routers
	requestRouter := mux.NewRouter()

	for _, requestHandler := range server.requestHandlers {
		requestRoute := requestHandler.Route
		requestHandler := requestHandler.Handler

		// add logging
		requestHandler = handlers.LogRequests(requestHandler)

		requestRouter.Handle(requestRoute, requestHandler)
	}

	return requestRouter
}

// Get the http binding if it is enabled.
func (server *Server) httpEndpoint() (httpEndpoint HTTPEndpoint, enabled bool) {

	if !server.config.Server.HTTP.Enabled {
		return HTTPEndpoint{}, false
	}

	return HTTPEndpoint{
		isSecure:    false,
		forceHTTPS:  server.config.Server.HTTPS.HTTPSIsForced(),
		tcpBindings: server.config.Server.HTTP.Bindings,
	}, true

}

// Get the https binding if it is enabled.tcpBinding
func (server *Server) httpsEndpoint() (httpsEndpoint HTTPSEndpoint, enabled bool) {

	if !server.config.Server.HTTPS.Enabled {
		return HTTPSEndpoint{}, false
	}

	httpEndpoint := HTTPEndpoint{
		domain:      server.config.Server.DomainName,
		isSecure:    true,
		tcpBindings: server.config.Server.HTTPS.Bindings,
	}

	certFilePath, keyFilePath, _ := server.config.CertificateFilePaths()

	httpsEndpoint = HTTPSEndpoint{
		HTTPEndpoint: httpEndpoint,
		certFilePath: certFilePath,
		keyFilePath:  keyFilePath,
	}

	return httpsEndpoint, true

}

// HTTPEndpoint contains HTTP server endpoint parameters such as a domain name and TCP bindings.
type HTTPEndpoint struct {
	domain      string
	isSecure    bool
	forceHTTPS  bool
	tcpBindings []*config.TCPBinding
}

func (endpoint *HTTPEndpoint) String() string {
	return fmt.Sprintf("Domain: %q, IsSecure: %v, ForeceHTTPs: %v", endpoint.domain, endpoint.isSecure, endpoint.forceHTTPS)
}

// IsSecure returns a flag indicating whether the current HTTPEndpoint is secure (HTTPS) or not.
func (endpoint *HTTPEndpoint) IsSecure() bool {
	return endpoint.isSecure
}

// Protocol returns the protocol of the current HTTPEndpoint. "https" if this endpoint is secure; otherwise "http".
func (endpoint *HTTPEndpoint) Protocol() string {
	if endpoint.isSecure {
		return "https"
	}
	return "http"
}

// ForceHTTPS returns a flag indicating whether a secure connection shall be preferred over insecure connections.
func (endpoint *HTTPEndpoint) ForceHTTPS() bool {
	return endpoint.forceHTTPS
}

// Bindings returns all TCP bindings of the current HTTP endpoint.
func (endpoint *HTTPEndpoint) Bindings() []*config.TCPBinding {
	return endpoint.tcpBindings
}

// DefaultURL return the default url for the current HTTP endpoint. It will include the domain name if one is configured.
// If none is configured it will use the IP address as the host name.
func (endpoint *HTTPEndpoint) DefaultURL() string {

	// no point in returning a url if there are no tcp bindings
	if len(endpoint.tcpBindings) == 0 {
		return ""
	}

	// use the first tcp binding as the default
	defaultBinding := *endpoint.tcpBindings[0]

	// create an URL from the tcp binding if no domain is configured
	if endpoint.domain == "" {
		return getURL(*endpoint, defaultBinding)
	}

	// determine the port suffix (e.g. ":8080")
	portSuffix := ""
	portNumber := defaultBinding.Port
	isDefaultPort := portNumber == 80 || portNumber == 443
	if !isDefaultPort {
		portSuffix = fmt.Sprintf(":%v", portNumber)
	}

	return fmt.Sprintf("%s://%s%s", endpoint.Protocol(), endpoint.domain, portSuffix)
}

// HTTPSEndpoint contains a secure version of a HTTPEndpoint with parameters for secure TLS connections such as the certificate paths.
type HTTPSEndpoint struct {
	HTTPEndpoint

	certFilePath string
	keyFilePath  string
}

// CertFilePath returns the SSL certificate file (e.g. "cert.pem") name of this HTTPSEndpoint.
func (endpoint *HTTPSEndpoint) CertFilePath() string {
	return endpoint.certFilePath
}

// KeyFilePath returns the SSL certificate key file name (e.g. "cert.key") of this HTTPSEndpoint.
func (endpoint *HTTPSEndpoint) KeyFilePath() string {
	return endpoint.keyFilePath
}

// getURL returns the formatted URL (e.g. "https://localhost:8080") for the given TCP binding,
// using the IP address as the hostname.
func getURL(endpoint HTTPEndpoint, tcpBinding config.TCPBinding) string {
	tcpAddress := tcpBinding.GetTCPAddress()
	hostname := tcpAddress.String()

	// don't use wildcard addresses for the URL, use localhost instead
	if isWildcardAddress(tcpAddress.IP) {
		if isIPv6Address(tcpAddress.IP) {
			// IPv6 addresses are wrapped in brackets
			hostname = strings.Replace(hostname, fmt.Sprintf("[%s]", tcpAddress.IP.String()), "localhost", 1)
		} else {
			hostname = strings.Replace(hostname, tcpAddress.IP.String(), "localhost", 1)
		}
	}

	return fmt.Sprintf("%s://%s", endpoint.Protocol(), hostname)
}

// isIPv6Address checks if the given IP address is a IPv6 address or not.
func isIPv6Address(ip net.IP) bool {
	// if the ip cannot be converted to IPv4 it must be an IPv6 address
	return ip.To4() == nil
}

// isWildcardAddress checks if the supplied ip is a "source address for this host on this network"
// (see: RFC 5735, Section 3, example: 0.0.0.0)
func isWildcardAddress(ip net.IP) bool {
	switch ip.String() {
	case "0.0.0.0":
		return true

	case "::":
		return true

	default:
		return false
	}
}
