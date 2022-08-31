package aggregatebuilder

import (
	"fmt"

	constants "ukgm/internal/app/ukgm/constant"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type (
	AggregationGroup struct {
		GroupIDs     []GroupID
		Accumulators []AggregationAccumulator
	}

	GroupID struct {
		SourceField      string
		DestinationField string
	}

	AggregationAccumulator struct {
		Type    string
		Field   string
		DBField string
	}
)

func (accumulator *AggregationAccumulator) ValidateType() (err error) {
	if accumulator.Type == "" || !constants.ValidateAggregationAccumulatorType(accumulator.Type) {
		err = fmt.Errorf("Type is invalid")
		zap.S().Error(err)
	}
	return
}

func (accumulator *AggregationAccumulator) ValidateField() (err error) {
	if accumulator.DBField == "" || accumulator.Field == "" {
		err = fmt.Errorf("Field is invalid")
		zap.S().Error(err)
	}
	return
}

func (accumulator *AggregationAccumulator) Validate() (err error) {
	err = accumulator.ValidateField()
	if err != nil {
		zap.S().Error(err)
		return
	}
	err = accumulator.ValidateType()
	if err != nil {
		zap.S().Error(err)
		return
	}
	return
}

func (accumulator *AggregationAccumulator) ToBSON() (result bson.E) {
	if accumulator.Validate() != nil {
		return
	}
	result = bson.E{
		accumulator.Field, bson.D{
			{
				fmt.Sprintf("$%s", accumulator.Type), fmt.Sprintf("$%s", accumulator.DBField),
			},
		},
	}
	return
}

func (groupID *GroupID) Validate() (err error) {
	if groupID.DestinationField == "" || groupID.SourceField == "" {
		err = fmt.Errorf("Field is invalid")
		zap.S().Error(err)
	}
	return
}

func (groupID *GroupID) ToBSON() (result bson.E) {
	if groupID.Validate() != nil {
		return
	}
	result = bson.E{
		groupID.DestinationField, fmt.Sprintf("$%s", groupID.SourceField),
	}
	return
}

func (aggregate *AggregationGroup) ToBSON() (result bson.M) {
	result = bson.M{}
	IDBsonD := bson.D{}
	for _, groupID := range aggregate.GroupIDs {
		IDBsonD = append(IDBsonD, groupID.ToBSON())
	}
	combinationBsonD := bson.D{
		{"_id", IDBsonD},
	}
	for _, accumulator := range aggregate.Accumulators {
		combinationBsonD = append(combinationBsonD, accumulator.ToBSON())
	}
	result["$group"] = combinationBsonD
	return
}
