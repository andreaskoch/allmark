// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"github.com/andreaskoch/allmark/common/config"
	"github.com/andreaskoch/allmark/common/logger"
	"github.com/andreaskoch/allmark/common/util/fsutil"
	"github.com/andreaskoch/allmark/web/header"
	"github.com/andreaskoch/allmark/web/orchestrator"
	"github.com/andreaskoch/allmark/web/view/templates"
	"fmt"
	"net/http"
)

var (

	// BasePath defines the path path for all requests.
	BasePath = "/"

	// TagPathPrefix defines the prefix for tag-routes.
	TagPathPrefix = "/tags.html#"

	// TagmapHandlerRoute defines the route for tagmap-handler requests.
	TagmapHandlerRoute = "/tags.html"

	// ThemeRoutePrefix defines the route-prefix for theme files.
	ThemeRoutePrefix = "/theme"

	// ThemeHandlerRoute defines the route for thumbnails.
	ThemeHandlerRoute = fmt.Sprintf("%s/{path:.*$}", ThemeRoutePrefix)

	// ThumbnailRoutePrefix defines the route-prefix for thumbnails.
	ThumbnailRoutePrefix = "/thumbnails"

	// ThumbnailHandlerRoute defines the route for thumbnails.
	ThumbnailHandlerRoute = fmt.Sprintf("%s/{path:.*$}", ThumbnailRoutePrefix)

	// PrintHandlerRoute defines the route for print-handler requests.
	PrintHandlerRoute = `/{path:.+\.print$|print$}`

	// JSONHandlerRoute defines the route for JSON-handler requests.
	JSONHandlerRoute = `/{path:.+\.json$|json$}`

	// MarkdownHandlerRoute defines the route for Markdown-handler requests.
	MarkdownHandlerRoute = `/{path:.+\.markdown$|markdown$}`

	// LatestHandlerRoute defines the route for latest-handler requests.
	LatestHandlerRoute = `/{path:.+\.latest$|latest$}`

	// DOCXHandlerRoute defines the route for rich-text-handler requests.
	DOCXHandlerRoute = `/{path:.+\.docx$|docx$}`

	// UpdateHandlerRoute defines the route for update-handler requests.
	UpdateHandlerRoute = `/{path:.+\.ws$|ws$}`

	// ItemHandlerRoute defines the route for item-handler requests.
	ItemHandlerRoute = "/{path:.*$}"

	// SitemapHandlerRoute defines the route for sitemap-handler requests.
	SitemapHandlerRoute = "/sitemap.html"

	// XMLSitemapHandlerRoute defines the route for xml-sitemap-handler requests.
	XMLSitemapHandlerRoute = "/sitemap.xml"

	// RSSHandlerRoute defines the route for RSS-feed-handler requests.
	RSSHandlerRoute = "/feed.rss"

	// RobotsTxtHandlerRoute defines the route for robotstxt-handler requests.
	RobotsTxtHandlerRoute = "/robots.txt"

	// SearchHandlerRoute defines the route for search-handler requests.
	SearchHandlerRoute = "/search"

	// OpenSearchDescriptionHandlerRoute defines the route for opensearch-descriptiption-handler requests.
	OpenSearchDescriptionHandlerRoute = "/opensearch.xml"

	// TypeAheadSearchHandlerRoute defines the route for typeahead-search-handler requests.
	TypeAheadSearchHandlerRoute = "/search.json"

	// TypeAheadTitlesHandlerRoute defines the route for typeahead-titles-handler requests.
	TypeAheadTitlesHandlerRoute = "/titles.json"

	// RedirectHandlerRoute defines the route for redirect-handler requests.
	RedirectHandlerRoute = "/{path:.*$}"

	// AliasLookupHandlerRoute defines the route for alias-lookup-handler requests.
	AliasLookupHandlerRoute = "/!{alias:.+$}"

	// AliasIndexHandlerRoute defines the route for alias-lookup-handler requests.
	AliasIndexHandlerRoute = "/!"
)

// RouteAndHandler combines routes and http-handlers.
type RouteAndHandler struct {
	Route   string
	Handler http.Handler
}

// HandlerList is a list of routes and http-handlers.
type HandlerList []RouteAndHandler

// Add the specified route and http handler to the current list.
func (list *HandlerList) Add(route string, handler http.Handler) {
	*list = append(*list, RouteAndHandler{route, handler})
}

// GetRedirectHandlers returns a list of redirect handlers.
func GetRedirectHandlers(logger logger.Logger, baseURITarget string, baseHandler http.Handler) HandlerList {
	handlers := make(HandlerList, 0)
	handlers.Add(RedirectHandlerRoute, Redirect(logger, baseURITarget))
	return handlers
}

