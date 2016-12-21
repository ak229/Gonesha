package models
type ModuleData struct {

	Module string
	Status bool
	Response []string
}

type HostData struct{

	Ip string
	ModuleDataList []ModuleData
}
