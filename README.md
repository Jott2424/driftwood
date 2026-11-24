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

## Ground rules

- **Not a batch project.** If it can be a cron job, it's probably not the point.
- **Learning over polish.** Trying things, breaking them, and writing down what I learned matters more than a finished product.
- **Language is open.** Python is familiar, so I'm free (and encouraged) to reach for something else where it teaches more.

## Roadmap

- [ ] Pick a concrete app idea to build around
- [ ] A basic HTTP API with a database behind it
- [ ] Add a realtime channel (WebSockets / SSE)
- [ ] Add a message broker or pub/sub layer
- [ ] Auth, rate limiting, and error handling
- [ ] Logging, metrics, and health checks
- [ ] Containerize and deploy

## Notes

Learnings and decisions will be captured as the project grows.
