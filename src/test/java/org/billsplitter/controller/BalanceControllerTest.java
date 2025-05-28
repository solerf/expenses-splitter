package org.billsplitter.controller;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.dropwizard.testing.junit5.DropwizardExtensionsSupport;
import io.dropwizard.testing.junit5.ResourceExtension;
import jakarta.ws.rs.client.Entity;
import jakarta.ws.rs.core.Response;
import jakarta.ws.rs.core.Response.Status.Family;
import org.billsplitter.model.Balance;
import org.billsplitter.model.Transaction;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;

import java.math.BigDecimal;
import java.util.Arrays;
import java.util.Comparator;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertIterableEquals;

@ExtendWith(DropwizardExtensionsSupport.class)
class BalanceControllerTest {

    private static final ResourceExtension CONTROLLER = ResourceExtension.builder()
            .addResource(new BalanceController(transactions -> Arrays.asList(
                    new Balance("A", BigDecimal.valueOf(30)),
                    new Balance("B", BigDecimal.valueOf(0)),
                    new Balance("C", BigDecimal.valueOf(-30))
            )))
            .build();

    @Test
    void test_null_payload() {
        try (Response response = CONTROLLER.target("/balance/calculate").request().post(Entity.json(null))) {
            assertEquals(Family.CLIENT_ERROR, Family.familyOf(response.getStatus()));
        }
    }

    @Test
    void test_invalid_payload() {
        try (Response response = CONTROLLER.target("/balance/calculate").request().post(Entity.json("{\"a\": 2}"))) {
            assertEquals(Family.CLIENT_ERROR, Family.familyOf(response.getStatus()));
        }
    }

    @Test
    void test_valid_payload() throws JsonProcessingException {
        var objectMapper = new ObjectMapper();
        var writer = objectMapper.writer().withDefaultPrettyPrinter();
        var json = writer.writeValueAsString(Arrays.asList(
                new Transaction("A", "B", BigDecimal.valueOf(40)),
                new Transaction("B", "C", BigDecimal.valueOf(40)),
                new Transaction("C", "A", BigDecimal.valueOf(10))
        ));
        try (Response response = CONTROLLER.target("/balance/calculate").request().post(Entity.json(json))) {
            var rawJson = response.readEntity(String.class);
            var actual = objectMapper.readValue(rawJson, new TypeReference<List<Balance>>() {
            });

            var expected = Arrays.asList(
                    new Balance("A", BigDecimal.valueOf(30)),
                    new Balance("B", BigDecimal.valueOf(0)),
                    new Balance("C", BigDecimal.valueOf(-30))
            );
            expected.sort(Comparator.comparing(Balance::getName));
            actual.sort(Comparator.comparing(Balance::getName));

            assertEquals(Family.SUCCESSFUL, Family.familyOf(response.getStatus()));
            assertEquals("[{\"name\":\"A\",\"amount\":30},{\"name\":\"B\",\"amount\":0},{\"name\":\"C\",\"amount\":-30}]", rawJson);

            assertIterableEquals(expected, actual);
        }
    }

}
