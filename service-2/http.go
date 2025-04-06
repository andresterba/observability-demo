package main

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type Store interface {
	Get(ctx context.Context, key string) string
	Set(ctx context.Context, key, value string)
}

type Controller struct {
	tracer trace.Tracer
}

func NewController(tracer trace.Tracer) *Controller {
	return &Controller{
		tracer: tracer,
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, span := c.tracer.Start(r.Context(), "in-controller-entry")
	defer span.End()

	time.Sleep(1 * time.Second)

	switch r.Method {
	case "GET":
		c.handleGet(w, r)
		time.Sleep(1 * time.Second)
	case "POST":
		c.handlePost(w, r)
		time.Sleep(1 * time.Second)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *Controller) handleGet(w http.ResponseWriter, r *http.Request) {
	_, span := c.tracer.Start(r.Context(), "in-handle-get")
	defer span.End()

	time.Sleep(1 * time.Second)

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}
}

func (c *Controller) handlePost(w http.ResponseWriter, r *http.Request) {
	_, span := c.tracer.Start(r.Context(), "in-handle-post")
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

	w.WriteHeader(http.StatusNoContent)
}
