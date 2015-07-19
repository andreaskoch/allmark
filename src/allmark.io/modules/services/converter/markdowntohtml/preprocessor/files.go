// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package preprocessor

import (
	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/pattern"
	"allmark.io/modules/common/route"
	"allmark.io/modules/model"
	"allmark.io/modules/services/converter/markdowntohtml/filetreerenderer"
	"regexp"
	"strings"
)

var (
	// files: [*description text*](*folder path*)
	filesMarkdownExtensionPattern = regexp.MustCompile(`files: \[([^\]]+)\]\(([^)]+)\)`)
)

func newFilesExtension(pathProvider paths.Pather, baseRoute route.Route, files []*model.File) *filesExtension {
	return &filesExtension{
		pathProvider:     pathProvider,
		base:             baseRoute,
		fileTreeRenderer: filetreerenderer.New(pathProvider, baseRoute, files),
	}
}

type filesExtension struct {
	pathProvider     paths.Pather
	base             route.Route
	fileTreeRenderer *filetreerenderer.FileTreeRenderer
}

func (converter *filesExtension) Convert(markdown string) (convertedContent string, converterError error) {

	convertedContent = markdown

	for {

		// search for files-extension code
		found, matches := pattern.IsMatch(convertedContent, filesMarkdownExtensionPattern)
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
		renderedCode := converter.fileTreeRenderer.Render(title, "filelinks", path)

		// replace markdown with link list
		convertedContent = strings.Replace(convertedContent, originalText, renderedCode, 1)

	}

	return convertedContent, nil
}
