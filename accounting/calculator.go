package accounting

import (
	"cmp"
	"math"
	"slices"
)

// CalculateBalance calculates the final balance for each person involved in the transactions
func CalculateBalance(transactions Transactions) Balances {
	b := make(map[string]float64, len(transactions))
	for _, t := range transactions {
		if t.From != t.To {
			b[t.From] = b[t.From] + t.Amount
			b[t.To] = b[t.To] - t.Amount
		}
	}

	finalB := make(Balances, 0, len(b))
	for n, a := range b {
		finalB = append(finalB, Balance{Name: n, Amount: a})
	}
	return finalB
}

// ReduceTransactions finds the minimum number of transactions to balance to 0 the amount of each person
func ReduceTransactions(balances Balances) ReductionResult {
	var traverseBalances func(b Balances, t Transactions) (Balances, Transactions)
	traverseBalances = func(b Balances, t Transactions) (Balances, Transactions) {
		slices.SortFunc(b, func(b1 Balance, b2 Balance) int {
			return cmp.Compare(b1.Amount, b2.Amount)
		})

		negIdx, posIdx := 0, len(b)-1
		neg := b[negIdx]
		pos := b[posIdx]

		if neg.Amount == 0.0 || pos.Amount == 0.0 {
			return b, t
		}

		diff := pos.Amount - math.Abs(neg.Amount)

		// tAmount is the lowest between two parties from transaction
		tAmount := math.Min(math.Abs(neg.Amount), pos.Amount)
		t = append(t, Transaction{From: neg.Name, To: pos.Name, Amount: tAmount})

		neg.Amount = math.Min(0.0, diff) // if < 0 still has debt
		pos.Amount = math.Max(0.0, diff) // if > 0 still has to receive

		b[negIdx] = neg
		b[posIdx] = pos
		return traverseBalances(b, t)
	}

	finalBalances, finalTransactions := traverseBalances(
		append(Balances{}, balances...),
		make(Transactions, 0, len(balances)),
	)
	return ReductionResult{finalBalances, finalTransactions}
}
