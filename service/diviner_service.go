package main

import (
	"context"
	"diviner/common/base58"
	"diviner/common/cast"
	"diviner/common/config"
	"diviner/fabsdk"
	pbl "diviner/protos/lmsr"
	pbm "diviner/protos/member"
	pbs "diviner/protos/service"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/hyperledger/fabric-sdk-go/api/apifabclient"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"github.com/hyperledger/fabric-sdk-go/def/fabapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type divinerService struct {
	Client       apifabclient.FabricClient
	FabricConfig string
	SDK          *fabapi.FabricSDK
	ChannelID    string
	Channel      apifabclient.Channel
	Chaincode    string
	User         string
	Expired      int64
	Wait         time.Duration
	Balance      float64
}

func (s *divinerService) queryFabric(client apitxn.ChannelClient, module, fcn string, data ...[]byte) ([]byte, error) {
	var args [][]byte
	args = append(args, []byte(module))
	args = append(args, data...)

	return fabsdk.QueryFabric(client, s.Chaincode, fcn, args...)

	/*qr := apitxn.QueryRequest{
		ChaincodeID: s.Chaincode,
		Fcn:         fcn,
		Args:        args,
	}

	return client.Query(qr)*/
}

func (s *divinerService) queryFabricByID(client apitxn.ChannelClient, module, fcn, id string) ([]byte, error) {
	return s.queryFabric(client, module, fcn, []byte(id))
}

func (s *divinerService) executeFabric(client apitxn.ChannelClient, module, fcn string, data ...[]byte) (apitxn.TransactionID, error) {
	var args [][]byte
	args = append(args, []byte(module))
	args = append(args, data...)

	return fabsdk.ExecuteFabric(client, s.Chaincode, fcn, args...)

	/*txr := apitxn.ExecuteTxRequest{
		ChaincodeID: s.Chaincode,
		Fcn:         fcn,
		Args:        args,
	}

	return client.ExecuteTx(txr)*/
}

func (s *divinerService) registerChaincodeEvent(client apitxn.ChannelClient, chaincode, name string) (chan *apitxn.CCEvent, apitxn.Registration) {
	/*id := name + "([a-zA-Z0-9]+)"
	notifier := make(chan *apitxn.CCEvent)
	rce := client.RegisterChaincodeEvent(notifier, s.Chaincode, id)
	return notifier, rce*/
	return fabsdk.RegisterChaincodeEventWithDefaultRegex(client, s.Chaincode, name)
}

/*func (s *divinerService) selectEvent(notifier chan *apitxn.CCEvent, timeout time.Duration) []byte {
	select {
	case evt := <-notifier:
		log.Println("get notifier")
		return evt.Payload
	case <-time.After(time.Second * timeout):
		log.Println("timeout")
		return nil
	}
}*/

func (s *divinerService) returnMemberInfoResponse(bytes []byte) (*pbs.MemberInfoResponse, error) {
	member, err := pbm.Unmarshal(bytes)
	if err != nil {
		return nil, err
	}

	return &pbs.MemberInfoResponse{
		Member: member,
		Time:   ptypes.TimestampNow(),
	}, nil
}

func (s *divinerService) QueryMember(ctx context.Context, req *pbs.QueryRequest) (*pbs.MemberInfoResponse, error) {
	if ok, err := pbs.CheckQueryRequest(req, s.Expired); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("data is illegal")
	}

	if base58.Encode(req.Check.PublicKey) != req.Id {
		return nil, fmt.Errorf("address not match")
	}

	// call chaincode
	client, err := s.SDK.NewChannelClient(s.ChannelID, s.User)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	bytes, err := s.queryFabricByID(client, "member", "query", req.Id)
	if err != nil {
		return nil, err
	}
	return s.returnMemberInfoResponse(bytes)
}

