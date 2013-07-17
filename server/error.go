// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"net/http"
)

var error404Handler func(w http.ResponseWriter, r *http.Request)
