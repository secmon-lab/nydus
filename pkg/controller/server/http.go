package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/secmon-as-code/nydus/pkg/domain/context/logging"
	"github.com/secmon-as-code/nydus/pkg/domain/interfaces"
	"github.com/secmon-as-code/nydus/pkg/domain/model"
	"github.com/secmon-as-code/nydus/pkg/usecase"
)

type Server struct {
	route *chi.Mux
}

func (x *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	x.route.ServeHTTP(w, r)
}

func New(uc interfaces.UseCase) *Server {
	route := chi.NewRouter()

	if v, ok := uc.(*usecase.UseCase); ok {
		logging.Default().Debug("[tmp] check client: new server", "client", v.Clients())
	}

	route.Route("/google/pubsub", func(r chi.Router) {
		r.Post("/cloud-storage", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotImplemented)
		})
	})
	route.Route("/aws/sqs", func(r chi.Router) {
		r.Post("/s3", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotImplemented)
		})
	})
	route.Route("/azure/event-grid", func(r chi.Router) {
		r.Options("/blob-storage", handleAzureEventGridValidate(uc))
		r.Post("/blob-storage", handleAzureEventGridMessage(uc))
	})

	return &Server{
		route: route,
	}
}

func handleAzureEventGridMessage(uc interfaces.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ev model.CloudEventSchema
		logger := logging.From(r.Context())

		// Do not use json.Decoder to avoid missing the request body for logging
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Warn("failed to read request body from Azure", "err", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &ev); err != nil {
			logger.Warn("failed to unmarshal request body from Azure", "err", err, "body", string(body))
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if v, ok := uc.(*usecase.UseCase); ok {
			logger.Debug("[tmp] check client: received Azure CloudEvent", "client", v.Clients())
		}

		if ev.Type == "Microsoft.Storage.BlobCreated" {
			if err := uc.HandleAzureCloudEvent(r.Context(), &ev); err != nil {
				logger.Warn("failed to handle Azure CloudEvent", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			logger.Warn("unexpected event type", "type", ev.Type)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func handleAzureEventGridValidate(uc interfaces.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Webhook-Request-Origin") != "eventgrid.azure.net" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if err := uc.ValidateAzureEventGrid(r.Context(), r.Header.Get("Webhook-Request-Callback")); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
