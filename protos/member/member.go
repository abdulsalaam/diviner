package member

import (
	"diviner/common/base58"

	proto "github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/bccsp"
	perrors "github.com/pkg/errors"
)

// NewMember ...
func NewMember(priv bccsp.Key, balance float64) (*Member, error) {
	if !priv.Private() {
		return nil, perrors.New("priv must be a private key")
	}

	pub, err := priv.PublicKey()
	if err != nil {
		return nil, err
	}

	addr, err := base58.Address(pub)
	if err != nil {
		return nil, err
	}

	return &Member{
		Id:      addr,
		Address: addr,
		Balance: balance,
		Assets:  make(map[string]float64),
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

// Marshal ...
func Marshal(m *Member) ([]byte, error) {
	return proto.Marshal(m)
}
