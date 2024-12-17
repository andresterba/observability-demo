package main

import (
	"context"
	"net/http"
	"time"
)

type Store interface {
	Get(ctx context.Context, key string) string
	Set(ctx context.Context, key, value string)
}

type Controller struct {
	// server is the HTTP server.
	store Store
}

func NewController(store Store) *Controller {
	return &Controller{
		store: store,
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "in-controller-entry")
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
	_, span := tracer.Start(r.Context(), "in-handle-get")
	defer span.End()

	time.Sleep(1 * time.Second)

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	value := c.store.Get(r.Context(), key)
	if value == "" {
		time.Sleep(time.Second)

		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	time.Sleep(time.Second)

	w.Write([]byte(value))
}

func (c *Controller) handlePost(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "in-handle-post")
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

	c.store.Set(r.Context(), key, value)
	w.WriteHeader(http.StatusNoContent)
}
