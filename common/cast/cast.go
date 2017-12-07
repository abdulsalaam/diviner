package cast

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// BytesToFloat64 cast byte array to float64
func BytesToFloat64(in []byte) (float64, error) {
	var ret float64
	buf := bytes.NewReader(in)
	err := binary.Read(buf, binary.LittleEndian, &ret)
	if err != nil {
		return 0.0, err
	}
	return ret, nil
}

// BytesToBool ...
func BytesToBool(in []byte) (bool, error) {
	if len(in) != 1 {
		return false, fmt.Errorf("in length error: %d", len(in))
	}

	return in[0] != 0, nil
}

// StringsToBytes ...
func StringsToBytes(in ...string) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, x := range in {
		if err := binary.Write(buf, binary.LittleEndian, []byte(x)); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// ToBytes cast to byte array
func ToBytes(in ...interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	var err error
	for _, x := range in {
		switch x.(type) {
		case string:
			err = binary.Write(buf, binary.LittleEndian, []byte(x.(string)))
		default:
			err = binary.Write(buf, binary.LittleEndian, x)
		}

		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// StringsToByteArray ...
func StringsToByteArray(args ...string) [][]byte {
	result := make([][]byte, len(args))

	for i, x := range args {
		result[i] = []byte(x)
	}

	return result
}

// ByteArrayToStrings ...
func ByteArrayToStrings(in [][]byte) []string {
	result := make([]string, len(in))
	for i, x := range in {
		result[i] = string(x)
	}

	return result
}
