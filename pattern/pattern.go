package pattern

import (
	"regexp"
)

// DocumentStructure contains a set of regular expression
// for parsing documents.
type DocumentStructure struct {
	EmptyLine      regexp.Regexp
	Title          regexp.Regexp
	Description    regexp.Regexp
	HorizontalRule regexp.Regexp
	MetaData       regexp.Regexp
}

// NewDocumentStructure returns a new DocumentStructure
// for parsing documents.
func NewDocumentStructure() DocumentStructure {
	// Lines which contain nothing but white space characters
	// or no characters at all.
	emptyLineRegexp := regexp.MustCompile("^\\s*$")

	// Lines which a start with a hash, followed by zero or more
	// white space characters, followed by text.
	titleRegexp := regexp.MustCompile("\\s*#\\s*(\\w.+)")

	// Lines which start with text
	descriptionRegexp := regexp.MustCompile("^\\w.+")

	// Lines which nothing but dashes
	horizontalRuleRegexp := regexp.MustCompile("^-{2,}")

	// Lines with a "key: value" syntax
	metaDataRegexp := regexp.MustCompile("^(\\w+):\\s*(\\w.+)$")

	return DocumentStructure{
		EmptyLine:      *emptyLineRegexp,
		Title:          *titleRegexp,
		Description:    *descriptionRegexp,
		HorizontalRule: *horizontalRuleRegexp,
		MetaData:       *metaDataRegexp,
	}
}
