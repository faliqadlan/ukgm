package entities

type Doctor struct {
	ID         interface{} `json:"id" bson:"_id,omitempty"`
	DistrictID string      `json:"district_id" bson:"district_id,omitempty"`
	Name       string      `json:"name" bson:"name"`

	// include
	Timestamp
}
