// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package files

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/tree/filetree"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/pattern"
	"regexp"
	"strings"
)

var (
	// files: [*description text*](*folder path*)
	markdownPattern = regexp.MustCompile(`files: \[([^\]]+)\]\(([^)]+)\)`)
)

func New(pathProvider paths.Pather, baseRoute route.Route, files []*model.File) *FilesExtension {
	return &FilesExtension{
		pathProvider: pathProvider,
		base:         baseRoute,
		files:        convertFilesToTree(files),
	}
}

type FilesExtension struct {
	pathProvider paths.Pather
	base         route.Route
	files        *filetree.FileTree
}

func (converter *FilesExtension) Convert(markdown string) (convertedContent string, conversionError error) {

	convertedContent = markdown

	for {

		// search for files-extension code
		found, matches := pattern.IsMatch(convertedContent, markdownPattern)
		if !found || (found && len(matches) != 3) {
			break // abort. no (more) files-extension code found
		}

		// extract the parameters from the pattern matches
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// normalize the path with the current path provider
		path = converter.pathProvider.Path(path)

		// get the code
		renderedCode := converter.getFileSystemCode(title, path)

		// replace markdown with link list
		convertedContent = strings.Replace(convertedContent, originalText, renderedCode, 1)

	}

	return convertedContent, nil
}

func (converter *FilesExtension) getFileSystemCode(title, path string) string {

	// create the base route from the path
	folderRoute, err := route.NewFromRequest(path)
	if err != nil {
		// abort. an error occured.
		// todo: log error
		return ""
	}

	fullFolderRoute, err := route.Combine(&converter.base, folderRoute)
	if err != nil {
		// abort. an error occured.
		// todo: log error
		return ""
	}

	// render the filesystem
	code := fmt.Sprintf(`<section class="filelinks"><h1>%s</h1>`, title)
	code += converter.renderFilesystemEntry(fullFolderRoute)
	code += `</section>`

	return code
}

func (converter *FilesExtension) renderFilesystemEntry(route *route.Route) string {

	filepath := converter.pathProvider.Path(route.Value())
	html := fmt.Sprintf(`<a href="%s" title="%s">%s</a>`, filepath, route.Value(), route.LastComponentName())

	childs := converter.files.GetChildFiles(route)

	if len(childs) > 0 {

		html += "<ul>\n"

		for _, child := range childs {
			html += fmt.Sprintf("<li>%s</li>\n", converter.renderFilesystemEntry(child.Route()))
		}

		html += "</ul>\n"
	}

	return html
}

func convertFilesToTree(files []*model.File) *filetree.FileTree {

	tree := filetree.New()

	for _, file := range files {
		tree.InsertFile(file)
	}

	return tree
}
