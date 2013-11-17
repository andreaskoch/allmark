// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"fmt"
	"github.com/andreaskoch/allmark/repository"
	"io"
	"path/filepath"
)

var (

	// sort the items by date and folder name
	dateAndFolder = func(item1, item2 *repository.Item) bool {

		if item1.MetaData.CreationDate.Equal(item2.MetaData.CreationDate) {
			// ascending by directory name
			return filepath.Base(item1.Directory()) < filepath.Base(item2.Directory())
		}

		// descending by date
		return item1.MetaData.CreationDate.After(item2.MetaData.CreationDate)
	}
)

func (renderer *Renderer) XMLSitemap(writer io.Writer, host string) {

	targetFile := "sitemap.xml"
	pathProvider := renderer.pathProvider
	rssRenderer := func(writer io.Writer, host string) {
		xmlsitemap(writer, host, renderer.root)
	}

	cacheReponse(targetFile, pathProvider, rssRenderer, host, writer)
}

func xmlsitemap(writer io.Writer, host string, rootItem *repository.Item) {

	fmt.Fprintln(writer, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintln(writer, `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

	// get all child items
	items := repository.GetAllChilds(rootItem, func(item *repository.Item) bool {
		isNotVirtual := !item.IsVirtual()
		return isNotVirtual
	})

	// sort the items by date and folder name
	dateAndFolder := func(item1, item2 *repository.Item) bool {

		if item1.MetaData.CreationDate.Equal(item2.MetaData.CreationDate) {
			// ascending by directory name
			return filepath.Base(item1.Directory()) < filepath.Base(item2.Directory())
		}

		// descending by date
		return item1.MetaData.CreationDate.After(item2.MetaData.CreationDate)
	}

	repository.By(dateAndFolder).Sort(items)

	for _, item := range items {
		fmt.Fprintln(writer, `<url>`)
		fmt.Fprintln(writer, fmt.Sprintf(`<loc>%s</loc>`, getItemLocation(host, item)))
		fmt.Fprintln(writer, fmt.Sprintf(`<lastmod>%s</lastmod>`, getItemDate(item)))
		fmt.Fprintln(writer, `</url>`)
	}

	fmt.Fprintln(writer, `</urlset>`)

}
