// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"allmark.io/modules/common/config"
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/util/fsutil"
	"allmark.io/modules/common/util/hashutil"
	"allmark.io/modules/dataaccess"
	"allmark.io/modules/services/converter"
	"allmark.io/modules/services/parser"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/server/handler"
	"allmark.io/modules/web/server/header"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/webpaths"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/net/websocket"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	BasePath      = "/"
	TagPathPrefix = fmt.Sprintf("%stags.html#", BasePath)

	// Static Routes
	StaticThemeFolderRoute     = "/theme"
	StaticThumbnailFolderRoute = "/thumbnails"

	// Dynamic Routes
	PrintHandlerRoute  = `/{path:.+\.print$|print$}`
	JsonHandlerRoute   = `/{path:.+\.json$|json$}`
	LatestHandlerRoute = `/{path:.+\.latest$|latest$}`
	RtfHandlerRoute    = `/{path:.+\.rtf$|rtf$}`
	UpdateHandlerRoute = `/{path:.+\.ws$|ws$}`

	ItemHandlerRoute  = "/{path:.*$}"
	ThemeHandlerRoute = fmt.Sprintf("%s/{path:.*$}", StaticThemeFolderRoute)

	TagmapHandlerRoute                = "/tags.html"
	SitemapHandlerRoute               = "/sitemap.html"
	XmlSitemapHandlerRoute            = "/sitemap.xml"
	RssHandlerRoute                   = "/feed.rss"
	RobotsTxtHandlerRoute             = "/robots.txt"
	SearchHandlerRoute                = "/search"
	OpenSearchDescriptionHandlerRoute = "/opensearch.xml"

	TypeAheadSearchHandlerRoute = "/search.json"
	TypeAheadTitlesHandlerRoute = "/titles.json"
)

func New(logger logger.Logger, config config.Config, repository dataaccess.Repository, parser parser.Parser, converter converter.Converter) (*Server, error) {

	// paths
	patherFactory := webpaths.NewFactory(logger, repository)
	itemPathProvider := patherFactory.Absolute(BasePath)
	tagPathProvider := patherFactory.Absolute(TagPathPrefix)
	webPathProvider := webpaths.NewWebPathProvider(patherFactory, itemPathProvider, tagPathProvider)

	// handlers
	templateProvider := templates.NewProvider(config.TemplatesFolder())
	orchestratorFactory := orchestrator.NewFactory(logger, config, repository, parser, converter, webPathProvider)
	reindexInterval := config.Indexing.IntervalInSeconds
	headerWriterFactory := header.NewHeaderWriterFactory(reindexInterval)
	handlerFactory := handler.NewFactory(logger, config, templateProvider, *orchestratorFactory, headerWriterFactory)

	return &Server{
		logger: logger,
		config: config,

		headerWriterFactory: headerWriterFactory,
		handlerFactory:      handlerFactory,
	}, nil

}

type Server struct {
	logger logger.Logger
	config config.Config

	headerWriterFactory header.WriterFactory
	handlerFactory      *handler.Factory

	standardRequestRouter *mux.Router
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
	redirectHandler := server.handlerFactory.NewRedirectHandler(baseUriTarget)
	redirectRouter.HandleFunc("/{path:.*$}", redirectHandler.Func())
	return redirectRouter
}

