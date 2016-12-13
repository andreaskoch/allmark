// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type OpenSearchDescription struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	FavIconURL  string `json:"favIconURL"`
	SearchURL   string `json:"searchURL"`
	Tags        string `json:"tags"`
}
