package handlers

import (
	"github.com/abbot/go-http-auth"
	"net/http"
)

// RequireDigestAuthentication forces digest access authentication for the given handler.
func RequireDigestAuthentication(baseHandler http.Handler, secretProvider auth.SecretProvider) http.Handler {

	authenticator := auth.NewBasicAuthenticator("", secretProvider)

	baseHandlerWithAuthentication := func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		baseHandler.ServeHTTP(w, &r.Request)
	}

	return authenticator.Wrap(baseHandlerWithAuthentication)
}
