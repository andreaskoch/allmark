// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdowntohtml

import (
	"github.com/elWyatt/allmark/common/logger"
	"github.com/elWyatt/allmark/common/paths"
	"github.com/elWyatt/allmark/model"
	"github.com/elWyatt/allmark/services/converter/markdowntohtml/imageprovider"
	"github.com/elWyatt/allmark/services/converter/markdowntohtml/postprocessor"
	"github.com/elWyatt/allmark/services/converter/markdowntohtml/preprocessor"
	"github.com/russross/blackfriday"
)

// Converter converts markdown to HTML
type Converter struct {
	logger        logger.Logger
	preprocessor  *preprocessor.Preprocessor
	postprocessor *postprocessor.Postprocessor
}

// New creates a new Markdown-to-HTML converter instance.
func New(logger logger.Logger, imageProvider *imageprovider.ImageProvider) *Converter {
	return &Converter{
		logger:        logger,
		preprocessor:  preprocessor.New(logger, imageProvider),
		postprocessor: postprocessor.New(logger, imageProvider),
	}
}

// Convert the supplied item with all paths relative to the supplied base route
func (converter *Converter) Convert(aliasResolver func(alias string) *model.Item, pathProvider paths.Pather, item *model.Item) (convertedContent string, converterError error) {

	converter.logger.Debug("Converting markdown for item %q.", item)

	// preprocessor
	rawMarkdownContent := item.Content
	preprocessedMarkdownContent, err := converter.preprocessor.Convert(aliasResolver, pathProvider, item.Route(), item.Files(), rawMarkdownContent)
	if err != nil {
		return "", err
	}

	// markdown to html
	htmlContent := markdownToHTML(preprocessedMarkdownContent)

	// postprocessing
	postProcessedHTMLContent, err := converter.postprocessor.Convert(pathProvider, item.Route(), item.Files(), htmlContent)
	if err != nil {
		return "", err
	}

	return postProcessedHTMLContent, nil
}

func markdownToHTML(markdown string) (html string) {
	// set up the HTML renderer
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	// set up the parser
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	extensions |= blackfriday.EXTENSION_HARD_LINE_BREAK

	return string(blackfriday.Markdown([]byte(markdown), renderer, extensions))
}
