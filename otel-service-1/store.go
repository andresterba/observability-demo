package main

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type MemoryStore struct {
	// store is the in-memory store.
	store map[string]string
}

func NewMemoryStore() Store {
	return &MemoryStore{
		store: make(map[string]string),
	}
}

func (s *MemoryStore) Get(ctx context.Context, key string) string {
	ctx, span := tracer.Start(ctx, "in-store-get")
	defer span.End()

	time.Sleep(1 * time.Second)

	if value, ok := s.store[key]; ok {
		return value
	}

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	// client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:4041/get", nil)
	if err != nil {
		panic(err)
	}
	client.Do(req)

	return ""
}

func (s *MemoryStore) Set(ctx context.Context, key, value string) {
	ctx, span := tracer.Start(ctx, "in-store-set")
	defer span.End()

	// s.store[key] = value

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	// client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:4041/set", nil)
	if err != nil {
		panic(err)
	}

	client.Do(req)
}
