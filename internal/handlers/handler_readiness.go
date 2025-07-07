package handlers
import (
    "net/http"
    "github.com/CodeEzard/RSSAggregator/internal/utils"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
    utils.RespondWithJSON(w, 200, struct{}{})
}