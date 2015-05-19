// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/common/util/fsutil"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/server/header"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
	"bufio"
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Rtf struct {
	logger                     logger.Logger
	conversionToolPath         string
	headerWriter               header.HeaderWriter
	converterModelOrchestrator *orchestrator.ConversionModelOrchestrator
	templateProvider           templates.Provider
	error404Handler            Handler
}

func (handler *Rtf) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		handler.headerWriter.Write(w, header.CONTENTTYPE_RTF)

		// get the path from the request variables
		vars := mux.Vars(r)
		path := vars["path"]

		// strip the "rtf" or ".rtf" suffix from the path
		path = strings.TrimSuffix(path, "rtf")
		path = strings.TrimSuffix(path, ".")

		// get the request route
		requestRoute := route.NewFromRequest(path)

		// make sure the request body is closed
		defer r.Body.Close()

		// get the conversion model
		hostname := getHostnameFromRequest(r)
		model, found := handler.converterModelOrchestrator.GetConversionModel(hostname, requestRoute)
		if !found {

			// display a 404 error page
			error404Handler := handler.error404Handler.Func()
			error404Handler(w, r)
			return
		}

		html := handler.convertToHtml(hostname, model)

		// write the html to a temp file
		htmlFilePath := fsutil.GetTempFileName("html-source") + ".html"
		htmlFile, err := os.Create(htmlFilePath)
		if err != nil {
			handler.logger.Error("Cannot open HTML file for writing. Error: %s", err.Error())
			return
		}

		htmlFile.WriteString(html)
		htmlFile.Sync()
		htmlFile.Close()

		// get a target file path
		targetFile := fsutil.GetTempFileName("rtf-target") + ".rtf"

		// call pandoc
		args := []string{
			"-s",
			fmt.Sprintf(`%s`, htmlFilePath),
			"-o",
			fmt.Sprintf(`%s`, targetFile),
		}

		cmd := exec.Command(handler.conversionToolPath, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			handler.logger.Error("Could not run pandoc: %v", err)
			return
		}

		// rtf file
		rtfFile, err := fsutil.OpenFile(targetFile)
		if err != nil {
			handler.logger.Error("Cannot open target file. Error: %s", err.Error())
			return
		}

		defer rtfFile.Close()

		w.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, getRichTextFilename(model)))

		io.Copy(w, rtfFile)

		return
	}
}

func (handler *Rtf) convertToHtml(hostname string, viewModel viewmodel.ConversionModel) string {

	// get a template
	template, err := handler.templateProvider.GetSubTemplate(hostname, templates.ConversionTemplateName)
	if err != nil {
		handler.logger.Error("No template for item of type %q.", viewModel.Type)
		return ""
	}

	// render template
	buffer := new(bytes.Buffer)
	writer := bufio.NewWriter(buffer)
	if err := renderTemplate(viewModel, template, writer); err != nil {
		handler.logger.Error("%s", err)
		return ""
	}

	writer.Flush()

	return buffer.String()
}

func getRichTextFilename(model viewmodel.ConversionModel) string {
	originalRoute := route.NewFromRequest(model.Route)
	fileNameRoute := route.NewFromRequest(originalRoute.LastComponentName())

	if model.Level == 0 {
		fileNameRoute = route.NewFromRequest(model.Title)
	}

	return fmt.Sprintf("%s.rtf", fileNameRoute.Value())
}

func execute(directory, commandText string) error {

	// get the command
	command := getCmd(directory, commandText)

	// execute the command
	if err := command.Start(); err != nil {
		return err
	}

	// wait for the command to finish
	return command.Wait()
}

func getCmd(directory, commandText string) *exec.Cmd {
	if commandText == "" {
		return nil
	}

	components := strings.Split(commandText, " ")

	// get the command name
	commandName := components[0]

	// get the command arguments
	arguments := make([]string, 0)
	if len(components) > 1 {
		arguments = components[1:]
	}

	// create the command
	command := exec.Command(commandName, arguments...)

	// set the working directory
	command.Dir = directory

	// redirect command io
	redirectCommandIO(command)

	return command
}

func redirectCommandIO(cmd *exec.Cmd) (*os.File, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	//direct. Masked passwords work OK!
	cmd.Stdin = os.Stdin
	return nil, err
}
