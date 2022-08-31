package aggregatebuilder_test

import (
	"testing"

	"github.com/faliqadlan/ukgm/internal/app/ukgm/model/request"
	aggregatebuilder "github.com/faliqadlan/ukgm/internal/app/ukgm/mongoRepository/aggregateBuilder"
)

func TestAggregationMatchToBSON(t *testing.T) {
	var MatchAggregate *aggregatebuilder.AggregationMatch = &aggregatebuilder.AggregationMatch{
		Pr: &request.PageRequest{
			Filters: []request.Filter{
				request.Filter{
					Option:   "document_number",
					Operator: "=",
					Value:    "OK",
				},
			},
		},
	}

	result := MatchAggregate.ToBSON()
	if result == nil {
		t.Fatal("NIL")
	}
}

func TestAggregationGeoNearToBSON(t *testing.T) {
	var MatchGeoNear *aggregatebuilder.AggregationGeoNear = &aggregatebuilder.AggregationGeoNear{
		Longitude:     100,
		Latitude:      100,
		MaxDistance:   5000,
		DistanceField: "distance",
	}

	result := MatchGeoNear.ToBSON()
	if result == nil {
		t.Fatal("NIL")
	}
}

func TestAggregationGroupToBSON(t *testing.T) {
	var Aggregate *aggregatebuilder.AggregationGroup = &aggregatebuilder.AggregationGroup{
		GroupIDs: []aggregatebuilder.GroupID{
			aggregatebuilder.GroupID{
				SourceField:      "sub_cob",
				DestinationField: "sub_cob",
			},
		},
		Accumulators: []aggregatebuilder.AggregationAccumulator{
			aggregatebuilder.AggregationAccumulator{
				Type:    "sum",
				Field:   "total_full_sum_insured",
				DBField: "total_full_sum_insured",
			},
		},
	}

	result := Aggregate.ToBSON()
	if result == nil {
		t.Fatal("NIL")
	}
}

func TestAggregationHelperToBSON(t *testing.T) {
	var Aggregate *aggregatebuilder.AggregationGroup = &aggregatebuilder.AggregationGroup{
		GroupIDs: []aggregatebuilder.GroupID{
			aggregatebuilder.GroupID{
				SourceField:      "sub_cob",
				DestinationField: "sub_cob",
			},
		},
		Accumulators: []aggregatebuilder.AggregationAccumulator{
			aggregatebuilder.AggregationAccumulator{
				Type:    "sum",
				Field:   "total_full_sum_insured",
				DBField: "total_full_sum_insured",
			},
		},
	}

	var MatchAggregate *aggregatebuilder.AggregationMatch = &aggregatebuilder.AggregationMatch{
		Pr: &request.PageRequest{
			Filters: []request.Filter{
				request.Filter{
					Option:   "document_number",
					Operator: "=",
					Value:    "OK",
				},
			},
		},
	}

	var MatchGeoNear *aggregatebuilder.AggregationGeoNear = &aggregatebuilder.AggregationGeoNear{
		Longitude:     100,
		Latitude:      100,
		MaxDistance:   5000,
		DistanceField: "distance",
	}

	var Helper *aggregatebuilder.AggregationHelper = &aggregatebuilder.AggregationHelper{}
	Helper.AppendRequest(&aggregatebuilder.AggregationRequest{AggregateMember: MatchAggregate})
	Helper.AppendRequest(&aggregatebuilder.AggregationRequest{AggregateMember: MatchGeoNear})
	Helper.AppendRequest(&aggregatebuilder.AggregationRequest{AggregateMember: Aggregate})
	result := Helper.ToBSON()
	if result == nil {
		t.Fatal("NIL")
	}
}
