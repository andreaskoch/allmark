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
	"strings"
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
	httpBinding, httpEnabled := server.getHttpBinding()
	httpsBinding, httpsEnabled := server.getHttpsBinding()

	// http
	if httpEnabled {

		go func() {
			server.logger.Info("HTTP Endpoint: %s", httpBinding.Url())

			if httpBinding.ForceHttps() {

				// Redirect HTTP → HTTPs
				redirectTarget := httpsBinding.Url()
				httpsRedirectRouter := server.getRedirectRouter(redirectTarget)

				if err := http.ListenAndServe(httpBinding.String(), httpsRedirectRouter); err != nil {
					result <- fmt.Errorf("Server failed with error: %v", err)
				} else {
					result <- nil
				}

			} else {

				// Standard HTTP Request Router
				if err := http.ListenAndServe(httpBinding.String(), standardRequestRouter); err != nil {
					result <- fmt.Errorf("Server failed with error: %v", err)
				} else {
					result <- nil
				}

			}

		}()
	}

	// https
	if httpsEnabled {

		go func() {
			server.logger.Info("HTTPs Endpoint: %s", httpsBinding.Url())

			// Standard HTTPs Request Router
			if err := http.ListenAndServeTLS(httpsBinding.String(), httpsBinding.CertFilePath(), httpsBinding.KeyFilePath(), standardRequestRouter); err != nil {
				result <- fmt.Errorf("Server failed with error: %v", err)
			} else {
				result <- nil
			}

		}()

	}

	// open url in browser
	repositoryUrl := httpBinding.Url()
	if httpsEnabled && httpBinding.ForceHttps() {
		repositoryUrl = httpsBinding.Url()
	}

	open.Run(repositoryUrl)

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

		// add authentication
		if httpsBinding, httpsEnabled := server.getHttpsBinding(); httpsEnabled && server.config.AuthenticationIsEnabled() {
			secretProvider := server.config.GetAuthenticationUserStore()
			if secretProvider == nil {
				panic("Authentication is enabled but the supplied secret provider is nil.")
			}

			requestHandler = handlers.RequireDigestAuthentication(requestHandler, httpsBinding.Hostname(), secretProvider)
		}

		requestRouter.Handle(requestRoute, requestHandler)
	}

	return requestRouter
}

// Get the http binding if it is enabled.
func (server *Server) getHttpBinding() (httpBinding HttpBinding, enabled bool) {

	if !server.config.Server.Http.Enabled {
		return HttpBinding{}, false
	}

	return HttpBinding{
		hostname:   server.getHostname(),
		portNumber: server.config.Server.Http.GetPortNumber(),
		isSecure:   false,
		forceHttps: server.config.Server.Https.ForceHttps(),
	}, true

}

// Get the https binding if it is enabled.
func (server *Server) getHttpsBinding() (httpsBinding HttpsBinding, enabled bool) {

	if !server.config.Server.Https.Enabled {
		return HttpsBinding{}, false
	}

	httpBinding := HttpBinding{
		hostname:   server.getHostname(),
		portNumber: server.config.Server.Https.GetPortNumber(),
		isSecure:   true,
	}

	certFilePath, keyFilePath := server.config.CertificateFilePaths()

	httpsBinding = HttpsBinding{
		HttpBinding:  httpBinding,
		certFilePath: certFilePath,
		keyFilePath:  keyFilePath,
	}

	return httpsBinding, true

}

// Get the configured hostname (default: "localhost")
func (server *Server) getHostname() string {
	hostname := strings.ToLower(strings.TrimSpace(server.config.Server.Hostname))
	if hostname == "" {
		return "localhost"
	}

	return hostname
}

type HttpBinding struct {
	hostname   string
	portNumber int
	isSecure   bool
	forceHttps bool
}

func (b *HttpBinding) String() string {
	return fmt.Sprintf("%s:%v", b.hostname, b.portNumber)
}

func (b *HttpBinding) Hostname() string {
	return b.hostname
}

func (b *HttpBinding) PortNumber() int {
	return b.portNumber
}

func (b *HttpBinding) IsSecure() bool {
	return b.isSecure
}

func (b *HttpBinding) Url() string {
	protocol := "http"
	if b.isSecure {
		protocol = "https"
	}

	isStandardPort := b.PortNumber() == 80 || b.PortNumber() == 443
	if isStandardPort {
		return fmt.Sprintf("%s://%s", protocol, b.Hostname())
	}

	return fmt.Sprintf("%s://%s:%v", protocol, b.Hostname(), b.PortNumber())
}

func (b *HttpBinding) ForceHttps() bool {
	return b.forceHttps
}

type HttpsBinding struct {
	HttpBinding

	certFilePath string
	keyFilePath  string
}

func (https *HttpsBinding) CertFilePath() string {
	return https.certFilePath
}

func (https *HttpsBinding) KeyFilePath() string {
	return https.keyFilePath
}
