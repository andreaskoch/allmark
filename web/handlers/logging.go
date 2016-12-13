package handlers

import (
	gorillahandlers "github.com/gorilla/handlers"
	"net/http"
	"os"
)

func LogRequests(baseHandler http.Handler) http.Handler {
	return gorillahandlers.LoggingHandler(os.Stdout, baseHandler)
}
