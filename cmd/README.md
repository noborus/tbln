# TBLN CLI tool

Import/Export TBLN file and RDBMS table.
Also sign and verify the TBLN file.

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
