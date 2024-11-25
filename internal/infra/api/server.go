package api

import "net/http"

type Config struct {
	Host string
}

func InitServer(cfg Config, h handlers) *http.Server {
	return &http.Server{
		Addr:    cfg.Host,
		Handler: newHandler(h),
	}
}
