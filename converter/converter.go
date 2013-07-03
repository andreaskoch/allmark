package converter

import (
	"github.com/andreaskoch/allmark/converter/html"
	"github.com/andreaskoch/allmark/repository"
)

func Convert(item *repository.Item) *repository.Item {
	return html.ToHtml(item)
}
