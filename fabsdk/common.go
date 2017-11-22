package fabsdk

import (
	"github.com/hyperledger/fabric-sdk-go/api/apiconfig"
	fab "github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	"github.com/hyperledger/fabric-sdk-go/def/fabapi"
	"github.com/hyperledger/fabric-sdk-go/pkg/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/errors"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/events"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/orderer"
)

// InitAPIConfig ...
func InitAPIConfig(file string) (apiconfig.Config, error) {
	return config.InitConfig(file)
}

// NewSDK ...
func NewSDK(file string) (*fabapi.FabricSDK, error) {
	sdkOptions := fabapi.Options{ConfigFile: file}

	return fabapi.NewSDK(sdkOptions)
}

// GetChannel ...
func GetChannel(client fab.FabricClient, channelID string, orgs []string) (fab.Channel, error) {
	channel, err := client.NewChannel(channelID)
	if err != nil {
		return nil, errors.WithMessage(err, "new channel failed")
	}

	ordererConfig, err := client.Config().RandomOrdererConfig()
	if err != nil {
		return nil, errors.WithMessage(err, "random orderer config failed")
	}

	serverHostOverride := ""
	if str, ok := ordererConfig.GRPCOptions["ssl-target-name-override"].(string); ok {
		serverHostOverride = str
	}

	orderer, err := orderer.NewOrderer(ordererConfig.URL, ordererConfig.TLSCACerts.Path, serverHostOverride, client.Config())
	if err != nil {
		return nil, errors.WithMessage(err, "new orderer failed")
	}

	err = channel.AddOrderer(orderer)
	if err != nil {
		return nil, errors.WithMessage(err, "adding orderer failed")
	}

	for _, org := range orgs {
		peerConfig, err := client.Config().PeersConfig(org)
		if err != nil {
			return nil, errors.WithMessage(err, "reading peer config failed")
		}

		for _, p := range peerConfig {
			serverHostOverride = ""
			if str, ok := p.GRPCOptions["ssl-target-name-override"].(string); ok {
				serverHostOverride = str
			}
			endorser, err := fabapi.NewPeer(p.URL, p.TLSCACerts.Path, serverHostOverride, client.Config())
			if err != nil {
				return nil, errors.WithMessage(err, "new peer failed")
			}

			err = channel.AddPeer(endorser)
			if err != nil {
				return nil, errors.WithMessage(err, "adding peer failed")
			}
		}
	}

	return channel, nil
}

// GetEventHub ...
func GetEventHub(client fab.FabricClient, org string) (fab.EventHub, error) {
	eventHub, err := events.NewEventHub(client)
	if err != nil {
		return nil, errors.WithMessage(err, "new event hub failed")
	}

	found := false

	peers, err := client.Config().PeersConfig(org)
	if err != nil {
		return nil, errors.WithMessage(err, "get peer configs failed")
	}

	for _, p := range peers {
		if p.URL != "" {
			serverHostOverride := ""
			if str, ok := p.GRPCOptions["ssl-target-name-override"].(string); ok {
				serverHostOverride = str
			}
			eventHub.SetPeerAddr(p.EventURL, p.TLSCACerts.Path, serverHostOverride)
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("event hub configuration not found")
	}

	return eventHub, nil
}
