package member

import (
	"diviner/common/base58"
	fmt "fmt"

	"github.com/golang/protobuf/ptypes"

	proto "github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/bccsp"
)

// NewMember ...
func NewMember(addr string, balance float64) *Member {
	return &Member{
		Address: addr,
		Balance: balance,
		Assets:  make(map[string]float64),
		Subsidy: 0.0,
		Blocked: false,
	}
}

// NewMemberWithPrivateKey ...
func NewMemberWithPrivateKey(priv bccsp.Key, balance float64) (*Member, error) {
	if !priv.Private() {
		return nil, fmt.Errorf("priv must be a private key")
	}

	pub, err := priv.PublicKey()
	if err != nil {
		return nil, err
	}

	addr, err := base58.Address(pub)
	if err != nil {
		return nil, err
	}

	return NewMember(addr, balance), nil
}

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

// Unmarshal ...
func Unmarshal(data []byte) (*Member, error) {
	mem := &Member{}
	if err := proto.Unmarshal(data, mem); err != nil {
		return nil, err
	}
	return mem, nil
}

// UnmarshalTransfer ...
func UnmarshalTransfer(data []byte) (*Transfer, error) {
	tx := &Transfer{}
	if err := proto.Unmarshal(data, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// Marshal ...
func Marshal(m *Member) ([]byte, error) {
	return proto.Marshal(m)
}

// MarshalTransfer ...
func MarshalTransfer(tx *Transfer) ([]byte, error) {
	return proto.Marshal(tx)
}
