package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type XKCDComic struct {
	Num       int    `json:"num"`
	Title     string `json:"title"`
	SafeTitle string `json:"safe_title"`
	Alt       string `json:"alt"`
	Img       string `json:"img"`
	Year      string `json:"year"`
	Month     string `json:"month"`
	Day       string `json:"day"`
}

type QuoteResponse struct {
	Quote  string `json:"quote"`
	Title  string `json:"title"`
	Number int    `json:"comic_number"`
	Date   string `json:"date"`
	URL    string `json:"url"`
}

func MagicHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger.Info("handling magic request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
		)

		// Get a random comic number between 1 and ~2800 (approximate current max)
		rand.Seed(time.Now().UnixNano())
		comicNum := rand.Intn(2800) + 1

		logger.Debug("fetching XKCD comic", zap.Int("comic_number", comicNum))

		// Fetch the XKCD comic!
		resp, err := http.Get(fmt.Sprintf("https://xkcd.com/%d/info.0.json", comicNum))
		if err != nil {
			logger.Error("failed to fetch XKCD comic",
				zap.Int("comic_number", comicNum),
				zap.Error(err),
			)
			http.Error(w, "Failed to fetch XKCD comic", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logger.Warn("XKCD comic not found",
				zap.Int("comic_number", comicNum),
				zap.Int("status_code", resp.StatusCode),
			)
			http.Error(w, "XKCD comic not found", http.StatusNotFound)
			return
		}

		var comic XKCDComic
		if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
			logger.Error("failed to decode XKCD comic",
				zap.Int("comic_number", comicNum),
				zap.Error(err),
			)
			http.Error(w, "Failed to decode XKCD comic", http.StatusInternalServerError)
			return
		}

		// Create the quote response using the alt text as the quote
		quote := QuoteResponse{
			Quote:  comic.Alt,
			Title:  comic.Title,
			Number: comic.Num,
			Date:   fmt.Sprintf("%s-%s-%s", comic.Year, comic.Month, comic.Day),
			URL:    fmt.Sprintf("https://xkcd.com/%d/", comic.Num),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(quote)

		duration := time.Since(start)
		logger.Info("magic request completed",
			zap.Int("comic_number", comicNum),
			zap.String("title", comic.Title),
			zap.Duration("duration", duration),
			zap.Int("status", http.StatusOK),
		)
	}
}
