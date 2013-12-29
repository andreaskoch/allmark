// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/debughandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/itemhandler"
	"github.com/andreaskoch/allmark2/ui/web/server/index"
)

func NewItemHandler(logger logger.Logger, index *index.Index, patherFactory paths.PatherFactory, converter conversion.Converter) Handler {
	return itemhandler.New(logger, index, patherFactory, converter)
}

func NewDebugHandler(logger logger.Logger, index *index.Index) Handler {
	return debughandler.New(logger, index)
}
