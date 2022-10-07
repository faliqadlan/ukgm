package request

import "fmt"

var (
	PageRequestCtxKey = &CtxKey{"PageRequest"}
	AuthCtxKey        = &CtxKey{"AuthUser"}
	UserCtxKey        = &CtxKey{"User"}
	DocumentCtxKey    = &CtxKey{"Document"}
	LimitCtxKey       = &CtxKey{"Limit"}
	LayerCtxKey       = &CtxKey{"Layer"}
	DistrictCtxKey    = &CtxKey{"District"}
	DiagnoseCtxKey    = &CtxKey{"Diagnose"}
)

type (
	CtxKey struct {
		name string
	}
)

func (k *CtxKey) String() string {
	return fmt.Sprintf("ukgm's context value: %s", k.name)
}

func GetDiagnoseCtxKey() *CtxKey {
	return DiagnoseCtxKey
}

func GetDistrictCtxKey() (result *CtxKey) {
	return DistrictCtxKey
}

func GetPageRequestCtxKey() (result *CtxKey) {
	result = PageRequestCtxKey
	return
}

func GetAuthCtxKey() (result *CtxKey) {
	result = AuthCtxKey
	return
}

func GetUserCtxKey() (result *CtxKey) {
	result = UserCtxKey
	return
}

func GetDocumentCtxKey() (result *CtxKey) {
	result = DocumentCtxKey
	return
}

func GetLimitCtxKey() (result *CtxKey) {
	result = LimitCtxKey
	return
}

func GetLayerCtxKey() (result *CtxKey) {
	result = LayerCtxKey
	return
}
