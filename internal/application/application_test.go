package application_test

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/IDK536/go-calc/internal/application"
)

func TestCalcHandlerSuccessCase(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		body           map[string]string
		statusCode     int
		expectedResult string
	}{
		{
			name: "valid expression",
			body: map[string]string{
				"expression": "2+2",
			},
			statusCode:     200,
			expectedResult: `result: 4.000000`,
		},
		{
			name: "valid expression",
			body: map[string]string{
				"expression": "2+2*(2+2)",
			},
			statusCode:     200,
			expectedResult: `result: 10.000000`,
		},
		{
			name: "invalid expression",
			body: map[string]string{
				"expression": "2+",
			},
			statusCode:     422,
			expectedResult: `err: invalid expression`,
		},
		{
			name: "invalid expression",
			body: map[string]string{
				"expression": "2+2h",
			},
			statusCode:     422,
			expectedResult: `err: invalid expression`,
		},
		{
			name: "invalid expression",
			body: map[string]string{
				"expression": "2++2",
			},
			statusCode:     422,
			expectedResult: `err: invalid expression`,
		},
		{
			name: "division by zero",
			body: map[string]string{
				"expression": "2/0",
			},
			statusCode:     422,
			expectedResult: `err: division by zero`,
		},
		{
			name: "division by zero",
			body: map[string]string{
				"expression": "2/(19-10-9)",
			},
			statusCode:     422,
			expectedResult: `err: division by zero`,
		},
		{
			name: "empty expression",
			body: map[string]string{
				"expression": "",
			},
			statusCode:     422,
			expectedResult: `err: empty expression`,
		},
		{
			name:       "internal server error",
			body:       nil,
			statusCode: 500,
			expectedResult: `EOF
err: internal server error`,
		},
	}
	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(testCase.body)
			var req *http.Request
			if testCase.body == nil {
				req = httptest.NewRequest("POSR", "/api/v1/calculate", nil)
			} else {
				req = httptest.NewRequest("POSR", "/api/v1/calculate", bytes.NewReader(bodyBytes))
			}
			w := httptest.NewRecorder()
			application.CalcHandler(w, req)
			resp := w.Result()
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			if resp.StatusCode != testCase.statusCode {
				t.Errorf("expected status code %d, got %d", testCase.statusCode, resp.StatusCode)
			}
			if string(body) != testCase.expectedResult {
				t.Errorf("expected body %s, got %s", testCase.expectedResult, string(body))
			}
		})
	}
}
