package entities

import (
	"time"
	"ukgm/internal/app/ukgm"
	"ukgm/internal/app/utils"
)

type Patient struct {
	ID           interface{} `json:"id" bson:"_id,omitempty"`
	UnitID       string      `json:"unit_id" bson:"_unit_id,omitempty"`
	DistrictID   string      `json:"district_id" bson:"_district_id,omitempty"`
	Name         string      `json:"name" bson:"name" part:"main" properties:"searchable"`
	NameUnit     string      `json:"name_unit" bson:"name_unit"`
	MrNumber     string      `json:"mr_number" bson:"mr_number"`
	Dob          *time.Time  `json:"dob" bson:"dob"`
	Gender       string      `json:"gender" bson:"gender"`
	LastDiagnose *Diagnose   `json:"last_diagnose" bson:"last_diagnose"`

	// include
	Timestamp
}

type PatientList struct {
	Id           string `json:"id"`
	MrNumber     string `json:"mr_number"`
	NamePatient  string `json:"name_patient"`
	NameUnit     string `json:"name_unit"`
	Dob          string `json:"dob"`
	LastDiagnose string `json:"last_diagnose"`
	Condition    string `json:"condition"`
}

func (p *Patient) ToPatientList() *PatientList {
	return &PatientList{
		Id:           utils.ObjectIDHex(p.ID),
		MrNumber:     p.MrNumber,
		NamePatient:  p.Name,
		NameUnit:     p.NameUnit,
		Dob:          p.Dob.Local().Format(ukgm.TimeFormat),
		LastDiagnose: p.LastDiagnose.UpdatedAt.Local().Format(ukgm.TimeFormat),
		Condition:    p.LastDiagnose.Condition,
	}
}
