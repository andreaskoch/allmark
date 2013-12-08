// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package audio

import (
	"fmt"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/converter/markdowntohtml/pattern"
	"mime"
	"regexp"
	"strings"
)

var (
	// audio: [*description text*](*a link to an audio file*)
	audioPattern = regexp.MustCompile(`audio: \[([^\]]+)\]\(([^)]+)\)`)
)

func Convert(markdown string, files []*model.File) string {

	for {

		found, matches := pattern.IsMatch(markdown, audioPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// get the code
		renderedCode := getAudioCode(title, path, files)

		// replace markdown with link list
		markdown = strings.Replace(markdown, originalText, renderedCode, 1)

	}

	return markdown
}

func getAudioCode(title, path string, files []*model.File) string {

	// html5 audio file
	if isAudioFile, mimeType := isAudioFileLink(path); isAudioFile {
		return getAudioFileLink(title, path, mimeType)
	}

	// fallback
	return fmt.Sprintf(`<a href="%s" target="_blank" title="%s">%s</a>`, path, title, title)
}

func isAudioFileLink(link string) (isAudioFile bool, mimeType string) {
	normalizedLink := strings.ToLower(link)
	fileExtension := normalizedLink[strings.LastIndex(normalizedLink, "."):]
	mimeType = mime.TypeByExtension(fileExtension)

	switch fileExtension {
	case ".mp3", ".ogg":
		return true, mimeType
	default:
		return false, ""
	}

	panic("Unreachable")
}

func getAudioFileLink(title, link, mimetype string) string {
	return fmt.Sprintf(`<section class="audio audio-file">
		<h1><a href="%s" target="_blank" title="%s">%s</a></h1>
		<audio controls>
			<source src="%s" type="%s">
		</audio>
	</section>`, link, title, title, link, mimetype)
}
