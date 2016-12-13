package handlers

import (
	"github.com/andreaskoch/allmark/common/logger"
	"github.com/abbot/go-http-auth"
	"net/http"
)

// RequireDigestAuthentication forces digest access authentication for the given handler.
func RequireDigestAuthentication(logger logger.Logger, baseHandler http.Handler, secretProvider auth.SecretProvider) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authenticator := auth.NewBasicAuthenticator("", secretProvider)

		baseHandlerWithAuthentication := func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
			baseHandler.ServeHTTP(w, &r.Request)
		}

		authHandler := authenticator.Wrap(baseHandlerWithAuthentication)
		authHandler.ServeHTTP(w, r)
	})

}
