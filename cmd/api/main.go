package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
"github.com/CodeEzard/RSSAggregator/internal/handlers"
"github.com/CodeEzard/RSSAggregator/internal/database" // Adjust the import path as necessary
"github.com/CodeEzard/RSSAggregator/internal/scraper"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	poststring := os.Getenv("PORT")
	if poststring == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	if !strings.Contains(dbURL, "sslmode=") {
    	if strings.Contains(dbURL, "?") {
        	dbURL += "&sslmode=disable"
    	} else {
        	dbURL += "?sslmode=disable"
    	}
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	db := database.New(conn)
	apiConfig := &handlers.APIConfig{
		DB: db,
	}

	go scraper.StartScraping(db, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE	", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlers.HandlerReadiness)
	v1Router.Get("/err", handlers.HandlerErr)
	v1Router.Post("/users", apiConfig.HandlerCreateUser)
	v1Router.Get("/users", apiConfig.MiddlewareAuth(apiConfig.HandlerGetUser))
	v1Router.Post("/feeds", apiConfig.MiddlewareAuth(apiConfig.HandlerCreateFeed))
	v1Router.Get("/feeds", apiConfig.HandlerGetFeeds)
	v1Router.Get("/posts", apiConfig.MiddlewareAuth(apiConfig.HandlerGetPostsForUser))
	v1Router.Post("/feed_follows", apiConfig.MiddlewareAuth(apiConfig.HandlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiConfig.MiddlewareAuth(apiConfig.HandlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiConfig.MiddlewareAuth(apiConfig.HandlerDeleteFeedFollow))

	router.Mount("/v1/", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + poststring,
	}

	log.Printf("Starting server on port %v", poststring)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port:", poststring)
}
