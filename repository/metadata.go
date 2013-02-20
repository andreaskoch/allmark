package repository

import (
	"strings"
	"time"
)

type MetaData struct {
	Language string
	Date     time.Time
	Tags     []string
}

func (metaData *MetaData) String() string {
	s := "Language: " + metaData.Language
	s += "\nDate: " + metaData.Date.String()
	s += "\nTags: " + strings.Join(metaData.Tags, ", ")

	return s
}
