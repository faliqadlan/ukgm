package entities

type Auth struct {
	User  AuthUser `json: "auth_user, omitempty" bson: "auth_user,omitempty"`
	Token string   `json:"token"`
}

type AuthUser struct {
	ID       string         `json:"id,omitempty" bson:"id,omitempty"`
	FullName string         `json:"full_name,omitempty" bson:"full_name,omitempty"`
	Role     []string       `json:"role"`
	RoleRef  *RoleReference `json:"role_reference" bson:"role_reference"`
}

type RoleReference struct {
	// MAIN
	ID   interface{} `json:"id" bson:"_id,omitempty"`
	Name string      `json:"name" bson:"name,omitempty"`
	// Modules     []string         `json:"modules" bson:"modules,omitempty"`
	// Units       []string         `json:"units" bson:"units,omitempty"`
	// COBs        []COBReference   `json:"cobs" bson:"cobs,omitempty"`
	Description string           `json:"description" bson:"description,omitempty"`
	IsActivated bool             `json:"is_activated" bson:"is_activated"`
	Users       []*UserReference `json:"users" bson:"users"`

	// Permission  Permission  `json:"permissions" bson:"permissions,omitempty"`
	// UserID      string      `json:"user_ids" bson:"user_ids"`
	// IsAdmin     bool        `json:"is_admin" bson:"is_admin,omitempty"`

	// RELATION

	// INCLUDE
	// Timestamp
}

type UserReference struct {
	ID       string   `json:"id,omitempty" bson:"id,omitempty"`
	FullName string   `json:"full_name,omitempty" bson:"full_name,omitempty"`
	Role     []string `json:"role,omitempty" bson:"role,omitempty"`
}

func (auth *AuthUser) ToUserReference() (result *UserReference) {
	result = &UserReference{
		ID:       auth.ID,
		FullName: auth.FullName,
		Role:     auth.Role,
	}
	return
}
