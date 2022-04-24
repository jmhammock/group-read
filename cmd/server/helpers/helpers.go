package helpers

import "strconv"

func StrToUint32(str string) (uint32, error) {
	v, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(v), nil
}
