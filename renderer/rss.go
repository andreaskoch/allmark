// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"fmt"
	"github.com/andreaskoch/allmark/converter/html"
	"github.com/andreaskoch/allmark/repository"
	"io"
)

func (renderer *Renderer) RSS(writer io.Writer, host string) {

	fmt.Fprintln(writer, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintln(writer, `<rss version="2.0">`)
	fmt.Fprintln(writer, `<channel>`)

	fmt.Fprintln(writer)
	fmt.Fprintln(writer, fmt.Sprintf(`<title><![CDATA[%s]]></title>`, renderer.root.Title))
	fmt.Fprintln(writer, fmt.Sprintf(`<description><![CDATA[%s]]></description>`, renderer.root.Description))
	fmt.Fprintln(writer, fmt.Sprintf(`<link>%s</link>`, getItemLocation(host, renderer.root)))
	fmt.Fprintln(writer, fmt.Sprintf(`<pubData>%s</pubData>`, getItemDate(renderer.root)))
	fmt.Fprintln(writer)

	// get all child items
	items := repository.GetAllChilds(renderer.root)

	// sort the items by date and folder name
	repository.By(dateAndFolder).Sort(items)

	for _, i := range items {

		// skip the root
		if i == renderer.root {
			continue
		}

		// render content for rss
		description := html.Convert(i, renderer.root.FilePathProvider())

		fmt.Fprintln(writer, `<item>`)
		fmt.Fprintln(writer, fmt.Sprintf(`<title><![CDATA[%s]]></title>`, i.Title))
		fmt.Fprintln(writer, fmt.Sprintf(`<description><![CDATA[%s]]></description>`, description))
		fmt.Fprintln(writer, fmt.Sprintf(`<link>%s</link>`, getItemLocation(host, i)))
		fmt.Fprintln(writer, fmt.Sprintf(`<pubData>%s</pubData>`, getItemDate(i)))
		fmt.Fprintln(writer, `</item>`)
		fmt.Fprintln(writer)
	}

	fmt.Fprintln(writer, `</channel>`)
	fmt.Fprintln(writer, `</rss>`)

}

func getItemDate(item *repository.Item) string {
	return item.Date
}

func getItemLocation(host string, item *repository.Item) string {
	route := item.AbsoluteRoute
	location := fmt.Sprintf(`http://%s/%s`, host, route)
	return location
}
