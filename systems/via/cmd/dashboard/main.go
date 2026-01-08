package main

import (
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-via/via"
	"github.com/go-via/via/h"
)

var (
	httpAddr = flag.String("addr", ":8084", "HTTP listen address")
	sseURL   = flag.String("sse-url", "http://localhost:8083", "nats2sse SSE endpoint URL")
	narunURL = flag.String("narun-url", "http://localhost:8082", "narun HTTP API URL")
)

func main() {
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	v := via.New()

	v.Config(via.Options{
		ServerAddress: *httpAddr,
		DocumentTitle: "UAV Swarm Dashboard",
	})

	// Health endpoint
	v.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// SSE proxy - avoids CORS issues by proxying nats2sse through same origin
	v.HandleFunc("GET /sse", func(w http.ResponseWriter, r *http.Request) {
		subject := r.URL.Query().Get("subject")
		if subject == "" {
			subject = "telemetry.>"
		}

		// Proxy to nats2sse
		resp, err := http.Get(*sseURL + "/events?subject=" + subject)
		if err != nil {
			http.Error(w, "SSE upstream error", http.StatusBadGateway)
			logger.Error("SSE proxy error", "error", err)
			return
		}
		defer resp.Body.Close()

		// Copy headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Stream data
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		buf := make([]byte, 4096)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				w.Write(buf[:n])
				flusher.Flush()
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}
		}
	})

	// Add CSS to head
	v.AppendToHead(h.Style(dashboardCSS()))

	// Add SSE script to foot
	v.AppendToFoot(
		h.Script(h.Attr("type", "module"), h.Attr("src", "https://cdn.jsdelivr.net/npm/@starfederation/datastar@1/dist/datastar.min.js")),
		h.Script(h.Text(sseScript(*sseURL))),
	)

	// Main dashboard page
	v.Page("/", func(c *via.Context) {
		c.View(func() h.H {
			return h.Div(
				h.Div(h.Class("header"),
					h.H1(h.Text("UAV Swarm Dashboard")),
				),
				h.Div(h.Class("container"),
					h.Div(h.ID("drones"), h.Class("grid"),
						h.Attr("data-signals", `{drones: [], sseConnected: false}`),
					),
				),
				h.Div(h.Class("sse-status"),
					h.Text("SSE: "),
					h.Span(
						h.Attr("data-class", `{connected: $sseConnected, disconnected: !$sseConnected}`),
						h.Attr("data-text", `$sseConnected ? 'Connected' : 'Disconnected'`),
					),
				),
			)
		})
	})

	logger.Info("Starting Via dashboard", "address", *httpAddr, "sse", *sseURL, "narun", *narunURL)
	v.Start()
}

func dashboardCSS() string {
	return `
* { box-sizing: border-box; margin: 0; padding: 0; }
body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: #1a1a2e;
    color: #eee;
    min-height: 100vh;
}
.header {
    background: #16213e;
    padding: 1rem 2rem;
    border-bottom: 1px solid #0f3460;
}
.header h1 { font-size: 1.5rem; }
.container { padding: 2rem; }
.grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1rem;
}
.card {
    background: #16213e;
    border-radius: 8px;
    padding: 1rem;
    border: 1px solid #0f3460;
}
.card h3 { color: #e94560; margin-bottom: 0.5rem; }
.status { display: flex; gap: 1rem; margin-bottom: 1rem; }
.status-item {
    background: #0f3460;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    font-size: 0.9rem;
}
.armed { color: #4ade80; }
.disarmed { color: #f87171; }
.telemetry { font-family: monospace; font-size: 0.85rem; }
.telemetry div { padding: 0.25rem 0; border-bottom: 1px solid #0f3460; }
.sse-status {
    position: fixed;
    bottom: 1rem;
    right: 1rem;
    padding: 0.5rem 1rem;
    background: #16213e;
    border-radius: 4px;
    font-size: 0.8rem;
}
.connected { color: #4ade80; }
.disconnected { color: #f87171; }
`
}

func sseScript(_ string) string {
	// Use local /sse proxy to avoid CORS issues
	return `
const evtSource = new EventSource('/sse?subject=telemetry.>');
evtSource.onopen = () => {
    const el = document.querySelector('[data-signals]');
    if (el && el.__datastar) el.__datastar.signals.sseConnected = true;
};
evtSource.onerror = () => {
    const el = document.querySelector('[data-signals]');
    if (el && el.__datastar) el.__datastar.signals.sseConnected = false;
};
evtSource.onmessage = (event) => {
    try {
        // nats2sse sends base64 encoded data with "processed: " prefix
        let rawData = event.data;
        try {
            rawData = atob(event.data);
            if (rawData.startsWith('processed: ')) {
                rawData = rawData.substring(11);
            }
        } catch (e) {
            // Not base64, use as-is
        }
        const data = JSON.parse(rawData);
        const el = document.querySelector('[data-signals]');
        if (el && el.__datastar) {
            const signals = el.__datastar.signals;
            const idx = signals.drones.findIndex(d => d.id === data.id);
            if (idx >= 0) {
                signals.drones[idx] = data;
            } else {
                signals.drones.push(data);
            }
        }
    } catch (e) {
        console.error('Failed to parse SSE data:', e);
    }
};
`
}
