# netlsr

netlsr is a simple tunneler between two VMs using a TUN interface over UDP. It forwards traffic from one VM (Client) to another (Server), routing traffic destined for a specified network CIDR.

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
  -ifname tun0 \
  -local-ip 10.99.0.2/16 \
  -peer-ip 10.99.0.1 \
  -port 5000
```

Enable IP forwarding and configure NAT:

```bash
sudo sysctl -w net.ipv4.ip_forward=1
sudo iptables -t nat -A POSTROUTING -s 10.100.0.0/16 -o eth0 -j MASQUERADE
sudo ip route add 10.100.0.0/16 dev tun0
```

### Client (VM A)

```bash
sudo ./netlsr \
  -mode client \
  -remote 192.168.2.128 \
  -ifname tun0 \
  -local-ip 10.99.0.1/16 \
  -peer-ip 10.99.0.2 \
  -port 5000
```

Configure route for the tunneled network:

```bash
sudo ip route add 10.100.0.0/16 dev tun0
```

## Flags

- `-mode`            Mode to run: `client` or `server` (default `client`)
- `-remote`          Server address (only for client mode)
- `-ifname`          Name of the TUN interface (default `tun0`)
- `-local-ip`        Local TUN IP (e.g., `10.100.0.1/16`)
- `-peer-ip`         Peer TUN IP (e.g., `10.100.0.2`)
- `-tun-cidr`        CIDR for the tunneled network (default `10.100.0.0/16`)
- `-port`            UDP port for the tunnel (default `5000`) 