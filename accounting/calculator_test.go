package accounting

import (
	"cmp"
	"reflect"
	"slices"
	"testing"
)

func Test_Calculate_Balance(t *testing.T) {
	scenarios := []struct {
		name     string
		input    Transactions
		expected Balances
	}{
		{
			name: "given example",
			input: Transactions{
				{"A", "B", 40.0},
				{"B", "C", 40.0},
				{"C", "A", 10.0},
			},
			expected: Balances{
				{"A", 30.0},
				{"B", 0.0},
				{"C", -30.0},
			},
		},
		{
			name: "when more transactions",
			input: Transactions{
				{"A", "B", 50.0},
				{"B", "C", 40.0},
				{"C", "B", 10.0},
				{"D", "A", 20.0},
				{"A", "C", 30.0},
				{"E", "A", 5.0},
			},
			expected: Balances{
				{"A", 55.0},
				{"B", -20.0},
				{"C", -60.0},
				{"D", 20.0},
				{"E", 5.0},
			},
		},
		{
			name: "when self transaction",
			input: Transactions{
				{"A", "B", 50.0},
				{"B", "B", 50.0},
			},
			expected: Balances{
				{"A", 50.0},
				{"B", -50.0},
			},
		},
		{
			name: "when single self transaction",
			input: Transactions{
				{"B", "B", 50.0},
			},
			expected: Balances{
				{"B", 0.0},
			},
		},
		{
			name:     "when empty transactions",
			input:    Transactions{},
			expected: Balances{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			actual := calculateBalance(s.input)

			slices.SortFunc(s.expected, func(a, b Balance) int {
				return cmp.Compare(a.Name, b.Name)
			})
			slices.SortFunc(actual, func(a, b Balance) int {
				return cmp.Compare(a.Name, b.Name)
			})

			if !reflect.DeepEqual(s.expected, actual) {
				t.Errorf("\nExpected:	%+v\nGot:		%+v", s.expected, actual)
			}
		})
	}
}

func Test_Minimize_Transactions(t *testing.T) {
	scenarios := []struct {
		name     string
		input    Balances
		expected Statement
	}{
		{
			name: "given example",
			input: Balances{
				{"A", 30.0},
				{"B", 0.0},
				{"C", -30.0},
			},
			expected: Statement{
				UpdatedBalances: Balances{
					{"A", 0.0},
					{"B", 0.0},
					{"C", 0.0},
				},
				Transactions: Transactions{
					{"C", "A", 30.0},
				},
			},
		},
		{
			name: "when more debts than credits",
			input: Balances{
				{"A", 30.0},
				{"B", -5.0},
				{"C", -5.0},
				{"D", -10.0},
				{"E", -5.0},
				{"F", -5.0},
			},
			expected: Statement{
				UpdatedBalances: Balances{
					{"A", 0.0},
					{"B", 0.0},
					{"C", 0.0},
					{"D", 0.0},
					{"E", 0.0},
					{"F", 0.0},
				},
				Transactions: Transactions{
					{"B", "A", 5.0},
					{"C", "A", 5.0},
					{"D", "A", 10.0},
					{"E", "A", 5.0},
					{"F", "A", 5.0},
				},
			},
		},
		{
			name: "when more credits than debts",
			input: Balances{
				{"A", -30.0},
				{"B", 5.0},
				{"C", 5.0},
				{"D", 10.0},
				{"E", 5.0},
				{"F", 5.0},
			},
			expected: Statement{
				UpdatedBalances: Balances{
					{"A", 0.0},
					{"B", 0.0},
					{"C", 0.0},
					{"D", 0.0},
					{"E", 0.0},
					{"F", 0.0},
				},
				Transactions: Transactions{
					{"A", "B", 5.0},
					{"A", "C", 5.0},
					{"A", "D", 10.0},
					{"A", "E", 5.0},
					{"A", "F", 5.0},
				},
			},
		},
		{
			name: "when more data in balances",
			input: Balances{
				{"A", 10.0},
				{"B", 5.0},
				{"C", -30.0},
				{"D", -50.0},
				{"E", 100.0},
				{"F", -35.0},
			},
			expected: Statement{
				UpdatedBalances: Balances{
					{"A", 0.0},
					{"B", 0.0},
					{"C", 0.0},
					{"D", 0.0},
					{"E", 0.0},
					{"F", 0.0},
				},
				Transactions: Transactions{
					{"C", "B", 5.0},
					{"C", "A", 10.0},
					{"C", "E", 15.0},
					{"F", "E", 35.0},
					{"D", "E", 50.0},
				},
			},
		},
		{
			name: "when incomplete balances, incomplete pay to owned",
			input: Balances{
				{"A", 30.0},
				{"B", -5.0},
				{"C", -5.0},
			},
			expected: Statement{
				UpdatedBalances: Balances{
					{"A", 20.0},
					{"B", 0.0},
					{"C", 0.0},
				},
				Transactions: Transactions{
					{"B", "A", 5.0},
					{"C", "A", 5.0},
				},
			},
		},
		{
			name: "when incomplete balances, complete pay to owned",
			input: Balances{
				{"A", 5.0},
				{"B", -5.0},
				{"C", -5.0},
			},
			expected: Statement{
				UpdatedBalances: Balances{
					{"A", 0.0},
					{"B", 0.0},
					{"C", -5.0},
				},
				Transactions: Transactions{
					{"B", "A", 5.0},
				},
			},
		},
		{
			name: "when single balance",
			input: Balances{
				{"A", 5.0},
			},
			expected: Statement{
				UpdatedBalances: Balances{
					{"A", 5.0},
				},
				Transactions: Transactions{},
			},
		},
		{
			name:  "when empty balances",
			input: Balances{},
			expected: Statement{
				UpdatedBalances: Balances{},
				Transactions:    Transactions{},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {

			actual := minimizeTransactions(s.input)

			slices.SortFunc(actual.UpdatedBalances, func(a, b Balance) int {
				return cmp.Compare(a.Name, b.Name)
			})
			slices.SortFunc(s.expected.UpdatedBalances, func(a, b Balance) int {
				return cmp.Compare(a.Name, b.Name)
			})

			slices.SortFunc(actual.Transactions, func(a, b Transaction) int {
				return cmp.Compare(a.To, b.To) + cmp.Compare(a.From, b.From)
			})
			slices.SortFunc(s.expected.Transactions, func(a, b Transaction) int {
				return cmp.Compare(a.To, b.To) + cmp.Compare(a.From, b.From)
			})

			if !reflect.DeepEqual(s.expected.UpdatedBalances, actual.UpdatedBalances) {
				t.Errorf("\nBalances:\nExpected:	%+v\nGot:		%+v", s.expected.UpdatedBalances, actual.UpdatedBalances)
			}
			if !reflect.DeepEqual(s.expected.Transactions, actual.Transactions) {
				t.Errorf("\nTransactions:\nExpected:	%+v\nGot:		%+v", s.expected.Transactions, actual.Transactions)
			}
		})
	}
}
