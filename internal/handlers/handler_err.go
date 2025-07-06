package main

import (
	"net/http"
	"github.com/yourusername/RSSAggregator/internal/utils"
)

func handlerErr(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, 400, "Something went wrong")
}