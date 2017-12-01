package fabsdk

import (
	"diviner/common/cast"
	"time"

	"github.com/hyperledger/fabric-sdk-go/api/apiconfig"
	fab "github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
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

// QueryFabric ...
func QueryFabric(client apitxn.ChannelClient, chaincode, fcn string, data ...[]byte) ([]byte, error) {
	qr := apitxn.QueryRequest{
		ChaincodeID: chaincode,
		Fcn:         fcn,
		Args:        data,
	}

	return client.Query(qr)
}

// QueryFabricByID ...
func QueryFabricByID(client apitxn.ChannelClient, chaincode, fcn, id string) ([]byte, error) {
	return QueryFabric(client, chaincode, fcn, []byte(id))
}

// ExecuteFabric ...
func ExecuteFabric(client apitxn.ChannelClient, chaincode, fcn string, data ...[]byte) (apitxn.TransactionID, error) {
	txr := apitxn.ExecuteTxRequest{
		ChaincodeID: chaincode,
		Fcn:         fcn,
		Args:        data,
	}

	return client.ExecuteTx(txr)
}

// ExecuteFabricWithStrings ...
func ExecuteFabricWithStrings(client apitxn.ChannelClient, chaincode, fcn string, data ...string) (apitxn.TransactionID, error) {
	tmp := cast.StringsToByteArray(data...)
	return ExecuteFabric(client, chaincode, fcn, tmp...)
}

// RegisterChaincodeEvent ...
func RegisterChaincodeEvent(client apitxn.ChannelClient, chaincode, name, regex string) (chan *apitxn.CCEvent, apitxn.Registration) {
	id := name + regex

	notifier := make(chan *apitxn.CCEvent)
	rce := client.RegisterChaincodeEvent(notifier, chaincode, id)
	return notifier, rce
}

// RegisterChaincodeEventWithDefaultRegex ...
func RegisterChaincodeEventWithDefaultRegex(client apitxn.ChannelClient, chaincode, name string) (chan *apitxn.CCEvent, apitxn.Registration) {
	return RegisterChaincodeEvent(client, chaincode, name, "([a-zA-Z0-9]+)")
}

// SelectEvent ...
func SelectEvent(notifier chan *apitxn.CCEvent, txid string, timeout time.Duration) []byte {
	for {
		select {
		case evt := <-notifier:
			if evt.TxID == txid {
				return evt.Payload
			}
		case <-time.After(time.Second * timeout):
			return nil
		}
	}
}
