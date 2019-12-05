package main

func convertInt16ToByte(a []int16) []byte {
	b := make([]byte, 2*len(a))
	bI := 0
	for i := 0; i < len(a); i++ {
		b[bI] = byte(uint(a[i]))
		b[bI+1] = byte(uint(a[i] >> 8))
		bI += 2
	}
	return b
}

func convertByteToInt16(b []byte) []int16 {
	a := make([]int16, len(b)/2)
	bI := 0
	for i := 0; i < len(a); i++ {
		a[i] = int16(uint16(b[bI]) | uint16(b[bI+1])<<8)
		bI += 2
	}
	return a
}
