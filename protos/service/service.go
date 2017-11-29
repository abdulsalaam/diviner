package service

import (
	pbc "diviner/protos/common"
	pbm "diviner/protos/member"
	"fmt"

	"diviner/common/base58"
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

func CheckQueryRequest(req *QueryRequest, expired int64) (bool, error) {
	return pbc.Verify(req.Check, []byte(req.Id), expired)
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

func CheckMemberCreateRequest(req *MemberCreateRequest, expired int64) (bool, error) {
	if req.Member.Id != req.Member.Address {
		return false, fmt.Errorf("member id and address not match: %s, %s", req.Member.Id, req.Member.Address)
	}

	if base58.Encode(req.Check.PublicKey) != req.Member.Address {
		return false, fmt.Errorf("member address and public not match")
	}

	bytes, err := pbm.Marshal(req.Member)
	if err != nil {
		return false, err
	}

	return pbc.Verify(req.Check, bytes, expired)

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

func CheckEventCreateRequest(req *EventCreateRequest, expired int64) (bool, error) {
	if req.User != base58.Encode(req.Check.PublicKey) {
		return false, fmt.Errorf("creator and caller are not match")
	}

	var tmp []string
	tmp = append(tmp, req.User)
	tmp = append(tmp, req.Title)
	tmp = append(tmp, req.Outcome...)

	if bytes, err := cast.StringsToBytes(tmp...); err != nil {
		return false, err
	} else {
		return pbc.Verify(req.Check, bytes, expired)
	}
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

func CheckMarketCreateRequest(req *MarketCreateRequest, expired int64) (bool, error) {
	if req.User != base58.Encode(req.Check.PublicKey) {
		return false, fmt.Errorf("creator and caller are not match")
	}

	if bytes, err := cast.ToBytes(req.User, req.Event, req.Num, req.IsFund); err != nil {
		return false, err
	} else {
		return pbc.Verify(req.Check, bytes, expired)
	}
}

func NewTxRequest(priv bccsp.Key, user string, buy bool, share string, volume float64) (*TxRequest, error) {
	bytes, err := cast.ToBytes(user, buy, share, volume)
	if err != nil {
		return nil, err
	}

	v, err := pbc.NewVerification(priv, bytes)
	if err != nil {
		return nil, err
	}

	return &TxRequest{
		User:   user,
		IsBuy:  buy,
		Share:  share,
		Volume: volume,
		Check:  v,
	}, nil
}

func CheckTxRequest(req *TxRequest, expired int64) (bool, error) {
	if req.User != base58.Encode(req.Check.PublicKey) {
		return false, fmt.Errorf("user are not match")
	}

	if bytes, err := cast.ToBytes(req.User, req.IsBuy, req.Share, req.Volume); err != nil {
		return false, err
	} else {
		return pbc.Verify(req.Check, bytes, expired)
	}
}
