// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/common/util/fsutil"
	"allmark.io/modules/web/header"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func DOCX(logger logger.Logger,
	conversionToolPath string,
	headerWriter header.HeaderWriter,
	fileOrchestrator *orchestrator.FileOrchestrator,
	converterModelOrchestrator *orchestrator.ConversionModelOrchestrator,
	templateProvider templates.Provider,
	error404Handler http.Handler) http.Handler {

	convertToHtml := func(baseURL string, viewModel viewmodel.ConversionModel) string {

		// get a template
		template, err := templateProvider.GetConversionTemplate(baseURL)
		if err != nil {
			logger.Error("No template for item of type %q.", viewModel.Type)
			return ""
		}

		// render template
		buffer := new(bytes.Buffer)
		writer := bufio.NewWriter(buffer)
		if err := renderTemplate(template, viewModel, writer); err != nil {
			logger.Error("%s", err)
			return ""
		}

		writer.Flush()

		return buffer.String()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_RTF)

		// strip the "docx" or ".docx" suffix from the path
		path := r.URL.Path
		path = strings.TrimSuffix(path, "docx")
		path = strings.TrimSuffix(path, ".")

		// get the request route
		requestRoute := route.NewFromRequest(path)

		// make sure the request body is closed
		defer r.Body.Close()

		// the temporary working directory
		targetDirectory := fsutil.GetTempDirectory()

		// get the conversion model
		baseURL := getBaseURLFromRequest(r)

		// make the baseURL HTTP because pandoc has problems with HTTPS
		baseURL = strings.Replace(baseURL, "https://", "http://", 1)

		// make sure pandoc only performs local requests
		originalHostname := r.Host
		newHostname := "localhost"
		_, port, _ := net.SplitHostPort(originalHostname)
		if port != "" {
			newHostname = newHostname + ":" + port
		}

		baseURL = strings.Replace(baseURL, originalHostname, newHostname, 1)

		model, found := converterModelOrchestrator.GetConversionModel(baseURL, requestRoute)
		if !found {

			// display a 404 error page
			error404Handler.ServeHTTP(w, r)
			return
		}

		html := convertToHtml(baseURL, model)

		// write the html to a temp file (Note: the file extension .html is important for pandoc)
		htmlFilePath := filepath.Join(targetDirectory, "source.html")
		htmlFile, err := os.Create(htmlFilePath)
		if err != nil {
			logger.Error("Cannot open HTML file for writing. Error: %s", err.Error())
			return
		}

		htmlFile.WriteString(html)
		htmlFile.Sync()

		// close and delete the file at the end of the function
		htmlFile.Close()

		// get a target file path (Note: the file extension .docx is important for pandoc)
		targetFilePath := filepath.Join(targetDirectory, "target.docx")

		// call pandoc
		args := []string{
			"-s",
			fmt.Sprintf(`%s`, htmlFilePath),
			"-o",
			fmt.Sprintf(`%s`, targetFilePath),
		}

		cmd := exec.Command(conversionToolPath, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = targetDirectory

		if err := cmd.Run(); err != nil {
			logger.Error("Could not run pandoc: %v", err)
			return
		}

		logger.Debug("Saving conversion files to directory: %q", targetDirectory)

		// docx file
		docxFile, err := fsutil.OpenFile(targetFilePath)
		if err != nil {
			logger.Error("Cannot open target file. Error: %s", err.Error())
			return
		}

		// close and delete the file at the end of the function
		defer docxFile.Close()

		// remove the temp directory at the end
		defer func() {
			logger.Debug("Deleting conversion file directory: %q", targetDirectory)
			if err := deleteFile(targetDirectory); err != nil {
				logger.Error("Could not delete the temporary working directory (%q) that has been created during the conversion. Error: %s", targetDirectory, err.Error())
			}
		}()

		w.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, getRichTextFilename(model)))

		io.Copy(w, docxFile)

		return
	})

}

func getRichTextFilename(model viewmodel.ConversionModel) string {
	originalRoute := route.NewFromRequest(model.Route)
	fileNameRoute := route.NewFromRequest(originalRoute.LastComponentName())

	if model.Level == 0 {
		fileNameRoute = route.NewFromRequest(model.Title)
	}

	return fmt.Sprintf("%s.docx", fileNameRoute.Value())
}

// deleteFile removes the file with the specified path.
func deleteFile(filepath string) error {
	return os.RemoveAll(filepath)
}