func (s *divinerService) CreateMember(ctx context.Context, req *pbs.MemberCreateRequest) (*pbs.MemberInfoResponse, error) {
	if ok, err := pbs.CheckMemberCreateRequest(req, s.Expired); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("data is illegal")
	}

	member := &pbm.Member{
		Id:      req.Member.Id,
		Address: req.Member.Address,
		Balance: s.Balance,
	}

	bytes, err := pbm.Marshal(member)
	if err != nil {
		return nil, err
	}

	request := apitxn.ChaincodeInvokeRequest{
		Targets:      []apitxn.ProposalProcessor{s.Channel.PrimaryPeer()},
		Fcn:          "create",
		Args:         [][]byte{[]byte("member"), bytes},
		TransientMap: nil,
		ChaincodeID:  s.Chaincode,
	}

	transactionProposalResponses, _, err := s.Channel.SendTransactionProposal(request)
	if err != nil {
		return nil, err
	}

	for _, v := range transactionProposalResponses {
		if v.Err != nil {
			return nil, fmt.Errorf("endorser %s error: %v", v.Endorser, v.Err)
		}
	}

	tx, err := s.Channel.CreateTransaction(transactionProposalResponses)
	if err != nil {
		return nil, err
	}

	response, err := s.Channel.SendTransaction(tx)
	if err != nil {
		return nil, err
	}

	if response.Err != nil {
		return nil, fmt.Errorf("orderer %s error: %v", response.Orderer, response.Err)
	}

	client, err := s.SDK.NewChannelClient(s.ChannelID, s.User)
	if err != nil {
		return nil, err
	}
	bytes, err = s.queryFabricByID(client, "member", "query", member.Id)
	if err != nil {
		return nil, err
	}
	return s.returnMemberInfoResponse(bytes)

	/*client, err := s.SDK.NewChannelClient(s.ChannelID, s.User)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	notifier, rce := s.registerChaincodeEvent(client, s.Chaincode, "member")
	defer client.UnregisterChaincodeEvent(rce)

	tx, err := s.executeFabric(client, "member", "create", bytes)
	if err != nil {
		return nil, err
	}

	bytes = fabsdk.SelectEvent(notifier, tx.ID, 5)
	if bytes != nil {
		return s.returnMemberInfoResponse(bytes)
	}

	log.Println("time out")
	bytes, err = s.queryFabricByID(client, "member", "query", member.Id)
	if err != nil {
		return nil, err
	}

	return s.returnMemberInfoResponse(bytes)
	*/
}

func (s *divinerService) returnEventInfoResponse(bytes []byte) (*pbs.EventInfoResponse, error) {
	msg, err := pbl.UnmarshalEvent(bytes)
	if err != nil {
		return nil, err
	}

	return &pbs.EventInfoResponse{
		Event: msg,
		Time:  ptypes.TimestampNow(),
	}, nil
}

func (s *divinerService) QueryEvent(ctx context.Context, req *pbs.QueryRequest) (*pbs.EventInfoResponse, error) {
	if ok, err := pbs.CheckQueryRequest(req, s.Expired); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("data is illegal")
	}

	client, err := s.SDK.NewChannelClient(s.ChannelID, s.User)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	bytes, err := s.queryFabricByID(client, "event", "query", req.Id)
	if err != nil {
		return nil, err
	}
	return s.returnEventInfoResponse(bytes)
}

func (s *divinerService) CreateEvent(ctx context.Context, req *pbs.EventCreateRequest) (*pbs.EventInfoResponse, error) {
	if ok, err := pbs.CheckEventCreateRequest(req, s.Expired); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("data is illegal")
	}

	client, err := s.SDK.NewChannelClient(s.ChannelID, s.User)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	event, err := pbl.NewEvent(req.User, req.Title, req.Outcome...)
	if err != nil {
		return nil, err
	}

	data, err := pbl.MarshalEvent(event)
	if err != nil {
		return nil, err
	}

	notifier, rce := s.registerChaincodeEvent(client, s.Chaincode, "event")
	defer client.UnregisterChaincodeEvent(rce)

	tx, err := s.executeFabric(client, "event", "create", data)
	if err != nil {
		fmt.Printf("execute fabric error: %v\n", err)
		return nil, err
	}

	bytes := fabsdk.SelectEvent(notifier, tx.ID, s.Wait)
	if bytes != nil {
		return s.returnEventInfoResponse(bytes)
	}

	bytes, err = s.queryFabricByID(client, "event", "query", event.Id)
	if err != nil {
		return nil, err
	}
	return s.returnEventInfoResponse(bytes)

}

func (s *divinerService) returnMarketInfoResponse(bytes []byte) (*pbs.MarketInfoResponse, error) {
	msg, err := pbl.UnmarshalMarket(bytes)
	if err != nil {
		return nil, err
	}

	return &pbs.MarketInfoResponse{
		Market: msg,
		Time:   ptypes.TimestampNow(),
	}, nil
}

func (s *divinerService) QueryMarket(ctx context.Context, req *pbs.QueryRequest) (*pbs.MarketInfoResponse, error) {
	if ok, err := pbs.CheckQueryRequest(req, s.Expired); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("data is illegal")
	}

	client, err := s.SDK.NewChannelClient(s.ChannelID, s.User)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	bytes, err := s.queryFabricByID(client, "market", "query", req.Id)
	if err != nil {
		return nil, err
	}
	return s.returnMarketInfoResponse(bytes)
}

