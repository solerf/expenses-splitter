package org.billsplitter.service;

import org.billsplitter.model.Balance;
import org.billsplitter.model.Statement;

import java.util.List;

/**
 * <code>TransactionCalculator</code> is the basic representation of a calculator to minimize transactions
 */
@FunctionalInterface
public interface TransactionCalculator {
    /**
     * Based on balances parameter equalize all amounts to 0, generating minimum transactions needed
     *
     * @param balances {@code List<Balance>}
     * @return <code>Statement</code>
     */
    Statement minimize(List<Balance> balances);
}
