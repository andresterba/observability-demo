# otel example

Run with:

```sh
    docker compose up -d
    make run
    curl -X GET "localhost:4040/?key=test"
    curl -X POST "localhost:4040/?key=test&value=test"
```

You can find the [jaeger UI here](http://localhost:16686/search).
