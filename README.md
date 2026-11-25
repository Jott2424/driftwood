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
- Realtime communication: WebSockets, pub/sub
- Several services that work together, the way production systems are run
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

## Architecture

Driftwood is a set of separate Go services, not one program. Each service is
built and proven on its own, then deployed to its own Proxmox VM.

```
browser / loadgen
        │
     gateway ── combined Swagger docs at /docs
        │
   ┌────┼──────────┐
 auth  chat     realtime  ←→  NATS
   │    │            │
   └────┴── Postgres (one database per service)
```

| Service | Job |
|---|---|
| **gateway** | The single entry point. Checks tokens, forwards requests, serves combined API docs. |
| **auth** | Registration, login, profiles. Issues short-lived JWTs and refresh tokens. |
| **chat** | Boards, groups, 1:1 chats, message history. |
| **realtime** | WebSocket connections, presence, typing, live delivery. |
| **web** | React client. Talks only to the gateway. |

Design decisions:

- **Clients call only the gateway.** Services can still talk to each other
  directly behind it. The gateway is not a hub for service-to-service calls.
- **Each service owns its data.** One Postgres server, one database per service,
  and no service reads another's tables.
- **Each service wraps its own database** with a small repository layer (for
  example `CreateUser`), so no raw SQL is spread through the code.
- **Auth** uses a short-lived JWT plus a server-stored refresh token that can be
  revoked. Passwords are stored as argon2id hashes.
- **Docs:** OpenAPI with Swagger UI for the REST APIs, written spec-first.
  AsyncAPI for the WebSocket messages and NATS events.
- **Tests:** unit and integration tests (against a real Postgres in Docker) from
  the first endpoint, plus CI on GitHub.

## Stack

Built from scratch, because this is what I want to learn:

- **Go** services: gateway, auth, chat, realtime (standard library `net/http`,
  `pgx` with hand-written SQL, `goose` migrations)
- **React** web client (TypeScript, Vite)
- **Go** load generator

Existing open source, configured and not rebuilt:

- **Postgres** for users, rooms, memberships, and messages
- **NATS** for events between services
- **HashiCorp Vault** for credentials and certificate management, with services
  re-fetching values so one update in Vault reaches all later requests
- **SeaweedFS** for object storage (avatars and attachments), through its S3 API
- **Docker Compose** for the local environment
- **Proxmox VMs** on an isolated internal network, so stress-test traffic stays
  off the home router, created with **Terraform**
- Later: Mailpit for email verification, TOTP, Prometheus and Grafana

## Layout

```
services/gateway/   API gateway
services/auth/      accounts, login, tokens
services/chat/      boards, groups, DMs, message history
services/realtime/  WebSocket connections, presence, typing
web/                React client
deploy/             Docker Compose and shared infrastructure config
docs/               decisions, roadmap, and learnings
tools/loadgen/      Go bot accounts for stress testing
```

## Ground rules

- **Build the app, use tools for the rest.** The chat application and its
  services are written from scratch. Supporting infrastructure is existing open
  source that I configure, like Postgres, NATS, Vault, and SeaweedFS.
- **Not a batch project.** If it can be a cron job, it's probably not the point.
- **Learning over polish.** Trying things, breaking them, and writing down what I learned matters more than a finished product.
- **Go is the main language.** Python is familiar, so Go is the deliberate stretch.
- **One service at a time.** Build it, prove it works, then give it its own VM.
- **Every session ends with something concrete**, however long the session is.

## Roadmap

The full list of session-sized milestones is in [docs/roadmap.md](docs/roadmap.md).

- [x] Pick a concrete app idea: a realtime chat platform
- [x] Proxmox internal network (`vmbr1`) with NAT; see [docs/infrastructure.md](docs/infrastructure.md)
- [ ] **Phase 0:** Go basics and local setup (Go, a small server, Postgres, migrations, CI)
- [ ] **Phase 1:** Proxmox base (VM template, Terraform for the VMs, Postgres VM)
- [ ] **Phase 2:** auth service
- [ ] **Phase 3:** chat service
- [ ] **Phase 4:** gateway
- [ ] **Phase 5:** NATS and the realtime service
- [ ] **Phase 6:** React website
- [ ] **Phase 7:** scale and hardening (load generator, Vault, metrics, SeaweedFS, extras)

## Notes

Learnings and decisions will be captured in [docs/](docs/) as the project grows.
