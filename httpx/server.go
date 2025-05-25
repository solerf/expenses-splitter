package httpx

import (
	"bill-splitter/accounting"
	"errors"
	"log"
	"net/http"
)

type BalanceService interface {
	Calculate(accounting.Transactions) accounting.Balances
}

type TransactionService interface {
	Minimize(accounting.Balances) accounting.Statement
}

type HttpServer struct {
	serverMux *http.ServeMux
}

func NewServer(balanceService BalanceService, transactionService TransactionService) *HttpServer {
	serverMux := &http.ServeMux{}

	e := newEndpointHandler(balanceService, transactionService)
	e.register(serverMux)

	hs := HttpServer{
		serverMux: serverMux,
	}
	return &hs
}

func (hs *HttpServer) Run() {
	log.Println("server started at :8000")
	if err := http.ListenAndServe(":8000", hs.serverMux); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %+v", err)
		return
	}
}
