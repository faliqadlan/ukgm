package aggregatebuilder

import (
	"math"

	"go.mongodb.org/mongo-driver/bson"
)

type (
	AggregationHelper struct {
		requests []*AggregationRequest
	}

	AggregateObject interface {
		ToBSON() (result bson.M)
	}

	AggregationRequest struct {
		index           int
		AggregateMember AggregateObject
	}
)

func (helper *AggregationHelper) GetRequests() (result []*AggregationRequest) {
	return helper.requests
}

func (helper *AggregationHelper) AppendRequest(request *AggregationRequest) {
	// if request.Validate() != nil {
	// 	return
	// }
	var startIndex int = len(helper.requests)
	request.index = startIndex
	helper.requests = append(helper.requests, request)
}

func (request *AggregationHelper) GetObjectByIndex(index int) (result *AggregationRequest) {
	var mapAggRequest map[int]*AggregationRequest = map[int]*AggregationRequest{}
	var maxIndex, minIndex int
	minIndex = math.MaxInt64
	for _, req := range request.requests {
		mapAggRequest[req.index] = req
		if req.index > maxIndex {
			maxIndex = req.index
		}
		if req.index < minIndex {
			minIndex = req.index
		}
	}
	if index > maxIndex {
		index = maxIndex
	} else if index < minIndex {
		index = minIndex
	}
	return mapAggRequest[index]
}

func (aggregate *AggregationHelper) ToBSON() (result bson.A) {
	for i := 0; i < len(aggregate.requests); i++ {
		result = append(result, aggregate.GetObjectByIndex(i).AggregateMember.ToBSON())
	}
	return
}

// func (aggregate *AggregationRequest) Validate() (err error) {
// 	if len(aggregate.AggregateMember) == 0 {
// 		err = fmt.Errorf("Request is invalid")
// 		zap.S().Error(err)
// 	}
// 	return
// }
