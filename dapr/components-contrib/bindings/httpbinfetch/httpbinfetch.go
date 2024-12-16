package httpbinfetch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/metadata"
	"github.com/dapr/kit/logger"
	kitmd "github.com/dapr/kit/metadata"
)

type Metadata struct {
	Endpoint string `mapstructure:"endpoint"`
	Profile  string `mapstructure:"profile"`
	ApiKey   string `mapstructure:"apikey"`
}

// CustomMiddleware is a struct that implements the middleware interface
type HttBinFetch struct {
	metadata Metadata
	logger   logger.Logger
}

// NewCustomMiddleware creates a new instance of CustomMiddleware
func NewHttBinFetch(logger logger.Logger) bindings.OutputBinding {
	return &HttBinFetch{logger: logger}
}

func (s *HttBinFetch) parseMetadata(meta bindings.Metadata) (Metadata, error) {
	httpBinMeta := Metadata{}
	err := kitmd.DecodeMetadata(meta.Properties, &httpBinMeta)
	// required metadata properties
	if err != nil {
		s.logger.Error("Error parsing input params:", err)
		return httpBinMeta, err
	}

	if httpBinMeta.Endpoint == "" || httpBinMeta.ApiKey == "" {
		s.logger.Debug("Config  params:", httpBinMeta)
		return httpBinMeta, errors.New("smtp binding error: host and port fields are required in metadata")
	}
	return httpBinMeta, nil
}

// Init smtp component (parse metadata).
func (s *HttBinFetch) Init(_ context.Context, metadata bindings.Metadata) error {
	// parse metadata
	meta, err := s.parseMetadata(metadata)
	if err != nil {
		return err
	}
	s.metadata = meta

	return nil
}

// Operations returns the allowed binding operations.
func (s *HttBinFetch) Operations() []bindings.OperationKind {
	return []bindings.OperationKind{bindings.CreateOperation}
}

// Invoke sends an email message.
func (m *HttBinFetch) Invoke(_ context.Context, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	// Merge config metadata with request metadata

	var requestData map[string]interface{}
	if err := json.Unmarshal(req.Data, &requestData); err != nil {
		m.logger.Error("Error unmarshalling JSON:", err)
		return nil, err
	}

	fmt.Printf("Request Body captured : %s\n", requestData)

	requestData["additional_property"] = "HttBinFetch_performed_task"
	requestData["apiKey"] = m.metadata.ApiKey
	requestData["profile"] = m.metadata.Profile

	modifiedBody, err := json.Marshal(requestData)
	if err != nil {
		m.logger.Error("Error marshalling modified JSON:", err)
		return nil, err
	}
	//ioutil.NopCloser(bytes.NewBuffer(modifiedBody))

	reqBody := []byte(modifiedBody)

	resp, err := http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		m.logger.Error("Error making POST request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		m.logger.Error("Error while parsing response :", err)
		return nil, err
	}

	fmt.Printf("Response from external service : %s\n", response)
	returnType := "application/json"
	status := strconv.Itoa(resp.StatusCode)
	invokeResponse := &bindings.InvokeResponse{
		Data:        response, // Modify as needed based on your logic
		Metadata:    map[string]string{"status": status},
		ContentType: &returnType,
	}

	return invokeResponse, nil

	//return nil, nil
}

// GetComponentMetadata returns the metadata of the component.
func (s *HttBinFetch) GetComponentMetadata() (metadataInfo metadata.MetadataMap) {
	metadataStruct := Metadata{}
	metadata.GetMetadataInfoFromStructType(reflect.TypeOf(metadataStruct), &metadataInfo, metadata.BindingType)
	return
}

func (s *HttBinFetch) Close() error {
	return nil
}
