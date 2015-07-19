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
	folderRoute := route.NewFromRequest(path)
	fullFolderRoute := route.Combine(r.base, folderRoute)

	// render the filesystem
	code := ""
	if strings.TrimSpace(title) != "" {
		code += fmt.Sprintf("\n**%s**\n\n", title)
	}

	if rootNode := r.files.GetNode(fullFolderRoute); rootNode != nil {

		// Render the childs of the root node
		code += "- " + r.renderFileNode(rootNode, 1)

	}

	return code
}

func (r *FileTreeRenderer) renderFileNode(node *FileNode, indentation int) string {

	html := ""

	if file := node.Value(); file != nil {
		fileRoute := file.Route()
		filepath := r.pathProvider.Path(fileRoute.Value())
		html = fmt.Sprintf("[%s](%s)\n", fileRoute.LastComponentName(), filepath)
	} else {
		html = node.Name() + "\n"
	}

	if childs := node.Childs(); len(childs) > 0 {

		for _, child := range childs {
			html += fmt.Sprintf("%s- %s", getIndentation(indentation, "\t"), r.renderFileNode(child, indentation+1))
		}

	}

	return html
}

func getIndentation(depth int, character string) string {
	indentation := ""
	for level := 1; level <= depth; level++ {
		indentation += character
	}
	return indentation
}

func convertFilesToTree(files []*model.File) *FileTree {

	tree := newTree()

	for _, file := range files {
		tree.Insert(file)
	}

	return tree
}
