# Documentation (1st part)

## How to run

Ensure to have `java` at least version `21`, build or test the application. `docker` is not mandatory to have installed.

All operations are listed in the [Makefile](./Makefile).

### Testing

To execute the tests, run:

```bash
 make tests
```

### Building

There are two options to build the application:

- With `maven`
  Run:

```bash
 make build
```

This will generate a `jar` file under `target` folder.

- With `docker`
  Run:

```bash
 # DOCKER_TAG is not mandatory if not informed `bill-splitter` will be assumed 
 make build-docker DOCKER_TAG=your-tag
```

This will generate a docker image of the application in your local repository identified by the `DOCKER_TAG`.

### Running

There are three options to run the application:

- With `maven`
  Run:

```bash
 make run
```

Or, if the application was built with the `make build` command from the previous section. Must have a `java`
installation.

```bash
 java -jar target/billsplitter-1.0-SNAPSHOT.jar server config.yml
```

- With `docker`
  Run:

```bash
 # DOCKER_TAG is not mandatory if not informed `bill-splitter` will be assumed 
 # HOST_PORT is not mandatory if not informed `8000` will be assumed 
 make run-docker DOCKER_TAG=your-tag HOST_PORT=your-port
```

Once application starts it will be accessible at port `:8000`, or the one defined with `HOST_PORT` if running from
`docker` image.

## Request Examples

### Calculating the balance

To calculate the balances we need to do a `POST` at `/balance/calculate` with a `JSON` of transactions.
Sample request:

```bash
 curl --header "Content-Type: application/json" \
      --request POST \
      --data '[{ "from": "A", "to": "B", "amount": 40 },{ "from": "B", "to": "C", "amount": 40 },{ "from": "C", "to": "A", "amount": 10 }]' \
      http://localhost:8000/balance/calculate
```

Sample response:

```json
[
  {
    "name": "A",
    "amount": 30
  },
  {
    "name": "B",
    "amount": 0
  },
  {
    "name": "C",
    "amount": -30
  }
]
```

### Minimizing the transactions

To minimize the transactions we need to do a `POST` at `/transaction/minimize` with a `JSON` of balances.
Sample request:

```bash
 curl --header "Content-Type: application/json" \
      --request POST \
      --data '[{ "name": "A", "amount": 30 }, { "name": "B", "amount": 0 }, { "name": "C", "amount": -30 }]' \
      http://localhost:8000/transaction/minimize
```

Sample response:

```json
{
  "updated_balances": [
    {
      "name": "C",
      "amount": 0
    },
    {
      "name": "B",
      "amount": 0
    },
    {
      "name": "A",
      "amount": 0
    }
  ],
  "transactions": [
    {
      "from": "C",
      "to": "A",
      "amount": 30
    }
  ]
}
```

## Assumptions

- Focusing on simplicity and ease of development, bootstrapped application using `dropwizard` archetype.
- It is a straight forward stateless API, no database involved.
- There are a lot of points for improvement, like:
    - Observability (metrics, logging, health).
    - Improve configuration support from outside via args and/or files.
    - If desired to effectively handle groups, Authentication/Authorization could be implemented.
    - Add database support to better track groups/persons balances and transactions.
    - Improve endpoints documentation adding `swagger` generating a static web for it.
- No `IntegrationTests` per se, once I assumed, due to application simplicity simple testing controllers with
  `dropwizard` support would be enough, even though `TestContainers` could be used for it.
