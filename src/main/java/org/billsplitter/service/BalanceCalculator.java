package org.billsplitter.service;

import org.billsplitter.model.Balance;
import org.billsplitter.model.Transaction;

import java.util.List;

/**
 * <code>BalanceCalculator</code> is the basic representation of a calculator to calculate balances from transactions
 */
@FunctionalInterface
public interface BalanceCalculator {
    /**
     * From transactions parameter calculate the final balance for all
     * people involved
     *
     * @param transactions {@code List<Transaction>}
     * @return
     */
    List<Balance> calculate(List<Transaction> transactions);
}

