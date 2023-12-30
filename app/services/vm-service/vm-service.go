package vm_service

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	ports_service "simple-hosting/compositor/app/services/ports-service"
	"simple-hosting/compositor/app/settings"
	tools_retry "simple-hosting/go-commons/tools/retry"
	tools_sequence "simple-hosting/go-commons/tools/sequence"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

var config = settings.GetServiceSettings()

func Init() *VmServiceContext {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	return &VmServiceContext{
		client: cli,
	}
}

func exposedPortsArrToMap(ports []string) map[nat.Port]struct{} {
	res := tools_sequence.ToMap(
		ports,
		func(p string) nat.Port {
			_, containerPort := findPortsInPortData(p)
			return nat.Port(strconv.Itoa(containerPort))
		},
		func(p string) struct{} {
			return struct{}{}
		},
	)
	return res
}

func findPortBeforeSlash(p string) int {
	portStr := strings.Split(p, "/")[0]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	return port
}

func findPortsInPortData(portData string) (int, int) { //f.e "25565/tcp" -> -1, 25565 "1234:5678/tcp" -> 1234, 5678
	ports := strings.Split(portData, ":")
	if len(ports) == 1 {
		return -1, findPortBeforeSlash(ports[0])
	}

	return findPortBeforeSlash(ports[0]), findPortBeforeSlash(ports[1])
}

func portsArrToPortBindings(ports []string) nat.PortMap {
	res := tools_sequence.ToMap(
		ports,
		func(p string) nat.Port {
			_, containerPort := findPortsInPortData(p)
			return nat.Port(strconv.Itoa(containerPort))
		},
		func(p string) []nat.PortBinding {
			hostPort, containerPort := findPortsInPortData(p)
			if hostPort != -1 {
				return []nat.PortBinding{{HostIP: "", HostPort: strconv.Itoa(hostPort)}}
			}
			red, err := ports_service.CreatePortRedirect(containerPort)
			if err != nil {
				panic(err)
			}
			return []nat.PortBinding{{HostIP: "", HostPort: strconv.Itoa(red.ExternalPort)}}
		},
	)
	return res
}

func (hypContext *VmServiceContext) pullImageFromOrigin(imageId string) error {
	cli := extractClient(hypContext)
	_, err := cli.ImagePull(context.Background(), imageId, types.ImagePullOptions{})
	return err
}

func (hypContext *VmServiceContext) GetAllVms() VmListAllResponse {
	cli := extractClient(hypContext)
	resp, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return VmListAllResponse{IsSuccess: false, Error: err}
	}
	units := tools_sequence.Mapper(resp, func(c types.Container) VmListUnit {
		return VmListUnit{Names: c.Names, Id: c.ID, State: c.State, Status: c.Status}
	})
	return VmListAllResponse{
		Vms:       units,
		IsSuccess: true,
		Error:     nil,
	}
}

func (hypContext *VmServiceContext) CreateVm(request VmCreateRequest) VmCreateResponse {
	cli := extractClient(hypContext)
	fmt.Printf("Client created\n")
	err := hypContext.pullImageFromOrigin(request.VmImage)
	if err != nil {
		fmt.Printf("Error when pulling image: %s\n", err)
		return VmCreateResponse{
			IsSuccess: false,
			Error:     err,
		}
	}

	delay, _ := time.ParseDuration("5s")
	resp, err := tools_retry.Retry(
		func() (container.CreateResponse, error) {
			return containerCreate(cli, request)
		},
		func(r container.CreateResponse, e error) bool {
			fmt.Printf("Error when creating container: %s\n", err)
			return r.ID != ""
		},
		config.Hypervisor.Services.ContainerCreateAttempts,
		delay,
		nil,
	)

	fmt.Printf("Error when creating container: %s\n", err)
	return VmCreateResponse{
		VmId:      resp.ID,
		IsSuccess: err == nil,
		Error:     err,
	}
}

func containerCreate(cli *client.Client, request VmCreateRequest) (container.CreateResponse, error) {
	portMap := portsArrToPortBindings(request.VmExposePorts)
	resp, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        request.VmImage,
			ExposedPorts: exposedPortsArrToMap(request.VmExposePorts),
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory:     request.VmAvailableRamBytes,
				MemorySwap: request.VmAvailableSwapBytes,
			},
			PortBindings: portMap,
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		request.VmName,
	)
	return resp, err
}

func getOutboundIP() string {
	return config.Hypervisor.Services.OutboundIP
}

func getHostPorts(cli *client.Client, vmId string) ([]string, error) {
	containerInfo, err := cli.ContainerInspect(context.Background(), vmId)
	if err != nil {
		return nil, err
	}
	l := list.New()
	for _, bindings := range containerInfo.HostConfig.PortBindings {
		for _, binding := range bindings {
			l.PushBack(binding.HostPort)
		}
	}
	res := make([]string, l.Len())
	for p, i := l.Front(), 0; p != nil; p, i = p.Next(), i+1 {
		v, ok := p.Value.(string)
		if ok {
			res[i] = v
		} else {
			return nil, errors.New("port list contains no-port value")
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

	if err != nil {
		return VmRunResponse{Error: err}
	}

	ports, err := getHostPorts(cli, request.VmId)
	if err != nil {
		return VmRunResponse{
			VmId:  request.VmId,
			Error: err,
		}
	}

	return VmRunResponse{
		VmId:          request.VmId,
		HostIp:        getOutboundIP(),
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
	con, err := cli.ContainerInspect(context.Background(), request.VmId)

	if err != nil {
		return VmDeleteResponse{
			VmId:      request.VmId,
			IsSuccess: err == nil,
			Error:     err,
		}
	}

	err = cli.ContainerRemove(context.Background(), request.VmId, types.ContainerRemoveOptions{})

	if err != nil {
		return VmDeleteResponse{
			VmId:      request.VmId,
			IsSuccess: err == nil,
			Error:     err,
		}
	}

	bind := con.HostConfig.PortBindings
	for conPort, hostPorts := range bind {
		for _, hostPort := range hostPorts {
			_, portInData := findPortsInPortData(hostPort.HostPort)
			if err = ports_service.ClosePortRedirect(ports_service.PortRedirect{
				InternalPort: conPort.Int(),
				ExternalPort: portInData,
			}); err != nil {
				return VmDeleteResponse{
					VmId:      request.VmId,
					IsSuccess: err == nil,
					Error:     err,
				}
			}
		}
	}

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
