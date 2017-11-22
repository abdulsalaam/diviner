package fabsdk

import (
	"math/rand"
	"time"

	fab "github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	"github.com/pkg/errors"
)

// RandomString ...
func RandomString(strlen int) string {
	rand.Seed(time.Now().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	charlen := len(chars)
	ret := make([]byte, strlen)

	for i := 0; i < strlen; i++ {
		ret[i] = chars[rand.Intn(charlen)]
	}

	return string(ret)
}

// GenerateRandomID ...
func GenerateRandomID() string {
	return RandomString(10)
}

// HasPrimaryPeerJoinedChannel checks whether the primary peer of a channel
// has already joined the channel. It returns true if it has, false otherwise,
// or an error
func HasPrimaryPeerJoinedChannel(client fab.FabricClient, channel fab.Channel) (bool, error) {
	foundChannel := false
	primaryPeer := channel.PrimaryPeer()
	response, err := client.QueryChannels(primaryPeer)
	if err != nil {
		return false, errors.WithMessage(err, "failed to query channel for primary peer")
	}
	for _, responseChannel := range response.Channels {
		if responseChannel.ChannelId == channel.Name() {
			foundChannel = true
		}
	}

	return foundChannel, nil
}
