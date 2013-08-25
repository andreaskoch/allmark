// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
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

func renderVideos(fileIndex *repository.FileIndex, pathProvider *path.Provider, markdown string) string {

	for {

		found, matches := util.IsMatch(markdown, videoPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// get a renderer
		renderer := getVideoRenderer(title, path, fileIndex, pathProvider)

		// execute the renderer
		renderedCode := renderer()

		// replace markdown with link list
		markdown = strings.Replace(markdown, originalText, renderedCode, 1)

	}

	return markdown
}

func getVideoRenderer(title, path string, fileIndex *repository.FileIndex, pathProvider *path.Provider) func() string {

	// youtube
	if isYouTube, videoId := isYouTubeLink(path); isYouTube {
		return func() string {
			return renderYouTubeVideo(title, videoId)
		}
	}

	// vimeo
	if isVimeo, videoId := isVimeoLink(path); isVimeo {
		return func() string {
			return renderVimeoVideo(title, videoId)
		}
	}

	// html5 video file
	if isVideoFile, mimeType := isVideoFileLink(path); isVideoFile {
		return func() string {
			return renderVideoFileLink(title, path, mimeType)
		}
	}

	// return the fallback handler
	return func() string {
		return fmt.Sprintf(`<a href="%s" target="_blank" title="%s">%s</a>`, path, title, title)
	}
}

func isYouTubeLink(link string) (isYouTubeLink bool, videoId string) {
	if found, matches := util.IsMatch(link, youTubeVideoPattern); found && len(matches) == 2 {
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
	if found, matches := util.IsMatch(link, vimeoVideoPattern); found && len(matches) == 2 {
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
