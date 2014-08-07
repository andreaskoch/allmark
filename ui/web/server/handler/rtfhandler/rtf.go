// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtfhandler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/errorhandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func New(logger logger.Logger, config *config.Config, itemIndex *index.Index, patherFactory paths.PatherFactory) *RtfHandler {

	// error
	error404Handler := errorhandler.New(logger, config, itemIndex, patherFactory)

	return &RtfHandler{
		logger:          logger,
		itemIndex:       itemIndex,
		config:          config,
		patherFactory:   patherFactory,
		error404Handler: error404Handler,
	}
}

type RtfHandler struct {
	logger          logger.Logger
	itemIndex       *index.Index
	config          *config.Config
	patherFactory   paths.PatherFactory
	error404Handler *errorhandler.ErrorHandler
}

func (handler *RtfHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// get the path from the request variables
		vars := mux.Vars(r)
		path := vars["path"]

		// strip the "rtf" or ".rtf" suffix from the path
		path = strings.TrimSuffix(path, "rtf")
		path = strings.TrimSuffix(path, ".")

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

		// check if the a conversion tool has been supplied
		conversionToolIsConfigured := len(handler.config.Conversion.Tool) > 0
		if !conversionToolIsConfigured {

			handler.logger.Warn("Cannot convert item %q to RTF. No conversion tool configured.", item.String())

			// display a 404 error page
			error404Handler := handler.error404Handler.Func()
			error404Handler(w, r)
			return

		}

		// prepare a path provider which includes the hostname
		hostname := handlerutil.GetHostnameFromRequest(r)
		addressPrefix := fmt.Sprintf("http://%s/", hostname)
		pathProvider := handler.patherFactory.Absolute(addressPrefix)

		// assemble the item url
		sourceUrl := pathProvider.Path(orchestrator.GetTypedItemUrl(item, "print"))

		// get a target file path
		targetFile := getTempFileName("rtf-target") + ".rtf"

		// call pandoc
		args := []string{
			"-s",
			fmt.Sprintf(`%s`, sourceUrl),
			"-o",
			fmt.Sprintf(`%s`, targetFile),
		}

		cmd := exec.Command(handler.config.Conversion.Tool, args...)
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

		w.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, getRichTextFilename(item)))

		io.Copy(w, rtfFile)

		return
	}
}

func getRichTextFilename(item *model.Item) string {
	fallback := "document"

	fileNameRoute, err := route.NewFromRequest(item.Route().LastComponentName())
	if err != nil {
		return fallback
	}

	if item.Route().Level() == 0 {
		fileNameRoute, err = route.NewFromRequest(item.Title)
		if err != nil {
			return fallback
		}
	}

	return fmt.Sprintf("%s.rtf", fileNameRoute.Value())
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
