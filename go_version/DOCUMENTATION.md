# Documentation

## How to run

Ensure to have `go` at least version `1.24`, build or test the application. `docker` is not mandatory to have installed.

All operations are listed in the [Makefile](./Makefile).

### Testing

To execute the tests, run:

```bash
 make tests
```

### Building

There are two options to build the application:

- With `go`
  Run:

```bash
 make build
```

This will generate a binary under `bin` in the root folder.

- With `docker`
  Run:

```bash
 # DOCKER_TAG is not mandatory if not informed `bill-splitter` will be assumed 
 make build-docker DOCKER_TAG=your-tag
```

This will generate a docker image of the application in your local repository identified by the `DOCKER_TAG`.

### Running

There are three options to run the application:

- With `go`
  Run:

```bash
 make run
```

Or, if the application was built with the `make build` command from the previous section.

```bash
 ./bin/bill-splitter
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

- Targeting simplicity and ease of development, `go` was used with no third party dependencies involved.
- It is a straight forward stateless API, no database involved.
- There are a lot of points for improvement, like:
    - Observability (metrics, logging, health).
    - Support accepting configuration from outside via args and/or files. 
    - If desired to effectively handle groups, implement Authentication/Authorization.
    - Add database support to better track groups/persons balances and transactions.
    - Improve endpoints documentation adding `swagger`.
- Using float for simplicity, even though it's known it's issues with precision
