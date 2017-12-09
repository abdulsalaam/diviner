package member

import (
	fmt "fmt"
	"strings"

	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

// NewMemberRegisterRequest ...
func NewMemberRegisterRequest(addr string) *MemberRegisterRequest {
	return &MemberRegisterRequest{
		Address: addr,
	}
}

// CheckMemberRegisterRequest ...
func CheckMemberRegisterRequest(in *MemberRegisterRequest) (bool, error) {
	if strings.TrimSpace(in.Address) == "" {
		return false, fmt.Errorf("address is empty")
	}

	return true, nil
}

// NewMember ...
func NewMember(req *MemberRegisterRequest, balance float64) (*Member, error) {
	if ok, err := CheckMemberRegisterRequest(req); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("request error")
	}

	return &Member{
		Address: req.Address,
		Balance: balance,
		Assets:  make(map[string]float64),
		Subsidy: 0.0,
		Blocked: false,
	}, nil
}

// NewTransfer ...
func NewTransfer(target string, amount float64) (*Transfer, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be larger than 0, but %v", amount)
	}

	return &Transfer{
		Target: target,
		Amount: amount,
		Time:   ptypes.TimestampNow(),
	}, nil
}

// CheckTransfer ...
func CheckTransfer(in *Transfer) (bool, error) {
	if in.Amount <= 0 {
		return false, fmt.Errorf("amount must be larger than 0, but %v", in.Amount)
	}

	if strings.TrimSpace(in.Target) == "" {
		return false, fmt.Errorf("target is empty")
	}

	return true, nil
}

// UnmarshalMember ...
func UnmarshalMember(data []byte) (*Member, error) {
	ret := new(Member)
	if err := proto.Unmarshal(data, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// UnmarshalMemberRegisterRequest ...
func UnmarshalMemberRegisterRequest(data []byte) (*MemberRegisterRequest, error) {
	ret := new(MemberRegisterRequest)
	if err := proto.Unmarshal(data, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// UnmarshalTransfer ...
func UnmarshalTransfer(data []byte) (*Transfer, error) {
	ret := new(Transfer)
	if err := proto.Unmarshal(data, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// MarshalMember ...
func MarshalMember(in *Member) ([]byte, error) {
	return proto.Marshal(in)
}

// MarshalMemberRegisterRequest ...
func MarshalMemberRegisterRequest(in *MemberRegisterRequest) ([]byte, error) {
	return proto.Marshal(in)
}

// MarshalTransfer ...
func MarshalTransfer(in *Transfer) ([]byte, error) {
	return proto.Marshal(in)
}
