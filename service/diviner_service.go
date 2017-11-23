package main

import (
	"context"
	"diviner/common/base58"
	"diviner/common/config"
	"diviner/fabsdk"
	pbm "diviner/protos/member"
	pbs "diviner/protos/service"
	"fmt"
	"log"
	"net"

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
	User         string
	Expired      int64
	Balance      float64
}

func (s *divinerService) queryFabric(client apitxn.ChannelClient, chaincode, fcn, id string) ([]byte, error) {
	qr := apitxn.QueryRequest{
		ChaincodeID: chaincode,
		Fcn:         fcn,
		Args:        [][]byte{[]byte(id)},
	}

	return client.Query(qr)
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

	if bytes, err := s.queryFabric(client, "member", "query", req.Id); err != nil {
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

	txr := apitxn.ExecuteTxRequest{
		ChaincodeID: "member",
		Fcn:         "create",
		Args:        [][]byte{bytes},
	}

	_, err = client.ExecuteTx(txr)
	if err != nil {
		return nil, err
	}

	if bytes, err = s.queryFabric(client, "member", "query", member.Id); err != nil {
		return nil, err
	} else {
		return s.returnMemberInfoResponse(bytes)
	}

}

func (s *divinerService) QueryEvent(ctx context.Context, req *pbs.QueryRequest) (*pbs.EventInfoResponse, error) {
	return nil, nil
}

func (s *divinerService) CreateEvent(ctx context.Context, req *pbs.EventCreateRequest) (*pbs.EventInfoResponse, error) {
	return nil, nil
}

func (s *divinerService) QueryMarket(ctx context.Context, req *pbs.QueryRequest) (*pbs.MarketInfoResponse, error) {
	return nil, nil
}

func (s *divinerService) CreateMarket(ctx context.Context, req *pbs.MarketCreateRequest) (*pbs.MarketInfoResponse, error) {
	return nil, nil
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
