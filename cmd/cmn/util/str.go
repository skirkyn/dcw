package util

func StrToByteSlice(str string) []byte {
	res := []byte(str)
	return res
}

func StrSliceToByteSlice(str []string) []byte {
	res := make([]byte, 0)
	for i := 0; i < len(str); i++ {
		res = append(res, []byte(str[i])...)
	}
	return res
}
