package accounting

// Service just represents a way to access calculate and minimize operations
type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Calculate(transactions Transactions) Balances {
	return calculateBalance(transactions)
}

func (s *Service) Minimize(balances Balances) Statement {
	return minimizeTransactions(balances)
}
