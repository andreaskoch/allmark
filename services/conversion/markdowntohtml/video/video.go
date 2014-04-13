// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package video

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
	// video: [*description text*](*a link to a youtube video or to a video file*)
	videoPattern = regexp.MustCompile(`video: \[([^\]]+)\]\(([^)]+)\)`)

	// youtube video link pattern
	youTubeVideoPattern = regexp.MustCompile(`http[s]?://www\.youtube\.com/watch\?v=([^&]+)`)

	// vimeo video link pattern
	vimeoVideoPattern = regexp.MustCompile(`http[s]?://vimeo\.com/([\d]+)`)
)

func New(pathProvider paths.Pather, files []*model.File) *VideoExtension {
	return &VideoExtension{
		pathProvider: pathProvider,
		files:        files,
	}
}

type VideoExtension struct {
	pathProvider paths.Pather
	files        []*model.File
}

func (converter *VideoExtension) Convert(markdown string) (convertedContent string, conversionError error) {

	convertedContent = markdown

	for {

		found, matches := pattern.IsMatch(convertedContent, videoPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// get the code
		renderedCode := converter.getVideoCode(title, path)

		// replace markdown
		convertedContent = strings.Replace(convertedContent, originalText, renderedCode, 1)

	}

	return convertedContent, nil
}

func (converter *VideoExtension) getMatchingFile(path string) *model.File {
	for _, file := range converter.files {
		if file.Route().IsMatch(path) && util.IsVideoFile(file) {
			return file
		}
	}

	return nil
}

func (converter *VideoExtension) getVideoCode(title, path string) string {

	// internal video file
	if audioFile := converter.getMatchingFile(path); audioFile != nil {

		if mimeType, err := util.GetMimeType(audioFile); err == nil {
			filepath := converter.pathProvider.Path(audioFile.Route().Value())
			return renderVideoFileLink(title, filepath, mimeType)
		}

	}

	// external: youtube
	if isYouTube, videoId := isYouTubeLink(path); isYouTube {
		return renderYouTubeVideo(title, videoId)
	}

	// external: vimeo
	if isVimeo, videoId := isVimeoLink(path); isVimeo {
		return renderVimeoVideo(title, videoId)
	}

	// external: html5 video file
	if isVideoFile, mimeType := isVideoFileLink(path); isVideoFile {
		return renderVideoFileLink(title, path, mimeType)
	}

	// return the fallback handler
	return fmt.Sprintf(`<a href="%s" target="_blank" title="%s">%s</a>`, path, title, title)
}

func isYouTubeLink(link string) (isYouTubeLink bool, videoId string) {
	if found, matches := pattern.IsMatch(link, youTubeVideoPattern); found && len(matches) == 2 {
		return true, matches[1]
	}

	return false, ""
}

func renderYouTubeVideo(title, videoId string) string {
	return fmt.Sprintf(`<section class="video video-external video-youtube">
		<h1><a href="http://www.youtube.com/watch?v=%s" target="_blank" title="%s">%s</a></h1>
		<iframe width="560" height="315" src="http://www.youtube.com/embed/%s" frameborder="0" allowfullscreen></iframe>
	</section>`, videoId, title, title, videoId)
}

func isVimeoLink(link string) (isVimeoLink bool, videoId string) {
	if found, matches := pattern.IsMatch(link, vimeoVideoPattern); found && len(matches) == 2 {
		return true, matches[1]
	}

	return false, ""
}

func renderVimeoVideo(title, videoId string) string {
	return fmt.Sprintf(`<section class="video video-external video-vimeo">
		<h1><a href="https://vimeo.com/%s" target="_blank" title="%s">%s</a></h1>
		<iframe src="http://player.vimeo.com/video/%s" width="560" height="315" frameborder="0" webkitAllowFullScreen mozallowfullscreen allowFullScreen></iframe>
	</section>`, videoId, title, title, videoId)
}

func isVideoFileLink(link string) (isVideoFile bool, mimeType string) {
	normalizedLink := strings.ToLower(link)
	fileExtension := normalizedLink[strings.LastIndex(normalizedLink, "."):]
	mimeType = mime.TypeByExtension(fileExtension)

	switch fileExtension {
	case ".mp4", ".ogg", ".ogv", ".webm", ".3gp":
		return true, mimeType
	default:
		return false, ""
	}

	panic("Unreachable")
}

func renderVideoFileLink(title, link, mimetype string) string {
	return fmt.Sprintf(`<section class="video video-file">
		<h1><a href="%s" target="_blank" title="%s">%s</a></h1>
		<video width="560" height="315" controls>
			<source src="%s" type="%s">
		</video>
	</section>`, link, title, title, link, mimetype)
}
