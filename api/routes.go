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
		_, err := rw.Write([]byte("pong!"))
		if err != nil {
			logger.L.Error().Err(err).Msg("Failed to write pong response")
		}
	})

	logger.L.Info().Msgf("Opened port %s for API calls", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		logger.L.Panic().Err(err).Msg("Failed to start http server")
	}
}
