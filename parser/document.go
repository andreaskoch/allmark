package parser

type DocumentParser struct {
	Patterns DocumentStructure
}

func NewDocumentParser(documentStructure DocumentStructure) DocumentParser {
	return DocumentParser{
		Patterns: documentStructure,
	}
}

func (parser DocumentParser) Parse(lines []string, metaData MetaData) (item ParsedItem, err error) {

	// assign meta data
	item.MetaData = metaData

	// title
	title, lines := getMatchingValue(lines, parser.Patterns.Title, parser.Patterns.EmptyLine)
	item.AddElement("title", title)

	// description
	description, lines := getMatchingValue(lines, parser.Patterns.Description, parser.Patterns.EmptyLine)
	item.AddElement("description", description)

	// content
	item.AddElement("content", getContent(lines))

	return item, nil
}
