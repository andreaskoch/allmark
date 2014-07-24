// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtfhandler

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/errorhandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/view/templates"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func New(logger logger.Logger, config *config.Config, itemIndex *index.Index, patherFactory paths.PatherFactory, converter conversion.Converter) *RtfHandler {

	// templates
	templateProvider := templates.NewProvider(config.TemplatesFolder())

	// error
	error404Handler := errorhandler.New(logger, config, itemIndex, patherFactory)

	// viewmodel
	conversionModelOrchestrator := orchestrator.NewConversionModelOrchestrator(itemIndex, converter)

	return &RtfHandler{
		logger:                      logger,
		itemIndex:                   itemIndex,
		config:                      config,
		patherFactory:               patherFactory,
		templateProvider:            templateProvider,
		error404Handler:             error404Handler,
		conversionModelOrchestrator: conversionModelOrchestrator,
	}
}

type RtfHandler struct {
	logger                      logger.Logger
	itemIndex                   *index.Index
	config                      *config.Config
	patherFactory               paths.PatherFactory
	templateProvider            *templates.Provider
	error404Handler             *errorhandler.ErrorHandler
	conversionModelOrchestrator orchestrator.ConversionModelOrchestrator
}

func (handler *RtfHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// get the path from the request variables
		vars := mux.Vars(r)
		path := vars["path"]

		// get the request route
		requestRoute, err := route.NewFromRequest(path)
		if err != nil {
			handler.logger.Error("Unable to get route from request. Error: %s", err.Error())
			return
		}

		// make sure the request body is closed
		defer r.Body.Close()

		// check if there is a item for the request
		item, found := handler.itemIndex.IsMatch(*requestRoute)
		if !found {

			// display a 404 error page
			error404Handler := handler.error404Handler.Func()
			error404Handler(w, r)
			return
		}

		// prepare a path provider which includes the hostname
		hostname := handlerutil.GetHostnameFromRequest(r)
		addressPrefix := fmt.Sprintf("http://%s/", hostname)
		pathProvider := handler.patherFactory.Absolute(addressPrefix)

		// render the view model
		viewModel := handler.conversionModelOrchestrator.GetConversionModel(pathProvider, item)
		html := handler.convertToHtml(viewModel)

		// write the html to a temp file
		htmlFilePath := getTempFileName("html-source") + ".html"
		htmlFile, err := fsutil.OpenFile(htmlFilePath)
		if err != nil {
			handler.logger.Error("Cannot open HTML file for writing. Error: %s", err.Error())
			return
		}

		defer htmlFile.Close()
		htmlFile.WriteString(html)

		// get a target file path
		targetFile := getTempFileName("rtf-target") + ".rtf"

		// call pandoc
		args := []string{
			"-s",
			fmt.Sprintf(`%s`, htmlFilePath),
			"-o",
			fmt.Sprintf(`%s`, targetFile),
		}

		cmd := exec.Command("pandoc.exe", args...)
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

		io.Copy(w, rtfFile)

		return
	}
}

func (handler *RtfHandler) convertToHtml(viewModel viewmodel.ConversionModel) string {

	// get a template
	template, err := handler.templateProvider.GetSubTemplate(templates.ConversionTemplateName)
	if err != nil {
		handler.logger.Error("No template for item of type %q.", viewModel.Type)
		return ""
	}

	// render template
	buffer := new(bytes.Buffer)
	writer := bufio.NewWriter(buffer)
	if err := handlerutil.RenderTemplate(viewModel, template, writer); err != nil {
		handler.logger.Error("%s", err)
		return ""
	}

	writer.Flush()

	return buffer.String()
}

func getTempFileName(prefix string) string {
	file, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("%s-rtf-conversion", prefix))
	if err != nil {
		panic(err)
	}

	defer file.Close()

	return file.Name()
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
