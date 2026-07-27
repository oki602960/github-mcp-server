package middleware

import (
	"errors"
	"fmt"
	"net/http"

	ghcontext "github.com/github/github-mcp-server/pkg/context"
	"github.com/github/github-mcp-server/pkg/http/oauth"
	"github.com/github/github-mcp-server/pkg/utils"
)

func ExtractUserToken(oauthCfg *oauth.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Check if token info already exists in context, if it does, skip extraction.
			// In remote setup, we may have already extracted token info earlier.
			if _, ok := ghcontext.GetTokenInfo(ctx); ok {
				// Token info already exists in context, skip extraction
				next.ServeHTTP(w, r)
				return
			}

			tokenType, token, err := utils.ParseAuthorizationHeader(r)
			if err != nil {
				if errors.Is(err, utils.ErrMissingAuthorizationHeader) {
					// Fall back to ?token= query parameter so clients that cannot set
					// custom HTTP headers can embed the GitHub token in the server URL.
					if qt := r.URL.Query().Get("token"); qt != "" {
						tokenType, token, err = utils.ParseBearerToken(qt)
						if err != nil {
							http.Error(w, err.Error(), http.StatusBadRequest)
							return
						}
						// continue below with token from query param
					} else {
						// No token anywhere — issue OAuth challenge per MCP spec.
						sendAuthChallenge(w, r, oauthCfg)
						return
					}
				} else {
					// Bad format or unsupported scheme in the Authorization header.
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

			ctx = ghcontext.WithTokenInfo(ctx, &ghcontext.TokenInfo{
				Token:     token,
				TokenType: tokenType,
			})
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// sendAuthChallenge sends a 401 Unauthorized response with WWW-Authenticate header
// containing the OAuth protected resource metadata URL as per RFC 6750 and MCP spec.
func sendAuthChallenge(w http.ResponseWriter, r *http.Request, oauthCfg *oauth.Config) {
	resourcePath := oauth.ResolveResourcePath(r, oauthCfg)
	resourceMetadataURL := oauth.BuildResourceMetadataURL(r, oauthCfg, resourcePath)
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Bearer resource_metadata=%q`, resourceMetadataURL))
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
