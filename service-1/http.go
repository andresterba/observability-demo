package main

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Client interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
}

type Controller struct {
	client Client

	tracer trace.Tracer
	log    *zap.SugaredLogger
}

func NewController(store Client, tracer trace.Tracer, log *zap.SugaredLogger) *Controller {
	return &Controller{
		client: store,
		tracer: tracer,
		log:    log,
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := c.tracer.Start(r.Context(), "in-controller-entry", trace.WithAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
		attribute.String("http.host", r.Host),
	))
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
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	value, err := c.client.Get(ctx, key)
	if err != nil {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	_, err = w.Write([]byte(value))
	if err != nil {
		c.log.Errorw("failed to write response", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
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

	err := c.client.Set(ctx, key, value)
	if err != nil {
		c.log.Errorw("failed to set value", "key", key, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
