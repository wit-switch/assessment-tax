package tax_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tax Handler Suite")
}

func compacJSON(jsonStr string) ([]byte, error) {
	buffer := bytes.Buffer{}
	if err := json.Compact(&buffer, []byte(jsonStr)); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func request(route string, body io.Reader, e *echo.Echo) (int, string) {
	req := httptest.NewRequest(http.MethodPost, route, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	// custom error handler only plays with real request
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}
