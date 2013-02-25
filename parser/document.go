package parser

func ParseDocument(lines []string, metaData MetaData) (item ParsedItem, err error) {

	// assign meta data
	item.MetaData = metaData

	// title
	title, lines := getMatchingValue(lines, TitlePattern)
	item.AddElement("title", title)

	// description
	description, lines := getMatchingValue(lines, DescriptionPattern)
	item.AddElement("description", description)

	// content
	item.AddElement("content", getContent(lines))

	return item, nil
}
