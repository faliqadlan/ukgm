package entities

type Doctor struct {
	ID   interface{} `json:"id" bson:"_id, omitempty"`
	Name string      `json:"name" bson:"name"`

	// include
	Timestamp
}
