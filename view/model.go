// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package view

type Model struct {
	Level                int                   `json:"level"`
	AbsoluteRoute        string                `json:"absoluteRoute"`
	RelativeRoute        string                `json:"relativeRoute"`
	Title                string                `json:"title"`
	Description          string                `json:"description"`
	Content              string                `json:"content"`
	LanguageTag          string                `json:"languageTag"`
	Type                 string                `json:"type"`
	CreationDate         string                `json:"creationdate"`
	LastModifiedDate     string                `json:"lastmodifieddate"`
	Tags                 []*Tag                `json:"tags"`
	Childs               []*Model              `json:"childs"`
	ToplevelNavigation   *ToplevelNavigation   `json:"toplevelNavigation"`
	BreadcrumbNavigation *BreadcrumbNavigation `json:"breadcrumbNavigation"`
	TagCloud             *TagCloud             `json:"tagCloud"`
	Locations            []*GeoLocation        `json:"locations"`
	GeoLocation          *GeoLocation          `json:"geoLocation"`
}

func Error(title, content, relativPath, absolutePath string) *Model {
	return &Model{
		Level:         0,
		Title:         title,
		RelativeRoute: relativPath,
		AbsoluteRoute: absolutePath,
		Content:       content,
		Type:          "error",
	}
}
