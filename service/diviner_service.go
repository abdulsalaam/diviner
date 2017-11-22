package main

import (
	"context"
	"diviner/common/base58"
	"diviner/common/config"
	pbm "diviner/protos/member"
	pbs "diviner/protos/service"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type divinerService struct{}

func (s *divinerService) QueryMember(ctx context.Context, req *pbs.QueryRequest) (*pbs.MemberInfoResponse, error) {
	addr := base58.Encode(req.Check.PublicKey)
	if addr != req.Id {
		return nil, fmt.Errorf("can not query others, need %s, but %s", addr, req.Id)
	}

	// call chaincode

	log.Println("query member: ", req)
	member := &pbm.Member{
		Id:      req.Id,
		Address: req.Id,
		Balance: 0.0,
	}
	return &pbs.MemberInfoResponse{
		Member: member,
		Time:   ptypes.TimestampNow(),
	}, nil
}

func (s *divinerService) CreateMember(ctx context.Context, req *pbs.MemberCreateRequest) (*pbs.MemberInfoResponse, error) {
	log.Println("create member: ", req)
	member := &pbm.Member{
		Id:      req.Member.Id,
		Address: req.Member.Address,
		Balance: 0.0,
	}
	return &pbs.MemberInfoResponse{
		Member: member,
		Time:   ptypes.TimestampNow(),
	}, nil
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

	pbs.RegisterDivinerSerivceServer(s, &divinerService{})

	reflection.Register(s)

	log.Println("serving...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Println("end")
}