// GetBaseHandlers returns a full-list of all http-handlers in this package.
func GetBaseHandlers(logger logger.Logger, config config.Config, templateProvider templates.Provider, orchestratorFactory orchestrator.Factory, headerWriterFactory header.WriterFactory) HandlerList {
	handlers := make(HandlerList, 0)

	// orchestrators
	navigationOrchestrator := orchestratorFactory.NewNavigationOrchestrator()
	viewModelOrchestrator := orchestratorFactory.NewViewModelOrchestrator()
	fileOrchestrator := orchestratorFactory.NewFileOrchestrator()

	// global handlers
	errorHandler := Error(headerWriterFactory.Static(), templateProvider, navigationOrchestrator)

	itemHandler := Item(
		logger,
		headerWriterFactory.Dynamic(),
		fileOrchestrator,
		viewModelOrchestrator,
		templateProvider, errorHandler)

	// theme
	if themeFolder := config.ThemeFolder(); fsutil.DirectoryExists(themeFolder) {
		requestPrefixToStripFromRequestURI := "/" + config.Server.ThemeFolderName

		handlers.Add(
			ThemeHandlerRoute,
			AddETAgToStaticFileHandler(
				Static(
					themeFolder,
					ThemeRoutePrefix),
				headerWriterFactory.Static(), themeFolder, requestPrefixToStripFromRequestURI))

	} else {

		handlers.Add(
			ThemeHandlerRoute,
			InMemoryTheme(
				"/"+config.Server.ThemeFolderName+"/",
				headerWriterFactory.Static(),
				errorHandler))
	}

	// alias lookup
	handlers.Add(
		AliasLookupHandlerRoute,
		AliasLookup(headerWriterFactory.Dynamic(),
			viewModelOrchestrator,
			itemHandler))

	// alias index
	handlers.Add(
		AliasIndexHandlerRoute,
		AliasIndex(
			headerWriterFactory.Dynamic(),
			navigationOrchestrator,
			orchestratorFactory.NewAliasIndexOrchestrator(),
			templateProvider))

	// thumbnails
	if thumbnailsFolder := config.ThumbnailFolder(); fsutil.DirectoryExists(thumbnailsFolder) {
		requestPrefixToStripFromRequestURI := "/" + config.Conversion.Thumbnails.FolderName

		handlers.Add(
			ThumbnailHandlerRoute,
			AddETAgToStaticFileHandler(Static(thumbnailsFolder,
				ThumbnailRoutePrefix),
				headerWriterFactory.Static(),
				thumbnailsFolder,
				requestPrefixToStripFromRequestURI))
	}

	// robots.txt
	handlers.Add(RobotsTxtHandlerRoute, RobotsTxt(headerWriterFactory.Static(), templateProvider))

	// sitemap.html
	handlers.Add(
		SitemapHandlerRoute,
		Sitemap(headerWriterFactory.Dynamic(),
			navigationOrchestrator,
			orchestratorFactory.NewSitemapOrchestrator(),
			templateProvider))

	// tags.html
	handlers.Add(
		TagmapHandlerRoute,
		Tags(headerWriterFactory.Dynamic(),
			navigationOrchestrator,
			orchestratorFactory.NewTagsOrchestrator(),
			templateProvider))

	// search
	handlers.Add(
		SearchHandlerRoute,
		Search(
			headerWriterFactory.Dynamic(),
			navigationOrchestrator,
			orchestratorFactory.NewSearchOrchestrator(),
			templateProvider,
			errorHandler))

	// sitemap.xml
	handlers.Add(
		XMLSitemapHandlerRoute,
		XMLSitemap(headerWriterFactory.Dynamic(),
			orchestratorFactory.NewXMLSitemapOrchestrator(),
			templateProvider))

	// opensearch.xml
	handlers.Add(
		OpenSearchDescriptionHandlerRoute,
		OpenSearchDescription(headerWriterFactory.Static(),
			orchestratorFactory.NewOpenSearchDescriptionOrchestrator(),
			templateProvider))

	// titles.json
	handlers.Add(
		TypeAheadTitlesHandlerRoute,
		Titles(headerWriterFactory.Dynamic(),
			orchestratorFactory.NewTitlesOrchestrator()))

	// search.json
	handlers.Add(
		TypeAheadSearchHandlerRoute,
		TypeAhead(headerWriterFactory.Dynamic(),
			orchestratorFactory.NewTypeAheadOrchestrator()))

	// latest.json
	handlers.Add(LatestHandlerRoute, Latest(logger, headerWriterFactory.Dynamic(), viewModelOrchestrator, itemHandler))

	// rss
	handlers.Add(
		RSSHandlerRoute,
		RSS(headerWriterFactory.Dynamic(),
			orchestratorFactory.NewFeedOrchestrator(),
			templateProvider,
			errorHandler))

	// json
	handlers.Add(JSONHandlerRoute,
		JSON(headerWriterFactory.Dynamic(),
			viewModelOrchestrator,
			itemHandler))

	// markdown
	handlers.Add(MarkdownHandlerRoute,
		Markdown(headerWriterFactory.Dynamic(),
			viewModelOrchestrator,
			itemHandler))

	// conversion
	conversionModelOrchestrator := orchestratorFactory.NewConversionModelOrchestrator()

	// print
	handlers.Add(
		PrintHandlerRoute,
		Print(logger,
			headerWriterFactory.Dynamic(),
			conversionModelOrchestrator,
			templateProvider,
			errorHandler))

	// docx
	conversionEndpointTCPAddress := config.Conversion.EndpointBinding().GetTCPAddress()
	conversionEndpointAddress := conversionEndpointTCPAddress.String()
	handlers.Add(
		DOCXHandlerRoute,
		DOCX(logger,
			config.Conversion.DOCX.Tool(),
			conversionEndpointAddress,
			headerWriterFactory.Dynamic(),
			conversionModelOrchestrator,
			templateProvider,
			errorHandler))

	// update
	handlers.Add(
		UpdateHandlerRoute,
		Update(
			logger,
			headerWriterFactory.Dynamic(),
			templateProvider,
			orchestratorFactory.NewUpdateOrchestrator()))

	// items
	handlers.Add(
		ItemHandlerRoute,
		itemHandler)

	return handlers
}
