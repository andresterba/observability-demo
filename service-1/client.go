package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"observability-demo/lib"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type StoreClient struct {
	log    *zap.SugaredLogger
	tracer trace.Tracer
}

func NewClient(tracer trace.Tracer, log *zap.SugaredLogger) Client {
	return &StoreClient{
		tracer: tracer,
		log:    log,
	}
}

func (s *StoreClient) Get(ctx context.Context, key string) (string, error) {
	ctx, span := s.tracer.Start(ctx, "in-client-get")
	defer span.End()

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	url := fmt.Sprintf("http://localhost:4041/get?key=%s", key)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get key %s: %s", key, resp.Status)
	}

	defer resp.Body.Close()

	var result lib.Result
	body, err := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		s.log.Errorf("failed to unmarshal response body: %v", err)

		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	s.log.Infof("Got value: %s for key: %s", result.Value, key)

	return result.Value, nil
}

func (s *StoreClient) Set(ctx context.Context, key, value string) error {
	ctx, span := s.tracer.Start(ctx, "in-client-set")
	defer span.End()

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	url := fmt.Sprintf("http://localhost:4041/set?key=%s&value=%s", key, value)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to set key %s: %s", key, resp.Status)
	}

	return nil
}
