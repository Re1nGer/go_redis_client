# Redis Client for Go

This is a lightweight, feature-rich Redis client for Go. It provides a simple and efficient way to interact with Redis servers, supporting various Redis commands and connection options with RESP v2.

## Features

- Simple and intuitive API
- Support for multiple Redis commands
- Connection pooling (placeholder, not yet implemented)
- Custom connection functions
- SSL/TLS support (placeholder, not yet implemented)
- Authentication support
- Flexible configuration options
- Socket timeout and keep-alive settings

## Installation

To install the Redis client, use the following command:

```
go get github.com/Re1nGer/go_redis_client@v1.2.0
```

## Usage

### Connecting to Redis

To connect to a Redis server, use the `NewClient` function:

```go
import "github.com/re1nger/redis-client-go"

client, err := redisclient.NewClient("localhost", 6379)
if err != nil {
    // Handle error
}
defer client.Close()
```

#### Connection Options

You can provide additional options when connecting:

```go
client, err := redisclient.NewClient("localhost", 6379,
    redisclient.WithUsername("myuser"),
    redisclient.WithPassword("mypassword"),
    redisclient.DecodeResponses(true),
    redisclient.WithClientName("my-app"),
    redisclient.WithSocketConnectTimeout(5 * time.Second),
    redisclient.WithSocketTimeout(10 * time.Second),
    redisclient.WithSocketKeepAlive(true, 30 * time.Second),
)
if err != nil {
    // Handle error
}
defer client.Close()
```

#### Custom Connection Function

You can also provide a custom connection function:

```go
customConnFunc := func() (net.Conn, error) {
    // Your custom connection logic here
    return net.Dial("tcp", "localhost:6379")
}

client, err := redisclient.NewClient("localhost", 6379,
    redisclient.WithCustomConnectFunc(customConnFunc),
)
if err != nil {
    // Handle error
}
defer client.Close()
```

### Executing Commands

To execute Redis commands, use corresponding command names as methods:

```go
// Using a specific command method
result, err := client.Set("mykey", "myvalue")
if err != nil {
    // Handle error
}

// Getting a value
value, err := client.Get("mykey")
if err != nil {
    // Handle error
}
```

### Using SET with Options

The client provides a flexible way to use the SET command with various options:

```go
opts := redisclient.NewSetOpts().WithNX().WithEX(60)
result, err := client.Set("mykey", "myvalue", opts)
if err != nil {
    // Handle error
}
```

To execute any command use `Do` method:

```go
// Using the general Do method
result, err := client.Do("SET", "mykey", "myvalue")
if err != nil {
    // Handle error
}
```

## Available Commands

The client supports the following Redis commands:

### String Operations
- GET
- SET (with options: NX, XX, EX, PX, EXAT, PXAT, KEEPTTL, GET)
- SETNX
- SETRANGE
- STRLEN
- APPEND
- DECR
- DECRBY
- GETDEL
- GETEX
- GETRANGE
- GETSET
- INCR
- INCRBY
- INCRYBYFLOAT
- LCS
- MGET
- MSET
- MSETNX
- PSETEX

### List Operations
- LPOP
- LINDEX
- LINSERT
- LLEN
- LMOVE
- LMPOP - pending
- LPOS
- LPUSH
- LPUSHX
- LRANGE
- LREM
- LSET
- LTRIM
- RPOP
- RPOPLPUSH - deprecated 
- RPUSH
- RPUSHX

### Set Operations
- SADD
- SCARD
- SDIFF
- SDIFFSTORE
- SINTER
- SINTERCARD
- SINTERSTORE
- SISMEMBER
- SMEMBERS
- SMISMEMBER
- SMOVE
- SPOP
- SRANDMEMBER
- SREM
- SSCAN
- SUNION
- SUNIONSTORE


### Hash Operations
- HDEL
- HSET
- HGET
- HEXISTS
- HGETALL
- HLEN
- HEXPIRETIME
- HINCRBY
- HINCRBYFLOAT
- HMGET
- HPERSIST
- HSTRLEN
- HEXPIRE

