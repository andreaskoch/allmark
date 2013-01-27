/*
	Package model defines the basic
	data structures of the docs engine.
*/
package model

import "time"

type Document struct {
	Path        string    // The documents folder
	Title       string    // The document title
	Description string    // A short description of the document content.
	Content     string    // The document content
	Language    string    // [optional] The ISO language code document (e.g. "en-GB", "de-DE")
	Date        time.Time // [optional] The date the document has been created
}
