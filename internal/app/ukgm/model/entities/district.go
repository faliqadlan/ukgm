package entities

import (
	"log"
	"ukgm/internal/app/ukgm"
	"ukgm/internal/app/utils"
)

type (
	District struct {
		ID     interface{}        `json:"id" bson:"_id,omitempty"`
		Name   string             `json:"name" bson:"name"`
		Units  map[string]*Unit   `json:"units" bson:"units"`
		Doctor map[string]*Doctor `json:"Doctor" bson:"Doctor"`
		Nurse  map[string]*Nurse  `json:"Nurse" bson:"Nurse"`

		// include
		Timestamp
	}

	DistrictList struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	DistrictResponse struct {
		ID                interface{} `json:"_id"`
		Name              string      `json:"name"`
		Units             []*Unit     `json:"units"`
		Doctor            []*Doctor   `json:"Doctor"`
		Nurse             []*Nurse    `json:"Nurse"`
		ConditionOverview `json:"condition_overview"`

		// include
		Timestamp
	}

	ConditionOverview struct {
		Baik       int `json:"baik"`
		CukupBaik  int `json:"cukup_baik"`
		KurangBaik int `json:"kurang_baik"`
		Buruk      int `json:"buruk"`
	}
)

func (d *District) ToDistrictList() *DistrictList {
	return &DistrictList{
		ID:   utils.ObjectIDHex(d.ID),
		Name: d.Name,
	}
}

func (d *District) ToDistrictResponse() *DistrictResponse {
	res := &DistrictResponse{}
	res.ID = utils.ObjectIDHex(d.ID)
	res.Name = d.Name

	if d.Units != nil {
		for _, v := range d.Units {
			if v != nil {
				res.Units = append(res.Units, v)
			}

			if v.Patients != nil {
				for _, pat := range v.Patients {
					if pat != nil {
						if pat.LastDiagnose != nil {

							switch pat.LastDiagnose.Condition {
							case ukgm.KondisiBaik:
								res.ConditionOverview.Baik += 1
							case ukgm.KondisiCukupBaik:
								res.ConditionOverview.CukupBaik += 1
							case ukgm.KondisiKurangBaik:
								res.ConditionOverview.KurangBaik += 1
							case ukgm.KondisiBuruk:
								res.ConditionOverview.Buruk += 1
							}

						}
					}
				}
			}
		}
	}

	if d.Doctor != nil {
		for _, v := range d.Doctor {
			if v != nil {
				res.Doctor = append(res.Doctor, v)
			}
		}
	}

	if d.Nurse != nil {
		for _, v := range d.Nurse {
			res.Nurse = append(res.Nurse, v)
		}
	}

	res.Timestamp = d.Timestamp
	return res
}

func (d *District) ToDistrictResponseOld() (res *DistrictResponse) {
	log.Println(d)
	log.Println(d.ID)
	res = &DistrictResponse{}
	res.ID = utils.ObjectIDHex(d.ID)
	res.Name = d.Name

	if d.Units != nil {
		for _, v := range d.Units {
			if v != nil {
				res.Units = append(res.Units, v)
			}

			if v.Patients != nil {
				for _, pat := range v.Patients {
					if pat != nil {
						if pat.LastDiagnose != nil {

							switch pat.LastDiagnose.Condition {
							case ukgm.KondisiBaik:
								res.ConditionOverview.Baik += 1
							case ukgm.KondisiCukupBaik:
								res.ConditionOverview.CukupBaik += 1
							case ukgm.KondisiKurangBaik:
								res.ConditionOverview.KurangBaik += 1
							case ukgm.KondisiBuruk:
								res.ConditionOverview.Buruk += 1
							}

						}
					}
				}
			}
		}
	}

	if d.Doctor != nil {
		for _, v := range d.Doctor {
			if v != nil {
				res.Doctor = append(res.Doctor, v)
			}
		}
	}

	if d.Nurse != nil {
		for _, v := range d.Nurse {
			res.Nurse = append(res.Nurse, v)
		}
	}

	res.Timestamp = d.Timestamp

	return
}
