package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Start(port string) {
	r := chi.NewRouter()
	r.Use(middleware.RedirectSlashes)

	r.Get("/ping", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("pong!"))
	})

	log.Print("loaded routes !")
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Panicf("failed to start http server: %v", err)
	}
}
