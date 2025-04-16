# piecehub
A Filecoin piece storage server for curio market.

## Installation

```bash
go install github.com/web3tea/piecehub/cmd/piecehub@latest
```

## Quick Start

### 1. Simple Directory Mode:
```bash
piecehub dir /data/pieces1 [/data/pieces2 ...]
```

For more options:
```bash
piecehub dir -h
```

### 2. Simple S3 Mode:
```bash
piecehub s3 --endpioint xx --ak xx --sk xx bucket1 [bucket2 ...]
```

For more options:
```bash
piecehub s3 -h
```

### 3. Hybrid Storage Mode:

Create a configuration file `config.toml`:

```toml
[server]
address = ":8080"
read_timeout = 30
write_timeout = 30
tokens = ["xxx"]

[[disks]]
name = "local1"
root_dir = "/data/pieces1"

[[disks]]
name = "local2"
root_dir = "/data/pieces2"

[[s3s]]
name = "remote1"
endpoint = "https://s3.amazonaws.com"
region = "us-east-1"
bucket = "my-pieces2"
access_key = "xxx"
secret_key = "xxx"
use_ssl = false
prefix = ""

[[s3s]]
name = "remote2"
endpoint = "https://s3.amazonaws.com"
region = "us-east-1"
bucket = "my-pieces2"
access_key = "xxx"
secret_key = "xxx"
use_ssl = false
prefix = ""
```

Start the server:

```bash
piecehub -c config.toml
```

### 4. Authentication

No authentication by default.

To enable authentication, set the `tokens` field in the configuration file or use the command line option.

1. Add token to the configuration file:

```toml
[server]
tokens = ["token1", "token2"]
```

2. Use the command line option:

```bash
piecehub --token token1 --token token2 -c config.toml
```



## API

### Check Piece Existence
```http
HEAD /pieces?id=<pieceCid>
```

### Get Piece Data
```http
GET /pieces?id=<pieceCid>
```

### List Storage Name
```http
GET /storages
```

### Examples

Using curl:

```bash
# Check if piece exists
curl -I "http://localhost:8080/pieces?id=<pieceCid>"

# Download piece
curl -O "http://localhost:8080/pieces?id=<pieceCid>"

# With token
curl -H "Authorization: your-token" -O "http://localhost:8080/pieces?id=<pieceCid>"

# Generatge car file
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"size":268435456,"storageName":"test-storage-name"}' \
  http://localhost:8080/debug/generate-car

# Response
{
    "pieceCid":"baga6ea4seaqb46zh6n4fig7nuf5lmfylxr4flmzu2tgfjm6k4werggcnp3fvspy",
    "pieceSize":536870912,
    "payloadSize":268445499,
    "carSize":268445499,
    "carCid":"bafkreibq4fevl27rgurgnxbp7adh42aqiyd6ouflxhj3gzmcxcxzbh6lla"
}
```
