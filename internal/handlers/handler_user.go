package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/CodeEzard/RSSAggregator/internal/database" // Adjust the import path as necessary
)

func (apiConfig *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON:%v", err))
		return
	}

	user, err := apiConfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:     params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseUserToUser(user))
}

func (apiConfig *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, databaseUserToUser(user))	
}

func (apiConfig *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {

	clean := r.URL.Query().Get("clean") == "true"

	posts, err := apiConfig.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting posts for user: %v", err))
		return
	}

	// Clean posts before sending response
    cleanPosts := make([]CleanPost, len(posts))
    for i, post := range posts {
        cleanPosts[i] = CleanPost{
            ID:             post.ID, 
            CreatedAt:      post.CreatedAt,
            UpdatedAt:      post.UpdatedAt,
            Title:          cleanTitle(post.Title),
            Description:    cleanPostDescription(post.Description),
            PublishedAt:    post.PublishedAt,
            ApplicationUrl: post.Url,
            CompanyID:      post.FeedID,
        }
    }

    if clean {
        respondWithJSON(w, 200, cleanPosts)
    } else {
        respondWithJSON(w, 200, databasePostsToPosts(posts))
    }
}