package aggregatebuilder

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	AggregationGeoNear struct {
		Longitude     float64
		Latitude      float64
		MaxDistance   float64
		DistanceField string
	}
)

func (aggregate *AggregationGeoNear) ToBSON() (result bson.M) {
	result = primitive.M{}
	// result["timestamp.deleted_at"] = nil
	result["$geoNear"] = bson.D{
		{"near",
			bson.D{
				{"type", "Point"},
				{"coordinates",
					bson.A{
						aggregate.Longitude,
						aggregate.Latitude,
					},
				},
			},
		},
		{
			"distanceField", aggregate.DistanceField,
		},
		{"maxDistance", aggregate.MaxDistance},
		{"query", bson.M{"timestamp.deleted_at": nil}},
	}
	return
}
