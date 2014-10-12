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
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/allmark2/services/converter"
	"github.com/andreaskoch/allmark2/services/parser"
	"github.com/andreaskoch/allmark2/web/orchestrator"
	"github.com/andreaskoch/allmark2/web/server/handler"
	"github.com/andreaskoch/allmark2/web/server/header"
	"github.com/andreaskoch/allmark2/web/webpaths"
	"github.com/gorilla/mux"
	"github.com/skratchdot/open-golang/open"
	"math"
	"net/http"
	"strings"
)

var (
	BasePath      = "/"
	TagPathPrefix = fmt.Sprintf("%stags.html#", BasePath)

	// Dynamic Routes
	PrintHandlerRoute  = `/{path:.+\.print$|print$}`
	JsonHandlerRoute   = `/{path:.+\.json$|json$}`
	LatestHandlerRoute = `/{path:.+\.latest$|latest$}`
	RtfHandlerRoute    = `/{path:.+\.rtf$|rtf$}`
	UpdateHandlerRoute = `/{path:.+\.ws$|ws$}`

	ItemHandlerRoute = "/{path:.*$}"

	TagmapHandlerRoute                = "/tags.html"
	SitemapHandlerRoute               = "/sitemap.html"
	XmlSitemapHandlerRoute            = "/sitemap.xml"
	RssHandlerRoute                   = "/feed.rss"
	RobotsTxtHandlerRoute             = "/robots.txt"
	SearchHandlerRoute                = "/search"
	OpenSearchDescriptionHandlerRoute = "/opensearch.xml"

	TypeAheadSearchHandlerRoute = "/search.json"
	TypeAheadTitlesHandlerRoute = "/titles.json"

	// Static Routes
	ThemeFolderRoute = "/theme"
)

func New(logger logger.Logger, config config.Config, repository dataaccess.Repository, parser parser.Parser, converter converter.Converter) (*Server, error) {

	// paths
	patherFactory := webpaths.NewFactory(logger, repository)
	itemPathProvider := patherFactory.Absolute(BasePath)
	tagPathProvider := patherFactory.Absolute(TagPathPrefix)
	webPathProvider := webpaths.NewWebPathProvider(patherFactory, itemPathProvider, tagPathProvider)

	// orchestrator
	orchestratorFactory := orchestrator.NewFactory(logger, config, repository, parser, converter, webPathProvider)

	// handlers
	handlerFactory := handler.NewFactory(logger, config, *orchestratorFactory)

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
	result := make(chan error)

	go func() {
		server.isRunning = true

		// register requst routers
		requestRouter := mux.NewRouter()

		// websocket update handler
		// updateHub := update.NewHub(server.logger, server.updateHub)
		// go updateHub.Run()

		updateHandler := server.handlerFactory.NewUpdateHandler()
		requestRouter.Handle(UpdateHandlerRoute, websocket.Handler(updateHandler.Func()))

		// serve auxiliary dynamic files
		requestRouter.HandleFunc(RobotsTxtHandlerRoute, header.GzipHandlerFunc(server.handlerFactory.NewRobotsTxtHandler().Func()))
		requestRouter.HandleFunc(XmlSitemapHandlerRoute, header.GzipHandlerFunc(server.handlerFactory.NewXmlSitemapHandler().Func()))
		requestRouter.HandleFunc(TagmapHandlerRoute, header.GzipHandlerFunc(server.handlerFactory.NewTagsHandler().Func()))
		requestRouter.HandleFunc(SitemapHandlerRoute, header.GzipHandlerFunc(server.handlerFactory.NewSitemapHandler().Func()))
		requestRouter.HandleFunc(RssHandlerRoute, header.GzipHandlerFunc(server.handlerFactory.NewRssHandler().Func()))
		requestRouter.HandleFunc(PrintHandlerRoute, header.GzipHandlerFunc(server.handlerFactory.NewPrintHandler().Func()))
		requestRouter.HandleFunc(SearchHandlerRoute, header.GzipHandlerFunc(server.handlerFactory.NewSearchHandler().Func()))
		requestRouter.HandleFunc(OpenSearchDescriptionHandlerRoute, header.GzipHandlerFunc(server.handlerFactory.NewOpenSearchDescriptionHandler().Func()))
		requestRouter.HandleFunc(TypeAheadSearchHandlerRoute, header.GzipHandlerFunc(server.handlerFactory.NewTypeAheadSearchHandler().Func()))
		requestRouter.HandleFunc(TypeAheadTitlesHandlerRoute, header.GzipHandlerFunc(server.handlerFactory.NewTypeAheadTitlesHandler().Func()))

		// serve static files
		if themeFolder := server.config.ThemeFolder(); fsutil.DirectoryExists(themeFolder) {
			s := http.StripPrefix(ThemeFolderRoute, maxAgeHandler(header.STATICCONTENT_CACHEDURATION_SECONDS, http.FileServer(http.Dir(themeFolder))))
			requestRouter.PathPrefix(ThemeFolderRoute).Handler(header.GzipHandler(s))
		}

		// rich text
		requestRouter.HandleFunc(RtfHandlerRoute, server.handlerFactory.NewRtfHandler().Func())

		// serve items
		requestRouter.HandleFunc(JsonHandlerRoute, header.GzipHandlerFunc(server.handlerFactory.NewJsonHandler().Func()))
		requestRouter.HandleFunc(LatestHandlerRoute, header.GzipHandlerFunc(server.handlerFactory.NewLatestHandler().Func()))
		requestRouter.HandleFunc(ItemHandlerRoute, header.GzipHandlerFunc(server.handlerFactory.NewItemHandler().Func()))

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

func maxAgeHandler(seconds int, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header.Cache(w, r, seconds)
		h.ServeHTTP(w, r)
	})
}
