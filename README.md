# own-redis

## Description

`own-redis` is a minimalist implementation of an in-memory key-value database inspired by Redis. It uses the **UDP protocol** for communication and supports basic commands: `PING`, `SET`, and `GET`, including expiration via `PX`.

This project is written in **Go**, without any third-party libraries. The main purpose is to gain a deeper understanding of **UDP networking**, **concurrent access handling**, and **key-value storage principles**.

## Features

- Communication over the **UDP protocol**
- In-memory key-value storage
- Supports commands:
  - `PING` — check if the server is alive
  - `SET` — store a value by key
  - `GET` — retrieve a value by key
  - `SET ... PX <ms>` — store with expiration in milliseconds
- Automatic removal of expired keys (via background cleanup)
- Thread-safe access using `sync.RWMutex`
- Command-line flags: `--port`, `--help`

### Build and Run

```bash
go build -o own-redis .
```

### Run with the default port 8080:
```bash
./own-redis
```

### Run with a custom port:
```bash
./own-redis --port 7828
```

### Display usage help:
```bash
./own-redis --help
```

## Usage Example

You can test the server using `netcat` (nc):
```bash
nc -u 127.0.0.1 8080
```

### Commands

## Supported Commands

| Description         | Command Format                     | Example                                    | Server Response |
|---------------------|------------------------------------|--------------------------------------------|-----------------|
| Check server status | `PING`                             | `PING`                                     | `PONG`          |
| Set a key's value   | `SET <key> <value>`                | `SET name Zako`                            | `OK`            |
| Set with TTL        | `SET <key> <value> PX <milliseconds>` | `SET temp 123 PX 5000`                 | `OK`            |
| Get a key's value   | `GET <key>`                        | `GET name`                                 | `Zako` or `(nil)` |

> Notes:
> - PX sets expiration in milliseconds. After that, the key is automatically deleted.
> - `(nil)` is returned if the key doesn't exist or expired.
> - The server responds and moves to a new line automatically after each command.

```text
SET foo                  --> (error) ERR wrong number of arguments for SET command
GET                      --> (error) ERR wrong number of arguments for GET command
SET a b PX               --> (error) ERR syntax error
SET a b PX abc           --> (error) ERR value is not an integer or out of range
SET a b PX 100 PX 200    --> (error) ERR PX already specified
SET a b PX 100 extra     --> (error) ERR syntax error after PX
```

## Project Structure

```tree
├── README.md
├── go.mod
├── internal
│   ├── flags
│   │   └── flags.go
│   ├── server
│   │   ├── handlers.go
│   │   └── server.go
│   └── utils
│       └── usage.go
└── main.go
```

## Internal Details
- The server stores data in a map[string]valueEntry where valueEntry contains the value and expiration timestamp.

- A sync.RWMutex ensures safe concurrent access.

- A background goroutine removes expired keys every second.

- Each UDP request is handled in a separate goroutine.

## Testing

You can test using `nc`:
```bash
nc -u 127.0.0.1 8080
```

## Example session:
```text
PING
PONG

SET key1 value1
OK

GET key1
value1

SET key2 temp PX 2000
OK

GET key2
temp

# Wait 2 seconds
GET key2
(nil)
```

# Made with ❤️ by [zaaripzha](https://platform.alem.school/git/zaaripzha) aka [1tzme](https://github.com/1tzme)