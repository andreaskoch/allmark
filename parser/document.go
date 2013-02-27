package parser

func (item *ParsedItem) ParseDocument(lines []string, metaData MetaData) *ParsedItem {

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

	return item
}
