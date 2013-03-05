package indexer

import (
	"time"
)

type MetaData struct {
	Language string
	Date     time.Time
	Tags     []string
}

func EmptyMetaData() MetaData {
	return MetaData{}
}
