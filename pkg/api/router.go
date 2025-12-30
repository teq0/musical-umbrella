package api

import (
	"net/http"

	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/magic", MagicHandler(logger))

	return mux
}
