# TBLN CLI tool

Import/Export TBLN file and RDBMS table.
Supports digital signatures and verification for TBLN files.

```
Usage:
  tbln [command]

Available Commands:
  export      Export database table or query
  genkey      Generate a new key pair
  hash        Add hash value to TBLN file
  help        Help about any command
  import      Import database table
  sign        Sign a TBLN file with a private key
  verify      Verify signature and checksum of TBLN file

Flags:
  -h, --help             help for tbln
  -k, --keyname string   key name
      --keypath string   key store path
      --pubfile string   public Key File
      --seckey string    Secret Key File
```

## example

Export the database table and output the TBLN file.

```sh
$ tbln export --db postgres --dsn "host=localhost dbname=sampletest" -t simple -o simple.tbln
```

Import the TBLN file into the database.

```sh
$ tbln import --db postgres --dsn "host=localhost dbname=sampletest" -t simple2 -f simple.tbln
```
