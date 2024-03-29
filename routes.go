package velocidad

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (v *Velocidad) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)

	if v.Debug {
		mux.Use(middleware.Logger)
	}

	mux.Use(middleware.Recoverer)

	mux.Get("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Welcome to Velocidad")
	})

	return mux
}
