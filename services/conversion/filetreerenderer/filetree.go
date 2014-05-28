// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filetreerenderer

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/tree/filetree"
	"github.com/andreaskoch/allmark2/model"
	"strings"
)

func New(pathProvider paths.Pather, baseRoute route.Route, files []*model.File) *FileTreeRenderer {
	return &FileTreeRenderer{
		pathProvider: pathProvider,
		base:         baseRoute,
		files:        convertFilesToTree(files),
	}
}

type FileTreeRenderer struct {
	pathProvider paths.Pather
	base         route.Route
	files        *filetree.FileTree
}

func (r *FileTreeRenderer) Render(title, cssClass, path string) string {

	// create the base route from the path
	folderRoute, err := route.NewFromRequest(path)
	if err != nil {
		// abort. an error occured.
		// todo: log error
		return ""
	}

	fullFolderRoute, err := route.Combine(&r.base, folderRoute)
	if err != nil {
		// abort. an error occured.
		// todo: log error
		return ""
	}

	// render the filesystem
	code := fmt.Sprintf(`<section class="%s">`, cssClass)
	if strings.TrimSpace(title) != "" {
		code += fmt.Sprintf("\n<header>%s</header>\n", title)
	}

	childs := r.files.GetChildFiles(fullFolderRoute)

	code += "<ul class=\"tree\">\n"

	if len(childs) > 1 {
		for _, child := range r.files.GetChildFiles(fullFolderRoute) {
			code += "<li>\n"
			code += r.renderFilesystemEntry(child.Route())
			code += "</li>\n"
		}
	} else {
		code += "<li>\n"
		code += r.renderFilesystemEntry(fullFolderRoute)
		code += "</li>\n"
	}

	code += "</ul>\n</section>"

	return code
}

func (r *FileTreeRenderer) renderFilesystemEntry(route *route.Route) string {

	filepath := r.pathProvider.Path(route.Value())
	html := fmt.Sprintf(`<a href="%s" title="%s">%s</a>`, filepath, route.Value(), route.LastComponentName())

	childs := r.files.GetChildFiles(route)

	if len(childs) > 0 {

		html += "<ul>\n"

		for _, child := range childs {
			html += fmt.Sprintf("<li>%s</li>\n", r.renderFilesystemEntry(child.Route()))
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
