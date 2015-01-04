// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/common/util/hashutil"
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/allmark2/services/converter"
	"github.com/andreaskoch/allmark2/services/parser"
	"github.com/andreaskoch/allmark2/web/orchestrator"
	"github.com/andreaskoch/allmark2/web/server/handler"
	"github.com/andreaskoch/allmark2/web/server/header"
	"github.com/andreaskoch/allmark2/web/view/templates"
	"github.com/andreaskoch/allmark2/web/webpaths"
	"github.com/gorilla/mux"
	"github.com/skratchdot/open-golang/open"
	"math"
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
	handlerFactory := handler.NewFactory(logger, config, templateProvider, *orchestratorFactory)

	return &Server{
		logger: logger,
		config: config,

		handlerFactory: handlerFactory,
	}, nil

}

type Server struct {
	isRunning bool

	logger logger.Logger
	config config.Config

	handlerFactory *handler.Factory
}

func (server *Server) IsRunning() bool {
	return server.isRunning
}

func (server *Server) Start() chan error {

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
		s := http.StripPrefix(StaticThemeFolderRoute, addStaticFileHeaders(themeFolder, "/"+config.ThemeFolderName, header.STATICCONTENT_CACHEDURATION_SECONDS, themeFolderHandler))
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
		s := http.StripPrefix(StaticThumbnailFolderRoute, addStaticFileHeaders(thumbnailsFolder, "/"+config.ThumbnailsFolderName, header.STATICCONTENT_CACHEDURATION_SECONDS, thumbnailsFolderHandler))
		requestRouter.PathPrefix(StaticThumbnailFolderRoute).Handler(s)
	}

	// serve items
	requestRouter.HandleFunc(RtfHandlerRoute, rtfHandler.Func())
	requestRouter.HandleFunc(JsonHandlerRoute, jsonHandler.Func())
	requestRouter.HandleFunc(LatestHandlerRoute, latestHanlder.Func())
	requestRouter.HandleFunc(ItemHandlerRoute, itemHandler.Func())

	result := make(chan error)

	go func() {
		server.isRunning = true

		// start http server: http
		httpBinding := server.getHttpBinding()
		server.logger.Info("Starting http server %q\n", httpBinding)

		if err := http.ListenAndServe(httpBinding, requestRouter); err != nil {
			result <- fmt.Errorf("Server failed with error: %v", err)
		} else {
			result <- nil
		}

		server.isRunning = false
	}()

	// open the repository in the browser
	open.Run(server.getAddress())

	return result
}

func (server *Server) getHttpBinding() string {

	hostname := server.getHostname()
	port := server.getPort()

	if strings.TrimSpace(hostname) == "" {
		fmt.Sprintf(":%v", port)
	}

	return fmt.Sprintf("%s:%v", hostname, port)
}

func (server *Server) getAddress() string {
	hostname := server.getHostname()
	port := server.getPort()

	switch port {
	case 80:
		return fmt.Sprintf("http://%s", hostname)
	default:
		return fmt.Sprintf("http://%s:%v", hostname, port)
	}

	panic("Unreachable")
}

func (server *Server) getHostname() string {
	hostname := strings.ToLower(strings.TrimSpace(server.config.Server.Http.Hostname))
	if hostname == "" {
		return "localhost"
	}

	return hostname
}

func (server *Server) getPort() int {
	port := server.config.Server.Http.Port
	if port < 1 || port > math.MaxUint16 {
		panic(fmt.Sprintf("%q is an invalid value for a port. Ports can only be in the range of %v to %v,", port, 1, math.MaxUint16))
	}

	return port
}

func addStaticFileHeaders(baseFolder, requestPrefixToStripFromRequestUri string, seconds int, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
			header.ETag(w, r, etag)
		}

		header.Cache(w, r, seconds)
		header.VaryAcceptEncoding(w, r)

		h.ServeHTTP(w, r)
	})
}

func stripPathFromRequest(r *http.Request, path string) string {
	return strings.TrimPrefix(r.RequestURI, path)
}
