golangci-lint-version := "v2.1"

lint:
    docker run --rm -t -v $(pwd):/app -w /app \
    --user $(id -u):$(id -g) \
    -v $(go env GOCACHE):/.cache/go-build -e GOCACHE=/.cache/go-build \
    -v $(go env GOMODCACHE):/.cache/mod -e GOMODCACHE=/.cache/mod \
    -v ~/.cache/golangci-lint:/.cache/golangci-lint -e GOLANGCI_LINT_CACHE=/.cache/golangci-lint \
    golangci/golangci-lint:{{ golangci-lint-version }} golangci-lint run

lint-fix:
    docker run --rm -t -v $(pwd):/app -w /app \
    --user $(id -u):$(id -g) \
    -v $(go env GOCACHE):/.cache/go-build -e GOCACHE=/.cache/go-build \
    -v $(go env GOMODCACHE):/.cache/mod -e GOMODCACHE=/.cache/mod \
    -v ~/.cache/golangci-lint:/.cache/golangci-lint -e GOLANGCI_LINT_CACHE=/.cache/golangci-lint \
    golangci/golangci-lint:{{ golangci-lint-version }} golangci-lint run --fix

typos:
    typos .

typos-fix:
    typos -w .
