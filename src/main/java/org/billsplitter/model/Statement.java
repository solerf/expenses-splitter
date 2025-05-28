package org.billsplitter.model;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Collections;
import java.util.List;

/**
 * POJO representation for a <code>Statement</code>, which is the result of a <code>TransactionCalculator.minimize</code> operation
 */
public class Statement {

    public static final Statement EMPTY = new Statement(Collections.emptyList(), Collections.emptyList());

    private List<Balance> updatedBalances;
    private List<Transaction> transactions;

    public Statement() {
    }

    public Statement(List<Balance> updatedBalances, List<Transaction> transactions) {
        this.updatedBalances = updatedBalances;
        this.transactions = transactions;
    }

    @JsonProperty
    public List<Balance> getUpdatedBalances() {
        return updatedBalances;
    }

    @JsonProperty
    public List<Transaction> getTransactions() {
        return transactions;
    }
}
