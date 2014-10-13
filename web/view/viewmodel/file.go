// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

import (
	"time"
)

type File struct {
	Parent string `json:"parent"`
	Path   string `json:"path"`
	Route  string `json:"route"`
	Name   string `json:"name"`

	Hash         string    `json:"hash"`
	LastModified time.Time `json:"lastModified"`
	MimeType     string    `json:"mimeType"`
}
