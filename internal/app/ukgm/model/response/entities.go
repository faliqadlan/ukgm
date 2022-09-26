package response

type District struct {
}

type ConditionOverview struct {
	Baik       int `json:"baik"`
	CukupBaik  int `json:"cukup_baik"`
	KurangBaik int `json:"kurang_baik"`
	Buruk      int `json:"buruk"`
}