### Bitmap
- BITCOUNT
- BITFIELD
- BITFIELD_RO
- BITOP
- GETBIT
- SETBIT

### Key Operations
- DEL
- DUMP

### Connection Management
- AUTH
- CLIENT CACHING
- CLIENT GETNAME
- CLIENT GETREDIR
- CLIENT ID
- CLIENT INFO
- CLIENT KILL
- CLIENT LIST
- CLIENT NO-NOEVICT
- CLIENT NO-TOUCH
- CLIENT PAUSE
- CLIENT REPLY
- CLIENT SETINFO
- CLIENT SETNAME
- CLIENT TRACKING
- CLIENT TRACKINGINFO
- CLIENT UNBLOCK
- CLIENT UNPAUSE
- ECHO
- PING
- QUIT (deprecated)

### Server Management 
- ACL CAT
- ACL DELUSER
- ACL DRYRUN
- ACL GENPASS
- ACL GETUSER
- ACL LIST
- ACL LOAD
- ACL SAVE
- ACL SETUSER
- ACL USERS
- ACL WHOAMI
- BGREWRITEAOF
- BGSAVE
- COMMAND
- COMMAND COUNT
- COMMAND DOCS
- COMMAND GETKEYS
- COMMAND GETKEYSANDFLAGS
- COMMAND INFO
- COMMAND LIST
- CONFIG GET
- CONFIG RESETSTAT
- CONFIG REWRITE
- CONFIG SET
- DBSIZE
- FAILOVER
- FLUSHALL
- FLUSHDB
- INFO
- LASTSAVE
- LATENCY DOCTOR
- LATENCY GRAPH
- LATENCY HISTOGRAM
- LATENCY HISTORY
- LATENCY LATEST
- LATENCY RESET
- MEMORY DOCTOR
- MEMORY MALLOC-STATS
- MEMORY PURGE
- MEMORY STATS
- MEMORY USAGE
- MODULE UNLOAD
- MONITOR
- PSYNC
- REPLCONF
- REPLICAOF
- ROLE
- SAVE
- SLOWLOG LEN
- SLOWLOG RESET
- SWAPDB
- SYNC
- TIME

More commands will be added in future updates.

## Configuration Options

The client supports various configuration options:

- `WithUsername(username string)`: Set the username for authentication
- `WithPassword(password string)`: Set the password for authentication
- `WithClientName(clientName string)`: Set a custom client name
- `WithCustomConnectFunc(connFunc RedisConnectFunc)`: Provide a custom connection function
- `WithSocketConnectTimeout(duration time.Duration)`: Set the timeout for connecting to the Redis server
- `WithSocketTimeout(duration time.Duration)`: Set the timeout for socket operations
- `WithSocketKeepAlive(enabled bool, interval time.Duration)`: Enable and set the interval for TCP keep-alive packets

## Error Handling

The client returns errors for various scenarios, including connection failures, authentication errors, and Redis command errors. Always check the returned error and handle it appropriately in your application.

## Contributing

Contributions to this Redis client are welcome! Here's how you can contribute:

1. Fork the repository
2. Create a new branch for your feature or bug fix
3. Write your code and tests
4. Ensure all tests pass
5. Submit a pull request with a clear description of your changes

### Under Development

- Implement connection pooling
- Add support for more Redis commands
- Implement SSL/TLS support
- Add more comprehensive error handling
- Add benchmarks tests
- Implement Context to add timeouts and command cancellations 

## License

This Redis client is released under the MIT License. See the [LICENSE](https://github.com/Re1nGer/go_redis_client/blob/master/LICENSE) file for details.

## Disclaimer

This client is a work in progress and may not be suitable for production use without further testing and development. Use at your own risk.

## Support

For bug reports, feature requests, or general questions, please open an issue on the GitHub repository.

---
