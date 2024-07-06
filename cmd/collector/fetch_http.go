package main

import (
	"context"
	"net/http"
	"slog"

	"github.com/google/uuid"
)

type HTTPFetcher struct {
	client *http.Client
}

// FetchHTTPSource retrieves src (protocols: http://, https://) using GET method and returns the bytes and nil error.
// Otherwise, an appropriate error is returned. It's possible to customize the request using the options.
func FetchHTTPSource(ctx context.Context, src string, options ...HTTPFetchOption) ([]byte, error) {
	// Try to make sense of the provided src
	if !(strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://")) {
		return nil, fmt.Errorf("invalid prefix, expected http:// or https://")
	}
	if *appname == "" { // flags
		*appname = "FIXME-to-be-nice"
	}
	options = append(options, WithUserAgent(*appname))

	// Generate a request ID (UUID), add it to the request and insert it into the context
	reqId := uuid.New() // may panic
	options = append(options, WithRequestIdAndAppname(reqId, *appname))

	ctx = NewContextWithRequestId(reqId)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, src, nil)
	if err != nil {
		slog.Error("FetchHTTPSource request creation failed", "error", err, slog.Group("request", "method", method, "src", src))
		return nil, fmt.Errorf("creating request failed: %v", err)
	}

	for opt := range options {
		err = opt(req)
		if err != nil {
			slog.Error("FetchHTTPSource option setting failed", "error", err, slog.Group("request", "method", method, "src", src))
			return nil, fmt.Errorf("setting HTTPFetchOption failed: %v", err)
		}
	}

	slog.Debug("request created", "request-id", &reqId, slog.Group("request", "method", req.Method, "url", &req.URL))

	return nil, fmt.Errorf("unimplemented")
}

// vim: cc=120:
