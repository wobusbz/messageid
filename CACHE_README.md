# Game Server Cache

This repository contains both a protobuf message ID generator plugin and a game server cache implementation that uses gRPC for communication instead of Redis.

## Architecture

The cache server follows a hierarchical structure based on database table relationships defined in JSON configuration. It supports:

- **Master Tables**: Primary data entities like players, guilds, etc.
- **Slave Tables**: Child data entities that belong to master tables
- **Dynamic Key Resolution**: Uses path-like keys with parameter substitution

## Key Features

- **gRPC Communication**: Fast, type-safe communication using Protocol Buffers
- **Hierarchical Caching**: Tree-like cache structure following directory-like paths
- **Key Validation**: Automatic validation against configured table patterns
- **Thread-Safe**: Concurrent read/write operations using mutex locks
- **No External Dependencies**: Self-contained in-memory cache (no Redis required)

## Configuration

The cache server uses JSON configuration to define table structures and their Statdir patterns:

```json
{
    "MasterTables": [
        {
            "SqlName": "players",
            "Statdir": "/players/{plydbid}",
            "Slaves": ["db1.rolcampwall", "db1.rolcampnpc"]
        }
    ],
    "SlaveTables": [
        {
            "SqlName": "rolcampwall", 
            "Statdir": "/players/{plydbid}/campwall/{wallpos}"
        }
    ]
}
```

## Cache Key Examples

Based on the Statdir patterns, the cache accepts keys like:

- `/players/12345` - Player data for player ID 12345
- `/players/12345/campwall/1` - Camp wall position 1 for player 12345
- `/guilds/67890` - Guild data for guild ID 67890
- `/guilds/67890/guildmembers/12345` - Guild member data

## Building

### Build the cache server:
```bash
go build -o cache-server ./cmd/cache-server
```

### Build the example client:
```bash
go build -o cache-client ./cmd/cache-client
```

### Build the original protoc plugin:
```bash
go build -o protoc-gen-messageid .
```

## Running

### Start the cache server:
```bash
./cache-server -port 8080 -config config/tables.json
```

### Run the example client:
```bash
./cache-client -server localhost:8080
```

## gRPC API

The cache service provides three main operations:

### Get
```protobuf
rpc Get(GetRequest) returns (GetResponse);
```

### Set
```protobuf
rpc Set(SetRequest) returns (SetResponse);
```

### Delete
```protobuf
rpc Delete(DeleteRequest) returns (DeleteResponse);
```

## Testing

Run the test suite:
```bash
go test ./cache/
```

## Usage Example

```go
// Connect to cache server
conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
client := cache.NewCacheServiceClient(conn)

// Set player data
setResp, err := client.Set(ctx, &cache.SetRequest{
    Key:  "/players/12345",
    Data: []byte(`{"name": "Player1", "level": 50}`),
})

// Get player data
getResp, err := client.Get(ctx, &cache.GetRequest{
    Key: "/players/12345",
})

// Delete player data
deleteResp, err := client.Delete(ctx, &cache.DeleteRequest{
    Key: "/players/12345",
})
```

## Key Resolution

The cache automatically resolves dynamic parameters in Statdir patterns:

- Pattern: `/players/{plydbid}/campwall/{wallpos}`
- Resolved: `/players/12345/campwall/1`

Keys are validated against configured patterns to ensure data integrity and prevent invalid cache entries.