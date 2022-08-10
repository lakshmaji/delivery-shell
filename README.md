# Delivery shell app

1. Computes discount
1. Computes discount along with est delivery time
1. Validate offers schema

[Show me how to run the app](#build)

## Structuring

### Services

Uses **go** way of designing services (DI) using Hexagonal architecture.

- Offer service : Offer service is responsible for applying offer discount amount when meets the expected criteria.
- Delivery Service: Delivery service responsible for computing **delivery cost** and applies **offer discount** if applicable by using `offer service`
- Total delivery cost can be computed using method available on **PackageDetails**, so that all computation logic will be at one place (The better way could be have its own service).

```txt
📦 services
 ┣ 📂 delivery_svc
 ┃ ┣ 📜 default_svc.go
 ┃ ┗ 📜 delivery_svc.go
 ┣ 📂 offers_svc
 ┃ ┣ 📜 default_svc.go
 ┃ ┗ 📜 offers_svc.go
 ┗ 📂 shell_io_svc
 ┃ ┣ 📜 default_svc.go
 ┃ ┗ 📜 shell_io_svc.go
```

### Models

The domain models

```txt
📦 models
 ┣ 📜 offers.go
 ┣ 📜 package_details.go
 ┣ 📜 package_stats.go
 ┗ 📜 vehicles.go
```

### Clients

Write your own client

```txt
📦 clients
 ┣ 📜 base_client.go
 ┣ 📜 shell_client.go
```

## Managing Offers

We have created a schema to validate whether give offer is valid or applicable.

```go
type Offer struct {
  Code       OfferCode
  Conditions []Condition
  Discount   float64
}
```

Maintains a list of offers in `offers.json`  file, which adheres to schema defined above. We can add any no of offers or remove existing ones from `offers.json`. The modifications to `offers.json` file wont require any other code changes. We can keep this `offers.json` in a database for maintainability and ease of deployments.

### Schema definition

The current implementation calculates discount when all conditions specified for the offer code are met.

#### Facts

|Fact|Data type|
|:--|--:|
| distance | decimal or integer|
| weight| decimal or integer|

#### Operators

|Operator|effect|
|:--|:--:|
|lessThan|< |
|greaterThanOrEqual|< |
|lessThanOrEqual|< |

### scripts directory

Validates `offers.json` schema

```txt
 📦scripts
 ┣ 📂 src
 ┣ 📜 package.json
 ┗ 📜 tsconfig.json
```

### Validating offers schema

1. Make sure to [setup](https://bun.sh/) 
2. Navigate to **scripts** directory
3. Install dependencies 

    ```bash
      yarn setup
    ```

4. Validate `offers.json` data (format)

    ```bash
      yarn validate
    ```

--- 

## Development & Testing

### Developing

```bash
make dev
```

### Build

This works for osx (Mac), for other operating system, refer to this [doc](https://go.dev/doc/tutorial/compile-install).

```bash
# Build
make build
# Run app
make start
```


#### Sample (1) Input & Output

```bash
    100 3
    pkg1 5 5 OFR001
    pkg2 15 5 OFR002
    pkg3 10 100 OFR003
    2 70 200
```

```bash
    pkg1 0.00 175.00
    pkg2 0.00 275.00
    pkg3 35.00 665.00
```

#### Sample (2) Input & Output

```bash
    100 5
    PKG1 50 30 OFR001
    PKG2 75 125 OFR002
    PKG3 175 100 OFR008
    PKG4 110 60 OFR002
    PKG5 155 95 NA
    2 70 200
```

#### Sample (3) Input & Output

```bash
    100 7
    PKG1 3 30 OFR001
    PKG2 2 125 OFR002
    PKG3 3 100 OFR008
    PKG4 4 60 OFR002
    PKG5 1 95 NA
    PKG6 5 95 NA
    PKG7 6 95 NA
    2 70 6
```

### Testing

```bash
make test
make coverage
```

#### Lint

Install **golangci-lint**. Refer [Local Installation](http://golangci-lint.run) for installing.

```bash
make lint
```

## TODO

- [ ] CI/CD
- [ ] There scope for improvement interns (bigO - knapsack)