// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdowntohtml

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/paths"
	"allmark.io/modules/model"
	"allmark.io/modules/services/converter/markdowntohtml/postprocessor"
	"allmark.io/modules/services/converter/markdowntohtml/preprocessor"
	"allmark.io/modules/services/thumbnail"
	"github.com/russross/blackfriday"
)

type Converter struct {
	logger        logger.Logger
	preprocessor  *preprocessor.Preprocessor
	postprocessor *postprocessor.Postprocessor
}

func New(logger logger.Logger, thumbnailIndex *thumbnail.Index) *Converter {
	return &Converter{
		logger:        logger,
		preprocessor:  preprocessor.New(logger, thumbnailIndex),
		postprocessor: postprocessor.New(logger, thumbnailIndex),
	}
}

// Convert the supplied item with all paths relative to the supplied base route
func (converter *Converter) Convert(aliasResolver func(alias string) *model.Item, rootPathProvider, itemContentPathProvider paths.Pather, item *model.Item) (convertedContent string, converterError error) {

	converter.logger.Debug("Converting markdown for item %q.", item)

	// preprocessor
	rawMarkdownContent := item.Content
	preprocessedMarkdownContent, err := converter.preprocessor.Convert(aliasResolver, rootPathProvider, itemContentPathProvider, item.Route(), item.Files(), rawMarkdownContent)
	if err != nil {
		return "", err
	}

	// markdown to html
	htmlContent := markdownToHTML(preprocessedMarkdownContent)

	// postprocessing
	postProcessedHTMLContent, err := converter.postprocessor.Convert(rootPathProvider, itemContentPathProvider, item.Route(), item.Files(), htmlContent)
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
