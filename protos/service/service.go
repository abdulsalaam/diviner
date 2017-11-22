package service

import (
	pbc "diviner/protos/common"
	pbm "diviner/protos/member"

	"diviner/common/cast"

	"github.com/hyperledger/fabric/bccsp"
)

func NewQueryRequest(priv bccsp.Key, id string) (*QueryRequest, error) {
	v, err := pbc.NewVerification(priv, []byte(id))
	if err != nil {
		return nil, err
	}

	return &QueryRequest{
		Id:    id,
		Check: v,
	}, nil
}

func NewMemberCreateRequest(priv bccsp.Key) (*MemberCreateRequest, error) {
	member, err := pbm.NewMember(priv, 0.0)
	if err != nil {
		return nil, err
	}

	bytes, err := pbm.Marshal(member)
	if err != nil {
		return nil, err
	}

	v, err := pbc.NewVerification(priv, bytes)
	if err != nil {
		return nil, err
	}

	return &MemberCreateRequest{
		Member: member,
		Check:  v,
	}, nil
}

func NewEventCreateRequest(priv bccsp.Key, user, title string, outcomes []string) (*EventCreateRequest, error) {
	var tmp []string
	tmp = append(tmp, user)
	tmp = append(tmp, title)
	tmp = append(tmp, outcomes...)

	bytes, err := cast.StringsToBytes(tmp...)
	if err != nil {
		return nil, err
	}

	v, err := pbc.NewVerification(priv, bytes)
	if err != nil {
		return nil, err
	}

	return &EventCreateRequest{
		User:    user,
		Title:   title,
		Outcome: outcomes,
		Check:   v,
	}, nil
}

func NewMarketCreateRequest(priv bccsp.Key, user, event string, num float64, fund bool) (*MarketCreateRequest, error) {
	bytes, err := cast.ToBytes(user, event, num, fund)
	if err != nil {
		return nil, err
	}

	v, err := pbc.NewVerification(priv, bytes)
	if err != nil {
		return nil, err
	}

	return &MarketCreateRequest{
		User:   user,
		Event:  event,
		Num:    num,
		IsFund: fund,
		Check:  v,
	}, nil
}
