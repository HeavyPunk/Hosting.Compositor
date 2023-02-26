package vm_service

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func Init() *VmServiceContext {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	return &VmServiceContext{
		client: cli,
	}
}

func (hypContext *VmServiceContext) CreateVm(request VmCreateRequest) VmCreateResponse {
	cli := hypContext.client
	if cli == nil {
		panic("Hypervisor client is nil")
	}
	resp, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: request.VmImage,
		},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		&v1.Platform{},
		request.VmName,
	)
	return VmCreateResponse{
		VmId:      resp.ID,
		IsSuccess: err == nil,
		Error:     err,
	}
}

func (hypContext *VmServiceContext) StopVm(request VmStopRequest) VmStopResponse {
	panic("Not implemented")
}

func (hypContext *VmServiceContext) SuspendVm(request VmSuspendRequest) VmSuspendResponse {
	panic("Not implemented")
}
