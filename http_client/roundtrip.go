package http_client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type LoggingTripper struct {
	logger io.Writer
	next   http.RoundTripper
}

func (l LoggingTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	fmt.Fprintf(l.logger, "[%s] %s %s\n", time.Now().Format(time.ANSIC), r.Method, r.URL)
	return l.next.RoundTrip(r)
}