// Get an instance of the standard request router for all repository related routes.
func (server *Server) getStandardRequestRouter() *mux.Router {

	if server.standardRequestRouter != nil {
		return server.standardRequestRouter
	}

	// instantiate handlers
	updateHandler := server.handlerFactory.NewUpdateHandler()
	robotsTxtHandler := server.handlerFactory.NewRobotsTxtHandler()
	xmlSitemapHandler := server.handlerFactory.NewXmlSitemapHandler()
	tagsHandler := server.handlerFactory.NewTagsHandler()
	sitemapHandler := server.handlerFactory.NewSitemapHandler()
	rssHandler := server.handlerFactory.NewRssHandler()
	printHandler := server.handlerFactory.NewPrintHandler()
	searchHandler := server.handlerFactory.NewSearchHandler()
	opensearchDescriptionHandler := server.handlerFactory.NewOpenSearchDescriptionHandler()
	typeAheadHandler := server.handlerFactory.NewTypeAheadSearchHandler()
	titlesHandler := server.handlerFactory.NewTypeAheadTitlesHandler()
	rtfHandler := server.handlerFactory.NewRtfHandler()
	jsonHandler := server.handlerFactory.NewJsonHandler()
	latestHanlder := server.handlerFactory.NewLatestHandler()
	itemHandler := server.handlerFactory.NewItemHandler()

	// register requst routers
	requestRouter := mux.NewRouter()
	requestRouter.Handle(UpdateHandlerRoute, websocket.Handler(updateHandler.Func()))

	// serve auxiliary dynamic files
	requestRouter.HandleFunc(RobotsTxtHandlerRoute, robotsTxtHandler.Func())
	requestRouter.HandleFunc(XmlSitemapHandlerRoute, xmlSitemapHandler.Func())
	requestRouter.HandleFunc(TagmapHandlerRoute, tagsHandler.Func())
	requestRouter.HandleFunc(SitemapHandlerRoute, sitemapHandler.Func())
	requestRouter.HandleFunc(RssHandlerRoute, rssHandler.Func())
	requestRouter.HandleFunc(PrintHandlerRoute, printHandler.Func())
	requestRouter.HandleFunc(SearchHandlerRoute, searchHandler.Func())
	requestRouter.HandleFunc(OpenSearchDescriptionHandlerRoute, opensearchDescriptionHandler.Func())
	requestRouter.HandleFunc(TypeAheadSearchHandlerRoute, typeAheadHandler.Func())
	requestRouter.HandleFunc(TypeAheadTitlesHandlerRoute, titlesHandler.Func())

	// theme
	if themeFolder := server.config.ThemeFolder(); fsutil.DirectoryExists(themeFolder) {

		// serve static
		themeFolderHandler := http.FileServer(http.Dir(themeFolder))
		requestPrefixToStripFromRequestUri := "/" + config.ThemeFolderName
		s := http.StripPrefix(StaticThemeFolderRoute, server.addStaticFileHeaders(themeFolder, requestPrefixToStripFromRequestUri, themeFolderHandler))
		requestRouter.PathPrefix(StaticThemeFolderRoute).Handler(s)

	} else {

		// serve dynamic
		server.logger.Info("Serving default theme-files from memory. If you want to serve the theme from disc please use the 'init' command on your repository or home folder.")
		themeHandler := server.handlerFactory.NewThemeHandler()
		requestRouter.HandleFunc(ThemeHandlerRoute, themeHandler.Func())

	}

	// serve thumbnails
	if thumbnailsFolder := server.config.ThumbnailFolder(); fsutil.DirectoryExists(thumbnailsFolder) {
		thumbnailsFolderHandler := http.FileServer(http.Dir(thumbnailsFolder))
		requestPrefixToStripFromRequestUri := "/" + config.ThumbnailsFolderName
		s := http.StripPrefix(StaticThumbnailFolderRoute, server.addStaticFileHeaders(thumbnailsFolder, requestPrefixToStripFromRequestUri, thumbnailsFolderHandler))
		requestRouter.PathPrefix(StaticThumbnailFolderRoute).Handler(s)
	}

	// serve items
	requestRouter.HandleFunc(RtfHandlerRoute, rtfHandler.Func())
	requestRouter.HandleFunc(JsonHandlerRoute, jsonHandler.Func())
	requestRouter.HandleFunc(LatestHandlerRoute, latestHanlder.Func())
	requestRouter.HandleFunc(ItemHandlerRoute, itemHandler.Func())

	server.standardRequestRouter = requestRouter
	return server.standardRequestRouter
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

func (server *Server) addStaticFileHeaders(baseFolder, requestPrefixToStripFromRequestUri string, h http.Handler) http.Handler {

	staticHeaderWriter := server.headerWriterFactory.Static()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		staticHeaderWriter.Write(w, "")

		// determine the hash
		etag := ""

		// prepare the request uri
		requestUri := r.RequestURI
		if requestPrefixToStripFromRequestUri != "" {
			requestUri = stripPathFromRequest(r, requestPrefixToStripFromRequestUri)
		}

		// assemble the filepath on disc
		filePath := filepath.Join(baseFolder, requestUri)

		// read the the hash
		if file, err := fsutil.OpenFile(filePath); err == nil {
			defer file.Close()
			if fileHash, hashErr := hashutil.GetHash(file); hashErr == nil {
				etag = fileHash
			}
		}
		if etag != "" {
			header.ETag(w, etag)
		}

		h.ServeHTTP(w, r)
	})
}

func stripPathFromRequest(r *http.Request, path string) string {
	return strings.TrimPrefix(r.RequestURI, path)
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
