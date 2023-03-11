package controller_server

type CreateServerRequest struct {
	VmImageUri           string   `json:"image"`
	VmName               string   `json:"name"`
	VmAvailableRamBytes  int64    `json:"ram"`
	VmAvailableSwapBytes int64    `json:"swap"`
	VmAvailableDiskBytes int64    `json:"disk"`
	VmExposePorts        []string `json:"internal-ports"` //format: ["8888/tcp", "2222"/udp]
}

type CreateServerResponse struct {
	VmId string `json:"vm-id"`
}

type StartServerRequest struct {
	VmId string `json:"vm-id"`
}

type StartServerResponse struct {
	VmId         string   `json:"vm-id"`
	VmWhiteIp    string   `json:"ip"`
	VmWhitePorts []string `json:"ports"`
}

type StopServerRequest struct {
	VmId string `json:"vm-id"`
}

type StopServerResponse struct {
	Success bool `json:"success"`
}

type RemoveServerRequest struct {
	VmId string `json:"vm-id"`
}

type RemoveServerResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
