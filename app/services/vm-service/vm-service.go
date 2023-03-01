package vm_service

import (
	"container/list"
	"context"
	"errors"
	"log"
	"net"

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

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func getHostPorts(cli *client.Client, vmId string) ([]string, error) {
	containerInfo, err := cli.ContainerInspect(context.Background(), vmId)
	if err != nil {
		return nil, err
	}
	l := list.New()
	for hostPort := range containerInfo.NetworkSettings.Ports {
		l.PushBack(hostPort.Port())
	}
	res := make([]string, l.Len())
	for p, i := l.Front(), 0; p != nil; p, i = p.Next(), i+1 {
		v, ok := p.Value.(string)
		if ok {
			res[i] = v
		} else {
			return nil, errors.New("Port list contains no-port value")
		}
	}
	return res, nil
}

func (hypContext *VmServiceContext) RunVm(request VmRunRequest) VmRunResponse {
	cli := extractClient(hypContext)
	err := cli.ContainerStart(
		context.Background(),
		request.VmId,
		types.ContainerStartOptions{},
	)

	ports, err := getHostPorts(cli, request.VmId)
	if err != nil {
		return VmRunResponse{
			VmId:  request.VmId,
			Error: err,
		}
	}

	return VmRunResponse{
		VmId:          request.VmId,
		HostIp:        getOutboundIP().String(),
		ExternalPorts: ports,
		Error:         err,
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
