package accounting

type Transaction struct {
	From   string
	To     string
	Amount float64
}

type Transactions []Transaction

type Balance struct {
	Name   string
	Amount float64
}

type Balances []Balance

type ReductionResult struct {
	UpdatedBalances Balances
	Transactions    Transactions
}
