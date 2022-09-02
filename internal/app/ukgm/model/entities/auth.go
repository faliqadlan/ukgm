package entities

type AuthUser struct {
	ID   string `json:"id" bson:"id,omitempty"`
	Name string `json:"name" bson:"name"`
}
