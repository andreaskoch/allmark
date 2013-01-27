package indexer

import (
	"andyk/docs/model"
	"fmt"
	"io/ioutil"
	"time"
)

func Index() model.Document {

	files, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}

	fmt.Println(files)

	for i := 0; i < len(files); i++ {
		file := files[i]
		fmt.Println(file.Name())
	}

	var doc model.Document
	doc.Path = "Test"
	doc.Title = "Test"
	doc.Description = "Description"
	doc.Content = "Content"
	doc.Language = "en-US"
	doc.Date = time.Date(2013, 1, 13, 0, 0, 0, 0, time.UTC)

	return doc
}
