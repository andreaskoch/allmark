// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package preprocessor

import (
	"github.com/andreaskoch/allmark/common/paths"
	"github.com/andreaskoch/allmark/model"
	"github.com/andreaskoch/allmark/services/converter/markdowntohtml/util"
	"fmt"
	"mime"
	"regexp"
	"strings"
)

var (
	// audio: [*description text*](*a link to an audio file*)
	audioMarkdownExtensionPattern = regexp.MustCompile(`audio: \[([^\]]+)\]\(([^)]+)\)`)
)

func newAudioExtension(pathProvider paths.Pather, files []*model.File) *audioExtension {
	return &audioExtension{
		pathProvider: pathProvider,
		files:        files,
	}
}

type audioExtension struct {
	pathProvider paths.Pather
	files        []*model.File
}

func (converter *audioExtension) Convert(markdown string) (convertedContent string, converterError error) {

	convertedContent = markdown

	for _, match := range audioMarkdownExtensionPattern.FindAllStringSubmatch(convertedContent, -1) {

		if len(match) != 3 {
			continue
		}

		// parameters
		originalText := strings.TrimSpace(match[0])
		title := strings.TrimSpace(match[1])
		path := strings.TrimSpace(match[2])

		// get the code
		renderedCode := converter.getAudioCode(title, path)

		// replace markdown
		convertedContent = strings.Replace(convertedContent, originalText, renderedCode, 1)

	}

	return convertedContent, nil
}

func (converter *audioExtension) getMatchingFile(path string) *model.File {
	for _, file := range converter.files {
		if file.Route().IsMatch(path) && model.IsAudioFile(file) {
			return file
		}
	}

	return nil
}

func (converter *audioExtension) getAudioCode(title, path string) string {

	fallback := util.GetHtmlLinkCode(title, path)

	// internal audio file
	if util.IsInternalLink(path) {

		if audioFile := converter.getMatchingFile(path); audioFile != nil {

			if mimeType, err := model.GetMimeType(audioFile); err == nil {
				filepath := converter.pathProvider.Path(audioFile.Route().Value())
				return getAudioFileLink(title, filepath, mimeType)
			}

		}

	} else {

		// external audio file
		if isAudioFile, mimeType := isAudioFileLink(path); isAudioFile {
			return getAudioFileLink(title, path, mimeType)
		}

	}

	// fallback
	return fallback
}

func getAudioFileLink(title, link, mimeType string) string {

	var code string
	if title != "" {
		code += fmt.Sprintf("**%s**\n\n", title)
	}

	code += "<audio controls>"
	code += fmt.Sprintf("<source src=\"%s\" type=\"%s\">", link, mimeType)
	code += "</audio>"

	return code
}

func isAudioFileLink(link string) (isAudioFile bool, mimeType string) {

	// abort if the link does not contain a dot
	if !strings.Contains(link, ".") {
		return false, ""
	}

	normalizedLink := strings.ToLower(link)
	fileExtension := normalizedLink[strings.LastIndex(normalizedLink, "."):]
	mimeType = mime.TypeByExtension(fileExtension)

	switch fileExtension {
	case ".mp3", ".ogg", ".wav":
		return true, mimeType
	default:
		return false, ""
	}

	panic("Unreachable")
}
