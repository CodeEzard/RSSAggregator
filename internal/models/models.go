package models

import (
	"database/sql"
	"time"
	"github.com/CodeEzard/RSSAggregator/internal/utils"
	"github.com/CodeEzard/RSSAggregator/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
}

func DatabaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		APIKey:    dbUser.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func DatabaseFeedToFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,	
		UpdatedAt: dbFeed.UpdatedAt,
		Name:      dbFeed.Name,
		Url:       dbFeed.Url,
		UserID:    dbFeed.UserID,
	}
}

func DatabaseFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	 feeds := []Feed{}
	 for _, dbFeed := range dbFeeds {
         feeds = append(feeds, DatabaseFeedToFeed(dbFeed))
	 }
	 return feeds
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	UserID   uuid.UUID   `json:"user_id"`
	FeedID   uuid.UUID   `json:"feed_id"`
}

func DatabaseFeedFollowToFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
		UserID:   dbFeedFollow.UserID,
		FeedID:   dbFeedFollow.FeedID,
	}
}

func DatabaseFeedFollowsToFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollow {
	 feedFollows := []FeedFollow{}
	 for _, dbFeedFollow := range dbFeedFollows {
         feedFollows = append(feedFollows, DatabaseFeedFollowToFeedFollow(dbFeedFollow))
	 }
	 return feedFollows
}

func DatabaseDeleteFeedFollowsToDeleteFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollow {
	 feedFollows := []FeedFollow{}
	 for _, dbFeedFollow := range dbFeedFollows {
         feedFollows = append(feedFollows, DatabaseFeedFollowToFeedFollow(dbFeedFollow))
	 }
	 return feedFollows
}

type Post struct {
	ID          uuid.UUID      `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Title       string         `json:"title"`
	Description *string        `json:"description"`
	PublishedAt time.Time      `json:"published_at"`
	Url         string         `json:"url"`
	FeedID      uuid.UUID      `json:"feed_id"`
}

func DatabasePostToPost(dbPost database.Post) Post {
	var description *string
	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}
	return Post{
		ID: dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Description: description,
		PublishedAt: dbPost.PublishedAt,
		Url:         dbPost.Url,
		FeedID:      dbPost.FeedID,
	}
}

func DatabasePostsToPosts(dbPosts []database.Post) []Post {
	posts := []Post{}
	for _, dbPost := range dbPosts {
		posts = append(posts, DatabasePostToPost(dbPost))
	}
	return posts
}	

type CleanPost struct {
    ID             uuid.UUID `json:"id"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
    Title          string    `json:"title"`
    Description    string    `json:"description"`
    PublishedAt    time.Time `json:"posted_at"`
    ApplicationUrl string    `json:"application_url"`
    CompanyID      uuid.UUID `json:"company_id"`
}

func CleanPostDescription(desc sql.NullString) string {
    if !desc.Valid {
        return ""
    }

	cleaned := utils.CleanDescription(desc.String)
	
	// Limit to first 200 characters for API response
	if len(cleaned) > 200 {
		return cleaned[:200] + "..."
	}
	return cleaned
}

// cleanDescription removes unwanted characters or HTML tags from the description string.
// cleanDescription removes unwanted characters or HTML tags from the description string.
// This function is defined in utils and should be used directly as utils.StripHTML.
