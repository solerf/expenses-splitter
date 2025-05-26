package accounting

// Balance holds the current amount for a person
type Balance struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

// Balances type alias for Balance slice
type Balances = []Balance
