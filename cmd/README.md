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
  keystore    keystore is a command to operate keystore
  sign        Sign a TBLN file with a private key
  verify      Verify signature and checksum of TBLN file

Flags:
  -h, --help              help for tbln
  -k, --keyname string    key name
      --keypath string    key store path
      --keystore string   keystore file
      --pubfile string    public Key file
      --secfile string    secret Key file
```

## Basic example

Export the database table and output the TBLN file.

```sh
$ tbln export --db postgres --dsn "host=localhost dbname=sampletest" -t simple -o simple.tbln
```

Import the TBLN file into the database.

```sh
$ tbln import --db postgres --dsn "host=localhost dbname=sampletest" -t simple2 -f simple.tbln
```

Data type and primary key are restored in this example.

```sh
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

```sh
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

```sh
$ tbln sign testdata/simple.tbln
```

You can sign the specified file by entering the previously entered password.

## Signature verification

Signature verification verifies signatures with the public key contained in the keystore.

```sh
$ tbln verify simple.tbln
2019/03/24 00:33:50 Signature verification successful
```

You can also treat a public key as a keystore instead of a keystore.

The signature verification included in this repository should be successful.

```sh
$ tbln verify --keystore testdata/test.pub testdata/simple.tbln
2019/03/24 00:33:50 Signature verification successful
```
