package parser

import (
	"errors"
	"fmt"
	"github.com/andreaskoch/docs/repository"
	"github.com/andreaskoch/docs/util"
	"github.com/russross/blackfriday"
	"os"
	"regexp"
	"strings"
)

var (
	// Lines which contain nothing but white space characters
	// or no characters at all.	
	EmptyLinePattern = regexp.MustCompile(`^\s*$`)

	// Lines which a start with a hash, followed by zero or more
	// white space characters, followed by text.
	TitlePattern = regexp.MustCompile(`\s*#\s*(\w.+)`)

	// Lines which start with text
	DescriptionPattern = regexp.MustCompile(`^\w.+`)

	// Lines which nothing but dashes
	HorizontalRulePattern = regexp.MustCompile(`^-{2,}`)

	// Lines with a "key: value" syntax
	MetaDataPattern = regexp.MustCompile(`^(\w+):\s*(\w.+)$`)
)

func Parse(item *repository.Item) (*repository.Item, error) {

	// open the file
	file, err := os.Open(item.SourcePath())
	if err != nil {
		return item, err
	}

	defer file.Close()

	// get the lines
	lines := util.GetLines(file)

	switch item.Type {
	case repository.DocumentItemType, repository.CollectionItemType, repository.RepositoryItemType:
		{
			return parseDocumentLikeItem(item, lines), nil
		}
	case repository.MessageItemType:
		{
			return parseMessage(item, lines), nil
		}
	}

	return item, errors.New(fmt.Sprintf("Items of type \"%v\" cannot be parsed.", item.Type))
}

func getMatchingValue(lines []string, matchPattern *regexp.Regexp) (string, []string) {

	// In order to be the "matching value" the line must
	// either be empty or match the supplied pattern.
	for lineNumber, line := range lines {

		lineMatchesTitlePattern, matches := util.IsMatch(line, matchPattern)
		if lineMatchesTitlePattern {
			nextLine := getNextLinenumber(lineNumber, lines)
			return util.GetLastElement(matches), lines[nextLine:]
		}

		lineIsEmpty := EmptyLinePattern.MatchString(line)
		if !lineIsEmpty {
			break
		}
	}

	return "", lines
}

func getContent(lines []string) string {

	// all remaining lines are the (raw markdown) content
	rawMarkdownContent := strings.TrimSpace(strings.Join(lines, "\n"))

	// html encode the markdown
	htmlEncodedContent := string(blackfriday.MarkdownBasic([]byte(rawMarkdownContent)))

	return htmlEncodedContent
}

func getNextLinenumber(lineNumber int, lines []string) int {
	nextLine := lineNumber + 1

	if nextLine <= len(lines) {
		return nextLine
	}

	return lineNumber
}
