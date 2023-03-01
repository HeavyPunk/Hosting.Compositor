package vm_service

import "github.com/docker/docker/client"

type VmServiceContext struct {
	client *client.Client
}

type VmCreateRequest struct {
	VmName               string
	VmImage              string
	VmAvailableDiskBytes int64
	VmAvailableRamBytes  int64
	VmAvailableSwapBytes int64
	VmExposePorts        []string //should has format like "25565/tcp" or "19132/udp"
}

type VmCreateResponse struct {
	VmId      string
	IsSuccess bool
	Error     error
}

type VmRunRequest struct {
	VmId string
}

type VmRunResponse struct {
	VmId          string
	Error         error
	HostIp        string
	ExternalPorts []string
}

type VmStopRequest struct {
	VmId string
}

type VmStopResponse struct {
	VmId      string
	IsSuccess bool
	Error     error
}

type VmDeleteResponse struct {
	VmId      string
	IsSuccess bool
	Error     error
}

type VmDeleteRequest struct {
	VmId string
}
