package collector

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func getFloat64FromString(s string) (float64, error) {
	n, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		fmt.Println(err)
		return 0.0, err
	}
	return math.Float64frombits(uint64(n)), err
}

func getFloatVersionFromString(s string) (float64, error) {
	value := strings.Replace(s, ".", "", -1)
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Println(err)
		return 0.0, err
	}
	return v, err
}
