package member

import (
	"diviner/common/base58"

	"github.com/gogo/protobuf/proto"
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
	}, nil
}

// Unmarshal ...
func Unmarshal(data []byte) (*Member, error) {
	mem := &Member{}
	err := proto.Unmarshal(data, mem)
	return mem, err
}

// Marshal ...
func Marshal(m *Member) ([]byte, error) {
	return proto.Marshal(m)
}
