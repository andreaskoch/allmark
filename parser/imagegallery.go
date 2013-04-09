package parser

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"strings"
)

func renderImageGalleries(item *repository.Item, lines []string) []string {

	for lineNumber, text := range lines {

		if found, matches := util.IsMatch(text, ImageGalleryPattern); found && len(matches) == 3 {

			// parameters
			originalText := matches[0]
			descriptionText := matches[1]
			path := matches[2]

			// create image gallery code
			files := item.Files.GetFilesByPath(path)
			imageLinks := getImageLinks(item, files)
			imageGalleryCode := fmt.Sprintf(`<div class="imagegallery">
				<header>
					<span>%s</span>
				</header>
				%s
			</div>`, descriptionText, strings.Join(imageLinks, "\n"))

			// replace markdown with image gallery
			lines[lineNumber] = strings.Replace(text, originalText, imageGalleryCode, 1)
		}
	}

	return lines

}

func getImageLinks(item *repository.Item, files []*repository.File) []string {
	pathProvider := path.NewProvider(item.Directory())
	imagelinks := make([]string, len(files), len(files))
	for index, file := range files {
		imagelinks[index] = getImageLink(pathProvider, file)
	}
	return imagelinks
}

func getImageLink(pathProvider *path.Provider, file *repository.File) string {
	fileRoute := pathProvider.GetFileRoute(file)
	return fmt.Sprintf(`<img src="%s" />`, fileRoute)
}
