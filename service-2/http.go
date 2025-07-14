package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"observability-demo/lib"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
}

type Controller struct {
	tracer trace.Tracer
	logger *zap.SugaredLogger
	store  Store
}

func NewController(tracer trace.Tracer, store Store, logger *zap.SugaredLogger) *Controller {
	return &Controller{
		tracer: tracer,
		store:  store,
		logger: logger,
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := c.tracer.Start(r.Context(), "in-controller-entry")
	defer span.End()

	switch r.Method {
	case "GET":
		c.handleGet(ctx, w, r)
	case "POST":
		c.handlePost(ctx, w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *Controller) handleGet(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx, span := c.tracer.Start(ctx, "in-handle-get")
	defer span.End()

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key in request", http.StatusBadRequest)
		return
	}

	value, err := c.store.Get(ctx, key)
	if err != nil {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	result := lib.Result{
		Key:   key,
		Value: value,
	}

	fmt.Printf("Returning value for key %s: %s\n", key, value)

	body, err := json.Marshal(result)
	if err != nil {
		c.logger.Errorw("failed to marshal result", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// write the JSON response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)

}

func (c *Controller) handlePost(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx, span := c.tracer.Start(ctx, "in-handle-post")
	defer span.End()

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	value := r.URL.Query().Get("value")
	if value == "" {
		http.Error(w, "missing value", http.StatusBadRequest)
		return
	}

	err := c.store.Set(ctx, key, value)
	if err != nil {
		http.Error(w, "failed to set value", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
