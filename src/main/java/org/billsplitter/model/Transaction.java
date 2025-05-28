package org.billsplitter.model;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.math.BigDecimal;
import java.util.Objects;

/**
 * POJO representation for a <code>Transaction</code> between two people
 */
public class Transaction {

    private String from;
    private String to;
    private BigDecimal amount;

    public Transaction() {
    }

    public Transaction(String from, String to, BigDecimal amount) {
        this.from = from;
        this.to = to;
        this.amount = amount;
    }

    @JsonProperty
    public String getFrom() {
        return from;
    }

    @JsonProperty
    public String getTo() {
        return to;
    }

    @JsonProperty
    public BigDecimal getAmount() {
        return amount;
    }

    public boolean notSelfTransaction() {
        return !this.from.equals(this.to);
    }

    @Override
    public String toString() {
        return "Transaction{" +
                "from='" + from + '\'' +
                ", to='" + to + '\'' +
                ", amount=" + amount +
                '}';
    }

    @Override
    public boolean equals(Object o) {
        if (o == null || getClass() != o.getClass()) return false;
        Transaction that = (Transaction) o;
        return Objects.equals(from, that.from) && Objects.equals(to, that.to) && Objects.equals(amount, that.amount);
    }

    @Override
    public int hashCode() {
        return Objects.hash(from, to, amount);
    }
}
