package main

import "fmt"

type RoomRequest struct {
	RoomName        string
	MaxParticipants int
	SubscriberID    string
	AssetReady      bool
}

type RoomDecision struct {
	Accepted bool
	Reason   string
	Event    string
}

// DecideAdmission keeps the business rule independent from the transport.
func DecideAdmission(req RoomRequest, current int) RoomDecision {
	if req.SubscriberID == "" {
		return RoomDecision{Reason: "subscriber identity is required"}
	}
	if !req.AssetReady {
		return RoomDecision{Reason: "content is still processing"}
	}
	if current >= req.MaxParticipants {
		return RoomDecision{Reason: "room is full"}
	}
	return RoomDecision{
		Accepted: true,
		Reason:   "subscriber admitted",
		Event:    fmt.Sprintf("subscriber.%s.ready", req.SubscriberID),
	}
}
