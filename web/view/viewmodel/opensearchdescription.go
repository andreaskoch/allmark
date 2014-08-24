// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type OpenSearchDescription struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	FavIconUrl  string `json:"favIconUrl"`
	SearchUrl   string `json:"searchUrl"`
	Tags        string `json:"tags"`
}
