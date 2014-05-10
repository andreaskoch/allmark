// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package audio

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/pattern"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/util"
	"mime"
	"regexp"
	"strings"
)

var (
	// audio: [*description text*](*a link to an audio file*)
	markdownPattern = regexp.MustCompile(`audio: \[([^\]]+)\]\(([^)]+)\)`)
)

func New(pathProvider paths.Pather, files []*model.File) *AudioExtension {
	return &AudioExtension{
		pathProvider: pathProvider,
		files:        files,
	}
}

type AudioExtension struct {
	pathProvider paths.Pather
	files        []*model.File
}

func (converter *AudioExtension) Convert(markdown string) (convertedContent string, conversionError error) {

	convertedContent = markdown

	for {

		found, matches := pattern.IsMatch(convertedContent, markdownPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// get the code
		renderedCode := converter.getAudioCode(title, path)

		// replace markdown
		convertedContent = strings.Replace(convertedContent, originalText, renderedCode, 1)

	}

	return convertedContent, nil
}

func (converter *AudioExtension) getMatchingFile(path string) *model.File {
	for _, file := range converter.files {
		if file.Route().IsMatch(path) && util.IsAudioFile(file) {
			return file
		}
	}

	return nil
}

func (converter *AudioExtension) getAudioCode(title, path string) string {

	fallback := util.GetFallbackLink(title, path)

	// internal audio file
	if util.IsInternalLink(path) {

		if audioFile := converter.getMatchingFile(path); audioFile != nil {

			if mimeType, err := util.GetMimeType(audioFile); err == nil {
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
	return fmt.Sprintf(`<section class="audio audio-file">
		<h1><a href="%s" target="_blank" title="%s">%s</a></h1>
		<audio controls>
			<source src="%s" type="%s">
		</audio>
	</section>`, link, title, title, link, mimeType)
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
