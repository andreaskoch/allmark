// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type Analytics struct {
	Enabled         bool            `json:"enabled"`
	GoogleAnalytics GoogleAnalytics `json:"googleAnalytics"`
}

type GoogleAnalytics struct {
	Enabled    bool   `json:"enabled"`
	TrackingID string `json:"trackingId"`
}
