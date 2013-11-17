// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package message

import (
	"github.com/andreaskoch/allmark/repository"
	"strings"
	"time"
)

func Parse(item *repository.Item, lines []string, fallbackTitle string) (sucess bool, err error) {

	// title
	item.Title = getTitle(item)
	if item.Title == "" {
		item.Title = fallbackTitle
	}

	messageContent := strings.TrimSpace(strings.Join(lines, "\n"))
	item.Description = messageContent

	return true, nil
}

func getTitle(item *repository.Item) string {
	time := item.MetaData.CreationDate.Format(time.RFC850)
	return time
}
