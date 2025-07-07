package handlers

import (
    "net/http"
    "github.com/CodeEzard/RSSAggregator/internal/database"
)

type ApiConfig struct {
    DB *database.Queries
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)