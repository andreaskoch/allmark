// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themes

import "github.com/andreaskoch/allmark2/ui/web/view/themes/themefiles"

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

			// pdf preview
			newFileFromText("pdf.js", themefiles.PdfJs),
			newFileFromText("pdf-preview.js", themefiles.PdfPreviewJs),

			// favicon
			newFileFromBase64("favicon.ico", themefiles.FaviconIco),

			// presentations
			newFileFromText("deck.js", themefiles.DeckJs),
			newFileFromText("deck.css", themefiles.DeckCss),
			newFileFromText("presentation.js", themefiles.PresentationJs),

			// auto-suggest
			newFileFromText("typeahead.js", themefiles.TypeAheadJs),
			newFileFromText("search.js", themefiles.SearchJs),

			// code highlighting
			newFileFromBase64("codehighlighting/highlight.js", themefiles.HighlightJs),
			newFileFromBase64("codehighlighting/highlight.css", themefiles.HighlightCss),
		},
	}

}
