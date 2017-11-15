package lmsr

import (
	proto "github.com/golang/protobuf/proto"
	perrors "github.com/pkg/errors"
)

func NewAsset(user, share string, volume float64) (*Asset, error) {
	if volume < 0 {
		return nil, perrors.Errorf("volume must be larger or equals 0: %v", volume)
	}

	id := AssetID(user, share)

	return &Asset{
		Id:     id,
		Volume: volume,
	}, nil
}

func UnmarshalAsset(data []byte) (*Asset, error) {
	asset := &Asset{}
	if err := proto.Unmarshal(data, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

func MarshalAsset(asset *Asset) ([]byte, error) {
	return proto.Marshal(asset)
}

func MarshalAssets(lst *Assets) ([]byte, error) {
	return proto.Marshal(lst)
}
