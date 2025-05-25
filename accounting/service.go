package accounting

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Calculate(transactions Transactions) Balances {
	return calculateBalance(transactions)
}

func (s *Service) Minimize(balances Balances) Statement {
	return minimizeTransactions(balances)
}
