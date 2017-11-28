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

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"github.com/hyperledger/fabric-sdk-go/def/fabapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type divinerService struct {
	FabricConfig string
	SDK          *fabapi.FabricSDK
	ChannelID    string
	Chaincode    string
	User         string
	Expired      int64
	Balance      float64
}

func (s *divinerService) queryFabric(client apitxn.ChannelClient, chaincode, fcn string, data ...[]byte) ([]byte, error) {
	var args [][]byte
	args = append(args, []byte(chaincode))
	args = append(args, data...)

	qr := apitxn.QueryRequest{
		ChaincodeID: s.Chaincode,
		Fcn:         fcn,
		Args:        args,
	}

	return client.Query(qr)
}

func (s *divinerService) queryFabricById(client apitxn.ChannelClient, chaincode, fcn, id string) ([]byte, error) {
	return s.queryFabric(client, chaincode, fcn, []byte(id))
}

func (s *divinerService) executeFabric(client apitxn.ChannelClient, chaincode, fcn string, data ...[]byte) (apitxn.TransactionID, error) {
	var args [][]byte
	args = append(args, []byte(chaincode))
	args = append(args, data...)

	txr := apitxn.ExecuteTxRequest{
		ChaincodeID: s.Chaincode,
		Fcn:         fcn,
		Args:        args,
	}

	return client.ExecuteTx(txr)
}

func (s *divinerService) registerChaincodeEvent(client apitxn.ChannelClient, chaincode, name string) (chan *apitxn.CCEvent, apitxn.Registration) {
	id := name + "([a-zA-Z0-9]+)"
	notifier := make(chan *apitxn.CCEvent)
	rce := client.RegisterChaincodeEvent(notifier, s.Chaincode, id)
	return notifier, rce
}

func (s *divinerService) selectEvent(notifier chan *apitxn.CCEvent, timeout time.Duration) []byte {
	select {
	case evt := <-notifier:
		log.Println("get notifier")
		return evt.Payload
	case <-time.After(time.Second * timeout):
		log.Println("timeout")
		return nil
	}
}

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

	if bytes, err := s.queryFabricById(client, "member", "query", req.Id); err != nil {
		return nil, err
	} else {
		return s.returnMemberInfoResponse(bytes)
	}
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

	client, err := s.SDK.NewChannelClient(s.ChannelID, s.User)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	notifier, rce := s.registerChaincodeEvent(client, s.Chaincode, "member")
	defer client.UnregisterChaincodeEvent(rce)

	_, err = s.executeFabric(client, "member", "create", bytes)
	if err != nil {
		return nil, err
	}

	if bytes := s.selectEvent(notifier, 5); bytes != nil {
		return s.returnMemberInfoResponse(bytes)
	} else {
		if bytes, err := s.queryFabricById(client, "member", "query", member.Id); err != nil {
			return nil, err
		} else {
			return s.returnMemberInfoResponse(bytes)
		}
	}

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

	if bytes, err := s.queryFabricById(client, "event", "query", req.Id); err != nil {
		return nil, err
	} else {
		return s.returnEventInfoResponse(bytes)
	}
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

	_, err = s.executeFabric(client, "event", "create", data)
	if err != nil {
		fmt.Printf("execute fabric error: %v\n", err)
		return nil, err
	}

	if bytes := s.selectEvent(notifier, 5); bytes != nil {
		return s.returnEventInfoResponse(bytes)
	} else {
		if bytes, err := s.queryFabricById(client, "event", "query", event.Id); err != nil {
			return nil, err
		} else {
			return s.returnEventInfoResponse(bytes)
		}
	}
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

	if bytes, err := s.queryFabricById(client, "market", "query", req.Id); err != nil {
		return nil, err
	} else {
		return s.returnMarketInfoResponse(bytes)
	}
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

	ccEventID := "market([a-zA-Z0-9]+)"
	notifier := make(chan *apitxn.CCEvent)
	rce := client.RegisterChaincodeEvent(notifier, s.Chaincode, ccEventID)
	defer client.UnregisterChaincodeEvent(rce)

	_, err = s.executeFabric(client, "market", "create", []byte(req.User), []byte(req.Event), num, flag)
	if err != nil {
		fmt.Printf("execute fabric error: %v\n", err)
		return nil, err
	}

	select {
	case evt := <-notifier:
		log.Println("get notifier")
		return s.returnMarketInfoResponse(evt.Payload)
	case <-time.After(time.Second * 5):
		log.Println("timeout")
		return nil, fmt.Errorf("can not get event notify")
	}
}

func (s *divinerService) Tx(ctx context.Context, req *pbs.TxRequest) (*pbs.TxResponse, error) {
	return nil, nil
}

func main() {
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
		Balance:      conf.GetFloat64("balance"),
	}

	service.SDK, err = fabsdk.NewSDK(conf.GetString("fabric"))
	if err != nil {
		log.Fatalf("init fab sdk error: %v", err)
	}

	log.Println(service)

	pbs.RegisterDivinerSerivceServer(s, service)

	reflection.Register(s)

	log.Println("serving...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Println("end")
}
