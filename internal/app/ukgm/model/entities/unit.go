package entities

type Unit struct {
	ID         interface{} `json:"id" bson:"_id, omitempty"`
	DistrictID string      `json:"district_id" bson:"district_id, omitempty"`
	Name       string      `json:"name" bson:"name"`
	Patients   []*Patient  `json:"patients" bson:"patients"`

	// include
	Timestamp
}
