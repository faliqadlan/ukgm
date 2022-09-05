package districtservice

import "ukgm/internal/app/ukgm/model/entities"

type DistrictServiceColl interface {
	ShowDistrictbyID(id string) (result *entities.District, err error)
	GetDistrictList() (result []entities.DistrictList, err error)
}
