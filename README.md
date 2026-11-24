# driftwood

A learning project for building **realtime services and APIs** — the kind of
long-running, event-driven software that sits at the other end of the spectrum
from batch data processing.

## Why

My background is batch workloads and Python: run a job, process a pile of data,
exit. Driftwood is a deliberate step away from that. The goal is to get hands-on
with how typical production software is built and operated:

- Long-lived services instead of scripts that run to completion
- Request/response APIs (REST) and design around contracts, versioning, and errors
- Realtime communication: WebSockets, server-sent events, pub/sub
- State, persistence, and concurrency under live traffic
- Observability, testing, deployment, and the rest of the engineering lifecycle

## The app

Driftwood is a **realtime chat platform**:

- **Accounts:** register, log in, and a profile page
- **Community boards:** public, topic-based channels with message history
- **Groups:** private rooms with members and roles
- **1:1 chats:** private direct messages
- Presence, typing indicators, unread counts, and reconnect with catch-up

It is also built to be **stress tested**: a load generator creates bot accounts
that register, join rooms, and hold scripted conversations, so I can join in and
watch the system under load.

## Planned stack

- **Go** for the API and WebSocket server
- **Postgres** for users, rooms, memberships, and messages
- **JWT** auth to start (Keycloak is a possible later swap)
- **Docker Compose** for the local environment
- **Proxmox VMs** on an isolated internal network, so stress-test traffic stays
  off the home router
- Later: Redis or NATS for fan-out between server instances, MinIO for avatars
  and attachments, a message broker, Mailpit for email verification, TOTP, and
  Prometheus + Grafana for metrics

## Layout

```
api/            Go service: REST API and WebSocket server
web/            frontend client
deploy/         Docker Compose and shared infrastructure config
docs/           decisions and learnings
tools/loadgen/  bot accounts and scripted conversations for stress testing
```

## Ground rules

- **Not a batch project.** If it can be a cron job, it's probably not the point.
- **Learning over polish.** Trying things, breaking them, and writing down what I learned matters more than a finished product.
- **Go is the main language.** Python is familiar, so Go is the deliberate stretch.
- **Use open source where it fits.** Configure existing services instead of rebuilding them.

## Roadmap

- [x] Pick a concrete app idea: a realtime chat platform
- [x] Proxmox internal network (`vmbr1`) with NAT; see [docs/infrastructure.md](docs/infrastructure.md)
- [ ] First VM and a cloud-init template for the rest
- [ ] Local environment: Docker Compose and Postgres
- [ ] API skeleton: health check, database connection, migrations
- [ ] Registration, login, and profile
- [ ] WebSocket messaging on community boards, with history
- [ ] Groups and membership permissions
- [ ] 1:1 chats
- [ ] Presence, typing indicators, and unread counts
- [ ] Multiple server instances with pub/sub fan-out
- [ ] Load generator: bot accounts and scripted conversations
- [ ] Avatars and attachments in object storage
- [ ] Email verification and optional TOTP
- [ ] Rate limiting, error handling, and tests
- [ ] Logging, metrics, and health checks
- [ ] Containerize and deploy

## Notes

Learnings and decisions will be captured in [docs/](docs/) as the project grows.
