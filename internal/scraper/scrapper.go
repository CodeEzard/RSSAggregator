package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"
	"github.com/CodeEzard/RSSAggregator/pkg/Rss"
	"github.com/CodeEzard/RSSAggregator/internal/database"
	"github.com/google/uuid"
)

func StartScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequests time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("Error fetching next feed:", err)
			continue
		}

		var wg sync.WaitGroup
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, &wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries,wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	rssFeed, err := Rss.UrlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching RSS feed:", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {

    // Aggressively clean all text data
    	cleanTitle := forceASCII(item.Title)
		cleanDesc := forceASCII(item.Description)
		cleanURL := forceASCII(item.Link)
    
    // Remove HTML from title
    	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
    	cleanTitle = htmlTagRegex.ReplaceAllString(cleanTitle, "")
    	cleanTitle = strings.TrimSpace(cleanTitle)
    
    // Clean description
    	cleanDesc = htmlTagRegex.ReplaceAllString(cleanDesc, "")
    	cleanDesc = strings.TrimSpace(cleanDesc)
    	if len(cleanDesc) > 500 {
    	    cleanDesc = cleanDesc[:500] + "..."
    	}

		description := sql.NullString{}
		if cleanDesc != "" {
        description.String = cleanDesc
        description.Valid = true
        }

		if cleanTitle == "" || cleanURL == "" {
    	  	log.Printf("Skipping post with empty title or URL after cleaning")
    	    continue
}


		pubAt, err := parsePublicationDate(item.PubDate)
        if err != nil {
            log.Printf("Error parsing date %v with err %v", item.PubDate, err)
            pubAt = time.Now() // Default to current time if parsing fails
        }

		_, err = db.CreatePost(context.Background(),
		 	database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       cleanTitle,
				Description: description,
				PublishedAt: pubAt,
				Url:         cleanURL,
				FeedID:      feed.ID,
			})	
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Error creating post:", err)
		}
	}
	log.Printf("Feed: %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
	// Implement the scraping logic here
}

func forceASCII(s string) string {
    var result strings.Builder
    for _, r := range s {
        if r < 128 && r > 31 { 
            result.WriteRune(r)
        } else if r == ' ' || r == '\t' || r == '\n' {
            result.WriteRune(' ')
        }
    }
    return strings.TrimSpace(result.String())
}

func parsePublicationDate(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC1123Z,                    // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC3339,                     // "2006-01-02T15:04:05Z07:00"  
		"2006-01-02T15:04:05-07:00",     // ISO 8601 variant
		"Mon, 2 Jan 2006 15:04:05 -0700", // RFC1123Z without leading zero
		"2 Jan 2006 15:04:05 -0700",     // Without day name
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}