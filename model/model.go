/*
	Package model defines the basic
	data structures of the docs engine.
*/
package model

type Document struct {
	Path        string // The documents folder
	Title       string // The document title
	Description string // A short description of the document content.
	Content     string // The document content
	Language    string // [optional] The ISO language code document (e.g. "en-GB", "de-DE")
	Date        string // [optional] The date the document has been created
}
