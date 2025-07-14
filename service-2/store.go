package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type MemoryStore struct {
	store map[string]string

	log    *zap.SugaredLogger
	tracer trace.Tracer
}

func NewMemoryStore(tracer trace.Tracer, logger *zap.SugaredLogger) Store {
	return &MemoryStore{
		store:  make(map[string]string),
		tracer: tracer,
		log:    logger,
	}
}

func (s *MemoryStore) Get(ctx context.Context, key string) (string, error) {
	_, span := s.tracer.Start(ctx, "in-store-get")
	defer span.End()

	if value, ok := s.store[key]; ok {
		s.log.Infof("found key %s with value %s", key, value)

		return value, nil
	}

	return "", fmt.Errorf("key %s not found in store", key)
}

func (s *MemoryStore) Set(ctx context.Context, key, value string) error {
	_, span := s.tracer.Start(ctx, "in-store-set")
	defer span.End()

	s.store[key] = value

	s.log.Infof("set key %s with value %s", key, value)

	return nil
}
