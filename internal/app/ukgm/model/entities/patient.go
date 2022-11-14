package entities

import (
	"time"
	"ukgm/internal/app/ukgm"
	"ukgm/internal/app/utils"
)

type Patient struct {
	ID           interface{} `json:"id" bson:"_id,omitempty"`
	UnitID       string      `json:"unit_id" bson:"unit_id,omitempty"`
	DistrictID   string      `json:"district_id" bson:"district_id,omitempty"`
	Name         string      `json:"name" bson:"name" part:"main" properties:"searchable"`
	NameUnit     string      `json:"name_unit" bson:"name_unit"`
	MrNumber     string      `json:"mr_number" bson:"mr_number"`
	Dob          *time.Time  `json:"dob" bson:"dob"`
	Gender       string      `json:"gender" bson:"gender"`
	NoHp         string      `json:"no_hp" bson:"no_hp"`
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
	NoHp         string `json:"no_hp"`
	LastDiagnose string `json:"last_diagnose"`
	Condition    string `json:"condition"`
}

func (p *Patient) ToPatientList() *PatientList {

	var dobS, ldS, cond string

	if p.Dob != nil {
		dobS = p.Dob.Local().Format(ukgm.TimeFormat)
	}

	if p.LastDiagnose != nil {
		ldS = p.LastDiagnose.CreatedAt.Local().Format(ukgm.TimeFormat)
	}

	if p.LastDiagnose != nil {
		cond = p.LastDiagnose.Condition
	}

	return &PatientList{
		Id:           utils.ObjectIDHex(p.ID),
		MrNumber:     p.MrNumber,
		NamePatient:  p.Name,
		NameUnit:     p.NameUnit,
		Dob:          dobS,
		LastDiagnose: ldS,
		Condition:    cond,
		NoHp:         p.NoHp,
	}
}

func (p *Patient) GetSearchableField() (result []string) {
	if p == nil {
		p = &Patient{}
	}
	result, err := utils.GetSearchableField(p)
	if err != nil {
		result = []string{}
	}
	return
}
