package collector

import (
	"fmt"
	"hash/fnv"
	"strconv"
)

func GetStakeFromString(s string) float64 {
	if len(s) == 1 {
		return 0
	}
	l := len(s) - 19 - 5
	v, err := strconv.ParseFloat(s[0:l], 64)
	if err != nil {
		fmt.Println(err)
	}
	return float64(v)
}

func GetFloatFromString(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return v
}

func HashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
