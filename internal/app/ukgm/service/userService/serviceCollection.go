package user_service

import "ukgm/internal/app/ukgm/model/entities"

type UserServiceColl interface {
	ShowUserByUid(uid string) (result *entities.User, err error)
}
