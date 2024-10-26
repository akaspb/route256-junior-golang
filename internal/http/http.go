package http

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewServer(
	address string,
) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	server := &http.Server{Addr: address, Handler: mux}
	return server
}
