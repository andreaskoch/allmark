package renderer

import (
	"andyk/docs/indexer"
	"andyk/docs/parser"
	"andyk/docs/templates"
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func RenderItem(item indexer.Item) {

	parsedItem, err := parser.ParseItem(item)
	if err != nil {
		log.Printf("Could not parse item \"%v\". Error: %v", item.Path, err)
		return
	}

	renderedItemFilePath := getRenderedItemPath(item)

	switch parsedItem.MetaData.ItemType {
	case parser.DocumentItemType:
		{
			file, err := os.Create(renderedItemFilePath)
			if err != nil {
				panic(err)
			}
			writer := bufio.NewWriter(file)

			defer func() {
				writer.Flush()
				file.Close()
			}()

			document := getDocument(parsedItem)
			template := template.New(parser.DocumentItemType)
			template.Parse(templates.DocumentTemplate)
			template.Execute(writer, document)
		}
	}
}

type Document struct {
	Title       string
	Description string
	Content     string
}

func getDocument(parsedItem parser.ParsedItem) Document {
	return Document{
		Title:       parsedItem.GetElementValue("title"),
		Description: parsedItem.GetElementValue("description"),
		Content:     parsedItem.GetElementValue("content"),
	}
}

// Get the filepath of the rendered repository item
func getRenderedItemPath(item indexer.Item) string {
	itemDirectory := filepath.Dir(item.Path)
	itemName := strings.Replace(filepath.Base(item.Path), filepath.Ext(item.Path), "", 1)

	renderedFilePath := filepath.Join(itemDirectory, itemName+".html")
	return renderedFilePath
}
