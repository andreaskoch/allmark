package handlers

import (
	"allmark.io/modules/common/logger"
	"github.com/abbot/go-http-auth"
	"net/http"
)

// RequireDigestAuthentication forces digest access authentication for the given handler.
func RequireDigestAuthentication(logger logger.Logger, baseHandler http.Handler, secretProvider auth.SecretProvider) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// don't require authentication if it's a local request
		if isLocalRequest(r) {
			logger.Debug("Skipping authentication for request %q because it's a local request", r.URL.String())
			baseHandler.ServeHTTP(w, r)
			return
		}

		authenticator := auth.NewBasicAuthenticator("", secretProvider)

		baseHandlerWithAuthentication := func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
			baseHandler.ServeHTTP(w, &r.Request)
		}

		authHandler := authenticator.Wrap(baseHandlerWithAuthentication)
		authHandler.ServeHTTP(w, r)
	})

}
