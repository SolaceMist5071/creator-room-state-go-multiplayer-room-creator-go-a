package main

import "testing"

func TestDecideAdmissionWaitsForProcessedContent(t *testing.T) {
	tests := []struct {
		name     string
		request  RoomRequest
		current  int
		accepted bool
		reason   string
	}{
		{
			name:    "ready subscriber gets a room event",
			request: RoomRequest{RoomName: "studio", MaxParticipants: 4, SubscriberID: "sub-7", AssetReady: true},
			current: 2, accepted: true, reason: "subscriber admitted",
		},
		{
			name:    "processing asset keeps subscriber out",
			request: RoomRequest{RoomName: "studio", MaxParticipants: 4, SubscriberID: "sub-7", AssetReady: false},
			current: 2, accepted: false, reason: "content is still processing",
		},
		{
			name:    "full room keeps subscriber out",
			request: RoomRequest{RoomName: "studio", MaxParticipants: 2, SubscriberID: "sub-7", AssetReady: true},
			current: 2, accepted: false, reason: "room is full",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DecideAdmission(tt.request, tt.current)
			if got.Accepted != tt.accepted || got.Reason != tt.reason {
				t.Fatalf("got %#v, want accepted=%v reason=%q", got, tt.accepted, tt.reason)
			}
		})
	}
}