func (s *divinerService) CreateMarket(ctx context.Context, req *pbs.MarketCreateRequest) (*pbs.MarketInfoResponse, error) {
	if ok, err := pbs.CheckMarketCreateRequest(req, s.Expired); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("data is illegal")
	}

	num, err := cast.ToBytes(req.Num)
	if err != nil {
		return nil, err
	}

	flag, err := cast.ToBytes(req.IsFund)
	if err != nil {
		return nil, err
	}

	client, err := s.SDK.NewChannelClient(s.ChannelID, s.User)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	notifier, rce := s.registerChaincodeEvent(client, s.Chaincode, "market")
	defer client.UnregisterChaincodeEvent(rce)

	tx, err := s.executeFabric(client, "market", "create", []byte(req.User), []byte(req.Event), num, flag)
	if err != nil {
		fmt.Printf("execute fabric error: %v\n", err)
		return nil, err
	}

	bytes := fabsdk.SelectEvent(notifier, tx.ID, s.Wait)
	if bytes != nil {
		return s.returnMarketInfoResponse(bytes)
	}
	return nil, fmt.Errorf("can not get event notify")
}

func (s *divinerService) returnTxResponse(data []byte) (*pbs.TxResponse, error) {
	price, err := cast.BytesToFloat64(data)
	if err != nil {
		return nil, err
	}
	return &pbs.TxResponse{
		Price: price,
		Time:  ptypes.TimestampNow(),
	}, nil

}

func (s *divinerService) Tx(ctx context.Context, req *pbs.TxRequest) (*pbs.TxResponse, error) {
	if ok, err := pbs.CheckTxRequest(req, s.Expired); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("data is illegal")
	}

	volume, err := cast.ToBytes(req.Volume)
	if err != nil {
		return nil, err
	}

	client, err := s.SDK.NewChannelClient(s.ChannelID, s.User)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	var cmd string
	if req.IsBuy {
		cmd = "buy"
	} else {
		cmd = "sell"
	}

	notifier, rce := s.registerChaincodeEvent(client, s.Chaincode, "tx")
	defer client.UnregisterChaincodeEvent(rce)
	tx, err := s.executeFabric(client, "tx", cmd, []byte(req.User), []byte(req.Share), volume)
	if err != nil {
		return nil, err
	}

	bytes := fabsdk.SelectEvent(notifier, tx.ID, s.Wait)
	if bytes != nil {
		return s.returnTxResponse(bytes)
	}
	return nil, fmt.Errorf("can not get event notify")
}

func (s *divinerService) monitor() {
	log.Println("monitor starting...")
	client, err := s.SDK.NewChannelClient(s.ChannelID, s.User)
	if err != nil {
		return
	}
	defer client.Close()

	notifier := make(chan *apitxn.CCEvent)
	rce := client.RegisterChaincodeEvent(notifier, s.Chaincode, "member([a-zA-z0-9]+)")
	defer client.UnregisterChaincodeEvent(rce)

	for {
		select {
		case ccEvent := <-notifier:
			log.Printf("chaincode: %s, name: %s, txid: %s\n", ccEvent.ChaincodeID, ccEvent.EventName, ccEvent.TxID)
		case <-time.After(time.Second * 5):
			log.Println("wait...")
		}
	}

}

func main() {
	// TODO: add trace chaincode and block event with gorouting.
	conf, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	lis, err := net.Listen(conf.GetString("protocol"), conf.GetString("listen"))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	service := &divinerService{
		FabricConfig: conf.GetString("fabric"),
		ChannelID:    conf.GetString("channel"),
		Chaincode:    conf.GetString("chaincode"),
		User:         conf.GetString("user"),
		Expired:      conf.GetInt64("expired"),
		Wait:         time.Duration(conf.GetInt64("wait")),
		Balance:      conf.GetFloat64("balance"),
	}

	service.SDK, err = fabsdk.NewSDK(conf.GetString("fabric"))
	if err != nil {
		log.Fatalf("init fab sdk error: %v", err)
	}

	session, err := service.SDK.NewPreEnrolledUserSession("Diviner", "User1")
	if err != nil {
		log.Fatalf("session error: %v\n", err)
	}
	service.Client, err = service.SDK.NewSystemClient(session)
	if err != nil {
		log.Fatalf("fabric client error: %v\n", err)
	}

	service.Channel, err = fabsdk.GetChannel(service.Client, service.ChannelID, []string{"Diviner"})
	if err != nil {
		log.Fatalf("get channel error: %v\n", err)
	}

	pbs.RegisterDivinerSerivceServer(s, service)

	reflection.Register(s)
	go service.monitor()

	log.Println("serving...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Println("end")
}
