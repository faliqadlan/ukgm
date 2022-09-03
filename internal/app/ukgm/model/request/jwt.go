package request

import (
	"ukgm/internal/app/ukgm/model/entities"

	"github.com/gbrlsnchs/jwt/v2"
)

type Token struct {
	*jwt.JWT
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (t *Token) ToAuth() entities.AuthUser {
	return entities.AuthUser{
		ID:   t.ID,
		Name: t.Name,
	}
}
