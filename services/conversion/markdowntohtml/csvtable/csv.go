// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvtable

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/pattern"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/util"
	"io"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// csv: [*description text*](*file path*)
	markdownPattern = regexp.MustCompile(`csv: \[([^\]]+)\]\(([^)]+)\)`)
)

func New(pathProvider paths.Pather, files []*model.File) *CsvTableExtension {
	return &CsvTableExtension{
		pathProvider: pathProvider,
		files:        files,
	}
}

type CsvTableExtension struct {
	pathProvider paths.Pather
	files        []*model.File
}

func (converter *CsvTableExtension) Convert(markdown string) (convertedContent string, conversionError error) {

	convertedContent = markdown

	for {

		found, matches := pattern.IsMatch(convertedContent, markdownPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// get the code
		renderedCode := converter.getTableCode(title, path)

		// replace markdown
		convertedContent = strings.Replace(convertedContent, originalText, renderedCode, 1)

	}

	return convertedContent, nil
}

func (converter *CsvTableExtension) getMatchingFile(path string) *model.File {
	for _, file := range converter.files {
		if file.Route().IsMatch(path) && isCSVFile(path) {
			return file
		}
	}

	return nil
}

func (converter *CsvTableExtension) getTableCode(title, path string) string {

	// internal csv file
	if csvFile := converter.getMatchingFile(path); csvFile != nil {

		filePath := converter.pathProvider.Path(csvFile.Route().Value())
		if tableData, err := readCSV(csvFile); err == nil {

			// table header
			tableCode := fmt.Sprintf(`<section class="csv">
					<h1><a href="%s" target="_blank">%s</a></h1>
					<table>
				`, filePath, title)

			// rows
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

			// table footer
			tableCode += `</tbody>
				</table>
			</section>`

			return tableCode
		} else {
			return fmt.Sprintf("<!-- Cannot read csv file %q (Error: %s) -->", path, err)
		}

	}

	// fallback
	return util.GetFallbackLink(title, path)
}

func readCSV(file *model.File) (data [][]string, err error) {
	contentProvider := file.ContentProvider()
	bytesBuffer := new(bytes.Buffer)

	// get the file content
	dataWriter := bufio.NewWriter(bytesBuffer)
	contentReader := func(content io.ReadSeeker) error {
		_, err := io.Copy(dataWriter, content)
		return err
	}

	if dataError := contentProvider.Data(contentReader); dataError != nil {
		return
	}

	// determine the separator
	separator := determineCSVColumnSeparator(bytesBuffer, ';')

	// read the csv
	csvReader := csv.NewReader(bytesBuffer)
	csvReader.Comma = separator

	return csvReader.ReadAll()
}

func determineCSVColumnSeparator(data io.Reader, fallback rune) rune {

	reader := bufio.NewReader(data)
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

func isCSVFile(path string) bool {
	fileExtension := strings.ToLower(filepath.Ext(path))
	switch fileExtension {
	case ".csv":
		return true
	default:
		return false
	}

	panic("Unreachable")
}
