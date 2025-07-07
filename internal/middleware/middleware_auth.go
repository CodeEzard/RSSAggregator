package middleware

import (
	"fmt"
	"net/http"
	"github.com/CodeEzard/RSSAggregator/internal/auth"
	"github.com/CodeEzard/RSSAggregator/internal/handlers"
	"github.com/CodeEzard/RSSAggregator/internal/utils"
	"github.com/CodeEzard/RSSAggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiConfig *handlers.apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
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