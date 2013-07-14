// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"net/http"
)

var robotsTxtHandler = func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `User-agent: *`)
	fmt.Fprintln(w, fmt.Sprintf(`Sitemap: http://%s/sitemap.xml`, r.Host))
}
