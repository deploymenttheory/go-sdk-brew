package client

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// applyOpenTelemetry wraps the HTTP transport with OpenTelemetry instrumentation.
// This is always enabled and uses the global OpenTelemetry providers set via:
//   - otel.SetTracerProvider()
//   - otel.SetMeterProvider()
//   - otel.SetTextMapPropagator()
//
// If no global providers are configured, the instrumentation is a no-op.
func (t *Transport) applyOpenTelemetry() {
	httpClient := t.client.Client()
	if httpClient == nil {
		return
	}

	transport := httpClient.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	instrumentedTransport := otelhttp.NewTransport(transport)
	httpClient.Transport = instrumentedTransport

	t.logger.Debug("OpenTelemetry HTTP instrumentation enabled (uses global providers)")
}
