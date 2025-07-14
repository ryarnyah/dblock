# DBLock [![Build Status](https://github.com/ryarnyah/dblock/actions/workflows/build.yml/badge.svg)](https://github.com/ryarnyah/dblock/)

```bash
________ __________.____                  __
\______ \\______   \    |    ____   ____ |  | __
 |    |  \|    |  _/    |   /  _ \_/ ___\|  |/ /
 |    `   \    |   \    |__(  <_> )  \___|    <
/_______  /______  /_______ \____/ \___  >__|_ \
        \/       \/        \/          \/     \/
```

Tool to maintain compatibility beetween multiple SQL database versions.

## Rules
| CODE          | Test                                         |
| :------------ | -------------------------------------------: |
| DCT02         | Add a column NOT NULL without default value  |
| DCT01         | Change Column type                           |
| DC001         | Delete Column                                |
| DT001         | Delete Table                                 |


## Installation

#### Binaries

- **linux** [amd64](https://github.com/ryarnyah/dblock/releases/download/0.4.4/dblock-linux-amd64) [386](https://github.com/ryarnyah/dblock/releases/download/0.4.4/dblock-linux-386) [arm](https://github.com/ryarnyah/dblock/releases/download/0.4.4/dblock-linux-arm) [arm64](https://github.com/ryarnyah/dblock/releases/download/0.4.4/dblock-linux-arm64)
- **windows** [amd64](https://github.com/ryarnyah/dblock/releases/download/0.4.4/dblock-windows-amd64) [386](https://github.com/ryarnyah/dblock/releases/download/0.4.4/dblock-windows-386)
- **darwin** [amd64](https://github.com/ryarnyah/dblock/releases/download/0.4.4/dblock-darwin-amd64)

```bash
sudo curl -L https://github.com/ryarnyah/dblock/releases/download/0.4.4/dblock-linux-amd64 -o /usr/local/bin/dblock && sudo chmod +x /usr/local/bin/dblock
```

#### Via Go

```bash
$ go get github.com/ryarnyah/dblock/cmd/dblock
```

#### From Source

```bash
$ mkdir -p $GOPATH/src/github.com/ryarnyah
$ git clone https://github.com/ryarnyah/dblock $GOPATH/src/github.com/ryarnyah/dblock
$ cd !$
$ make
```

#### Running with Docker
```bash
docker run ryarnyah/dblock-linux-amd64:0.4.4 <option>
```

## Usage

```bash
________ __________.____                  __
\______ \\______   \    |    ____   ____ |  | __
 |    |  \|    |  _/    |   /  _ \_/ ___\|  |/ /
 |    `   \    |   \    |__(  <_> )  \___|    <
/_______  /______  /_______ \____/ \___  >__|_ \
        \/       \/        \/          \/     \/
 Check db schema compatibility.
 Version: 0.4.4
 Build: a6d4ec3-dirty

  -alsologtostderr
        log to standard error as well as files
  -database-lock-file string
        file where database schemas will be persisted (default ".dblock.lock")
  -error-json-file string
        JSON file to write all errors
  -file-source string
        New schema in a json file (default ".new-schema.json")
  -log_backtrace_at value
        when logging hits line file:N, emit a stack trace
  -log_dir string
        If non-empty, write log files in this directory
  -logtostderr
        log to standard error instead of files
  -mssql-conn-info string
        Mssql connetion info (default "sqlserver://sa@localhost/SQLExpress?database=master&connection+timeout=30")
  -mssql-schema-regexp string
        Reex to filter schema to process (default ".*")
  -mysql-conn-info string
        MysqlQL connetion info (default "user:password@/dbname")
  -mysql-schema-regexp string
        Regex to filter schema to process (default ".*")
  -pg-conn-info string
        PostgreSQL connetion info (default "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable password=postgres")
  -pg-schema-regexp string
        Reex to filter schema to process (default ".*")
  -provider string
        DB provider (supported values: postgres, mysql, mssql, file) (default "postgres")
  -stderrthreshold value
        logs at or above this threshold go to stderr
  -v value
        log level for V logs
  -version
        Print version
  -vmodule value
        comma-separated list of pattern=N settings for file-filtered logging
```

## About

### Supported Providers

#### File (-provider file)
```bash
  -file-source string
        New schema in a json file (default ".new-schema.json")
```

#### PostgreSQL (-provider postgres)
```bash
  -pg-conn-info string
        PostgreSQL connetion info (default "host=localhost port=5432 user=postgres dbname=postgres sslmo
  -pg-schema-regexp string
        Reex to filter schema to process (default ".*")
```

#### MySQL (-provider mysql)
```bash
  -mysql-conn-info string
        MysqlQL connetion info (default "user:password@/dbname")
  -mysql-schema-regexp string
        Regex to filter schema to process (default ".*")
```

#### MSSQL (-provider mssql)
```bash
  -mssql-conn-info string
        Mssql connetion info (default "sqlserver://sa@localhost/SQLExpress?database=master&connection+timeout=30")
  -mssql-schema-regexp string
        Reex to filter schema to process (default ".*")
```
