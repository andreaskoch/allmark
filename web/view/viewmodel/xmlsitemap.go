// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type XMLSitemap struct {
	Entries []XmlSitemapEntry
}

type XmlSitemapEntry struct {
	Loc          string                 `json:"loc"`
	LastModified string                 `json:"lastModified"`
	Images       []XmlSitemapEntryImage `json:"image:image"`
}

type XmlSitemapEntryImage struct {
	Loc string `json:"image:loc"`
}
