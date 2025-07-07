package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/CodeEzard/RSSAggregator/internal/database" // Adjust the import path as necessary
	"github.com/CodeEzard/RSSAggregator/internal/utils"
    "github.com/CodeEzard/RSSAggregator/internal/models"
)

func (apiConfig *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON:%v", err))
		return
	}

	feed, err := apiConfig.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:     params.Name,
		Url:       params.URL,
		UserID:   user.ID,
	})
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Error creating feed: %v", err))
		return
	}

	utils.RespondWithJSON(w, 201, models.DatabaseFeedToFeed(feed))
}

func (apiConfig *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	

	feeds, err := apiConfig.DB.GetFeeds(r.Context())
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Error creating feed: %v", err))
		return
	}

	utils.RespondWithJSON(w, 201, models.DatabaseFeedsToFeeds(feeds))
}

