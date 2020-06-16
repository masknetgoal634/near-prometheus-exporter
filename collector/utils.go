package collector

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
)

func StringToFloat64(s string) float64 {
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

func GetFloatVersionFromString(s string) (float64, error) {
	value := strings.Replace(s, ".", "", -1)
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Println(err)
		return 0.0, err
	}
	return v, err
}

func HashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
