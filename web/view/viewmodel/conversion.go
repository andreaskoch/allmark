// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type ConversionModel struct {
	Base

	Content string `json:"content"`

	Files []File `json:"files"`
}
