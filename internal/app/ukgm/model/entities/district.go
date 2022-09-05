package entities

type District struct {
	ID     interface{} `json:"id" bson:"_id,omitempty"`
	Email  string      `json:"email" bson:"email"`
	Name   string      `json:"name" bson:"name"`
	Uid    string      `json:"uid" bson:"uid"`
	Units  []*Unit     `json:"units" bson:"units"`
	Doctor []*Doctor   `json:"Doctor" bson:"Doctor"`
	Nurse  []*Nurse    `json:"Nurse" bson:"Nurse"`

	// include
	Timestamp
}

// func (d *District) ToAuthUser() *AuthUser {
// 	return &AuthUser{
// 		ID:   d.ID.(string),
// 		Name: d.Name,
// 	}
// }
