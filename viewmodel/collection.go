package viewmodel

type Collection struct {
	Title       string
	Description string
	Content     string
	LanguageTag string
	Entries     []CollectionEntry
}

type CollectionEntry struct {
	Title       string
	Description string
	Path        string
}
