# netlsr (netleisure)

> "Networking shouldn't be WORK, it should be LEISURE. That's why we have netleisure - making network tunneling as relaxing as a day at the beach! 🏖️"

A simple and elegant tunneler that creates a virtual network bridge between VMs using TUN interfaces over UDP. Perfect for when you want to connect your virtual machines without the complexity of traditional VPNs.

## Architecture

```
[VM A] <---> [TUN Interface] <---> [UDP Tunnel] <---> [TUN Interface] <---> [VM B]
(Client)      (10.99.0.1/16)        (Port 5000)        (10.99.0.2/16)      (Server)
```

## Features

- 🚀 Simple point-to-point tunneling
- 🔒 TUN interface for layer 3 packet forwarding
- 🌐 UDP-based transport for maximum compatibility
- 🛠️ Cross-platform support (Linux, macOS)
- ⚡ Low latency and high performance

## Prerequisites

- Go installed (>=1.16)
- Linux with iproute2 and iptables
- Root privileges (required to create TUN interfaces and modify routing/iptables)

## Installation

```bash
git clone https://github.com/mnlx/netlsr.git
cd netlsr
go build -o netlsr
```

## Usage

### Server (VM B)

```bash
sudo ./netlsr \
  -mode server \
  -local-ip 10.99.0.2/16 
```


### Client (VM A)

```bash
sudo ./netlsr \
  -mode client \
  -remote <server-address> \
  -local-ip 10.99.0.1/16 \
  -port 5000
```

Configure route for the tunneled network:

```bash
sudo ip route add 10.100.0.0/16 dev tun0
```

## Network Flow

```
[Application] → [TUN Interface] → [netlsr] → [UDP Socket] → [Internet] → [UDP Socket] → [netlsr] → [TUN Interface] → [Application]
```

## Configuration Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-mode` | Mode to run: `client` or `server` | `client` |
| `-remote` | Server address (client mode only) | - |
| `-ifname` | TUN interface name | `tun0` |
| `-local-ip` | Local TUN IP (e.g., `10.100.0.1/16`) | - |
| `-peer-ip` | Peer TUN IP (e.g., `10.100.0.2`) | - |
| `-port` | UDP port for the tunnel | `5000` |

## Troubleshooting

1. **Interface Creation Fails**
   - Ensure you have root privileges
   - Check if the interface name is available

2. **Connection Issues**
   - Verify firewall rules allow UDP traffic
   - Check network connectivity between VMs
   - Ensure correct IP addressing

## Contributing

Feel free to submit issues and enhancement requests!

## License

MIT License - feel free to use this for your networking leisure time! 🎉 