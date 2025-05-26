package httpx

import (
	"bill-splitter/accounting"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
)

var invalidContentType = errors.New("Invalid `Content-Type` header. Expected `application/json`")

type customHandler func(http.ResponseWriter, *http.Request) error

// register register in the http.ServerMux all endpoints
func register(mux *http.ServeMux, balanceService BalanceService, transactionService TransactionService) {
	mux.HandleFunc(
		"POST /balance/calculate",
		mainHandlerFunc(validateContentType(balanceCalculate(balanceService))),
	)

	mux.HandleFunc(
		"POST /transaction/minimize",
		mainHandlerFunc(validateContentType(minimizeTransaction(transactionService))),
	)

	mux.HandleFunc("/", http.NotFound)
}

// mainHandlerFunc generic handler for ServerMux
func mainHandlerFunc(innerHandler customHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if err := innerHandler(writer, request); err != nil {
			switch {
			case errors.Is(err, invalidContentType):
				http.Error(writer, err.Error(), http.StatusBadRequest)
			default:
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

// validateContentType handler to just validate if the request has the correct content type
func validateContentType(next customHandler) customHandler {
	return func(writer http.ResponseWriter, request *http.Request) error {
		if !slices.Contains(request.Header.Values("Content-Type"), "application/json") {
			return invalidContentType
		}

		if next != nil {
			return next(writer, request)
		}
		return nil
	}
}

// balanceCalculate entry point for calculate balance
// accepts a JSON representation of a transactions array
// and returns an array of balances
func balanceCalculate(service BalanceService) customHandler {
	return func(writer http.ResponseWriter, request *http.Request) error {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println("failed to close request body: ", err)
			}
		}(request.Body)

		var t accounting.Transactions
		if err := json.NewDecoder(request.Body).Decode(&t); err != nil {
			return fmt.Errorf("failed to read request body: %+v", err)
		}

		balances := service.Calculate(t)

		if err := json.NewEncoder(writer).Encode(balances); err != nil {
			return fmt.Errorf("failed to write response body: %+v", err)
		}

		writer.Header().Set("Content-Type", "application/json")
		return nil
	}
}

func minimizeTransaction(service TransactionService) customHandler {
	return func(writer http.ResponseWriter, request *http.Request) error {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println("failed to close request body: ", err)
			}
		}(request.Body)

		var b accounting.Balances
		if err := json.NewDecoder(request.Body).Decode(&b); err != nil {
			return fmt.Errorf("failed to read request body: %+v", err)
		}

		statement := service.Minimize(b)

		if err := json.NewEncoder(writer).Encode(statement); err != nil {
			return fmt.Errorf("failed to write response body: %+v", err)
		}

		writer.Header().Set("Content-Type", "application/json")
		return nil
	}
}
