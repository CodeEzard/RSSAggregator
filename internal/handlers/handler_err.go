package handlers

import (
	"net/http"
	"github.com/CodeEzard/RSSAggregator/internal/utils"
)

func HandlerErr(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, http.StatusBadRequest, "Something went wrong")
}