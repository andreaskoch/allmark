// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"bufio"
	"fmt"
	"github.com/andreaskoch/allmark/converter/html"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"io"
	"os"
)

func (renderer *Renderer) RSS(writer io.Writer, host string) {

	// assemble the file path on disk
	targetFile := "rss.xml"
	renderTargetPath := renderer.pathProvider.GetRenderTargetPathForGiven(targetFile)

	// read the file from disk if it exists
	if util.FileExists(renderTargetPath) {
		file, err := os.Open(renderTargetPath)
		if err != nil {
			fmt.Fprintln(writer, fmt.Sprintf(`Unable to read %q.`, renderTargetPath))
			return
		}

		defer file.Close()

		data := make([]byte, 100)
		reader := bufio.NewReader(file)
		for {
			n, err := reader.Read(data)
			if n == 0 && err == io.EOF {
				return
			}

			if err != nil && err != io.EOF {
				fmt.Fprintln(writer, fmt.Sprintf(`Unable to read %q.`, renderTargetPath))
				return
			}

			writer.Write(data[:n])
		}

		return
	}

	// write to the file
	file, err := os.OpenFile(renderTargetPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Fprintln(writer, fmt.Sprintf(`Unable to open %q for writing. Err: %s`, renderTargetPath, err))
		return
	}

	defer file.Close()

	fileWriter := bufio.NewWriter(file)
	renderer.rss(fileWriter, host)

	fileWriter.Flush()

	// try to read from the cache again
	renderer.RSS(writer, host)
}

func (renderer *Renderer) rss(writer io.Writer, host string) {

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

		// render content for rss
		filePathProvider := renderer.root.FilePathProvider()
		httpRouteProvider := filePathProvider.NewHttpPathProvider(host)
		description := html.Convert(i, httpRouteProvider)

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
