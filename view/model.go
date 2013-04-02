package view

type Model struct {
	Path        string
	Title       string
	Description string
	Content     string
	LanguageTag string
	Entries     []Model
}
