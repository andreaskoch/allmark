// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package files

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/pattern"
	"regexp"
	"strings"
)

var (
	// files: [*description text*](*folder path*)
	filesPattern = regexp.MustCompile(`files: \[([^\]]+)\]\(([^)]+)\)`)
)

func New(pathProvider paths.Pather, files []*model.File) *FilesExtension {
	return &FilesExtension{
		pathProvider: pathProvider,
		files:        files,
	}
}

type FilesExtension struct {
	pathProvider paths.Pather
	files        []*model.File
}

func (converter *FilesExtension) Convert(markdown string) (convertedContent string, conversionError error) {

	convertedContent = markdown

	for {

		found, matches := pattern.IsMatch(convertedContent, filesPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// fix the path
		path = converter.pathProvider.Path(path)

		// create link list code
		files := converter.getFilesByPath(path)
		filePaths := converter.getFilePaths(title, files)
		rootFolderTitle := getLastPathComponent(path)
		rootFolder := NewRootFilesystemEntry(rootFolderTitle)
		rootFolder = getFileSystemFromLinks(rootFolder, filePaths)

		fileLinksCode := fmt.Sprintf(`<section class="filelinks"><h1>%s</h1>`, title)
		fileLinksCode += renderFilesystemEntry(rootFolder, 0)
		fileLinksCode += `</section>`

		// replace markdown with link list
		convertedContent = strings.Replace(convertedContent, originalText, fileLinksCode, 1)

	}

	return convertedContent, nil
}

func (converter *FilesExtension) getFilePaths(title string, files []*model.File) []string {

	numberOfFiles := len(files)
	fileLinks := make([]string, numberOfFiles, numberOfFiles)

	for index, file := range files {
		filePath := converter.pathProvider.Path(file.Route().Value())
		fileLinks[index] = filePath
	}

	return fileLinks
}

func (converter *FilesExtension) getFilesByPath(path string) []*File {

	if strings.Index(path, FilesDirectoryName) == 0 {
		path = path[len(FilesDirectoryName):]
	}

	matchingFiles := make([]*File, 0)

	for _, file := range converter.files {

		filePath := file.Path()
		indexPath := fileIndex.Path()

		if strings.Index(filePath, indexPath) != 0 {
			continue
		}

		relativeFilePath := filePath[len(indexPath):]
		if strings.HasPrefix(relativeFilePath, path) {
			matchingFiles = append(matchingFiles, file)
		}
	}

	return matchingFiles
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
