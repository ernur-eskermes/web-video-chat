package room

import (
	"time"

	"github.com/livekit/protocol/auth"
	lksdk "github.com/livekit/server-sdk-go"
)

type authBase struct {
	apiKey    string
	apiSecret string
}

type Room struct {
	*lksdk.RoomServiceClient
	authBase
}

func NewRoom(host, apiKey, secretKey string) Room {
	return Room{
		RoomServiceClient: lksdk.NewRoomServiceClient(host, apiKey, secretKey),
		authBase: authBase{
			apiKey:    apiKey,
			apiSecret: secretKey,
		},
	}
}

func (s Room) GetJoinToken(room, identity string) (string, error) {
	at := auth.NewAccessToken(s.apiKey, s.apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     room,
	}
	at.AddGrant(grant).
		SetIdentity(identity).
		SetValidFor(time.Hour)

	return at.ToJWT()
}
