# filegate

A CLI tool to instantly share files from your current directory over the network.

## Features

- **WebDAV Server** - Expose any directory as a WebDAV share, compatible with Windows, macOS, and Linux file managers
- **Public URLs** - Get a public URL instantly via the relay server (no port forwarding needed)
- **Local Mode** - Run on your LAN without external dependencies
- **DLNA Server** - Stream media files to smart TVs and media players
- **Basic Auth** - Password protection with auto-generated secure passwords

## Installation

```bash
go install github.com/filegate/filegate/cmd/filegate@latest
```

Or build from source:

```bash
make build
```

## Usage

### WebDAV (Public URL)

Share the current directory with a public URL:

```bash
filegate
```

This connects to the relay server and gives you a URL like `https://brave-tiger.filegate.app` that anyone can access.

### WebDAV (Local Network)

Share on your local network only:

```bash
filegate webdav --local
```

Or specify a custom port:

```bash
filegate webdav --local --port 9000
```

### DLNA (Smart TV)

Stream media to DLNA-compatible devices:

```bash
filegate dlna
```

Your smart TV should automatically discover the server.

## Options

### WebDAV Command

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--local` | `-l` | Run in local mode (LAN only) | `false` |
| `--port` | `-p` | Port to listen on (local mode) | `8080` |
| `--user` | `-u` | Username for Basic Auth | `admin` |
| `--pass` | | Password (auto-generated if omitted) | |

### DLNA Command

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--port` | `-p` | Port to listen on | `8080` |
| `--name` | `-n` | Server name | hostname |

## Connecting to WebDAV

### Windows
1. Open File Explorer
2. Right-click "This PC" → "Add a network location"
3. Enter the URL provided by filegate

### macOS
1. In Finder, press `Cmd+K`
2. Enter the URL provided by filegate

### Linux
Most file managers support WebDAV. In Nautilus, press `Ctrl+L` and enter:
```
davs://brave-tiger.filegate.app
```

## Architecture

```
┌─────────────┐         ┌─────────────────┐         ┌─────────────┐
│   Client    │ ──────► │  Relay Server   │ ◄────── │  filegate   │
│  (Browser)  │  HTTPS  │  filegate.app    │   WSS   │    CLI      │
└─────────────┘         └─────────────────┘         └─────────────┘
```

- **CLI** runs on your machine, serving files via WebDAV
- **Relay Server** provides public URLs and tunnels HTTP requests over WebSocket
- **Clients** connect via the public URL, requests are forwarded to your CLI

## Self-Hosting the Relay

Run your own relay server:

```bash
# Using Docker
docker-compose up -d

# Or directly
go run ./cmd/relay --domain yourdomain.com
```

Then point your CLI to it:

```bash
filegate webdav --relay wss://yourdomain.com/tunnel
```

## Development

```bash
# Build everything
make build

# Run tests
make test

# Build for all platforms
make build-all
```

## License

MIT
