# Infrastructure checklist

Running Driftwood as VMs on a single Proxmox node, on an isolated internal
network so test traffic never touches the home router.

## Design

```
Internet ── Router ── vmbr0 (LAN, physical NIC)
                          │
                    Proxmox host (NAT gateway)
                          │
                      vmbr1 (internal only, 10.10.0.0/24)
                          │
   gateway · auth · chat · realtime · web · postgres · nats ·
          vault · seaweedfs · loadgen
```

- `vmbr0`: the default bridge, attached to the physical NIC and the LAN.
- `vmbr1`: internal bridge with no physical NIC. The host is `10.10.0.1` and acts
  as the gateway and NAT for the VMs.
- VM-to-VM traffic on `vmbr1` stays inside the host. Only traffic to the
  internet or from my own machine crosses the router.

## Network setup (Proxmox host)

- [x] Create `vmbr1` with address `10.10.0.1/24`, no bridge ports, no gateway
- [x] Add `post-up` / `post-down` lines to `/etc/network/interfaces` for IP forwarding and MASQUERADE
- [x] Apply with `ifreload -a`
- [x] Verify `ip addr show vmbr1` shows `10.10.0.1/24`
- [x] Verify `cat /proc/sys/net/ipv4/ip_forward` prints `1`
- [x] Verify `iptables -t nat -L POSTROUTING -n -v` lists MASQUERADE for `10.10.0.0/24`

`vmbr1` block in `/etc/network/interfaces`:

```
auto vmbr1
iface vmbr1 inet static
        address 10.10.0.1/24
        bridge-ports none
        bridge-stp off
        bridge-fd 0

        post-up   echo 1 > /proc/sys/net/ipv4/ip_forward
        post-up   iptables -t nat -A POSTROUTING -s 10.10.0.0/24 -o vmbr0 -j MASQUERADE
        post-down iptables -t nat -D POSTROUTING -s 10.10.0.0/24 -o vmbr0 -j MASQUERADE
```

## First VM and template

- [ ] Create a test VM (Debian or Ubuntu server) with its NIC on `vmbr1`
- [ ] Set a temporary static IP `10.10.0.200/24`, gateway `10.10.0.1`, DNS `1.1.1.1` (the template keeps this range so it never clashes with a real service address)
- [ ] `ping 10.10.0.1` works (host)
- [ ] `ping 1.1.1.1` works (NAT to the internet)
- [ ] `ping google.com` works (DNS)
- [ ] Install Docker and add my SSH key
- [ ] Convert it to a cloud-init template so new VMs are quick to clone
- [ ] Clone a second VM and confirm VM-to-VM connectivity

## Address plan

Static addresses, kept here as the source of truth.

| Address | Name | Purpose | Status |
|---|---|---|---|
| 10.10.0.1 | proxmox | host, gateway, NAT | done |
| 10.10.0.10 | auth | accounts, login, tokens | planned |
| 10.10.0.11 | postgres | database, one database per service | planned |
| 10.10.0.12 | broker | NATS | planned |
| 10.10.0.13 | loadgen | bot accounts and stress tests | planned |
| 10.10.0.14 | gateway | API gateway, the single entry point | planned |
| 10.10.0.15 | vault | credentials and certificate management | planned |
| 10.10.0.16 | chat | boards, groups, DMs, message history | planned |
| 10.10.0.17 | realtime | WebSocket connections, presence, typing | planned |
| 10.10.0.18 | web | React client | planned |
| 10.10.0.19 | seaweedfs | object storage for avatars and attachments | planned |
| 10.10.0.200 | template | temporary address for the cloud-init template | planned |

## Reaching the VMs

- [ ] Gateway VM (`10.10.0.14`) as the single web entry point, with the host forwarding 80/443 to it
- [ ] Admin access (SSH, Grafana, Postgres) through a static route on my laptop for `10.10.0.0/24` via the Proxmox host, or Tailscale
- [ ] Do not forward database or broker ports to the LAN
- [ ] Name resolution: `/etc/hosts` entries on each VM, a small DNS server later

## Gotchas to remember

- Services must listen on their `10.10.0.x` address, not just `localhost` (for
  Postgres, set `listen_addresses` and allow the API's address in `pg_hba.conf`).
- VM firewalls (ufw) may need to allow `10.10.0.0/24`.
- Rules added by hand at the shell disappear on reboot, so keep them in
  `/etc/network/interfaces`.
- If the main bridge is not `vmbr0`, change `-o vmbr0` in the NAT lines.
- The load generator should call the API's internal address directly so test
  traffic stays on `vmbr1`.
