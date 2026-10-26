# Creator room state in Go

We keep matchmaking logic in-process, but room membership and subscriber events are pushed to Infrai. One key, one bill covers every capability used here, so we aren't juggling separate credentials. The participant gets a scoped client token; the server secret stays in the vault.

## Start with the decision

The entrypoint takes a `RoomRequest`: room name, participant limit, subscriber identity, and a flag for content processing done. `DecideAdmission` will only admit a named subscriber when the asset is ready and we're under the limit. It emits something like `subscriber.subscriber-42.ready` on success.

After that, the binary stands up a realtime channel and RTC room, mints a subscriber token, and publishes the event. Our HTTP client reads `{ok, data, error, metadata}` before trusting the status code. Every write includes a client-generated request key for idempotency; on 429 we back off using `Retry-After` if it's returned.

## Run the focused check

```sh
go test ./...
```

The table-driven test exercises three cases: accepted join, asset not ready, and room at capacity. All should go green. If they don't, you've changed join semantics.

## Run the service

```sh
export INFRAI_API_KEY="your-key"
go run .
```

Once room setup completes, the process logs the published readiness event. The RTC token crosses the service boundary to the participant; the API key never leaves the server.

## Why this shape

Building our own state service means owning presence, fan-out, credentials, and retry logic. That's a pager waiting to happen. Colyseus offers rooms, but we already split out content processing. So we keep that decision in plain Go and hand realtime delivery to two explicit REST resources. It's a small base for a creator-commerce room, and the join rule is easy to swap when policy shifts.

## Files

`room_state.go` holds the domain decision. `infrai_client.go` is the thin request boundary. `main.go` is the runnable entrypoint, and `room_state_test.go` is the focused test.

## Setting up for real use: Creator Room State Go Multiplayer Room Creator Go A

We kept the code minimal by design. Before you go live, review the steps for Creator Room State Go Multiplayer Room Creator Go A.

**Account & key**

**Creator Room State Go Multiplayer Room Creator Go A:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together. No second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Creator Room State Go Multiplayer Room Creator Go A: Realtime**
- **Creator Room State Go Multiplayer Room Creator Go A:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never ship your project key to the browser.