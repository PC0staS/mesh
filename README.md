# MESH - Monitor Each Server Health

MESH is a CLI tool to monitor servers and services easily and extensibly, using a local daemon and Unix socket communication.

## Features

- Add, remove, and list monitored servers
- Supports monitoring types: ping, http, tcp, ssh
- Configurable intervals and timeouts
- Configurable alert webhooks (json or string)
- Local daemon for background execution
- Efficient communication via Unix sockets

## Installation

### Install with Snap (Recommended)

You can install Mesh easily using Snap:

```sh
sudo snap install mesh
```

This is the easiest way to get started on most Linux distributions that support Snap.

### Manual Installation

1.  Clone the repository:

```sh
git clone https://github.com/PC0staS/mesh.git
cd mesh
```

2.  Build the project:

```sh
go build -o mesh
```

Alternatively, you can use Makefile targets:

      ```sh
      make build   # Build the binary
      make install # Install mesh to /usr/local/bin (may require sudo)
      make install-daemon # Install the mesh.service - You can use this one for a complete install.
      ```

This will copy the compiled mesh binary to /usr/local/bin so you can run mesh from anywhere.

## Quick Start

1. Start the daemon:
   ```sh
   ./mesh start
   ```
2. Add a server:
   ```sh
   ./mesh add
   ```
3. List servers:
   ```sh
   ./mesh list
   ```
4. Remove a server:
   ```sh
   ./mesh remove
   ```
5. Stop the daemon:
   ```sh
   ./mesh stop
   ```

## Project Structure

- `cmd/` — CLI commands (add, list, remove, start, stop, monitor)
- `internal/client/` — Client for daemon communication
- `internal/daemon/` — Daemon logic and monitoring
- `internal/config/` — Configuration and persistence
- `internal/monitor/` — Monitoring types and logic

## Configuration Example

The configuration file is stored at `~/.config/mesh/config.json` and contains the monitored servers.

## Requirements

- Go 1.18+
- Linux (requires Unix sockets)

## Contributing

Pull requests and suggestions are welcome!

## License

MIT
