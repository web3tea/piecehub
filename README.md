# piecehub
A high-performance Filecoin piece storage server.

## Installation

```bash
go install github.com/strahe/piecehub/cmd/piecehub@latest
```

## Quick Start

1. Simple Directory Mode:
```bash
# piecehub dir path1 path2 ...
# Example:

piecehub dir /data/pieces1 /data/pieces2 /data/pieces3
```

2. Mutiple Storage Mode:

Create a configuration file `config.toml`:

```toml
[server]
address = ":8080"
read_timeout = 30
write_timeout = 30

[[disks]]
name = "local1"
root_dir = "/data/pieces1"
max_size = 1073741824  # 1GB
direct_io = true

[[disks]]
name = "local2"
root_dir = "/data/pieces2"
max_size = 1073741824  # 1GB
direct_io = true

[[s3s]]
name = "remote1"
endpoint = "https://s3.amazonaws.com"
region = "us-east-1"
bucket = "my-pieces2"
access_key = "xxx"
secret_key = "xxx"

[[s3s]]
name = "remote2"
endpoint = "https://s3.amazonaws.com"
region = "us-east-1"
bucket = "my-pieces2"
access_key = "xxx"
secret_key = "xxx"
```

Start the server:

```bash
piecehub -c config.toml
```

## API

### Check Piece Existence
```http
HEAD /pieces?id=<pieceID>
```
or
```http
GET /pieces?id=<pieceID>
```

### Get Piece Data
```http
GET /data?id=<pieceID>
```

### Download with curl
```bash
# Check existence
curl -I http://localhost:8080/pieces?id=xxx

# Download piece
curl -O http://localhost:8080/data?id=xxx
```
