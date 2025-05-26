package httpx

import (
	"bill-splitter/accounting"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type BalanceService interface {
	Calculate(accounting.Transactions) accounting.Balances
}

type TransactionService interface {
	Minimize(accounting.Balances) accounting.Statement
}

type HttpServer struct {
	server    *http.Server
	osSigChan chan os.Signal
}

func NewServer(balanceService BalanceService, transactionService TransactionService) *HttpServer {
	serverMux := &http.ServeMux{}
	register(serverMux, balanceService, transactionService)

	server := &http.Server{
		Addr:    ":8000",
		Handler: serverMux,
	}

	hs := HttpServer{
		server:    server,
		osSigChan: make(chan os.Signal, 1),
	}
	return &hs
}

func (hs *HttpServer) Run() {
	go func() {
		log.Println("server started at :8000")
		if err := hs.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %+v", err)
			return
		}
	}()

	signal.Notify(hs.osSigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-hs.osSigChan

	// should use a timed context
	if err := hs.server.Shutdown(context.Background()); err != nil {
		log.Fatalf("server failed to shutdown: %+v", err)
	}
	log.Println("server stopped")
}

func (hs *HttpServer) Close() {
	log.Println("server shutting down...")
	hs.osSigChan <- syscall.SIGTERM
}
