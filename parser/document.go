package parser

func (item *ParsedItem) ParseDocument(lines []string) *ParsedItem {

	// meta data
	metaData, metaDataLocation, lines := ParseMetaData(lines)
	if metaDataLocation.Found {
		item.MetaData = metaData
	}

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
