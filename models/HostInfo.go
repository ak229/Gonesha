package models
// --- hostInfo ---
type HostInfo struct {

        Ip string `json:"ip"`
        Tag string `json:"tag"`
	Module []string  `json:"module"`
	User string `json:"user"`
	Password string `json:"password"`
	Test []string `json:"test"`
}
