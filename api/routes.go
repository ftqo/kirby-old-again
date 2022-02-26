package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ftqo/kirby/logger"
)

func Start(port string) {
	r := chi.NewRouter()
	r.Use(middleware.RedirectSlashes)

	r.Get("/ping", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("pong!"))
	})

	logger.L.Info().Msg("Loaded routes")
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		logger.L.Panic().Err(err).Msg("Failed to start http server")
	}
}
