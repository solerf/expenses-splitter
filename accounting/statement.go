package accounting

// Statement contains information of updated balances and transactions after minimize operation
type Statement struct {
	UpdatedBalances Balances     `json:"updated_balances"`
	Transactions    Transactions `json:"transactions"`
}
