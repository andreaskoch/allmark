// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package conversion

import (
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/model"
)

type Converter interface {

	// Convert the supplied item with all paths relative to the supplied base route
	Convert(pathProvider paths.Pather, item *model.Item, embedImages bool) (convertedContent string, conversionError error)
}
