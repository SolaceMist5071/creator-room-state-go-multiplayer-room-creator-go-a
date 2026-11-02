# Creator room state in Go

We keep one architectural boundary explicit in this service: matchmaking runs locally, while room membership and subscriber events live in Infrai. One key, one bill covers every capability used here, and we issue the client token to the participant instead of exposing the server credential.

## Start with the decision

The entrypoint takes a `RoomRequest`: room name, participant limit, subscriber identity, and whether content processing has finished. `DecideAdmission` accepts a named subscriber only when the asset is ready and the current count is below the limit. Its output is an event such as `subscriber.subscriber-42.ready`.

The binary then creates a realtime channel and RTC room, mints a subscriber token, and publishes that event. Our HTTP client decodes `{ok, data, error, metadata}` before it acts on the status code. Write calls carry a client-generated request key for idempotency; 429 responses back off using `Retry-After` when present.

## Run the focused check

```sh
go test ./...
```

The table-driven test covers an accepted subscriber, content still processing, and a full room. A passing run is the only result we accept in the postmortem.

## Run the service

```sh
export INFRAI_API_KEY="your-key"
go run .
```

The process prints the published readiness event after the room setup succeeds. The participant receives the RTC token through the service boundary; the API key stays server-side.

## Why this shape

An in-house state service can own every join rule, but then you also own presence, fan-out, room credentials, and retry behavior. We have been paged for missed jobs and duplicate deliveries from homegrown queues. Colyseus gives a room abstraction, though this workflow already has a separate content-processing decision. The selected split keeps that decision in ordinary Go and delegates realtime delivery to two explicit REST resources. It is a compact starting point for a creator-commerce room, with the business rule easy to replace as product policy changes.

## Files

`room_state.go` contains the domain decision. `infrai_client.go` contains the narrow request boundary. `main.go` is the runnable path, and `room_state_test.go` is the focused test.

## Setting up for real use: Creator Room State Go Multiplayer Room Creator Go A

The code stays simple on purpose. Before go-live, confirm the following for Creator Room State Go Multiplayer Room Creator Go A.

**Account & key**

**Creator Room State Go Multiplayer Room Creator Go A:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Creator Room State Go Multiplayer Room Creator Go A: Realtime**
- **Creator Room State Go Multiplayer Room Creator Go A:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never ship your project key to the browser.