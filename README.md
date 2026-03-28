# [zhaymm](https://github.com/isomorfisma/zhaymm)

<b>zhaymm : handle & automate your multitable mocks! <br></b>
Lightweight, DAG-based database seeder and data anonymizer written in Go. Generate massive relational datasets or sanitize production data for local environments in seconds.

## Overview

[zhaymm](https://github.com/isomorfisma/zhaymm) aims to solve the bottleneck of database seeding and ETL subsetting for developers. 

Traditional ORM seeders (like those in Laravel or Node.js) suffer from severe memory bloat and high network latency when generating hundreds of thousands of relational records, while zhaymm utilizes **zero-allocation chunking**, **bulk inserts**, and a **Directed Acyclic Graph (DAG)** execution model to process data natively at the speed of your hardware.

We currently only support PostgreSQL.

## Usage

The `zhaymm` CLI tool is designed to be standalone. You do not need Go installed on your machine to run the pre-compiled binaries.

### Installation

Download the latest binary for your OS (Windows, macOS, Linux) from the [Releases](https://github.com/isomorfisma/zhaymm/releases) page. Extract it and ensure it's in your system's `PATH`.

Alternatively, if you have Go installed, you can compile it directly:

```bash
$ go install github.com/isomorfisma/zhaymm/cmd/zhaymm@latest
$ zhaymm --help
```

### 1. Initialization
To generate a boilerplate schema configuration, use the `init` command:

```bash
$ zhaymm init
```
This will create a `schema.yaml` file in your current directory. You can define your table structures, row counts, dependencies, and faker functions here.

### 2. Seeding data (mock generation)
To generate 100% synthetic mock data based on your schema and push it to a local database:

```bash
$ zhaymm seed --config schema.yaml --db "postgres://user:pass@localhost:5432/local_db?sslmode=disable"
```

### 3. Subsetting & anonymizing (ETL)
To pull real data from a production/staging database, anonymize sensitive columns on-the-fly (as defined in your schema), and load it into a local database:

```bash
$ zhaymm pull --config schema.yaml \
    --source-db "postgres://prod_user:pass@remote-server/prod_db" \
    --target-db "postgres://local_user:pass@localhost:5432/local_db"
```

## Supported Faker Functions

zhaymm ships with a powerful expression engine. You can use the following functions directly inside the `columns` definition of your `schema.yaml`:

- **Identity:** `uuid()`, `person_name()`, `gender()`, `ssn()`, `job_title()`
- **Internet:** `email()`, `phone()`, `password()`, `ipv4()`, `mac_address()`
- **Location:** `city()`, `country()`, `street()`, `latitude()`, `longitude()`
- **Finance:** `credit_card()`, `company()`, `currency()`, `price()`
- **Relational:** `random_ref('table_name')`, `random_int(min, max)`

## Purpose

- **Eradicate CPU/RAM Bottlenecks**: Generate 150,000+ relational rows in ~20 seconds using < 60MB of RAM.
- **Safe Local Development**: Subsetting allows developers to debug real production data without violating GDPR/PDP by automatically masking sensitive PII (Personally Identifiable Information).
- **Zero-Dependency DX**: Ship as a single binary. QA engineers and frontend devs don't need to install Node modules or PHP dependencies just to populate their local DB.

## Related Projects

zhaymm stands on the shoulders of giants. Check out the core engines powering this tool:
- [brianvoe/gofakeit](https://github.com/brianvoe/gofakeit): The comprehensive dummy data generator for Go.
- [expr-lang/expr](https://github.com/expr-lang/expr): The ridiculously fast expression evaluator.

