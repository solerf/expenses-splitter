package accounting

type Balance struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

type Balances = []Balance
