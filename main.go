package main

import (
	"bill-splitter/accounting"
	"bill-splitter/httpx"
)

func main() {
	accService := accounting.NewService()
	s := httpx.NewServer(accService, accService)
	s.Run()
}
