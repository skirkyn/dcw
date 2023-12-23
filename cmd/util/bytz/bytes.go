package bytz

func IntToByteSlice(x int) []byte {
	buf := make([]byte, 4)
	mask := 0xFF
	for i := 0; i < len(buf); i++ {
		buf[i] = byte((x >> (i * 8)) & mask)
	}
	return buf
}

func StringToByteSlice(str string) []byte {
	res := []byte(str)
	return res
}

func SliceToByteSlice(str []string) []byte {
	res := make([]byte, 0)
	for i := 0; i < len(str); i++ {
		res = append(res, []byte(str[i])...)
	}
	return res
}
