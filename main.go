package main

import (
	"fmt"
	"log"
)

func main() {
	client, err := NewInfraiClient()
	if err != nil {
		log.Fatal(err)
	}
	request := RoomRequest{RoomName: "creator-studio", MaxParticipants: 8, SubscriberID: "subscriber-42", AssetReady: true}
	decision := DecideAdmission(request, 0)
	if !decision.Accepted {
		log.Fatal(decision.Reason)
	}
	if err := client.CreateChannel(request.RoomName, "public", "tencent_im"); err != nil {
		log.Fatal(err)
	}
	if err := client.CreateRoom(request.RoomName, request.MaxParticipants, 300, "global"); err != nil {
		log.Fatal(err)
	}
	if err := client.IssueRoomToken(request.RoomName, request.SubscriberID, "Subscriber 42", false, true); err != nil {
		log.Fatal(err)
	}
	if err := client.Publish(request.RoomName, decision.Event, request.SubscriberID, map[string]any{"room": request.RoomName, "asset_ready": request.AssetReady}); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: %s\n", decision.Event, decision.Reason)
}
