package vm_service

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
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
	cli := extractClient(hypContext)
	resp, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        request.VmImage,
			ExposedPorts: nat.PortSet{},
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory:     request.VmAvailableRamBytes,
				MemorySwap: request.VmAvailableSwapBytes,
			}},
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

func (hypContext *VmServiceContext) RunVm(request VmRunRequest) VmRunResponse {
	cli := extractClient(hypContext)
	err := cli.ContainerStart(
		context.Background(),
		request.VmId,
		types.ContainerStartOptions{},
	)
	return VmRunResponse{
		VmId:  request.VmId,
		Error: err,
	}
}

func (hypContext *VmServiceContext) StopVm(request VmStopRequest) VmStopResponse {
	cli := extractClient(hypContext)
	err := cli.ContainerStop(context.Background(), request.VmId, container.StopOptions{})
	return VmStopResponse{
		VmId:      request.VmId,
		IsSuccess: err == nil,
		Error:     err,
	}
}

func (hypContext *VmServiceContext) DeleteVm(request VmDeleteRequest) VmDeleteResponse {
	cli := extractClient(hypContext)
	err := cli.ContainerRemove(context.Background(), request.VmId, types.ContainerRemoveOptions{})
	return VmDeleteResponse{
		VmId:      request.VmId,
		IsSuccess: err == nil,
		Error:     err,
	}
}

func extractClient(hypContext *VmServiceContext) *client.Client {
	cli := hypContext.client
	if cli == nil {
		panic("Hypervisor client is nil")
	}
	return cli
}
