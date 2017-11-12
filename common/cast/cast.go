package cast

import (
	"bytes"
	"encoding/binary"
)

// ToFloat64 cast byte array to float64
func ToFloat64(in []byte) (float64, error) {
	var ret float64
	buf := bytes.NewReader(in)
	err := binary.Read(buf, binary.LittleEndian, &ret)
	if err != nil {
		return 0.0, err
	}
	return ret, nil
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
