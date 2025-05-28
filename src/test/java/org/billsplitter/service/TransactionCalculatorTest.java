package org.billsplitter.service;

import org.billsplitter.model.Balance;
import org.billsplitter.model.Statement;
import org.billsplitter.model.Transaction;
import org.junit.jupiter.api.Test;

import java.math.BigDecimal;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Comparator;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertIterableEquals;

class TransactionCalculatorTest {

    Calculator calculatorInstance = new Calculator();

    void runAssertion(List<Balance> input, Statement expected) {
        var statement = calculatorInstance.minimize(input);

        var expectedBalances = expected.getUpdatedBalances();
        var expectedTransactions = expected.getTransactions();

        var actualBalances = statement.getUpdatedBalances();
        var actualTransactions = statement.getTransactions();

        expectedBalances.sort(Comparator.comparing(Balance::getName));
        expectedTransactions.sort((a, b) -> a.getFrom().compareTo(b.getFrom()) + a.getTo().compareTo(b.getTo()));

        actualBalances.sort(Comparator.comparing(Balance::getName));
        actualTransactions.sort((a, b) -> a.getFrom().compareTo(b.getFrom()) + a.getTo().compareTo(b.getTo()));

        assertEquals(expectedBalances.size(), actualBalances.size());
        assertEquals(expectedTransactions.size(), actualTransactions.size());
        assertIterableEquals(expectedBalances, actualBalances);
        assertIterableEquals(expectedTransactions, actualTransactions);
    }

    @Test
    void test_given_example() {
        var input = Arrays.asList(
                new Balance("A", BigDecimal.valueOf(30)),
                new Balance("B", BigDecimal.valueOf(0)),
                new Balance("C", BigDecimal.valueOf(-30))
        );
        var expected = new Statement(
                Arrays.asList(
                        new Balance("A", BigDecimal.valueOf(0)),
                        new Balance("B", BigDecimal.valueOf(0)),
                        new Balance("C", BigDecimal.valueOf(0))
                ),
                Arrays.asList(
                        new Transaction("C", "A", new BigDecimal(30))
                )
        );
        runAssertion(input, expected);
    }

    @Test
    void test_when_more_debts_than_credits() {
        var input = Arrays.asList(
                new Balance("A", new BigDecimal(30)),
                new Balance("B", new BigDecimal(-5)),
                new Balance("C", new BigDecimal(-5)),
                new Balance("D", new BigDecimal(-10)),
                new Balance("E", new BigDecimal(-5)),
                new Balance("F", new BigDecimal(-5))
        );
        var expected = new Statement(
                Arrays.asList(
                        new Balance("A", new BigDecimal(0)),
                        new Balance("B", new BigDecimal(0)),
                        new Balance("C", new BigDecimal(0)),
                        new Balance("D", new BigDecimal(0)),
                        new Balance("E", new BigDecimal(0)),
                        new Balance("F", new BigDecimal(0))
                ),
                Arrays.asList(
                        new Transaction("B", "A", new BigDecimal(5)),
                        new Transaction("C", "A", new BigDecimal(5)),
                        new Transaction("D", "A", new BigDecimal(10)),
                        new Transaction("E", "A", new BigDecimal(5)),
                        new Transaction("F", "A", new BigDecimal(5))
                )
        );
        runAssertion(input, expected);
    }

    @Test
    void test_more_credits_than_debts() {
        var input = Arrays.asList(
                new Balance("A", new BigDecimal(-30)),
                new Balance("B", new BigDecimal(5)),
                new Balance("C", new BigDecimal(5)),
                new Balance("D", new BigDecimal(10)),
                new Balance("E", new BigDecimal(5)),
                new Balance("F", new BigDecimal(5))
        );
        var expected = new Statement(
                Arrays.asList(
                        new Balance("A", new BigDecimal(0)),
                        new Balance("B", new BigDecimal(0)),
                        new Balance("C", new BigDecimal(0)),
                        new Balance("D", new BigDecimal(0)),
                        new Balance("E", new BigDecimal(0)),
                        new Balance("F", new BigDecimal(0))
                ),
                Arrays.asList(
                        new Transaction("A", "B", new BigDecimal(5)),
                        new Transaction("A", "C", new BigDecimal(5)),
                        new Transaction("A", "D", new BigDecimal(10)),
                        new Transaction("A", "E", new BigDecimal(5)),
                        new Transaction("A", "F", new BigDecimal(5))
                )
        );
        runAssertion(input, expected);
    }

    @Test
    void test_when_more_data_in_balances() {
        var input = Arrays.asList(
                new Balance("A", new BigDecimal(10)),
                new Balance("B", new BigDecimal(5)),
                new Balance("C", new BigDecimal(-30)),
                new Balance("D", new BigDecimal(-50)),
                new Balance("E", new BigDecimal(100)),
                new Balance("F", new BigDecimal(-35))
        );
        var expected = new Statement(
                Arrays.asList(
                        new Balance("A", new BigDecimal(0)),
                        new Balance("B", new BigDecimal(0)),
                        new Balance("C", new BigDecimal(0)),
                        new Balance("D", new BigDecimal(0)),
                        new Balance("E", new BigDecimal(0)),
                        new Balance("F", new BigDecimal(0))
                ),
                Arrays.asList(
                        new Transaction("C", "B", new BigDecimal(5)),
                        new Transaction("C", "A", new BigDecimal(10)),
                        new Transaction("C", "E", new BigDecimal(15)),
                        new Transaction("F", "E", new BigDecimal(35)),
                        new Transaction("D", "E", new BigDecimal(50))
                )
        );
        runAssertion(input, expected);
    }

    @Test
    void test_when_incomplete_balances_incomplete_debtors() {
        var input = Arrays.asList(
                new Balance("A", new BigDecimal(30)),
                new Balance("B", new BigDecimal(-5)),
                new Balance("C", new BigDecimal(-5))
        );
        var expected = new Statement(
                Arrays.asList(
                        new Balance("A", new BigDecimal(20)),
                        new Balance("B", new BigDecimal(0)),
                        new Balance("C", new BigDecimal(0))
                ),
                Arrays.asList(
                        new Transaction("B", "A", new BigDecimal(5)),
                        new Transaction("C", "A", new BigDecimal(5))
                )
        );
        runAssertion(input, expected);
    }

    @Test
    void test_when_incomplete_balances_incomplete_creditors() {
        var input = Arrays.asList(
                new Balance("A", new BigDecimal(5)),
                new Balance("B", new BigDecimal(-5)),
                new Balance("C", new BigDecimal(-5))
        );
        var expected = new Statement(
                Arrays.asList(
                        new Balance("A", new BigDecimal(0)),
                        new Balance("B", new BigDecimal(0)),
                        new Balance("C", new BigDecimal(-5))
                ),
                Arrays.asList(
                        new Transaction("B", "A", new BigDecimal(5))
                )
        );
        runAssertion(input, expected);
    }

    @Test
    void test_when_single_balance() {
        var input = Arrays.asList(
                new Balance("A", new BigDecimal(5))
        );
        var expected = new Statement(
                Arrays.asList(
                        new Balance("A", new BigDecimal(5))
                ),
                Arrays.asList()
        );
        runAssertion(input, expected);
    }

    @Test
    void test_when_empty_balances() {
        var input = new ArrayList<Balance>();
        var expected = new Statement(
                Arrays.asList(),
                Arrays.asList()
        );
        runAssertion(input, expected);
    }

}
