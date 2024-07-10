package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/energostack/bisquitt-psk/pkg/config"

	"github.com/Lavalier/zchi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

// CustomRouter is a wrapper around the chi router.
type CustomRouter struct {
	router  *chi.Mux
	handler *CustomHandler
	config  *config.Config
}

// NewCustomRouter creates a new router with the specified handler.
func NewCustomRouter(handler *CustomHandler, cfg *config.Config) *CustomRouter {
	r := CustomRouter{
		router:  chi.NewRouter(),
		handler: handler,
		config:  cfg,
	}
	r.setupRoutes()
	return &r
}

func (cr *CustomRouter) setupRoutes() {
	cr.router.Use(middleware.RealIP)
	cr.router.Use(zchi.Logger(log.Logger))
	cr.router.Use(middleware.Recoverer)

	duration, err := time.ParseDuration(fmt.Sprintf("%ds", cr.config.APITimeoutSeconds))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse timeout duration")
	}

	cr.router.Use(middleware.Timeout(duration))

	cr.router.Route("/clients", func(r chi.Router) {
		r.Use(middleware.BasicAuth("auth", map[string]string{
			cr.config.BasicAuthUsername: cr.config.BasicAuthPassword,
		}))

		r.Get("/{id}", cr.handler.GetClient)
	})

	if cr.config.DocsEnabled {
		cr.router.Get("/docs/json", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./docs/swagger.json")
		})
		cr.router.Get("/docs/*", httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("http://%s:%s/docs/json", cr.config.Host, cr.config.Port)),
		))

		httpSwagger.UIConfig(map[string]string{
			"onComplete": fmt.Sprintf(`() => {
			window.ui.setHost('%s:%s');
	  	}`, cr.config.Host, cr.config.Port),
		})
	}
}

// ListenAndServe starts the HTTP server on the specified port.
func (cr *CustomRouter) ListenAndServe(port string) error {
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), cr.router)
	if err != nil {
		return err
	}
	return nil
}
