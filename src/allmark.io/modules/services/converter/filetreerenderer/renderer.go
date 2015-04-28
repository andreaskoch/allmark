// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filetreerenderer

import (
	"fmt"
	"strings"

	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/route"
	"allmark.io/modules/model"
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
	files        *FileTree
}

func (r *FileTreeRenderer) Render(title, cssClass, path string) string {

	// create the base route from the path
	folderRoute, err := route.NewFromRequest(path)
	if err != nil {
		// abort. an error occured.
		// todo: log error
		return ""
	}

	fullFolderRoute, err := route.Combine(r.base, folderRoute)
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

	if rootNode := r.files.GetNode(fullFolderRoute); rootNode != nil {

		// Render the childs of the root node
		code += "<ul class=\"tree\">\n"
		code += "<li>\n"
		code += r.renderFileNode(rootNode)
		code += "</li>\n"
		code += "</ul>\n</section>"

	}

	return code
}

func (r *FileTreeRenderer) renderFileNode(node *FileNode) string {

	html := ""

	if file := node.Value(); file != nil {
		fileRoute := file.Route()
		filepath := r.pathProvider.Path(fileRoute.Value())
		html = fmt.Sprintf(`<a href="%s" title="%s">%s</a>`, filepath, fileRoute.Value(), fileRoute.LastComponentName())
	} else {
		html = node.Name()
	}

	if childs := node.Childs(); len(childs) > 0 {

		html += "<ul>\n"

		for _, child := range childs {
			html += fmt.Sprintf("<li>%s</li>\n", r.renderFileNode(child))
		}

		html += "</ul>\n"
	}

	return html
}

func convertFilesToTree(files []*model.File) *FileTree {

	tree := newTree()

	for _, file := range files {
		tree.Insert(file)
	}

	return tree
}
