// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"net/http"
)

type Handler interface {
	Func() func(w http.ResponseWriter, r *http.Request)
}
