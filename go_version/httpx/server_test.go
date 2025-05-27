package httpx

import (
	"bill-splitter/accounting"
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func Test_Server_Integration(t *testing.T) {
	baseUrl := "http://localhost:8000"

	scenarios := []struct {
		name         string
		makeRequest  func() *http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name: "should process transactions and return balances",
			makeRequest: func() *http.Request {
				b := bytes.NewBuffer([]byte(`[{ "from": "A", "to": "B", "amount": 40 },{ "from": "B", "to": "C", "amount": 40 },{ "from": "C", "to": "A", "amount": 10 }]`))
				r, _ := http.NewRequest("POST", baseUrl+"/balance/calculate", b)
				r.Header = map[string][]string{"Content-Type": {"application/json"}}
				return r
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"name":"A","amount":30},{"name":"B","amount":0},{"name":"C","amount":-30}]`,
		},
		{
			name: "should process transactions and return balances",
			makeRequest: func() *http.Request {
				b := bytes.NewBuffer([]byte(`[{ "from": "A", "to": "B", "amount": 40.0 },{ "from": "A", "to": "C", "amount": 10.0 }]`))
				r, _ := http.NewRequest("POST", baseUrl+"/balance/calculate", b)
				return r
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid `Content-Type` header. Expected `application/json`",
		},
		{
			name: "should process balances and reduce transactions",
			makeRequest: func() *http.Request {
				b := bytes.NewBuffer([]byte(`[{"name":"A","amount":30.0},{"name":"B","amount":0.0},{"name":"C","amount":-30.0}]`))
				r, _ := http.NewRequest("POST", baseUrl+"/transaction/minimize", b)
				r.Header = map[string][]string{"Content-Type": {"application/json"}}
				return r
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"updated_balances":[{"name":"C","amount":0},{"name":"B","amount":0},{"name":"A","amount":0}],"transactions":[{"from":"C","to":"A","amount":30}]}`,
		},
		{
			name: "should return error when no content type",
			makeRequest: func() *http.Request {
				b := bytes.NewBuffer([]byte(`[{"name":"A","amount":30.0},{"name":"B","amount":0.0},{"name":"C","amount":-30.0}]`))
				r, _ := http.NewRequest("POST", baseUrl+"/transaction/minimize", b)
				return r
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid `Content-Type` header. Expected `application/json`",
		},
		{
			name: "should return error when no body",
			makeRequest: func() *http.Request {
				b := bytes.NewBuffer([]byte(""))
				r, _ := http.NewRequest("POST", baseUrl+"/transaction/minimize", b)
				r.Header = map[string][]string{"Content-Type": {"application/json"}}
				return r
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "failed to read request body: EOF",
		},
	}

	server := setup()
	defer server.Close()

	extractCodeAndBody := func(response *http.Response) (int, []byte) {
		defer response.Body.Close()
		actualCode := response.StatusCode
		actualBody, _ := io.ReadAll(response.Body)
		return actualCode, actualBody
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {

			req := s.makeRequest()

			// default client is enough for testing
			response, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("unexpected error doing request: %+v", err)
			}

			actualCode, actualBody := extractCodeAndBody(response)

			if s.expectedCode != actualCode {
				t.Errorf("\nExpected:	%+v\nGot:		%+v", s.expectedCode, actualCode)
			}

			// order of objects in json might change
			if strings.TrimSpace(s.expectedBody) != strings.TrimSpace(string(actualBody)) {
				t.Errorf("\nExpected:	%+v\nGot:		%+v", s.expectedBody, string(actualBody))
			}
		})
	}
}

func setup() *HttpServer {
	accService := accounting.NewService()
	s := NewServer(accService, accService)
	go func() {
		s.Run()
	}()
	return s
}
