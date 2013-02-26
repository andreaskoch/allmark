package renderer

import (
	"andyk/docs/indexer"
	"andyk/docs/parser"
	"andyk/docs/templates"
	"andyk/docs/util"
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
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
	LanguageTag string
}

func getDocument(parsedItem parser.ParsedItem) Document {
	return Document{
		Title:       parsedItem.GetElementValue("title"),
		Description: parsedItem.GetElementValue("description"),
		Content:     parsedItem.GetElementValue("content"),
		LanguageTag: getTwoLetterLanguageCode(parsedItem.MetaData.Language),
	}
}

// Get ISO 639-1 language code from a given language string (e.g. "en-US" => "en", "de-DE" => "de")
func getTwoLetterLanguageCode(languageString string) string {

	fallbackLangueCode := "en"
	if languageString == "" {
		// default value
		return fallbackLangueCode
	}

	// Check if the language string already matches
	// the ISO 639-1 language code pattern (e.g. "en", "de").
	iso6391TwoLetterLanguageCodePattern := regexp.MustCompile(`^[a-z]$`)
	if len(languageString) == 2 && iso6391TwoLetterLanguageCodePattern.MatchString(languageString) {
		return strings.ToLower(languageString)
	}

	// Check if the language string matches the
	// IETF language tag pattern (e.g. "en-US", "de-DE").
	ietfLanguageTagPattern := regexp.MustCompile(`^(\w\w)-\w{2,3}$`)
	matchesIETFPattern, matches := util.IsMatch(languageString, ietfLanguageTagPattern)
	if matchesIETFPattern {
		return matches[1]
	}

	// use fallback
	return fallbackLangueCode
}

// Get the filepath of the rendered repository item
func getRenderedItemPath(item indexer.Item) string {
	itemDirectory := filepath.Dir(item.Path)
	itemName := strings.Replace(filepath.Base(item.Path), filepath.Ext(item.Path), "", 1)

	renderedFilePath := filepath.Join(itemDirectory, itemName+".html")
	return renderedFilePath
}
