# Configuration

This directory contains configuration files for NetLSR.

## Configuration Files

- `config.yaml` - Complete configuration with both client and server settings
- `client.yaml` - Client-specific configuration
- `server.yaml` - Server-specific configuration

## Usage

### 1. Using Configuration Files

```bash
# Run with default configuration
./netlsr -mode client

# Run with custom configuration file
./netlsr -mode client -config ./configs/config.yaml

# Run server with custom configuration
./netlsr -mode server -config ./configs/server.yaml
```

### 2. Using Environment Variables

You can override any configuration value using environment variables:

```bash
# Client environment variables
export NETLSR_CLIENT_TUN_NAME="tun1"
export NETLSR_CLIENT_LOCAL_IP="10.0.0.10"
export NETLSR_CLIENT_PEER_IP="10.0.0.20"
export NETLSR_CLIENT_SERVER_ADDR="192.168.1.200"
export NETLSR_CLIENT_PORT="9090"

# Server environment variables
export NETLSR_SERVER_TUN_NAME="tun1"
export NETLSR_SERVER_LOCAL_IP="10.0.0.20"
export NETLSR_SERVER_PEER_IP="10.0.0.10"
export NETLSR_SERVER_PORT="9090"
export NETLSR_SERVER_EXT_IFACE="wlan0"
export NETLSR_SERVER_DEBUG="true"

# Run the application
./netlsr -mode client
```

### 3. Priority Order

Configuration values are loaded in this order (later values override earlier ones):

1. Default values (hardcoded in the application)
2. Configuration file values
3. Environment variable values

## Generating Default Configuration

To generate a default configuration file:

```bash
go run cmd/genconfig/main.go ./configs
```

This will create a `config.yaml` file with default values that you can customize.

## Configuration Options

### Client Configuration

| Field | Default | Description |
|-------|---------|-------------|
| `tun_name` | `tun0` | Name of the TUN interface |
| `local_ip` | `10.0.0.1` | Local IP address for TUN interface |
| `peer_ip` | `10.0.0.2` | Server's TUN IP address |
| `server_addr` | `""` | Server address to connect to |
| `port` | `8080` | Server port |

### Server Configuration

| Field | Default | Description |
|-------|---------|-------------|
| `tun_name` | `tun0` | Name of the TUN interface |
| `local_ip` | `10.0.0.1` | Local IP address for TUN interface |
| `peer_ip` | `10.0.0.2` | Client's TUN IP address |
| `port` | `8080` | Port to listen on |
| `ext_iface` | `""` | External interface for NAT |
| `debug` | `false` | Enable debug mode |

## Example Configuration

```yaml
client:
  tun_name: "tun0"
  local_ip: "10.0.0.1"
  peer_ip: "10.0.0.2"
  server_addr: "192.168.1.100"
  port: 8080

server:
  tun_name: "tun0"
  local_ip: "10.0.0.1"
  peer_ip: "10.0.0.2"
  port: 8080
  ext_iface: "eth0"
  debug: false
``` 