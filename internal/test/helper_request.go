package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
	"test-project/internal/api"
)

type GenericPayload map[string]interface{}
type GenericArrayPayload []interface{}

func (g GenericPayload) Reader(t TestingT) *bytes.Reader {
	t.Helper()

	b, err := json.Marshal(g)
	if err != nil {
		t.Fatalf("failed to serialize payload: %v", err)
	}

	return bytes.NewReader(b)
}

func (g GenericArrayPayload) Reader(t TestingT) *bytes.Reader {
	t.Helper()

	b, err := json.Marshal(g)
	if err != nil {
		t.Fatalf("failed to serialize payload: %v", err)
	}

	return bytes.NewReader(b)
}

func PerformRequestWithParams(t TestingT, s *api.Server, method string, path string, body GenericPayload, headers http.Header, queryParams map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	if body == nil {
		return PerformRequestWithRawBody(t, s, method, path, nil, headers, queryParams)
	}

	return PerformRequestWithRawBody(t, s, method, path, body.Reader(t), headers, queryParams)
}

func PerformRequestWithArrayAndParams(t TestingT, s *api.Server, method string, path string, body GenericArrayPayload, headers http.Header, queryParams map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	if body == nil {
		return PerformRequestWithRawBody(t, s, method, path, nil, headers, queryParams)
	}

	return PerformRequestWithRawBody(t, s, method, path, body.Reader(t), headers, queryParams)
}

func PerformRequestWithRawBody(t TestingT, s *api.Server, method string, path string, body io.Reader, headers http.Header, queryParams map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, body)

	if headers != nil {
		req.Header = headers
	}
	if body != nil && len(req.Header.Get(echo.HeaderContentType)) == 0 {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	if queryParams != nil {
		q := req.URL.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}

	res := httptest.NewRecorder()

	s.Echo.ServeHTTP(res, req)

	return res
}

func PerformRequest(t TestingT, s *api.Server, method string, path string, body GenericPayload, headers http.Header) *httptest.ResponseRecorder {
	t.Helper()

	return PerformRequestWithParams(t, s, method, path, body, headers, nil)
}

func PerformRequestWithArray(t TestingT, s *api.Server, method string, path string, body GenericArrayPayload, headers http.Header) *httptest.ResponseRecorder {
	t.Helper()

	return PerformRequestWithArrayAndParams(t, s, method, path, body, headers, nil)
}

func ParseResponseBody(t TestingT, res *httptest.ResponseRecorder, v interface{}) {
	t.Helper()

	if err := json.NewDecoder(res.Result().Body).Decode(&v); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
}

func ParseResponseAndValidate(t TestingT, res *httptest.ResponseRecorder, v runtime.Validatable) {
	t.Helper()

	ParseResponseBody(t, res, &v)

	if err := v.Validate(strfmt.Default); err != nil {
		t.Fatalf("Failed to validate response: %v", err)
	}
}

func HeadersWithAuth(t TestingT, token string) http.Header {
	t.Helper()

	return HeadersWithConfigurableAuth(t, "Bearer", token)
}

func HeadersWithConfigurableAuth(t TestingT, scheme string, token string) http.Header {
	t.Helper()

	headers := http.Header{}
	headers.Set(echo.HeaderAuthorization, fmt.Sprintf("%s %s", scheme, token))

	return headers
}

func HeadersWithAPIKeyAuth(t TestingT, token string) http.Header {
	t.Helper()

	return HeadersWithConfigurableAuth(t, "APIKey", token)
}
