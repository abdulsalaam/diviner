package common

import (
	"diviner/common/csp"
	"diviner/common/util"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric/bccsp"
	perrors "github.com/pkg/errors"
)

// NewVerification ...
func NewVerification(priv bccsp.Key, in []byte) (*Verification, error) {
	if !priv.Private() {
		return nil, perrors.New("priv must be a private key")
	}

	pubKey, err := csp.GetPublicKeyBytes(priv)
	if err != nil {
		return nil, err
	}

	hash, err := csp.Hash(in)
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
		Time:      ptypes.TimestampNow(),
	}, nil

}

// Verify ...
func Verify(v *Verification, in []byte, expired int64) (bool, error) {
	curr := ptypes.TimestampNow()
	if curr.Seconds-v.Time.Seconds >= expired {
		return false, perrors.New("data is expired")
	}

	hash, err := csp.Hash(in)
	if err != nil {
		return false, err
	}

	if util.CmpByteArray(hash, v.Hash) != 0 {
		return false, perrors.New("hash content not equal")
	}

	key, err := csp.ImportPubFromRaw(v.PublicKey)
	if err != nil {
		return false, err
	}

	return csp.Verify(key, v.Signature, hash)
}
