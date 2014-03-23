// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package files

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
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

		// search for files-extension code
		found, matches := pattern.IsMatch(convertedContent, filesPattern)
		if !found || (found && len(matches) != 3) {
			break // abort. no (more) files-extension code found
		}

		// extract the parameters from the pattern matches
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// normalize the path with the current path provider
		path = converter.pathProvider.Path(path)

		// create the base route from the path
		baseRoute, err := route.NewFromRequest(path)
		if err != nil {
			// abort. an error occured.
			return markdown, fmt.Errorf("Could not create a route from the path %q. Error: %s", path, err)
		}

		// get all matching files
		matchingFiles := getMatchingFiles(baseRoute, converter.files)
		pathsOfMatchingFiles := getFilePaths(matchingFiles, converter.pathProvider)

		// create a file system from the file paths
		rootFolderTitle := baseRoute.FolderName()
		rootFolder := NewRootFilesystemEntry(rootFolderTitle)
		rootFolder = getFileSystemFromLinks(rootFolder, pathsOfMatchingFiles)

		// render the filesystem
		fileLinksCode := fmt.Sprintf(`<section class="filelinks"><h1>%s</h1>`, title)
		fileLinksCode += renderFilesystemEntry(rootFolder, 0)
		fileLinksCode += `</section>`

		// replace markdown with link list
		convertedContent = strings.Replace(convertedContent, originalText, fileLinksCode, 1)

	}

	return convertedContent, nil
}

// Get all files that match (are childs of) the supplied path.
func getMatchingFiles(baseRoute *route.Route, files []*model.File) []*model.File {

	matchingFiles := make([]*model.File, 0)
	for _, file := range files {

		// check if the file is a child of the supplied path
		if !file.Route().IsChildOf(baseRoute) {
			continue
		}

		matchingFiles = append(matchingFiles, file)
	}

	return matchingFiles
}

// Get the files paths for the supplied File models.
func getFilePaths(files []*model.File, pathProvider paths.Pather) []string {

	numberOfFiles := len(files)
	fileLinks := make([]string, numberOfFiles, numberOfFiles)

	for index, file := range files {
		filePath := pathProvider.Path(file.Route().Value())
		fileLinks[index] = filePath
	}

	return fileLinks
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
		components = trimSlice(components)

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
