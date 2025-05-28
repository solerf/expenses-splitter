package org.billsplitter.service;

import org.billsplitter.model.Balance;
import org.billsplitter.model.Transaction;
import org.junit.jupiter.api.Test;

import java.math.BigDecimal;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Comparator;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertIterableEquals;

class BalanceCalculatorTest {

    Calculator calculatorInstance = new Calculator();

    void runAssertion(List<Transaction> input, List<Balance> expected) {
        var actual = calculatorInstance.calculate(input);

        expected.sort(Comparator.comparing(Balance::getName));
        actual.sort(Comparator.comparing(Balance::getName));

        assertEquals(expected.size(), actual.size());
        assertIterableEquals(expected, actual);
    }

    @Test
    void test_given_example() {
        var input = Arrays.asList(
                new Transaction("A", "B", BigDecimal.valueOf(40)),
                new Transaction("B", "C", BigDecimal.valueOf(40)),
                new Transaction("C", "A", BigDecimal.valueOf(10))
        );
        var expected = Arrays.asList(
                new Balance("A", BigDecimal.valueOf(30)),
                new Balance("B", BigDecimal.valueOf(0)),
                new Balance("C", BigDecimal.valueOf(-30))
        );
        runAssertion(input, expected);
    }

    @Test
    void test_when_more_transactions() {
        var input = Arrays.asList(
                new Transaction("A", "B", BigDecimal.valueOf(50)),
                new Transaction("B", "C", BigDecimal.valueOf(40)),
                new Transaction("C", "B", BigDecimal.valueOf(10)),
                new Transaction("D", "A", BigDecimal.valueOf(20)),
                new Transaction("A", "C", BigDecimal.valueOf(30)),
                new Transaction("E", "A", BigDecimal.valueOf(5))
        );
        var expected = Arrays.asList(
                new Balance("A", BigDecimal.valueOf(55)),
                new Balance("B", BigDecimal.valueOf(-20)),
                new Balance("C", BigDecimal.valueOf(-60)),
                new Balance("D", BigDecimal.valueOf(20)),
                new Balance("E", BigDecimal.valueOf(5))
        );
        runAssertion(input, expected);
    }

    @Test
    void test_when_self_transaction() {
        var input = Arrays.asList(
                new Transaction("A", "B", BigDecimal.valueOf(50)),
                new Transaction("B", "B", BigDecimal.valueOf(50))
        );
        var expected = Arrays.asList(
                new Balance("A", BigDecimal.valueOf(50)),
                new Balance("B", BigDecimal.valueOf(-50))
        );
        runAssertion(input, expected);
    }

    @Test
    void test_when_single_self_transaction() {
        var input = Arrays.asList(
                new Transaction("B", "B", BigDecimal.valueOf(50))
        );
        var expected = Arrays.asList(
                new Balance("B", BigDecimal.valueOf(0))
        );
        runAssertion(input, expected);
    }

    @Test
    void test_when_empty_transactions() {
        var input = new ArrayList<Transaction>();
        var expected = new ArrayList<Balance>();

        runAssertion(input, expected);
    }
}
