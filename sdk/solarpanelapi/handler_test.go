package solarpanelapi_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WilliamDeKeyzer/bdd-workshop-sdk/sdk/solarpanelapi"
	"github.com/stretchr/testify/require"
)

func TestHandler_Generation(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		method string
		body   string
		status int
		want   string
	}{
		"new panel":        {http.MethodGet, "", 200, `{"watts":0}`},
		"valid reading":    {http.MethodPut, `{"watts":700}`, 200, `{"watts":700}`},
		"zero reading":     {http.MethodPut, `{"watts":0}`, 400, `{"error":"invalid watt generation"}`},
		"negative reading": {http.MethodPut, `{"watts":-1}`, 400, `{"error":"invalid watt generation"}`},
		"malformed JSON":   {http.MethodPut, `{`, 400, `{"error":"invalid JSON body"}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			api := solarpanelapi.NewHandler()
			request := httptest.NewRequest(test.method, "/generation", strings.NewReader(test.body))
			response := httptest.NewRecorder()

			api.ServeHTTP(response, request)

			require.Equal(t, test.status, response.Code)
			require.Equal(t, "application/json", response.Header().Get("Content-Type"))
			require.JSONEq(t, test.want, response.Body.String())
			fetched := httptest.NewRecorder()
			api.ServeHTTP(fetched, httptest.NewRequest(http.MethodGet, "/generation", nil))
			wantReading := `{"watts":0}`
			if test.status == http.StatusOK && test.method == http.MethodPut {
				wantReading = test.want
			}
			require.JSONEq(t, wantReading, fetched.Body.String())
		})
	}
}

func TestHandler_RejectedUpdatePreservesReading(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"zero": `{"watts":0}`, "negative": `{"watts":-1}`, "malformed": `{`,
	}
	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			api := solarpanelapi.NewHandler()
			initial := httptest.NewRecorder()
			api.ServeHTTP(initial, httptest.NewRequest(http.MethodPut, "/generation", strings.NewReader(`{"watts":700}`)))
			require.Equal(t, http.StatusOK, initial.Code)

			rejected := httptest.NewRecorder()
			api.ServeHTTP(rejected, httptest.NewRequest(http.MethodPut, "/generation", strings.NewReader(body)))

			require.Equal(t, http.StatusBadRequest, rejected.Code)
			fetched := httptest.NewRecorder()
			api.ServeHTTP(fetched, httptest.NewRequest(http.MethodGet, "/generation", nil))
			require.JSONEq(t, `{"watts":700}`, fetched.Body.String())
		})
	}
}
