package org.billsplitter.controller;


import jakarta.validation.constraints.NotNull;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import org.billsplitter.model.Balance;
import org.billsplitter.model.Transaction;
import org.billsplitter.service.BalanceCalculator;

import java.util.List;

@Path("/balance")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
public class BalanceController {

    private final BalanceCalculator calculator;

    public BalanceController(BalanceCalculator calculator) {
        this.calculator = calculator;
    }

    /**
     * Endpoint for calculating the balances
     *
     * @param transactions {@code List<Transaction>} in JSON format
     * @return {@code List<Balance>} in JSON format
     */
    @POST
    @Path("/calculate")
    public List<Balance> calculate(@NotNull(message = "transactions can't not be null") List<Transaction> transactions) {
        return calculator.calculate(transactions);
    }

}
