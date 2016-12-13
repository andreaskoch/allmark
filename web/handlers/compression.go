package handlers

import (
	gorillahandlers "github.com/gorilla/handlers"
	"net/http"
)

func CompressResponses(baseHandler http.Handler) http.Handler {
	return gorillahandlers.CompressHandler(baseHandler)
}
