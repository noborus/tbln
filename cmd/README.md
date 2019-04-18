# TBLN CLI tool

Import/Export TBLN file and RDBMS table.
MERGE and EXCEPT (difference set) are possible from DB tables and files.
Supports digital signatures and verification for TBLN files.

## Install

```bash
$ go get -u github.com/noborus/tbln/cmd/tbln
```

### Note

* Requires a version of Go that supports modules. e.g. Go 1.11+
* The C compiler is required because the driver for SQLIte3 is included.

## Usage

```
Usage:
  tbln [command]

Available Commands:
  diff        Diff two TBLNs
  except      Except for other TBLN rows from self TBLN
  export      Export database table or query
  genkey      Generate a new key pair
  hash        Add hash value to TBLN file
  help        Help about any command
  import      Import database table
  keystore    keystore is a command to operate keystore
  merge       Merge two TBLNs
  sign        Sign a TBLN file with a private key
  verify      Verify signature and checksum of TBLN file
  version     Print the version number of tbln

Flags:
      --debug             debug output
  -h, --help              help for tbln
  -k, --keyname string    key name
      --keypath string    key store path
      --keystore string   keystore file
      --pubfile string    public Key file
      --secfile string    secret Key file
```

## Basic example

Export the database table and output the TBLN file.

The database driver name(--db) and DSN(--dsn) are required
to import/export database tables.

```bash
$ tbln export --db postgres --dsn "host=localhost dbname=sampletest" \
   -t simple -o simple.tbln
```

Import the TBLN file into the database.

```bash
$ tbln import --db postgres --dsn "host=localhost dbname=sampletest" \
  -t simple2 -f simple.tbln
```

Data type and primary key are restored in this example.

```bash
$ psql sampletest
# \d simple2
              Table "public.simple2"
 Column |  Type   | Collation | Nullable | Default
--------+---------+-----------+----------+---------
 id     | integer |           | not null |
 name   | text    |           |          |
Indexes:
    "simple2_pkey" PRIMARY KEY, btree (id)
```

## Signature example

First generate the private key and the public key.

```bash
$ tbln genkey
```

The above will generate ***\$(HOMEDIR)/.Tbln/\$(USER).key***
and ***\$(HOME)/.Tbln/\$(USER).pub***.

The file location and name can be changed by options.
```
  -k, --keyname string    key name
      --keypath string    key store path
      --pubfile string    public Key file
      --secfile string    secret Key file
```
You will be prompted to enter a password.

When generating a key, the private key is encrypted with a password and stored.

Passwords can be empty.

The generated public key is simultaneously registered in the keystore.
The keystore location is ***$(HOME)/.tbln/keystore.tbln***.
he keystore file can be changed optionally.
```
      --keypath string    key store path
```

Signing with a private key is possible, if you have generated a key.

```bash
$ tbln sign testdata/simple.tbln
```

You can sign the specified file by entering the previously entered password.

## Signature verification

Signature verification verifies signatures with the public key contained in the keystore.

```bash
$ tbln verify simple.tbln
2019/03/24 00:33:50 Signature verification successful
```

You can also treat a public key as a keystore instead of a keystore.

The signature verification included in this repository should be successful.

```bash
$ tbln verify --keystore testdata/test.pub testdata/simple.tbln
2019/03/24 00:33:50 Signature verification successful
```

## Merge

Merge other table or file to a self table or file.
Merging other files into its self table is the same as import.

Merge tables from another database.

This is an example of synchronizing MySQL tables with PostgreSQL.
With the --delete option, extra rows are deleted.

```
tbln merge --self-db "postgres" --self-dsn "database=test_db" --self-table simple \
           --other-db "mysql" --other-dsn "root@/test_db" --other-table simple \
           --delete
```

## Except

Except extracts differences as SQL except and outputs it as a TBLN file.
Remove the other rows from self rows and output the remaining rows.