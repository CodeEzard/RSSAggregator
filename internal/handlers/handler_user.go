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

func (apiConfig *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON:%v", err))
		return
	}

	user, err := apiConfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:     params.Name,
	})
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	utils.RespondWithJSON(w, 201, models.DatabaseUserToUser(user))
}

func (apiConfig *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	utils.RespondWithJSON(w, 200, models.DatabaseUserToUser(user))
}

func (apiConfig *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {

	clean := r.URL.Query().Get("clean") == "true"

	posts, err := apiConfig.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Error getting posts for user: %v", err))
		return
	}

	// Clean posts before sending response
    cleanPosts := make([]models.CleanPost, len(posts))
    for i, post := range posts {
        cleanPosts[i] = models.CleanPost{
            ID:             post.ID,
            CreatedAt:      post.CreatedAt,
            UpdatedAt:      post.UpdatedAt,
            Title:          utils.CleanTitle(post.Title),
            Description:    utils.CleanDescription(post.Description),
            PublishedAt:    post.PublishedAt,
            ApplicationUrl: post.Url,
            CompanyID:      post.FeedID,
        }
    }

    if clean {
        utils.RespondWithJSON(w, 200, cleanPosts)
    } else {
        utils.RespondWithJSON(w, 200, models.DatabasePostsToPosts(posts))
    }
}