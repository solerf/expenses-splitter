package accounting

import (
	"cmp"
	"math"
	"slices"
)

// calculateBalance calculates the final balance for each person involved in the transactions
func calculateBalance(transactions Transactions) Balances {
	b := make(map[string]float64, len(transactions))
	for _, t := range transactions {
		if t.From != t.To {
			b[t.From] += t.Amount
			b[t.To] -= t.Amount
		} else {
			// if self just record it as 0, for person to be accounted in balance
			b[t.To] += 0.0
		}
	}

	finalB := make(Balances, 0, len(b))
	for n, a := range b {
		finalB = append(finalB, Balance{Name: n, Amount: a})
	}
	return finalB
}

// minimizeTransactions finds the minimum number of transactions to balance to 0 the amount of each person
func minimizeTransactions(balances Balances) Statement {
	finalBalances := append(Balances{}, balances...)
	slices.SortFunc(finalBalances, func(b1 Balance, b2 Balance) int {
		return cmp.Compare(b1.Amount, b2.Amount)
	})

	finalTransactions := make(Transactions, 0, len(balances))

	negIdx, posIdx := 0, len(finalBalances)-1
	for posIdx > negIdx {
		neg := finalBalances[negIdx]
		pos := finalBalances[posIdx]

		diff := pos.Amount - math.Abs(neg.Amount)

		// tAmount is the lowest between two parties from transaction
		tAmount := math.Min(math.Abs(neg.Amount), pos.Amount)
		finalTransactions = append(finalTransactions, Transaction{From: neg.Name, To: pos.Name, Amount: tAmount})

		neg.Amount = math.Min(0.0, diff) // if < 0 still has debt
		pos.Amount = math.Max(0.0, diff) // if > 0 still has to receive

		finalBalances[negIdx], finalBalances[posIdx] = neg, pos

		if neg.Amount == 0.0 {
			negIdx++
		}

		if pos.Amount == 0.0 {
			posIdx--
		}
	}
	return Statement{finalBalances, finalTransactions}
}
