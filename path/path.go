// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"fmt"
	"github.com/andreaskoch/allmark/util"
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

func NewProvider(basePath string, useTempDir bool) *Provider {
	return &Provider{
		basePath:   basePath,
		tempDir:    os.TempDir(),
		useTempDir: useTempDir,
	}
}

type Provider struct {
	basePath   string
	tempDir    string
	useTempDir bool
}

func (provider *Provider) UseTempDir() bool {
	return provider.useTempDir
}

func (provider *Provider) TempDir() string {
	return provider.tempDir
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
		return provider.GetRenderTargetPath(pather)
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
	absoluteTargetFilesystemPath := provider.GetRenderTargetPath(pather)
	return provider.GetRouteFromFilepath(absoluteTargetFilesystemPath)
}

func (provider *Provider) GetFileRoute(pather Pather) string {
	absoluteFilepath := provider.GetRouteFromFilepath(pather.Path())
	return provider.GetRouteFromFilepath(absoluteFilepath)
}

func (provider *Provider) GetRouteFromFilepath(path string) string {
	relativeFilepath := provider.GetRelativePath(path)

	// remove temp dir
	if provider.UseTempDir() {
		relativeFilepath = strings.TrimPrefix(relativeFilepath, provider.TempDir())
	}

	// filepath to route
	route := filepath.ToSlash(relativeFilepath)

	// Trim leading slash
	route = StripLeadingUrlDirectorySeperator(route)

	return route
}

func (provider *Provider) GetRenderTargetPath(pather Pather) string {

	itemDirectoryRelative := provider.GetRelativePath(pather.Directory())
	relativeRenderTargetPath := filepath.Join(itemDirectoryRelative, WebServerDefaultFilename)

	var renderTargetPath string
	if provider.UseTempDir() {

		renderTargetPath = filepath.Join(provider.TempDir(), relativeRenderTargetPath)

		// make sure the directory exists
		util.CreateDirectory(filepath.Dir(renderTargetPath))

	} else {

		renderTargetPath = filepath.Join(provider.basePath, relativeRenderTargetPath)

	}

	return renderTargetPath
}
