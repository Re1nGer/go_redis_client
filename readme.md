
# Redis Client for Go

This is a lightweight, feature-rich Redis client for Go. It provides a simple and efficient way to interact with Redis servers, supporting various Redis commands and connection options with RESP v2

## Features

- Simple and intuitive API
- Support for multiple Redis commands
- Connection pooling (placeholder, not yet implemented)
- Custom connection functions
- SSL/TLS support (placeholder, not yet implemented)
- Authentication support
- Flexible configuration options

## Installation

To install the Redis client, use the following command:

```
go get github.com/re1nger/redis-client-go
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

You can also provide additional options:

```go
client, err := redisclient.NewClient("localhost", 6379,
    redisclient.WithUsername("myuser"),
    redisclient.WithPassword("mypassword"),
    redisclient.DecodeResponses(true),
    redisclient.WithClientName("my-app"),
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
result, err := client.SetWithOptions("mykey", "myvalue", opts)
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

The client currently supports the following Redis commands:

- GET
- SET (with options: NX, XX, EX, PX, EXAT, PXAT, KEEPTTL, GET)
- SETNX
- SETRANGE
- STRLEN
- ECHO
- EXISTS
- LPOP (partially implemented)
- SADD
- SCARD
- SDIFF
- SDIFFSTORE
- SINTER
- SINTERCARD

More commands will be added in future updates.

## Configuration Options

The client supports various configuration options:

- `WithUsername(username string)`: Set the username for authentication
- `WithPassword(password string)`: Set the password for authentication
- `DecodeResponses(shouldDecode bool)`: Enable/disable automatic response decoding
- `WithClientName(clientName string)`: Set a custom client name
- `WithCustomConnectFunc(connFunc RedisConnectFunc)`: Provide a custom connection function

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

This Redis client is released under the MIT License. See the [LICENSE](LICENSE) file for details.

## Disclaimer

This client is a work in progress and may not be suitable for production use without further testing and development. Use at your own risk.

## Support

For bug reports, feature requests, or general questions, please open an issue on the GitHub repository.

---