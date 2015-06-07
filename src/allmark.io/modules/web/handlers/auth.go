package handlers

import (
	"github.com/abbot/go-http-auth"
	"net/http"
)

// RequireDigestAuthentication forces digest access authentication for the given handler.
func RequireDigestAuthentication(baseHandler http.Handler, realm string, secretProvider auth.SecretProvider) http.Handler {

	authenticator := auth.NewBasicAuthenticator(realm, secretProvider)

	baseHandlerWithAuthentication := func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		baseHandler.ServeHTTP(w, &r.Request)
	}

	return authenticator.Wrap(baseHandlerWithAuthentication)
}
