package org.billsplitter.controller;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.dropwizard.testing.junit5.DropwizardExtensionsSupport;
import io.dropwizard.testing.junit5.ResourceExtension;
import jakarta.ws.rs.client.Entity;
import jakarta.ws.rs.core.Response;
import jakarta.ws.rs.core.Response.Status.Family;
import org.billsplitter.model.Balance;
import org.billsplitter.model.Statement;
import org.billsplitter.model.Transaction;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;

import java.math.BigDecimal;
import java.util.Arrays;
import java.util.Comparator;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertIterableEquals;

@ExtendWith(DropwizardExtensionsSupport.class)
class TransactionControllerTest {

    private static final ResourceExtension CONTROLLER = ResourceExtension.builder()
            .addResource(new TransactionController(transactions ->
                    new Statement(
                            Arrays.asList(
                                    new Balance("A", BigDecimal.valueOf(30)),
                                    new Balance("B", BigDecimal.valueOf(0)),
                                    new Balance("C", BigDecimal.valueOf(-30))
                            ),
                            Arrays.asList(
                                    new Transaction("A", "B", BigDecimal.valueOf(40)),
                                    new Transaction("B", "C", BigDecimal.valueOf(40)),
                                    new Transaction("C", "A", BigDecimal.valueOf(10))
                            )
                    )))
            .build();

    @Test
    void test_null_payload() {
        try (Response response = CONTROLLER.target("/transaction/minimize").request().post(Entity.json(null))) {
            assertEquals(Family.CLIENT_ERROR, Family.familyOf(response.getStatus()));
        }
    }

    @Test
    void test_invalid_payload() {
        try (Response response = CONTROLLER.target("/transaction/minimize").request().post(Entity.json("{\"a\": 2}"))) {
            assertEquals(Family.CLIENT_ERROR, Family.familyOf(response.getStatus()));
        }
    }

    @Test
    void test_valid_payload() throws JsonProcessingException {
        var objectMapper = new ObjectMapper();
        var writer = objectMapper.writer().withDefaultPrettyPrinter();
        var json = writer.writeValueAsString(Arrays.asList(
                new Balance("A", BigDecimal.valueOf(30)),
                new Balance("B", BigDecimal.valueOf(0)),
                new Balance("C", BigDecimal.valueOf(-30))
        ));
        try (Response response = CONTROLLER.target("/transaction/minimize").request().post(Entity.json(json))) {
            var rawJson = response.readEntity(String.class);
            var actual = objectMapper.readValue(rawJson, Statement.class);

            var expectedTransactions = Arrays.asList(
                    new Transaction("A", "B", BigDecimal.valueOf(40)),
                    new Transaction("B", "C", BigDecimal.valueOf(40)),
                    new Transaction("C", "A", BigDecimal.valueOf(10))
            );
            var expectedBalance = Arrays.asList(
                    new Balance("A", BigDecimal.valueOf(30)),
                    new Balance("B", BigDecimal.valueOf(0)),
                    new Balance("C", BigDecimal.valueOf(-30))
            );

            expectedBalance.sort(Comparator.comparing(Balance::getName));
            actual.getUpdatedBalances().sort(Comparator.comparing(Balance::getName));

            expectedTransactions.sort((a, b) -> a.getFrom().compareTo(b.getFrom()) + a.getTo().compareTo(b.getTo()));
            actual.getTransactions().sort((a, b) -> a.getFrom().compareTo(b.getFrom()) + a.getTo().compareTo(b.getTo()));

            assertEquals(Family.SUCCESSFUL, Family.familyOf(response.getStatus()));
            assertEquals(
                    "{\"updatedBalances\":[{\"name\":\"A\",\"amount\":30},{\"name\":\"B\",\"amount\":0},{\"name\":\"C\",\"amount\":-30}],\"transactions\":[{\"from\":\"A\",\"to\":\"B\",\"amount\":40},{\"from\":\"B\",\"to\":\"C\",\"amount\":40},{\"from\":\"C\",\"to\":\"A\",\"amount\":10}]}",
                    rawJson
            );

            assertIterableEquals(expectedBalance, actual.getUpdatedBalances());
            assertIterableEquals(expectedTransactions, actual.getTransactions());
        }
    }

}
