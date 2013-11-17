// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"regexp"
	"strings"
)

var (
	// files: [*description text*](*folder path*)
	fileLinksPattern = regexp.MustCompile(`files: \[([^\]]+)\]\(([^)]+)\)`)
)

func renderFileLinks(fileIndex *repository.FileIndex, pathProvider *path.Provider, markdown string) string {

	for {

		found, matches := util.IsMatch(markdown, fileLinksPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// create link list code
		files := fileIndex.FilesByPath(path, allFiles)
		filePaths := getFilePaths(title, files, pathProvider)
		rootFolderTitle := getLastPathComponent(path)
		rootFolder := NewRootFilesystemEntry(rootFolderTitle)
		rootFolder = getFileSystemFromLinks(rootFolder, filePaths)

		fileLinksCode := fmt.Sprintf(`<section class="filelinks"><h1>%s</h1>`, title)
		fileLinksCode += renderFilesystemEntry(rootFolder, 0)
		fileLinksCode += `</section>`

		// replace markdown with link list
		markdown = strings.Replace(markdown, originalText, fileLinksCode, 1)

	}

	return markdown
}

func getLastPathComponent(path string) string {
	if !strings.Contains(path, "/") {
		return path
	}

	components := strings.Split(path, "/")
	return components[len(components)-1]
}

func renderFilesystemEntry(fsEntry *FileSystemEntry, level int) string {

	html := ""

	// don't print the "root-folder" name (-> "files")
	if level > 0 {

		// assemble the file or folder link
		href := fsEntry.Path()
		title := fsEntry.Name()
		name := fsEntry.Name()

		if fsEntry.IsDirectory() {
			html = fmt.Sprintf(`%s`, name)
		} else {
			html = fmt.Sprintf(`<a href="%s" title="%s">%s</a>`, href, title, name)
		}
	}

	// recurse if the filesystem entry has childs
	if len(fsEntry.Childs) > 0 {

		html += "<ul>\n"
		for _, child := range fsEntry.Childs {
			html += fmt.Sprintf("<li>%s</li>\n", renderFilesystemEntry(child, level+1))
		}
		html += "</ul>\n"

	}

	return html
}

type FileSystemEntry struct {
	name   string
	parent *FileSystemEntry
	Childs []*FileSystemEntry
}

func NewRootFilesystemEntry(name string) *FileSystemEntry {
	return &FileSystemEntry{
		name:   name,
		Childs: make([]*FileSystemEntry, 0),
	}
}

func NewFilesystemEntry(parent *FileSystemEntry, name string) *FileSystemEntry {
	return &FileSystemEntry{
		name:   name,
		parent: parent,
		Childs: make([]*FileSystemEntry, 0),
	}
}

func (fsEntry *FileSystemEntry) Path() string {
	path := fsEntry.name

	if fsEntry.Parent() != nil {
		path = fsEntry.Parent().Path() + "/" + path
	}

	return path
}

func (fsEntry *FileSystemEntry) IsDirectory() bool {
	return len(fsEntry.Childs) > 0
}

func (fsEntry *FileSystemEntry) Parent() *FileSystemEntry {
	return fsEntry.parent
}

func (fsEntry *FileSystemEntry) Name() string {
	return util.DecodeUrl(fsEntry.name)
}

func (fsEntry *FileSystemEntry) GetChild(name string) *FileSystemEntry {
	for _, entry := range fsEntry.Childs {
		if entry.name == name {
			return entry
		}
	}

	return nil
}

func getFileSystemFromLinks(root *FileSystemEntry, filePaths []string) *FileSystemEntry {

	for _, filePath := range filePaths {

		current := root

		// split the path into components
		components := strings.Split(filePath, "/")

		// strip "files" component
		if len(components) > 1 {
			components = components[1:]
		}

		// trim empty components
		components = util.TrimSlice(components)

		// get the childs for each component
		for _, component := range components {
			child := current.GetChild(component)
			if child == nil {
				current.Childs = append(current.Childs, NewFilesystemEntry(current, component))
				child = current.GetChild(component)
			}
			current = child
		}
	}

	return root
}

func getFilePaths(title string, files []*repository.File, pathProvider *path.Provider) []string {

	numberOfFiles := len(files)
	fileLinks := make([]string, numberOfFiles, numberOfFiles)

	for index, file := range files {
		filePath := pathProvider.GetWebRoute(file)
		fileLinks[index] = filePath
	}

	return fileLinks
}
