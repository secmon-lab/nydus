package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/secmon-lab/nydus/pkg/domain/context/logging"
	"github.com/secmon-lab/nydus/pkg/domain/interfaces"
	"github.com/secmon-lab/nydus/pkg/domain/model"
)

type Server struct {
	route *chi.Mux
}

func (x *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	x.route.ServeHTTP(w, r)
}

func New(uc interfaces.UseCase) *Server {
	route := chi.NewRouter()
	route.Use(middlewareLogging)

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

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (x *statusWriter) WriteHeader(code int) {
	x.status = code
	x.ResponseWriter.WriteHeader(code)
}

func middlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.NewString()
		ctx := r.Context()
		logger := logging.From(ctx).With("request_id", reqID)

		ctx = logging.Inject(ctx, logger)

		sw := &statusWriter{ResponseWriter: w}
		next.ServeHTTP(sw, r.WithContext(ctx))

		logger.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", sw.status,
			"remote_addr", r.RemoteAddr,
			"header", r.Header,
			"user_agent", r.UserAgent(),
		)
	})
}

func handleAzureEventGridMessage(uc interfaces.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logging.From(r.Context())

		// Do not use json.Decoder to avoid missing the request body for logging
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Warn("failed to read request body from Azure", "err", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if r.Header.Get("Aeg-Event-Type") == "SubscriptionValidation" {
			var msgs []model.CloudEventValidation

			if err := json.Unmarshal(body, &msgs); err != nil {
				logger.Warn("failed to unmarshal request body from Azure", "err", err, "body", string(body))
				http.Error(w, "bad request", http.StatusBadRequest)
			}

			for _, msg := range msgs {
				if err := uc.ValidateAzureEventGrid(r.Context(), msg.Data.ValidationURL); err != nil {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}
			}
			return
		}

		var ev model.CloudEventSchema
		if err := json.Unmarshal(body, &ev); err != nil {
			logger.Warn("failed to unmarshal request body from Azure", "err", err, "body", string(body))
			http.Error(w, "bad request", http.StatusBadRequest)
			return
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
