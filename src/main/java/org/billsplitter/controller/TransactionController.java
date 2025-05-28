package org.billsplitter.controller;

import jakarta.validation.constraints.NotNull;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import org.billsplitter.model.Balance;
import org.billsplitter.model.Statement;
import org.billsplitter.service.TransactionCalculator;

import java.util.List;

@Path("/transaction")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
public class TransactionController {

    private final TransactionCalculator calculator;

    public TransactionController(TransactionCalculator calculator) {
        this.calculator = calculator;
    }

    /**
     * Endpoint for minimizing the transactions
     *
     * @param balances {@code List<Balance>} in JSON format
     * @return {@code Statement} in JSON format
     */
    @POST
    @Path("minimize")
    public Statement calculate(@NotNull(message = "transactions can't not be null") List<Balance> balances) {
        return calculator.minimize(balances);
    }

}
