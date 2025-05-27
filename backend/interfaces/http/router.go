package http

import (
	"net/http"
)

func SetupRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", PingHandler)

	return mux
}
