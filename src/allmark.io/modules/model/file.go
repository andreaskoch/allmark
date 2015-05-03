// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"allmark.io/modules/dataaccess"
)

// A File represents a file ressource that is associated with an Item.
type File struct {
	dataaccess.File
}
