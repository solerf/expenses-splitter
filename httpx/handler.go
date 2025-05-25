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

type httpHandler func(http.ResponseWriter, *http.Request) error

type endpoints struct {
	balanceService     BalanceService
	transactionService TransactionService
}

func newEndpointHandler(balanceService BalanceService, transactionService TransactionService) *endpoints {
	return &endpoints{
		balanceService:     balanceService,
		transactionService: transactionService,
	}
}

func (e *endpoints) register(mux *http.ServeMux) {
	mux.HandleFunc(
		"POST /balance/calculate",
		mainHandlerFunc(validateContentType(balanceCalculate(e.balanceService))),
	)

	mux.HandleFunc(
		"POST /transaction/minimize",
		mainHandlerFunc(validateContentType(minimizeTransaction(e.transactionService))),
	)

	mux.HandleFunc("/", http.NotFound)
}

func mainHandlerFunc(innerHandler httpHandler) http.HandlerFunc {
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

func validateContentType(next httpHandler) httpHandler {
	return func(writer http.ResponseWriter, request *http.Request) error {
		if !slices.Contains(request.Header.Values("Content-Type"), "application/json") {
			return invalidContentType
		}
		return next(writer, request)
	}
}

func balanceCalculate(service BalanceService) httpHandler {
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

func minimizeTransaction(service TransactionService) httpHandler {
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
