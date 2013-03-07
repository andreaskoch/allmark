package viewmodel

type Repository struct {
	Title       string
	Description string
	Content     string
	LanguageTag string
	Entries     []RepositoryEntry
}

type RepositoryEntry struct {
	Title       string
	Description string
	Path        string
}
