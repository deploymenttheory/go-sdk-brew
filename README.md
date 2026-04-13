# Go SDK for Homebrew Formulae API

[![Go Report Card](https://goreportcard.com/badge/github.com/deploymenttheory/go-sdk-brew)](https://goreportcard.com/report/github.com/deploymenttheory/go-sdk-brew)
[![GoDoc](https://pkg.go.dev/badge/github.com/deploymenttheory/go-sdk-brew)](https://pkg.go.dev/github.com/deploymenttheory/go-sdk-brew)
[![License](https://img.shields.io/github/license/deploymenttheory/go-sdk-brew)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/deploymenttheory/go-sdk-brew)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/deploymenttheory/go-sdk-brew)](https://github.com/deploymenttheory/go-sdk-brew/releases)
![Status: Experimental](https://img.shields.io/badge/status-experimental-yellow)

A Go client library for the [Homebrew Formulae API](https://formulae.brew.sh/docs/api/). No authentication required — the API is public and read-only. Includes a production-ready transport with retries, adaptive rate limiting, concurrency control, structured logging, and optional OpenTelemetry tracing.

## Quick Start

```bash
go get github.com/deploymenttheory/go-sdk-brew
```

```go
import (
    "context"
    "fmt"
    "log"

    brew "github.com/deploymenttheory/go-sdk-brew/brew"
    "github.com/deploymenttheory/go-sdk-brew/brew/services/analytics"
)

func main() {
    client, err := brew.NewClient("") // defaults to https://formulae.brew.sh
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // List all formulae
    formulae, _, err := client.Formulae.List(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Total formulae: %d\n", len(formulae))

    // Get a formula by name
    wget, _, err := client.Formulae.GetByName(ctx, "wget")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("wget %s: %s\n", wget.Versions.Stable, wget.Desc)

    // Install analytics for the last 30 days
    stats, _, err := client.Analytics.ListByCategory(ctx, analytics.CategoryInstall, analytics.Period30d)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Top formula: %s (%s installs)\n", stats.Items[0].Formula, stats.Items[0].Count)
}
```

## API Coverage

The SDK covers all endpoints of the [Homebrew Formulae API](https://formulae.brew.sh/docs/api/):

### Formulae (`client.Formulae`)

| Method | Endpoint | Description |
|---|---|---|
| `List(ctx)` | `GET /api/formula.json` | All Homebrew formulae |
| `GetByName(ctx, name)` | `GET /api/formula/{name}.json` | Single formula by name |

### Casks (`client.Casks`)

| Method | Endpoint | Description |
|---|---|---|
| `List(ctx)` | `GET /api/cask.json` | All Homebrew casks |
| `GetByName(ctx, name)` | `GET /api/cask/{name}.json` | Single cask by name |

### Analytics (`client.Analytics`)

| Method | Endpoint | Description |
|---|---|---|
| `ListByCategory(ctx, category, period)` | `GET /api/analytics/{category}/{period}.json` | Ranked install analytics (items array) |
| `ListCoreByCategory(ctx, category, period)` | `GET /api/analytics/{category}/homebrew-core/{period}.json` | homebrew-core analytics (formulae map) |
| `ListCaskInstalls(ctx, period)` | `GET /api/analytics/cask-install/homebrew-cask/{period}.json` | Cask install analytics (formulae map) |

**Analytics categories:** `analytics.CategoryInstall`, `analytics.CategoryInstallOnRequest`, `analytics.CategoryBuildError`

**Analytics periods:** `analytics.Period30d`, `analytics.Period90d`, `analytics.Period365d`

## Examples

The [examples directory](examples/) contains a working example demonstrating all three services:

- **[examples/basic/main.go](examples/basic/main.go)** — Lists formulae, gets wget, lists casks, gets iterm2, and fetches install analytics

## Configuration Options

### Creating a client

```go
import brew "github.com/deploymenttheory/go-sdk-brew/brew"

// Default base URL (https://formulae.brew.sh)
client, err := brew.NewClient("")

// Custom base URL
client, err := brew.NewClient("https://formulae.brew.sh")

// With options
client, err := brew.NewClient("",
    brew.WithTimeout(30*time.Second),
    brew.WithRetryCount(3),
    brew.WithLogger(logger),
)
```

### Available options

All configuration is passed via functional options to `NewClient`.

#### Basic Configuration

```go
brew.WithBaseURL("https://formulae.brew.sh")          // Override base URL
brew.WithTimeout(30*time.Second)                       // Request timeout (default: 30s)
brew.WithRetryCount(3)                                 // Retry attempts on transient errors (default: 3)
brew.WithRetryWaitTime(2*time.Second)                  // Initial retry wait (default: 2s)
brew.WithRetryMaxWaitTime(30*time.Second)              // Max retry wait with backoff (default: 30s)
brew.WithTotalRetryDuration(2*time.Minute)             // Total retry budget across all attempts
```

#### TLS / Security

```go
brew.WithTLSClientConfig(tlsConfig)                    // Custom *tls.Config
brew.WithInsecureSkipVerify()                          // Skip TLS verification (dev only)
```

#### Network

```go
brew.WithProxy("http://proxy.example.com:8080")        // HTTP/HTTPS/SOCKS5 proxy
brew.WithTransport(customTransport)                    // Custom net/http Transport
```

#### Headers

```go
brew.WithUserAgent("MyApp/1.0")                        // Set User-Agent
brew.WithGlobalHeader("X-Custom", "value")             // Add a single global header
brew.WithGlobalHeaders(map[string]string{...})         // Add multiple global headers
```

#### Observability

```go
brew.WithLogger(zapLogger)                             // Structured logging with go.uber.org/zap
brew.WithDebug()                                       // Verbose request/response logging (dev only)
```

OpenTelemetry tracing is configured at the transport level via `brew.WithTransport`. Wrap your existing transport with `otelhttp.NewTransport` before passing it in, or use `brew.WithOpenTelemetry` if enabled in your build.

#### Concurrency & Rate Limiting

```go
brew.WithMaxConcurrentRequests(10)                     // Max parallel in-flight requests (default: 10)
brew.WithMandatoryRequestDelay(50*time.Millisecond)    // Fixed delay between requests
```

The transport also applies an EMA-based adaptive delay automatically when API response times degrade.

### Example: Production Configuration

```go
import (
    "time"

    "go.uber.org/zap"
    brew "github.com/deploymenttheory/go-sdk-brew/brew"
)

logger, _ := zap.NewProduction()

client, err := brew.NewClient("",
    brew.WithTimeout(30*time.Second),
    brew.WithRetryCount(3),
    brew.WithRetryWaitTime(2*time.Second),
    brew.WithRetryMaxWaitTime(30*time.Second),
    brew.WithMaxConcurrentRequests(10),
    brew.WithLogger(logger),
    brew.WithUserAgent("my-app/1.0"),
)
```

## Response Shape

All service methods return `(result, *resty.Response, error)`:

```go
formula, resp, err := client.Formulae.GetByName(ctx, "wget")
if err != nil {
    // Handle error — includes API errors (404, 500, etc.) and network errors
}

fmt.Println(resp.StatusCode()) // 200
fmt.Println(formula.Name)      // "wget"
fmt.Println(formula.Versions.Stable)
fmt.Println(formula.Dependencies)
```

### Analytics response shapes

The three analytics endpoints return different Go types because the Homebrew API uses different JSON shapes for each:

```go
// General analytics: ranked items array
result, _, _ := client.Analytics.ListByCategory(ctx, analytics.CategoryInstall, analytics.Period30d)
for _, item := range result.Items {
    fmt.Printf("#%d %s — %s installs (%s%%)\n", item.Number, item.Formula, item.Count, item.Percent)
}

// homebrew-core: formulae map (formula name → []CoreFormulaItem)
core, _, _ := client.Analytics.ListCoreByCategory(ctx, analytics.CategoryInstall, analytics.Period30d)
for name, entries := range core.Formulae {
    fmt.Printf("%s: %s installs\n", name, entries[0].Count)
}

// Cask installs: formulae map (cask token → []CaskInstallItem)
casks, _, _ := client.Analytics.ListCaskInstalls(ctx, analytics.Period30d)
for token, entries := range casks.Formulae {
    fmt.Printf("%s: %s installs\n", token, entries[0].Count)
}
```

## Testing

```bash
# Unit tests (no network required — fixture-based mocks)
go test ./brew/services/...

# Acceptance tests (hits live https://formulae.brew.sh)
go test -v ./brew/acceptance/... -timeout 120s

# Optional env vars for acceptance tests
BREW_BASE_URL=https://formulae.brew.sh   # override API base URL
BREW_REQUEST_TIMEOUT=30s                 # per-test context timeout
BREW_VERBOSE=true                        # enable detailed logging
```

## Documentation

- [Homebrew Formulae API Reference](https://formulae.brew.sh/docs/api/)
- [GoDoc](https://pkg.go.dev/github.com/deploymenttheory/go-sdk-brew)

## Contributing

Contributions are welcome. Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting pull requests.

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

## Support

- **Issues:** [GitHub Issues](https://github.com/deploymenttheory/go-sdk-brew/issues)
- **Homebrew API docs:** [formulae.brew.sh/docs/api](https://formulae.brew.sh/docs/api/)

## Disclaimer

This is a community SDK and is not affiliated with or endorsed by the Homebrew project.
