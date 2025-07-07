package handlers

import (
    "net/http"
	"fmt"
	"github.com/CodeEzard/RSSAggregator/internal/auth"
	"github.com/CodeEzard/RSSAggregator/internal/utils"
    "github.com/CodeEzard/RSSAggregator/internal/database"
)

type APIConfig struct {
    DB *database.Queries
}

func (apiConfig APIConfig) MiddlewareAuth(handler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			utils.RespondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiConfig.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			utils.RespondWithError(w, 400, fmt.Sprintf("User not found: %v", err))
			return
		}
    	handler(w, r, user)
	}
}

type AuthedHandler func(http.ResponseWriter, *http.Request, database.User)