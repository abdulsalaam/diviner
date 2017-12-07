package member

import (
	"diviner/common/base58"
	fmt "fmt"

	proto "github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/bccsp"
)

// NewMember ...
func NewMember(pub []byte, balance float64) *Member {
	addr := base58.Encode(pub)

	return &Member{
		Address: addr,
		Balance: balance,
		Assets:  make(map[string]float64),
		Subsidy: 0.0,
		Blocked: false,
	}
}

// NewMember ...
func NewMember(priv bccsp.Key, balance float64) (*Member, error) {
	if !priv.Private() {
		return nil, fmt.Errorf("priv must be a private key")
	}

	pub, err := priv.PublicKey()
	if err != nil {
		return nil, err
	}

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
