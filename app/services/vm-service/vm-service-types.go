package vm_service

import "github.com/docker/docker/client"

type VmServiceContext struct {
	client *client.Client
}

type VmCreateRequest struct {
	VmName        string
	VmImage       string
	VmExposePorts []int
}

type VmCreateResponse struct {
	VmId      string
	IsSuccess bool
	Error     error
}

type VmRunRequest struct {
	VmId          string
	AvailableRam  uint
	AvailableSwap uint
}

type VmRunResponse struct {
}

type VmStopRequest struct {
	VmId string
}

type VmStopResponse struct {
	VmId      string
	IsSuccess bool
	Error     error
}

type VmSuspendResponse struct {
	VmId      string
	IsSuccess bool
	Error     error
}

type VmSuspendRequest struct {
	VmId string
}
