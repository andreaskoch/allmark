// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdowntohtml

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/converter/markdowntohtml/audio"
	"github.com/andreaskoch/allmark2/services/converter/markdowntohtml/markdown"
)

type Converter struct {
	logger logger.Logger
}

func New(logger logger.Logger) (*Converter, error) {
	return &Converter{
		logger: logger,
	}, nil
}

func (converter *Converter) Convert(item *model.Item) (*model.Item, error) {
	converter.logger.Info("Converting item %q.", item)

	// markdown extensions
	item.Content = audio.Convert(item.Content, item.Files())

	// markdown
	item.Content = markdown.Convert(item.Content)

	return item, nil
}
