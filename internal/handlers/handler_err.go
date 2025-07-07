package handlers

import (
	"net/http"
	"github.com/CodeEzard/RSSAggregator/internal/utils"
)

func HandlerErr(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, "Something went wrong", http.StatusBadRequest)
}