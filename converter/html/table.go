// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// csv: [*description text*](*file path*)
	tablePattern = regexp.MustCompile(`csv: \[([^\]]+)\]\(([^)]+)\)`)
)

func renderCSVTables(item *repository.Item, rawContent string) string {
	return convertCSVMarkdownExtension(rawContent, item.Files, item.FilePathProvider())
}

func convertCSVMarkdownExtension(markdown string, fileIndex *repository.FileIndex, pathProvider *path.Provider) string {

	for {

		found, matches := util.IsMatch(markdown, tablePattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// create image gallery code
		files := fileIndex.FilesByPath(path, isCSVFile)

		if len(files) == 0 {
			// file not found remove entry
			msg := fmt.Sprintf("<!-- Cannot render table. The file %q could not be found -->", path)
			markdown = strings.Replace(markdown, originalText, msg, 1)
			continue
		}

		matchedFile := files[0]
		realFilePath := matchedFile.Path()
		tableData, err := readCSV(realFilePath)
		if err != nil {
			// file not found remove entry
			msg := fmt.Sprintf("<!-- Cannot read csv file %q (Error: %s) -->", path, err)
			markdown = strings.Replace(markdown, originalText, msg, 1)
			continue
		}

		relativeFilePath := pathProvider.GetWebRoute(matchedFile)
		tableCode := fmt.Sprintf(`<section class="csv">
			<h1><a href="%s" target="_blank">%s</a></h1>
			<table>
		`, relativeFilePath, title)

		for rowNumber := range tableData {
			row := tableData[rowNumber]

			if rowNumber == 0 {
				tableCode += `<thead>`
			}

			if rowNumber == 1 {
				tableCode += `<tbody>`
			}

			tableCode += `<tr>`

			for columnNumber := range row {
				value := row[columnNumber]
				tableCode += fmt.Sprintf(`<td>%s</td>`, value)
			}

			tableCode += `</tr>`

			if rowNumber == 0 {
				tableCode += `</thead>`
			}
		}

		tableCode += `</tbody>
			</table>
		</section>`

		// replace markdown with image gallery
		markdown = strings.Replace(markdown, originalText, tableCode, 1)

	}

	return markdown
}

func readCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	// determine the separator
	separator := determineCSVColumnSeparator(path, ';')

	// read the csv
	csvReader := csv.NewReader(file)
	csvReader.Comma = separator

	return csvReader.ReadAll()
}

func determineCSVColumnSeparator(path string, fallback rune) rune {

	file, err := os.Open(path)
	if err != nil {
		return fallback
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	line, _, err := reader.ReadLine()
	if err != nil {
		return fallback
	}

	for _, character := range line {
		switch character {
		case ',':
			return ','
		case ';':
			return ';'
		case '\t':
			return '\t'
		}
	}

	return fallback
}

func isCSVFile(pather path.Pather) bool {
	fileExtension := strings.ToLower(filepath.Ext(pather.Path()))
	switch fileExtension {
	case ".csv":
		return true
	default:
		return false
	}

	panic("Unreachable")
}
