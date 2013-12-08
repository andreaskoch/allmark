// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdowntohtml

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/audio"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/markdown"
)

type Converter struct {
	logger logger.Logger
}

func New(logger logger.Logger) (*Converter, error) {
	return &Converter{
		logger: logger,
	}, nil
}

func (converter *Converter) Convert(item *model.Item) (convertedContent string, conversionError error) {

	converter.logger.Debug("Converting item %q.", item)

	content := item.Content

	// markdown extension: audio
	content, audioConversionError := audio.Convert(content, item.Files())
	if audioConversionError != nil {
		converter.logger.Warn("Error while converting audio extensions. Error: %s", audioConversionError)
	}

	// markdown to html
	content = markdown.Convert(content)

	return content, nil
}
