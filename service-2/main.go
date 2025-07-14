package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"observability-demo/lib"
	"os"
	"sync"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const HTTPPort = "4041"

func NewServer(controller *Controller) http.Handler {
	mux := http.NewServeMux()

	// handleFunc is a replacement for mux.HandleFunc
	// which enriches the handler's HTTP instrumentation with the pattern as the http.route.
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

	handleFunc("/", controller.ServeHTTP)

	// Add HTTP instrumentation for the whole server.
	handler := otelhttp.NewHandler(mux, "/")
	return handler

}

func run(ctx context.Context) error {
	lib.SetRuntimeSettings("service-2")
	traceProvider, err := lib.GetTracer(context.Background(), lib.Backend)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	log := lib.CreateProductionLogger("service-2")
	defer func() {
		err := log.Sync()
		if err != nil {
			log.Fatal("Error syncing logger")
		}
	}()
	logs := log.Sugar()

	httpSrvLogger := lib.CreateChildLogger(log, "http-server")
	storeLogger := lib.CreateChildLogger(log, "store")

	store := NewMemoryStore(traceProvider.Tracer("store"), storeLogger)
	controller := NewController(traceProvider.Tracer("controller"), store, httpSrvLogger)
	srv := NewServer(controller)

	// Handle SIGINT (CTRL+C) gracefully.
	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer stop()

	httpServer := &http.Server{
		Addr:         net.JoinHostPort("0.0.0.0", HTTPPort),
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      srv,
	}
	go func() {
		logs.Infof("listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Fatalf("error listening and serving: %s", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logs.Errorf("error shutting down http server: %s", err)
		}
	}()
	wg.Wait()
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
