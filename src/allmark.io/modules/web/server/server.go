// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"allmark.io/modules/common/config"
	"allmark.io/modules/common/logger"
	"allmark.io/modules/dataaccess"
	"allmark.io/modules/services/converter"
	"allmark.io/modules/services/parser"
	"allmark.io/modules/web/handlers"
	"allmark.io/modules/web/header"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/webpaths"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/skratchdot/open-golang/open"
	"net/http"
)

func New(logger logger.Logger, config config.Config, repository dataaccess.Repository, parser parser.Parser, converter converter.Converter) (*Server, error) {

	// create the request handlers
	patherFactory := webpaths.NewFactory(logger, repository)
	itemPathProvider := patherFactory.Absolute(handlers.BasePath)
	tagPathProvider := patherFactory.Absolute(handlers.TagPathPrefix)
	webPathProvider := webpaths.NewWebPathProvider(patherFactory, itemPathProvider, tagPathProvider)
	templateProvider := templates.NewProvider(config.TemplatesFolder())
	orchestratorFactory := orchestrator.NewFactory(logger, config, repository, parser, converter, webPathProvider)
	reindexInterval := config.Indexing.IntervalInSeconds
	headerWriterFactory := header.NewHeaderWriterFactory(reindexInterval)
	requestHandlers := handlers.GetBaseHandlers(logger, config, templateProvider, *orchestratorFactory, headerWriterFactory)

	return &Server{
		logger: logger,
		config: config,

		headerWriterFactory: headerWriterFactory,
		requestHandlers:     requestHandlers,
	}, nil

}

type Server struct {
	logger logger.Logger
	config config.Config

	headerWriterFactory header.WriterFactory

	requestHandlers handlers.HandlerList
}

func (server *Server) Start() chan error {

	result := make(chan error)

	standardRequestRouter := server.getStandardRequestRouter()

	// bindings
	httpEndpoint, httpEnabled := server.httpEndpoint()
	httpsEndpoint, httpsEnabled := server.httpsEndpoint()

	uniqueUrls := make(map[string]string)

	// http
	if httpEnabled {

		for _, tcpBinding := range httpEndpoint.Bindings() {

			tcpBinding.AssignFreePort()

			tcpAddr := tcpBinding.GetTCPAddress()
			address := tcpAddr.String()

			// start listening
			go func() {
				server.logger.Info("HTTP Endpoint: %s", address)

				if httpEndpoint.ForceHttps() {

					// Redirect HTTP → HTTPs
					redirectTarget := httpsEndpoint.DefaultURL()
					httpsRedirectRouter := server.getRedirectRouter(redirectTarget)

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
				uniqueUrls[endpointURL] = endpointURL
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
				server.logger.Info("HTTPs Endpoint: %s", address)

				// Standard HTTPs Request Router
				if err := http.ListenAndServeTLS(address, httpsEndpoint.CertFilePath(), httpsEndpoint.KeyFilePath(), standardRequestRouter); err != nil {
					result <- fmt.Errorf("Server failed with error: %v", err)
				} else {
					result <- nil
				}

			}()

			// store the URL for later opening
			endpointURL := httpsEndpoint.DefaultURL()
			uniqueUrls[endpointURL] = endpointURL
		}

	}

	// open HTTP URL(s) in a browser
	for _, url := range uniqueUrls {
		server.logger.Info("Open Url: %s", url)
		go open.Run(url)
	}

	return result
}

// Get a redirect router which redirects all requests to the url with the given base.
func (server *Server) getRedirectRouter(baseUriTarget string) *mux.Router {
	redirectRouter := mux.NewRouter()

	for _, requestHandler := range handlers.GetRedirectHandlers(baseUriTarget) {
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

			requestHandler = handlers.RequireDigestAuthentication(requestHandler, secretProvider)
		}

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
		forceHttps:  server.config.Server.HTTPs.HTTPsIsForced(),
		tcpBindings: server.config.Server.HTTP.Bindings,
	}, true

}

// Get the https binding if it is enabled.tcpBinding
func (server *Server) httpsEndpoint() (httpsEndpoint HTTPsEndpoint, enabled bool) {

	if !server.config.Server.HTTPs.Enabled {
		return HTTPsEndpoint{}, false
	}

	httpEndpoint := HTTPEndpoint{
		domain:      server.config.Server.DomainName,
		isSecure:    true,
		tcpBindings: server.config.Server.HTTPs.Bindings,
	}

	certFilePath, keyFilePath := server.config.CertificateFilePaths()

	httpsEndpoint = HTTPsEndpoint{
		HTTPEndpoint: httpEndpoint,
		certFilePath: certFilePath,
		keyFilePath:  keyFilePath,
	}

	return httpsEndpoint, true

}

type HTTPEndpoint struct {
	domain      string
	isSecure    bool
	forceHttps  bool
	tcpBindings []*config.TCPBinding
}

func (endpoint *HTTPEndpoint) IsSecure() bool {
	return endpoint.isSecure
}

func (endpoint *HTTPEndpoint) Protocol() string {
	if endpoint.isSecure {
		return "https"
	}
	return "http"
}

func (endpoint *HTTPEndpoint) ForceHttps() bool {
	return endpoint.forceHttps
}

func (endpoint *HTTPEndpoint) Bindings() []*config.TCPBinding {
	return endpoint.tcpBindings
}

func (endpoint *HTTPEndpoint) URL(tcpBinding config.TCPBinding) string {
	tcpAddress := tcpBinding.GetTCPAddress()
	return fmt.Sprintf("%s://%s", endpoint.Protocol(), tcpAddress.String())
}

func (endpoint *HTTPEndpoint) DefaultURL() string {

	if endpoint.domain != "" {
		return fmt.Sprintf("%s://%s", endpoint.Protocol(), endpoint.domain)
	}

	if len(endpoint.tcpBindings) == 0 {
		return ""
	}

	return endpoint.URL(*endpoint.tcpBindings[0])
}

type HTTPsEndpoint struct {
	HTTPEndpoint

	certFilePath string
	keyFilePath  string
}

func (endpoint *HTTPsEndpoint) CertFilePath() string {
	return endpoint.certFilePath
}

func (endpoint *HTTPsEndpoint) KeyFilePath() string {
	return endpoint.keyFilePath
}
