// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themes

var defaultTheme *Theme

func init() {

	defaultTheme = &Theme{
		Name: "default",
		Files: []*ThemeFile{

			// styles
			newFileFromText("screen.css", screenCss),
			newFileFromText("print.css", printCss),

			// javascript libraries
			newFileFromText("modernizr.js", modernizr),
			newFileFromText("jquery.js", jqueryJs),

			// web socket auto update
			newFileFromText("autoupdate.js", autoupdateJs),

			// pdf preview
			newFileFromText("pdf.js", pdfJs),
			newFileFromText("pdf-preview.js", pdfPreviewJs),

			// favicon
			newFileFromBase64("favicon.ico", faviconIco),

			// presentations
			newFileFromText("deck.js", deckJs),
			newFileFromText("deck.css", deckCss),

			// code highlighting
			newFileFromBase64("codehighlighting/highlight.js", highlightJs),
			newFileFromBase64("codehighlighting/highlight.css", highlightCss),
		},
	}

}
