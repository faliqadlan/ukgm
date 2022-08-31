package aggregatebuilder

import (
	"fmt"
	"ukgm/internal/app/ukgm/model/request"

	"go.mongodb.org/mongo-driver/bson"
)

type (
	AggregationMatch struct {
		Pr *request.PageRequest
	}
)

func (aggregate *AggregationMatch) ToBSON() (result bson.M) {
	filter, _ := aggregate.Pr.QueryPipeline([]string{})
	result = bson.M{}
	result["$match"] = filter
	fmt.Println(result)
	return
}
