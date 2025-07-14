package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"observability-demo/lib"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

const SERVER_ADDRESS = "http://localhost:4040"

var tmpl = template.Must(template.New("ui").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Key-Value Service</title>
</head>
<body>
	<h1>Key-Value Service</h1>
	<h2>Set Key-Value Pair</h2>
	<form method="POST" action="/set">
		<label for="key">Key:</label>
		<input type="text" id="key" name="key" required>
		<br>
		<label for="value">Value:</label>
		<input type="text" id="value" name="value" required>
		<br>
		<button type="submit">Set</button>
	</form>
	<h2>Get Value by Key</h2>
	<form method="GET" action="/get">
		<label for="key">Key:</label>
		<input type="text" id="key" name="key" required>
		<br>
		<button type="submit">Get</button>
	</form>
	{{if .Response}}
		<h3>Response:</h3>
		<p>{{.Response}}</p>
	{{end}}
</body>
</html>
`))

var traceClient trace.Tracer

func main() {
	lib.SetRuntimeSettings("ui")
	traceProvider, err := lib.GetTracer(context.Background(), lib.Backend)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	traceClient = traceProvider.Tracer("ui")

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)

	fmt.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	_, span := traceClient.Start(context.Background(), "set")
	defer span.End()

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	key := r.FormValue("key")
	value := r.FormValue("value")

	// Make POST request to external service
	externalURL := fmt.Sprintf("%s?key=%s&value=%s", SERVER_ADDRESS, key, value)
	resp, err := client.Post(externalURL, "application/json", nil)
	if err != nil {
		http.Error(w, "Failed to make POST request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	tmpl.Execute(w, map[string]string{"Response": string(body)})
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	_, span := traceClient.Start(context.Background(), "get")
	defer span.End()

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")

	// Make GET request to external service
	externalURL := fmt.Sprintf("%s?key=%s", SERVER_ADDRESS, key)
	resp, err := client.Get(externalURL)
	if err != nil {
		http.Error(w, "Failed to make GET request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	tmpl.Execute(w, map[string]string{"Response": string(body)})
}
