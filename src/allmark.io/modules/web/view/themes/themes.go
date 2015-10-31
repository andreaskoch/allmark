// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themes

import "allmark.io/modules/web/view/themes/themefiles"

var defaultTheme *Theme

// initialize themes
func init() {

	defaultTheme = &Theme{
		Name: "default",
		Files: []*ThemeFile{

			// styles
			newFileFromText("screen.css", themefiles.ScreenCss),
			newFileFromText("print.css", themefiles.PrintCss),

			// assets
			newFileFromBase64("tree-node.png", themefiles.NodePng),
			newFileFromBase64("tree-vertical-line.png", themefiles.VerticalLinePng),
			newFileFromBase64("tree-last-node.png", themefiles.LastNodePng),

			// javascript libraries
			newFileFromText("modernizr.js", themefiles.Modernizr),
			newFileFromText("jquery.js", themefiles.JqueryJs),

			// web socket auto update
			newFileFromText("autoupdate.js", themefiles.AutoupdateJs),

			// favicon
			newFileFromBase64("favicon.ico", themefiles.FaviconIco),

			// Github Ribbon
			newFileFromBase64("github-ribbon.png", themefiles.GithubRibbonPNG),

			// presentations
			newFileFromText("deck.js", themefiles.DeckJs),
			newFileFromText("deck.css", themefiles.DeckCss),
			newFileFromText("presentation.js", themefiles.PresentationJs),

			// auto-suggest
			newFileFromText("typeahead.js", themefiles.TypeAheadJs),
			newFileFromText("search.js", themefiles.SearchJs),

			// code highlighting
			newFileFromBase64("codehighlighting/highlight.js", themefiles.HighlightJs),
			newFileFromText("codehighlighting/highlight.css", themefiles.HighlightCss),

			// latest/preview
			newFileFromText("latest.js", themefiles.LatestJs),
			newFileFromText("jquery.tmpl.js", themefiles.JqueryTempl),

			// lazy-loading
			newFileFromText("lazysizes.js", themefiles.LazySizesJs),

			// global
			newFileFromText("site.js", themefiles.SiteJs),
		},
	}

}
