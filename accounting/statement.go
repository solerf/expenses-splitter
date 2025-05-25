package accounting

type Statement struct {
	UpdatedBalances Balances     `json:"updated_balances"`
	Transactions    Transactions `json:"transactions"`
}
