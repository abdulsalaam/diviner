package common

import (
	"diviner/common/csp"
	"diviner/common/util"
	"errors"

	"github.com/golang/protobuf/ptypes/timestamp"

	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric/bccsp"
)

func concate(in ...[]byte) []byte {
	sum := 0

	for _, x := range in {
		sum += len(x)
	}

	tmp := make([]byte, sum)

	idx := 0

	for _, x := range in {
		copy(tmp[idx:], x)
		idx += len(x)
	}

	return tmp
}

func hash(in []byte, ts *timestamp.Timestamp) ([]byte, error) {
	bytes, err := proto.Marshal(ts)
	if err != nil {
		return nil, err
	}

	tmp := concate(in, bytes)

	return csp.Hash(tmp)
}

// NewVerification ...
func NewVerification(priv bccsp.Key, in []byte) (*Verification, error) {
	if !priv.Private() {
		return nil, errors.New("priv must be a private key")
	}

	pubKey, err := csp.GetPublicKeyBytes(priv)
	if err != nil {
		return nil, err
	}

	curr := ptypes.TimestampNow()
	hash, err := hash(in, curr)

	if err != nil {
		return nil, err
	}

	sign, err := csp.Sign(priv, hash)
	if err != nil {
		return nil, err
	}

	return &Verification{
		PublicKey: pubKey,
		Hash:      hash,
		Signature: sign,
		Time:      curr,
	}, nil

}

// Verify ...
func Verify(v *Verification, in []byte, expired int64) (bool, error) {
	curr := ptypes.TimestampNow()
	if curr.Seconds-v.Time.Seconds >= expired {
		return false, errors.New("data is expired")
	}

	hash, err := hash(in, v.Time)

	if err != nil {
		return false, err
	}

	if util.CmpByteArray(hash, v.Hash) != 0 {
		return false, errors.New("hash content not equal")
	}

	key, err := csp.ImportPubFromRaw(v.PublicKey)
	if err != nil {
		return false, err
	}

	return csp.Verify(key, v.Signature, hash)
}

// Unmarshal ...
func Unmarshal(data []byte) (*Verification, error) {
	v := &Verification{}
	if err := proto.Unmarshal(data, v); err != nil {
		return nil, err
	}
	return v, nil
}

// Marshal ...
func Marshal(v *Verification) ([]byte, error) {
	return proto.Marshal(v)
}
