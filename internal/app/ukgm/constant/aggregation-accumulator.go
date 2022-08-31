package constants

import (
	"github.com/iman231/go-help"
)

const (
	AggregationAccumulatorTypeFirst   = "first"
	AggregationAccumulatorTypeLast    = "last"
	AggregationAccumulatorTypeMax     = "max"
	AggregationAccumulatorTypeMin     = "min"
	AggregationAccumulatorTypeCount   = "count"
	AggregationAccumulatorTypeSum     = "sum"
	AggregationAccumulatorTypeAverage = "avg"
)

var (
	aggregationAccumulatorType = []string{
		"first",
		"avg",
		"count",
		"last",
		"max",
		"min",
		"sum",
	}
)

func ValidateAggregationAccumulatorType(request string) bool {
	return help.InArray(request, aggregationAccumulatorType)
}
