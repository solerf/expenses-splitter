package httpx

import (
	"bill-splitter/accounting"
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Main_Handler_Func(t *testing.T) {
	scenarios := []struct {
		name          string
		customHandler customHandler
		expectedCode  int
	}{
		{
			name: "when no error",
			customHandler: func(writer http.ResponseWriter, request *http.Request) error {
				return nil
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "when internal server error",
			customHandler: func(writer http.ResponseWriter, request *http.Request) error {
				return errors.New("internal server error")
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "when bad request",
			customHandler: func(writer http.ResponseWriter, request *http.Request) error {
				return invalidContentType
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			mainHandlerFunc(s.customHandler)(
				recorder,
				httptest.NewRequest("POST", "/any", http.NoBody),
			)

			if s.expectedCode != recorder.Code {
				t.Errorf("\nExpected:	%+v\nGot:		%+v", s.expectedCode, recorder.Code)
			}
		})
	}
}

func Test_Validate_Content_Type(t *testing.T) {
	scenarios := []struct {
		name          string
		modifyRequest func(r *http.Request) *http.Request
		expectedError error
	}{
		{
			name: "when valid content type",
			modifyRequest: func(r *http.Request) *http.Request {
				r.Header.Set("Content-Type", "application/json")
				return r
			},
			expectedError: nil,
		},
		{
			name: "when invalid content type",
			modifyRequest: func(r *http.Request) *http.Request {
				r.Header.Set("Content-Type", "application/octet-stream")
				return r
			},
			expectedError: invalidContentType,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			err := validateContentType(nil)(
				httptest.NewRecorder(),
				s.modifyRequest(httptest.NewRequest("POST", "/any", http.NoBody)),
			)
			if !errors.Is(err, s.expectedError) {
				t.Errorf("\nExpected:	%+v\nGot:		%+v", s.expectedError, err)
			}
		})
	}
}

func Test_Balance_Calculate(t *testing.T) {
	scenarios := []struct {
		name          string
		bodyReader    io.Reader
		expectedError error
	}{
		{
			name:          "when valid json",
			bodyReader:    bytes.NewBuffer([]byte(`[{ "from": "A", "to": "B", "amount": 40.0 },{ "from": "A", "to": "C", "amount": 10.0 }]`)),
			expectedError: nil,
		},
		{
			name:          "when invalid json",
			bodyReader:    bytes.NewBuffer([]byte(`{ "from": "A", "to": "B", "amount": 40.0  }`)),
			expectedError: errors.New("failed to read request body: json: cannot unmarshal object into Go value of type []accounting.Transaction"),
		},
	}

	stubService := balanceServiceStub(func(_ accounting.Transactions) accounting.Balances {
		return accounting.Balances{
			{"A", 30.0},
			{"B", 0.0},
			{"C", -30.0},
		}
	})

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			err := balanceCalculate(stubService)(
				httptest.NewRecorder(),
				httptest.NewRequest("POST", "/any", s.bodyReader),
			)

			if !errors.Is(err, s.expectedError) && s.expectedError.Error() != err.Error() {
				t.Errorf("\nExpected:	%+v\nGot:		%+v", s.expectedError, err)
			}
		})
	}
}

func Test_Transactions_Minimize(t *testing.T) {
	scenarios := []struct {
		name          string
		bodyReader    io.Reader
		expectedError error
	}{
		{
			name:          "when valid json",
			bodyReader:    bytes.NewBuffer([]byte(`[{ "name": "A", "amount": 30.0 },{ "name": "B", "amount": 0.0 },{ "name": "C", "amount": -30.0 }]`)),
			expectedError: nil,
		},
		{
			name:          "when invalid json",
			bodyReader:    bytes.NewBuffer([]byte(`{ "from": "A", "to": "B", "amount": 40.0  }`)),
			expectedError: errors.New("failed to read request body: json: cannot unmarshal object into Go value of type []accounting.Balance"),
		},
	}

	stubService := transactionServiceStub(func(_ accounting.Balances) accounting.Statement {
		return accounting.Statement{
			UpdatedBalances: accounting.Balances{
				{"A", 0.0},
				{"B", 0.0},
				{"C", .0},
			},
			Transactions: accounting.Transactions{
				{From: "C", To: "A", Amount: 30.0},
			},
		}
	})

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			err := minimizeTransaction(stubService)(
				httptest.NewRecorder(),
				httptest.NewRequest("POST", "/any", s.bodyReader),
			)

			if !errors.Is(err, s.expectedError) && s.expectedError.Error() != err.Error() {
				t.Errorf("\nExpected:	%+v\nGot:		%+v", s.expectedError, err)
			}
		})
	}
}

type balanceServiceStub func(a accounting.Transactions) accounting.Balances

func (bss balanceServiceStub) Calculate(a accounting.Transactions) accounting.Balances {
	return bss(a)
}

type transactionServiceStub func(b accounting.Balances) accounting.Statement

func (tss transactionServiceStub) Minimize(b accounting.Balances) accounting.Statement {
	return tss(b)
}
