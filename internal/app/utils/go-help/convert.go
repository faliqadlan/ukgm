package help

import (
	"math"
	"strconv"
	"strings"
)

// FloatToAmountString func
func FloatToAmountString(f float64) string {
	f = ToFixed(f)

	fString := strconv.FormatFloat(math.Round(f*100000)/100000, 'f', -1, 64)

	splitDecimals := strings.Split(fString, ".")

	tempCount := len(splitDecimals[0]) % 3

	if tempCount == 0 {
		tempCount = 3
	}

	ret := make([]byte, 0)

	for i := 0; i < len(splitDecimals[0]); i++ {
		ret = append(ret, splitDecimals[0][i])
		tempCount = tempCount - 1
		if tempCount == 0 && i != len(splitDecimals[0])-1 {
			ret = append(ret, ',')
			tempCount = 3
		}
	}

	splitDecimals[0] = string(ret)

	return strings.Join(splitDecimals, ".")
}

// IntToAmountString func
func IntToAmountString(i int) string {
	fString := strconv.Itoa(i)

	tempCount := len(fString) % 3

	ret := make([]byte, 0)

	for i := 0; i < len(fString); i++ {
		ret = append(ret, fString[i])
		tempCount = tempCount - 1
		if tempCount == 0 && i != len(fString)-1 {
			ret = append(ret, ',')
			tempCount = 3
		}
	}

	return string(ret)
}

// RenderRow func
func RenderRow(row map[string][][]byte) map[string][]string {
	ret := make(map[string][]string)

	max := 0
	for _, v := range row {
		if len(v) > max {
			max = len(v)
		}
	}

	for k, v := range row {
		ret[k] = make([]string, max)
		for i := 0; i < max; i++ {
			if len(v)-i > 0 {
				ret[k][i] = string(v[i])
			} else {
				ret[k][i] = ""
			}
		}
	}

	return ret
}
