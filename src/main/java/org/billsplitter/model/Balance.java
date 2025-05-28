package org.billsplitter.model;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.math.BigDecimal;
import java.util.Objects;

/**
 * POJO representation for a person <code>Balance</code>
 */
public class Balance {

    private String name;
    private BigDecimal amount;

    public Balance() {
    }

    public Balance(String name, BigDecimal amount) {
        this.name = name;
        this.amount = amount;
    }

    @JsonProperty
    public String getName() {
        return name;
    }

    public void setAmount(BigDecimal amount) {
        this.amount = amount;
    }

    @JsonProperty
    public BigDecimal getAmount() {
        return amount;
    }

    @Override
    public String toString() {
        return "Balance{" +
                "name='" + name + '\'' +
                ", amount=" + amount +
                '}';
    }

    @Override
    public boolean equals(Object o) {
        if (o == null || getClass() != o.getClass()) return false;
        Balance balance = (Balance) o;
        return Objects.equals(name, balance.name) && Objects.equals(amount, balance.amount);
    }

    @Override
    public int hashCode() {
        return Objects.hash(name, amount);
    }
}
