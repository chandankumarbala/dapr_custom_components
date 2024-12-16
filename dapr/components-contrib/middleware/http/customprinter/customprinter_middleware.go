package customprinter

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/dapr/kit/logger"

	"github.com/dapr/components-contrib/middleware"
)

// CustomMiddleware struct
type CustomMiddleware struct {
	log logger.Logger
}

func NewCustomMiddleware(logger logger.Logger) middleware.Middleware {
	m := &CustomMiddleware{
		log: logger,
	}
	return m
}

func (m *CustomMiddleware) GetHandler(_ context.Context, metadata middleware.Metadata) (func(next http.Handler) http.Handler, error) {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Unable to read request body", http.StatusBadRequest)
				return
			}
			fmt.Printf("Request Body: %s\n", body)
			m.log.Debugf("Request Body: %s\n", body)
			// Restore the io.ReadCloser to allow further reading
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			// Call the next handler in the chain
			next.ServeHTTP(w, r)
		})
	}, nil
}

// GetHandler implements the middleware interface
