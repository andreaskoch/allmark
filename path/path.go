// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// Filesystem directory seperator
	FilesystemDirectorySeperator = string(os.PathSeparator)

	// Url directory seperator
	UrlDirectorySeperator = "/"

	// Web server default file
	WebServerDefaultFilename = "index.html"
)

func NewProvider(basePath string) *Provider {
	return &Provider{
		basePath: basePath,
	}
}

type Provider struct {
	basePath string
}

func (provider *Provider) GetWebRoute(pather Pather) string {

	switch pathType := pather.PathType(); pathType {
	case PatherTypeItem:
		return provider.GetItemRoute(pather)
	case PatherTypeFile:
		return provider.GetFileRoute(pather)
	default:
		panic(fmt.Sprintf("Unknown pather type %q", pathType))
	}

	panic("Unreachable. Unknown pather type")
}

func (provider *Provider) GetFilepath(pather Pather) string {

	switch pathType := pather.PathType(); pathType {
	case PatherTypeItem:
		return GetRenderTargetPath(pather)
	case PatherTypeFile:
		return pather.Path()
	default:
		panic(fmt.Sprintf("Unknown pather type %q", pathType))
	}

	panic("Unreachable. Unknown pather type")
}

func (provider *Provider) GetRelativePath(filepath string) string {
	return strings.Replace(filepath, provider.basePath, "", 1)
}

func (provider *Provider) GetItemRoute(pather Pather) string {
	absoluteTargetFilesystemPath := GetRenderTargetPath(pather)
	return provider.GetRouteFromFilepath(absoluteTargetFilesystemPath)
}

func (provider *Provider) GetFileRoute(pather Pather) string {
	absoluteFilepath := provider.GetRouteFromFilepath(pather.Path())
	return provider.GetRouteFromFilepath(absoluteFilepath)
}

func (provider *Provider) GetRouteFromFilepath(path string) string {
	relativeFilepath := provider.GetRelativePath(path)

	// filepath to route
	route := filepath.ToSlash(relativeFilepath)

	// Trim leading slash
	route = StripLeadingUrlDirectorySeperator(route)

	return route
}

func GetRenderTargetPath(pather Pather) string {
	sourceItemPath := pather.Path()
	renderTargetPath := sourceItemPath[0:strings.LastIndex(sourceItemPath, FilesystemDirectorySeperator)] + FilesystemDirectorySeperator + WebServerDefaultFilename
	return renderTargetPath
}
