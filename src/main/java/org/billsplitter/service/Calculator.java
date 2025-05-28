package org.billsplitter.service;

import org.billsplitter.model.Balance;
import org.billsplitter.model.Statement;
import org.billsplitter.model.Transaction;

import java.math.BigDecimal;
import java.util.*;
import java.util.function.BiFunction;
import java.util.function.BinaryOperator;
import java.util.stream.Collectors;
import java.util.stream.Stream;

/**
 * A Calculator implementation which extends from specifics {@code TransactionCalculator} and {@code BalanceCalculator}
 */
public class Calculator implements TransactionCalculator, BalanceCalculator {

    @Override
    public List<Balance> calculate(List<Transaction> transactions) {
        if ((transactions == null) || transactions.isEmpty()) {
            return Collections.emptyList();
        }

        Map<String, BigDecimal> balances = new HashMap<>();
        return transactions
                .stream()
                .reduce(balances, TransactionToBalance.accumulator, TransactionToBalance.combiner)
                .entrySet()
                .stream()
                .map(e -> new Balance(e.getKey(), e.getValue()))
                .collect(Collectors.toList());
    }

    @Override
    public Statement minimize(List<Balance> balances) {
        if ((balances == null) || balances.isEmpty()) {
            return Statement.EMPTY;
        }

        var balancesSorted = new ArrayList<>(balances);
        balancesSorted.sort(Comparator.comparing(Balance::getAmount));
        var negIdx = 0;
        var posIdx = balancesSorted.size() - 1;

        var transactions = new ArrayList<Transaction>();
        while (negIdx < posIdx) {
            var neg = balancesSorted.get(negIdx);
            var pos = balancesSorted.get(posIdx);

            // transactionAmount is the lowest between two parties from transaction
            var transactionAmount = neg.getAmount().abs().min(pos.getAmount());
            transactions.add(new Transaction(neg.getName(), pos.getName(), transactionAmount));

            var diff = pos.getAmount().subtract(neg.getAmount().abs());
            neg.setAmount(BigDecimal.ZERO.min(diff)); // if < 0 still has debt
            pos.setAmount(BigDecimal.ZERO.max(diff)); // if > 0 still has to receive

            if (neg.getAmount().equals(BigDecimal.ZERO)) {
                negIdx++;
            }

            if (pos.getAmount().equals(BigDecimal.ZERO)) {
                posIdx--;
            }
        }
        return new Statement(balancesSorted, transactions);
    }

    private static final class TransactionToBalance {
        static final BiFunction<Map<String, BigDecimal>, Transaction, Map<String, BigDecimal>> accumulator =
                (acc, t) -> {
                    BigDecimal fromAmount = acc.getOrDefault(t.getFrom(), BigDecimal.ZERO);
                    BigDecimal toAmount = acc.getOrDefault(t.getTo(), BigDecimal.ZERO);
                    if (t.notSelfTransaction()) {
                        acc.put(t.getFrom(), fromAmount.add(t.getAmount()));
                        acc.put(t.getTo(), toAmount.subtract(t.getAmount()));
                    } else {
                        acc.put(t.getTo(), toAmount.add(BigDecimal.ZERO));
                    }
                    return acc;
                };

        static final BinaryOperator<Map<String, BigDecimal>> combiner =
                (m1, m2) ->
                        Stream.concat(m1.entrySet().stream(), m2.entrySet().stream())
                                .collect(
                                        Collectors.toMap(
                                                Map.Entry::getKey,
                                                Map.Entry::getValue,
                                                BigDecimal::add
                                        )
                                );

    }
}
