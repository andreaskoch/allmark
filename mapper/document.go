// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"fmt"
	"github.com/andreaskoch/allmark/converter"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
)

func createDocumentMapperFunc(pathProvider *path.Provider, targetFormat string) Mapper {

	return func(item *repository.Item) view.Model {

		parsed, err := converter.Convert(item, targetFormat)
		if err != nil {
			return view.Error(fmt.Sprintf("%s", err), pathProvider.GetWebRoute(item))
		}

		return view.Model{
			Path:        pathProvider.GetWebRoute(item),
			Title:       parsed.Title,
			Description: parsed.Description,
			Content:     parsed.ConvertedContent,
			LanguageTag: getTwoLetterLanguageCode(parsed.MetaData.Language),
		}
	}

}
