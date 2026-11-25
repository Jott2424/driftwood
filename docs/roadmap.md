# Roadmap

Each milestone is one session and ends with something I can run or see. Some
sessions will be longer than others. The order is chosen so every service is
built and proven before it gets its own VM.

## Phase 0: Go basics and local setup

- [x] 1. Install Go and run a first program
- [x] 2. A small web server with a `/health` endpoint, and a first unit test
- [x] 3. Postgres in Docker Compose, with Go connecting through `pgx`
- [x] 4. A first database table, created with `goose` migrations
- [ ] 5. CI on GitHub that runs `go vet`, the tests, and the build, then require it in the `main` ruleset

## Phase 1: Proxmox base

- [ ] 6. First VM, turned into a cloud-init template (see [infrastructure.md](infrastructure.md))
- [ ] 7. Terraform the VMs: clone the template into service VMs from a map of VMs, reusing the pattern from `selfhosted/proxmox-vms`
- [ ] 8. Postgres VM at `10.10.0.11`, created by Terraform and reachable from my laptop

## Phase 2: auth service

- [ ] 9. Register endpoint with argon2id, a repository layer, and an integration test
- [ ] 10. Login that returns a JWT
- [ ] 11. Refresh tokens, logout, and revocation
- [ ] 12. Profile endpoints, reusable token checking, and the auth OpenAPI spec with Swagger UI at `/docs`
- [ ] 13. Deploy auth to its own VM

## Phase 3: chat service

- [ ] 14. Community boards and membership
- [ ] 15. Sending messages and paginated history
- [ ] 16. Private groups with roles
- [ ] 17. 1:1 chats
- [ ] 18. Deploy chat to its own VM

Each chat milestone also adds its endpoints to the chat OpenAPI spec.

## Phase 4: gateway

- [ ] 19. The gateway forwards `/auth` requests to the auth service
- [ ] 20. The gateway checks tokens and forwards chat routes
- [ ] 21. A combined Swagger page at `/docs`
- [ ] 22. Deploy the gateway at `10.10.0.14`

## Phase 5: NATS and the realtime service

- [ ] 23. NATS VM, and chat publishes a "message sent" event
- [ ] 24. WebSocket basics, starting with an echo server
- [ ] 25. Authenticated connections that deliver new messages from NATS
- [ ] 26. Presence and typing indicators
- [ ] 27. Unread counts and reconnect with catch-up
- [ ] 28. Deploy realtime, with the WebSocket connection going through the gateway, plus AsyncAPI docs for the messages and events

## Phase 6: React website

- [ ] 29. Login and registration
- [ ] 30. Board list and message history
- [ ] 31. Live chat
- [ ] 32. Deploy web, calling only the gateway

## Phase 7: scale and hardening

- [ ] 33. Go load generator with bot accounts, calling only the gateway
- [ ] 34. Several realtime instances sharing messages through NATS
- [ ] 35. Vault for credentials, with re-fetch instead of load-once
- [ ] 36. Vault as a certificate authority, with TLS between services
- [ ] 37. Rate limiting, structured logging, and Prometheus and Grafana
- [ ] 38. SeaweedFS VM, with avatar upload in auth
- [ ] 39. Attachments in chat, using SeaweedFS
- [ ] 40. Later extras: Mailpit email verification and TOTP

## Open choices to revisit

- NATS is the starting choice for the broker. Look at it again at milestone 23.
- Whether the realtime service's WebSocket traffic goes through the gateway or
  straight to the realtime service. Decide at milestone 28.
- Milestone 7 (Terraform), decide when we get there:
  - Whether the Proxmox host is the same one used in `selfhosted`, and whether
    to create a separate API token for driftwood (I'd suggest yes).
  - Whether to reuse the Python configuration scripts from `selfhosted`, or
    start with plain setup commands and bring in Ansible later.
  - Keep `terraform.tfvars` and the state file out of git. Vault can hold the
    API token later.
  - The provider can hang for up to 15 minutes waiting for the guest agent, so
    the template needs the agent installed.
