# [zhaymm](https://github.com/isomorfisma/zhaymm)

<b>zhaymm: handle & automate your multitable mocks</b><br>
A ridiculously lightweight, DAG-based database seeder and data anonymizer written in Go. I built this for fun to generate massive relational datasets or sanitize production data for local environments in seconds.

## Overview

[zhaymm](https://github.com/isomorfisma/zhaymm) started because I have a potato device and I use traditional ORM seeders (looking at you, Node.js and Laravel) to generate test data. They are great tools, but they suffer from severe memory bloat and high network latency when you need hundreds of thousands of relational records.

So, I rewrote the logic in Go. By utilizing **zero-allocation chunking**, **bulk inserts**, and a **Directed Acyclic Graph (DAG)** execution model, zhaymm processes data natively at the speed of your hardware. 

*Note: We currently only support PostgreSQL because that's what I use daily.*

## Usage

The `zhaymm` CLI tool is completely standalone. You don't even need Go installed on your machine to run it.

### Installation

**Option 1: Download the binary (Easiest)**
Download the latest binary for your OS (Windows, macOS, Linux) from the [Releases](https://github.com/isomorfisma/zhaymm/releases) page and extract it. 

To use it globally without typing `./zhaymm` every time, export it to your system's PATH. Add this line to your `~/.bashrc` or `~/.zshrc`:
```bash
export PATH=$PATH:/path/to/your/extracted/zhaymm/folder
```
Then, reload your terminal:
```bash
source ~/.bashrc  # or source ~/.zshrc
```

**Option 2: Build it from source**
If you're a Go developer, just compile it directly:
```bash
$go install [github.com/isomorfisma/zhaymm/cmd/zhaymm@latest$](https://github.com/isomorfisma/zhaymm/cmd/zhaymm@latest$) zhaymm --help
```

### 1. Initialization
Generate a boilerplate schema configuration so you don't have to write YAML from scratch:
```bash
$ zhaymm init
```
This creates a `schema.yaml` file. Define your table structures, row counts, dependencies, and faker functions here.

### 2. Seeding data (mock generation)
Generate 100% synthetic mock data based on your schema and blast it into your local database:
```bash
$ zhaymm seed --config schema.yaml --db "postgres://user:pass@localhost:5432/local_db?sslmode=disable"
```

### 3. Subsetting & anonymizing (ETL)
Pull real data from a staging/prod database, anonymize sensitive columns on-the-fly, and load it into your local DB safely:
```bash
$ zhaymm pull --config schema.yaml \
    --source-db "postgres://prod_user:pass@remote-server/prod_db" \
    --target-db "postgres://local_user:pass@localhost:5432/local_db"
```

## Supported Faker Functions

zhaymm ships with a pretty sweet expression engine. You can use these functions directly inside the `columns` definition of your `schema.yaml`:

- **Identity:** `uuid()`, `person_name()`, `gender()`, `ssn()`, `job_title()`
- **Internet:** `email()`, `phone()`, `password()`, `ipv4()`, `mac_address()`
- **Location:** `city()`, `country()`, `street()`, `latitude()`, `longitude()`
- **Finance:** `credit_card()`, `company()`, `currency()`, `price()`
- **Relational:** `random_ref('table_name')`, `random_int(min, max)`

## Why did I build this?

- **Just for fun (and speed):** I wanted to see how fast Go could actually go. Turns out, it can generate 150,000+ relational rows in ~20 seconds using < 60MB of RAM. 
- **Safe Local Debugging:** Subsetting allows me to debug real production data without violating GDPR/PDP by automatically masking sensitive PII.
- **Zero-Dependency DX:** QA engineers and frontend devs shouldn't need to run `npm install` or `composer install` just to populate a local test database. Just run the binary and go.

## Shoutouts

zhaymm stands on the shoulders of giants. Check out the core engines powering this tool:
- [brianvoe/gofakeit](https://github.com/brianvoe/gofakeit): The comprehensive dummy data generator for Go.
- [expr-lang/expr](https://github.com/expr-lang/expr): The ridiculously fast expression evaluator.
