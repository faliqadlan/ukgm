package entities

import (
	"log"
	"ukgm/internal/app/config/env"
	"ukgm/internal/app/ukgm"
	"ukgm/internal/app/utils"

	"firebase.google.com/go/auth"
)

type Officer struct {
	ID           interface{} `json:"id" bson:"_id,omitempty"`
	DistrictID   string      `json:"district_id" bson:"district_id,omitempty"`
	Name         string      `json:"name" bson:"name" part:"main" properties:"searchable"`
	Email        string      `json:"email" bson:"email"`
	Level        string      `json:"level" bson:"level"`
	Profesion    string      `json:"profesion" bson:"profesion"`
	DistrictName string      `json:"district_name" bson:"districtname"`

	// include
	Timestamp
}

type OfficerList struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Profesion    string `json:"profesion"`
	DistrictName string `json:"district_name"`
}

func (o *Officer) ToOfficerList() OfficerList {
	return OfficerList{
		Name:         o.Name,
		ID:           utils.ObjectIDHex(o.ID),
		Profesion:    o.Profesion,
		DistrictName: o.DistrictName,
	}
}

func (o *Officer) ToCreateUser(auth *auth.Token) (u User) {
	u.Name = o.Name
	u.OfficerId = utils.ObjectIDHex(o.ID)
	u.Email = o.Email
	u.DistrictID = o.DistrictID

	if o.DistrictID == "kecamatan" {
		distrctList := []DistrictList{}

		for i := 0; i < len(ukgm.DistrictNameList); i++ {

			var id string

			switch {
			case env.Get("DB_DATABASE") == "ukgm_dev":

				id = ukgm.DistrictIdListDev[i]

			case env.Get("DB_DATABASE") == "ukgm":
				id = ukgm.DistrictIdListLive[i]
			default:
				log.Println("error get id")
				return
			}

			distrctList = append(distrctList, DistrictList{
				ID:   id,
				Name: ukgm.DistrictNameList[i],
			})
		}

		u.ListDistrict = distrctList
	}

	u.Timestamp.CreateNew(auth)

	return
}
