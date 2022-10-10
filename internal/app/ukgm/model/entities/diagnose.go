package entities

import (
	"ukgm/internal/app/ukgm"
	"ukgm/internal/app/utils"
)

type Diagnose struct {
	ID          interface{} `json:"id" bson:"_id,omitempty"`
	PatientID   string      `json:"patient_id" bson:"patient_id,omitempty"`
	UnitID      string      `json:"unit_id" bson:"unit_id,omitempty"`
	DistrictID  string      `json:"district_id" bson:"district_id,omitempty"`
	KondisiGigi [8][]string `json:"kondisi_gigi" bson:"kondisi_gigi"`
	KarangGigi  KarangGigi  `json:"karang_gigi" bson:"karang_gigi"`
	Note        string      `json:"note" bson:"note"`
	Condition   string      `json:"condition" bson:"condition"`
	DoctorName  string      `json:"doctor_name" bson:"doctor_name"`
	NurseName   string      `json:"nurse_name" bson:"nurse_name"`

	// include
	Timestamp
}

type KarangGigi struct {
	KiriAtas   bool `json:"kiri_atas" bson:"kiri_atas"`
	KiriBawah  bool `json:"kiri_bawah" bson:"kiri_bawah"`
	KananAtas  bool `json:"kanan_atas" bson:"kanan_atas"`
	KananBawah bool `json:"kanan_bawah" bson:"kanan_bawah"`
	DepanAtas  bool `json:"depan_atas" bson:"depan_atas"`
	DepanBawah bool `json:"depan_bawah" bson:"depan_bawah"`
}

func (d *Diagnose) ValdiateCondition() error {
	return utils.ValidateInArray(true, ukgm.KondisiGigi[:], "kondisi diagnosa", d.Condition)
}

type DiagnoseProgress struct {
	Date string `json:"date"`
	ConditionOverview
}
