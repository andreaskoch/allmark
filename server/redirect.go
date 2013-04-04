package server

import (
	"net/http"
)

func redirect(w http.ResponseWriter, r *http.Request, route string) {
	http.Redirect(w, r, route, http.StatusMovedPermanently)
}
