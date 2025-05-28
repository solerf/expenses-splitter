package org.billsplitter;

import io.dropwizard.core.setup.Bootstrap;
import io.dropwizard.core.setup.Environment;
import io.dropwizard.jersey.setup.JerseyEnvironment;
import org.billsplitter.controller.BalanceController;
import org.billsplitter.controller.TransactionController;
import org.billsplitter.service.Calculator;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * <code>Application</code> to bootstrap and run the application
 */
public class Application extends io.dropwizard.core.Application<AppConfiguration> {

    public static void main(final String[] args) throws Exception {
        new Application().run(args);
    }

    @Override
    public String getName() {
        return "Bill Splitter";
    }

    @Override
    public void initialize(final Bootstrap<AppConfiguration> bootstrap) {
    }

    /**
     * Basic bootstrap for application
     *
     * @param configuration {@code AppConfiguration}
     * @param environment {@code Environment}
     */
    @Override
    public void run(final AppConfiguration configuration,
                    final Environment environment) {
        JerseyEnvironment jersey = environment.jersey();

        var calculator = new Calculator();
        jersey.register(new BalanceController(calculator));
        jersey.register(new TransactionController(calculator));
    }

}
