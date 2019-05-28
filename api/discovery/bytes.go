package discovery

func int64toBytes(i int64) []byte {
	return []byte{byte(i >> 56), byte(i >> 48), byte(i >> 40), byte(i >> 32), byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
}

func bytesToInt64(b []byte) int64 {
	var i int64 = 0
	for _, v := range b {
		i = (i << 8) | (int64(v) & 0xFF)
	}
	return i
}

func int32toBytes(i int32) []byte {
	return []byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
}

func bytesToInt32(b []byte) int32 {
	var i int32 = 0
	for _, v := range b {
		i = (i << 8) | (int32(v) & 0xFF)
	}
	return i
}
