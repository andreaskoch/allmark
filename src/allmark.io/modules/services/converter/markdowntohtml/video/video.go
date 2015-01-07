// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package video

import (
	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/pattern"
	"allmark.io/modules/model"
	"allmark.io/modules/services/converter/markdowntohtml/util"
	"fmt"
	"mime"
	"regexp"
	"strings"
)

var (
	// video: [*description text*](*a link to a youtube video or to a video file*)
	markdownPattern = regexp.MustCompile(`video: \[([^\]]+)\]\(([^)]+)\)`)

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

func (converter *VideoExtension) Convert(markdown string) (convertedContent string, converterError error) {

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

	fallback := util.GetHtmlLinkCode(title, path)

	// internal video file
	if util.IsInternalLink(path) {

		if videoFile := converter.getMatchingFile(path); videoFile != nil {

			if mimeType, err := util.GetMimeType(videoFile); err == nil {
				filepath := converter.pathProvider.Path(videoFile.Route().Value())
				return renderVideoFileLink(title, filepath, mimeType)
			}

		}

	} else {

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

	}

	// return the fallback handler
	return fallback
}

func isYouTubeLink(link string) (isYouTubeLink bool, videoId string) {
	if found, matches := pattern.IsMatch(link, youTubeVideoPattern); found && len(matches) == 2 {
		return true, matches[1]
	}

	return false, ""
}

func renderYouTubeVideo(title, videoId string) string {
	return fmt.Sprintf(`<section class="video video-external video-youtube">
		<header><a href="https://www.youtube.com/watch?v=%s" target="_blank" title="%s">%s</a></header>
		<iframe width="560" height="315" src="https://www.youtube.com/embed/%s" frameborder="0" allowfullscreen></iframe>
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
		<header><a href="https://vimeo.com/%s" target="_blank" title="%s">%s</a></header>
		<iframe src="https://player.vimeo.com/video/%s" width="560" height="315" frameborder="0" webkitAllowFullScreen mozallowfullscreen allowFullScreen></iframe>
	</section>`, videoId, title, title, videoId)
}

func isVideoFileLink(link string) (isVideoFile bool, mimeType string) {

	// abort if the link does not contain a dot
	if !strings.Contains(link, ".") {
		return false, ""
	}

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
		<header><a href="%s" target="_blank" title="%s">%s</a></header>
		<video width="560" height="315" controls src="%s" type="%s"></video>
	</section>`, link, title, title, link, mimetype)
}
