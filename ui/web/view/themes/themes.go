// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themes

import "github.com/andreaskoch/allmark/themes/themefiles"

var defaultTheme *Theme

// initialize themes
func init() {

	defaultTheme = &Theme{
		Name: "default",
		Files: []*ThemeFile{

			// styles
			newFileFromText("screen.css", themefiles.ScreenCss),
			newFileFromText("print.css", themefiles.PrintCss),

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

			// code highlighting
			newFileFromBase64("codehighlighting/highlight.js", themefiles.HighlightJs),
			newFileFromBase64("codehighlighting/highlight.css", themefiles.HighlightCss),
		},
	}

}
