package log_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/stretchr/testify/require"
)

func TestRequestLoggerGeneratesCorrelationIDWhenNotFoundInHeaders(t *testing.T) {
	response := httptest.NewRecorder()

	testUrl, err := url.Parse("http://localhost:8080")
	require.NoError(t, err)
	request := &http.Request{
		Method: http.MethodPost,
		URL:    testUrl,
		Header: map[string][]string{},
	}

	handler := log.RequestLogger()
	handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		entry := log.C(request.Context())

		correlationIDFromLogger, exists := entry.Data[log.FieldRequestID]
		require.True(t, exists)
		require.NotEmpty(t, correlationIDFromLogger)
	})).ServeHTTP(response, request)
}

func TestRequestLoggerUseCorrelationIDFromHeaderIfProvided(t *testing.T) {
	correlationID := "test-correlation-id"
	response := httptest.NewRecorder()

	testUrl, err := url.Parse("http://localhost:8080")
	require.NoError(t, err)
	request := &http.Request{
		Method: http.MethodPost,
		URL:    testUrl,
		Header: map[string][]string{},
	}
	request.Header.Set("x-request-id", correlationID)
	request.Header.Set("X-Real-IP", "127.0.0.1")

	handler := log.RequestLogger()
	handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		entry := log.C(request.Context())

		correlationIDFromLogger, exists := entry.Data[log.FieldRequestID]
		require.True(t, exists)
		require.Equal(t, correlationID, correlationIDFromLogger)
	})).ServeHTTP(response, request)
}
