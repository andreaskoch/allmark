package server

import (
	"strings"
)

func CombineUrlComponents(baseUrl string, components ...string) string {
	url := StripTrailingDirectorySeperator(baseUrl)

	for _, component := range components {
		url += UrlDirectorySeperator + StripTrailingDirectorySeperator(component)
	}

	return url
}

func StripTrailingDirectorySeperator(urlComponent string) string {

	url := urlComponent

	for strings.LastIndex(urlComponent, UrlDirectorySeperator)+1 == len(url) {
		url = strings.TrimRight(url, UrlDirectorySeperator)
	}

	return url
}
