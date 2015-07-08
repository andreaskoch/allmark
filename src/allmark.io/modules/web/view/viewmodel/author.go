// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url"`

	GooglePlusHandle string `json:"googlePlusHandle"`
	TwitterHandle    string `json:"twitterHandle"`
	FacebookHandle   string `json:"facebookHandle"`
}
