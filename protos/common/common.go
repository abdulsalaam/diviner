package common

import (
	"diviner/common/csp"
	"diviner/common/util"
	"fmt"

	"github.com/golang/protobuf/ptypes"
)

// Verify ...
func Verify(v *Verification, in []byte, expired int64) (bool, error) {
	curr := ptypes.TimestampNow()
	if curr.Seconds-v.Time.Seconds >= expired {
		return false, fmt.Errorf("data is expired")
	}

	hash, err := csp.Hash(in)
	if err != nil {
		return false, nil
	}

	if util.CmpByteArray(hash, v.Hash) != 0 {
		return false, nil
	}

	key, err := csp.ImportPubFromRaw(v.PublicKey)
	if err != nil {
		return false, err
	}

	return csp.Verify(key, v.Signature, hash)
}
