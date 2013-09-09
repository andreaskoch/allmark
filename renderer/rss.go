// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"bufio"
	"fmt"
	"github.com/andreaskoch/allmark/converter/html"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"io"
	"os"
)

func cacheReponse(targetFile string, pathProvider *path.Provider, responseWriter ResponseWriter, host string, writer io.Writer) {

	// assemble the file path on disk
	renderTargetPath := pathProvider.GetRenderTargetPathForGiven(targetFile)

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
	responseWriter(fileWriter, host)

	fileWriter.Flush()

	// try to read from the cache again
	cacheReponse(targetFile, pathProvider, responseWriter, host, writer)
}

func (renderer *Renderer) RSS(writer io.Writer, host string) {

	targetFile := "rss.xml"
	pathProvider := renderer.pathProvider
	rssRenderer := func(writer io.Writer, host string) {
		rss(writer, host, renderer.root)
	}

	cacheReponse(targetFile, pathProvider, rssRenderer, host, writer)
}

func rss(writer io.Writer, host string, rootItem *repository.Item) {

	fmt.Fprintln(writer, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintln(writer, `<rss version="2.0">`)
	fmt.Fprintln(writer, `<channel>`)

	fmt.Fprintln(writer)
	fmt.Fprintln(writer, fmt.Sprintf(`<title><![CDATA[%s]]></title>`, rootItem.Title))
	fmt.Fprintln(writer, fmt.Sprintf(`<description><![CDATA[%s]]></description>`, rootItem.Description))
	fmt.Fprintln(writer, fmt.Sprintf(`<link>%s</link>`, getItemLocation(host, rootItem)))
	fmt.Fprintln(writer, fmt.Sprintf(`<pubData>%s</pubData>`, getItemDate(rootItem)))
	fmt.Fprintln(writer)

	// get all child items
	items := repository.GetAllChilds(rootItem)

	// sort the items by date and folder name
	repository.By(dateAndFolder).Sort(items)

	for _, i := range items {

		// render content for rss
		filePathProvider := rootItem.FilePathProvider()
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
