package httpbinfetch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dapr/components-contrib/middleware"
	"github.com/dapr/kit/logger"
)

// CustomMiddleware is a struct that implements the middleware interface
type HttBinFetchMiddleware struct {
	logger logger.Logger
}

// NewCustomMiddleware creates a new instance of CustomMiddleware
func NewHttBinFetchMiddleware(logger logger.Logger) middleware.Middleware {
	m := &HttBinFetchMiddleware{
		logger: logger,
	}
	return m
}

// GetHandler returns a handler that processes incoming requests
func (m *HttBinFetchMiddleware) GetHandler(_ context.Context, metadata middleware.Metadata) (func(next http.Handler) http.Handler, error) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Inbound logic: Log incoming request

			// Prepare the POST request to the external service
			if r.Method == "POST" {

				var jsonData map[string]interface{}
				requestBodyFromIncoming, err := io.ReadAll(r.Body)
				if err != nil {
					m.logger.Error("Error capturing incoming request :", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				fmt.Printf("Request Body captured : %s\n", requestBodyFromIncoming)

				if err := json.Unmarshal(requestBodyFromIncoming, &jsonData); err != nil {
					m.logger.Error("Error unmarshalling JSON:", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				jsonData["additional_property"] = "HttBinFetchMiddleware_performed_task"

				modifiedBody, err := json.Marshal(jsonData)
				if err != nil {
					m.logger.Error("Error marshalling modified JSON:", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				//ioutil.NopCloser(bytes.NewBuffer(modifiedBody))

				reqBody := []byte(modifiedBody)

				resp, err := http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(reqBody))
				if err != nil {
					m.logger.Error("Error making POST request:", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				defer resp.Body.Close()

				// Read response from the external service
				//body, err := ioutil.ReadAll(resp.Body)
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					m.logger.Error("Error reading response body:", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				fmt.Printf("Response from external service : %s\n", body)

				// Forward the response to the client
				w.WriteHeader(resp.StatusCode)
				w.Write(body)
			}

			// Call the next handler in the chain
			next.ServeHTTP(w, r)
		})
	}, nil
}
