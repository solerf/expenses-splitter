package accounting

// Transaction holds the current amount for a person
type Transaction struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

// Transactions type alias for Transaction slice
type Transactions = []Transaction
