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

	// create a unique temp directory
	baseDirHash := util.GetHash(basePath)
	tempDir := filepath.Join(os.TempDir(), baseDirHash)
	if useTempDir {
		util.CreateDirectory(tempDir)
	}

	return &Provider{
		basePath:   basePath,
		tempDir:    tempDir,
		useTempDir: useTempDir,
	}
}

type Provider struct {
	basePath   string
	tempDir    string
	useTempDir bool
}

func (provider *Provider) New(basePath string) *Provider {
	return NewProvider(basePath, provider.UseTempDir())
}

func (provider *Provider) BasePath() string {
	return provider.basePath
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
		return util.EncodeUrl(provider.getItemRoute(pather))
	case PatherTypeFile, PatherTypeIndex:
		return util.EncodeUrl(provider.getFileRoute(pather))
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

func (provider *Provider) GetRenderTargetPath(pather Pather) string {

	itemDirectoryRelative := provider.getRelativePath(pather.Directory())
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

func (provider *Provider) getItemRoute(pather Pather) string {
	absoluteTargetFilesystemPath := provider.GetRenderTargetPath(pather)
	itemRoute := provider.getRouteFromFilepath(absoluteTargetFilesystemPath)

	return itemRoute
}

func (provider *Provider) getFileRoute(pather Pather) string {
	absoluteFilepath := provider.getRouteFromFilepath(pather.Path())
	return provider.getRouteFromFilepath(absoluteFilepath)
}

func (provider *Provider) getRouteFromFilepath(path string) string {
	relativeFilepath := provider.getRelativePath(path)

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

func (provider *Provider) getRelativePath(filepath string) string {
	return strings.Replace(filepath, provider.basePath, "", 1)
}
