package collector

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
)

func StringToInt64(s string) int64 {
	v, err := strconv.ParseInt(s[0:6], 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	return int64(v)
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
