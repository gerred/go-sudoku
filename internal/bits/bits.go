package bits

import (
	"strconv"
)

func HasSingleBit(a uint) bool {
	if a == 0 {
		return false
	}
	return a&(a-1) == 0
}

func GetNumberOfSetBits(x uint) uint {
	var count uint
	for count = 0; x != 0; count++ {
		x &= x - 1
	}
	return count
}

func GetSingleBitValue(val uint) uint {
	// TODO: naive approach
	var m uint
	for m = 1; m <= 9; m++ {
		if val == 1 {
			break
		}
		val >>= 1
	}
	return m
}

func GetString(val uint) string {
	var msg string
	for m := 1; m <= 9; m++ {
		if val&0x01 == 1 {
			if msg != "" {
				msg += ","
			}
			msg += strconv.Itoa(m)
		}
		val >>= 1
	}
	return msg
}

func GetBitList(val uint) []uint {
	var list []uint
	for m := uint(1); m <= 9; m++ {
		if val&0x01 == 1 {
			list = append(list, 1<<(m-1))
		}
		val >>= 1
	}
	return list
}
