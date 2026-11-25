# Roadmap

Each milestone is one session and ends with something I can run or see. Some
sessions will be longer than others. The order is chosen so every service is
built and proven before it gets its own VM.

## Phase 0: Go basics and local setup

- [ ] 1. Install Go and run a first program
- [ ] 2. A small web server with a `/health` endpoint, and a first unit test
- [ ] 3. Postgres in Docker Compose, with Go connecting through `pgx`
- [ ] 4. A first database table, created with `goose` migrations
- [ ] 5. CI on GitHub that runs `go vet`, the tests, and the build

## Phase 1: Proxmox base

- [ ] 6. First VM, turned into a cloud-init template (see [infrastructure.md](infrastructure.md))
- [ ] 7. Postgres VM at `10.10.0.11`, reachable from my laptop

## Phase 2: auth service

- [ ] 8. Register endpoint with argon2id, a repository layer, and an integration test
- [ ] 9. Login that returns a JWT
- [ ] 10. Refresh tokens, logout, and revocation
- [ ] 11. Profile endpoints, reusable token checking, and the auth OpenAPI spec with Swagger UI at `/docs`
- [ ] 12. Deploy auth to its own VM

## Phase 3: chat service

- [ ] 13. Community boards and membership
- [ ] 14. Sending messages and paginated history
- [ ] 15. Private groups with roles
- [ ] 16. 1:1 chats
- [ ] 17. Deploy chat to its own VM

Each chat milestone also adds its endpoints to the chat OpenAPI spec.

## Phase 4: gateway

- [ ] 18. The gateway forwards `/auth` requests to the auth service
- [ ] 19. The gateway checks tokens and forwards chat routes
- [ ] 20. A combined Swagger page at `/docs`
- [ ] 21. Deploy the gateway at `10.10.0.14`

## Phase 5: NATS and the realtime service

- [ ] 22. NATS VM, and chat publishes a "message sent" event
- [ ] 23. WebSocket basics, starting with an echo server
- [ ] 24. Authenticated connections that deliver new messages from NATS
- [ ] 25. Presence and typing indicators
- [ ] 26. Unread counts and reconnect with catch-up
- [ ] 27. Deploy realtime, with the WebSocket connection going through the gateway, plus AsyncAPI docs for the messages and events

## Phase 6: React website

- [ ] 28. Login and registration
- [ ] 29. Board list and message history
- [ ] 30. Live chat
- [ ] 31. Deploy web, calling only the gateway

## Phase 7: scale and hardening

- [ ] 32. Go load generator with bot accounts, calling only the gateway
- [ ] 33. Several realtime instances sharing messages through NATS
- [ ] 34. Vault for credentials, with re-fetch instead of load-once
- [ ] 35. Vault as a certificate authority, with TLS between services
- [ ] 36. Rate limiting, structured logging, and Prometheus and Grafana
- [ ] 37. SeaweedFS VM, with avatar upload in auth
- [ ] 38. Attachments in chat, using SeaweedFS
- [ ] 39. Later extras: Mailpit email verification and TOTP

## Open choices to revisit

- NATS is the starting choice for the broker. Look at it again at milestone 22.
- Whether the realtime service's WebSocket traffic goes through the gateway or
  straight to the realtime service. Decide at milestone 27.
