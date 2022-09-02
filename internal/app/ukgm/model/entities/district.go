package entities

type District struct {
	ID    interface{} `json:"id" bson:"_id, omitempty"`
	Email string      `json:"email" bson:"email"`
	Name  string      `json:"name" bson:"name"`
	// Unit
	Doctor []*Doctor `json:"Doctor" bson:"Doctor"`
	Nurse  []*Nurse  `json:"Nurse" bson:"Nurse"`
}

func (d *District) ToAuthUser() *AuthUser {
	return &AuthUser{
		ID:   d.ID.(string),
		Name: d.Name,
	}
}
