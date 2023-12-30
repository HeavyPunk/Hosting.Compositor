package controller_server

import (
	"fmt"
	"net/http"
	vm_service "simple-hosting/compositor/app/services/vm-service"
	error_utils "simple-hosting/go-commons/tools/errors"
	tools_sequence "simple-hosting/go-commons/tools/sequence"

	"github.com/gin-gonic/gin"
)

func GetServersList(c *gin.Context) {
	vmService := vm_service.Init()
	allVms := vmService.GetAllVms()
	resp := GetAllResponse{
		IsSuccess: allVms.IsSuccess,
		Error:     error_utils.GetErrorStringOrDefault(allVms.Error, ""),
	}
	resp.Vms = tools_sequence.Mapper(allVms.Vms, func(u vm_service.VmListUnit) vmListUnit {
		return vmListUnit{Names: u.Names, Id: u.Id, State: u.State, Status: u.Status}
	})
	c.JSON(http.StatusOK, resp)
}

func CreateServer(c *gin.Context) {
	var request CreateServerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Printf("Error binding request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vmService := vm_service.Init()
	response := vmService.CreateVm(vm_service.VmCreateRequest{
		VmName:               request.VmName,
		VmImage:              request.VmImageUri,
		VmExposePorts:        request.VmExposePorts,
		VmAvailableDiskBytes: request.VmAvailableDiskBytes,
		VmAvailableRamBytes:  request.VmAvailableRamBytes,
		VmAvailableSwapBytes: request.VmAvailableSwapBytes,
	})
	if !response.IsSuccess {
		fmt.Printf("Error creating vm: %v", error_utils.GetErrorStringOrDefault(response.Error, ""))
	}
	result := CreateServerResponse{
		VmId:    response.VmId,
		Success: response.IsSuccess,
		Error:   error_utils.GetErrorStringOrDefault(response.Error, ""),
	}
	if result.Success {
		c.JSON(http.StatusCreated, result)
	} else {
		c.JSON(http.StatusInternalServerError, result)
	}
}

func StartServer(c *gin.Context) {
	var request StartServerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Printf("Error binding request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vmService := vm_service.Init()
	resp := vmService.RunVm(vm_service.VmRunRequest{
		VmId: request.VmId,
	})
	result := StartServerResponse{
		VmId:         resp.VmId,
		VmWhiteIp:    resp.HostIp,
		VmWhitePorts: resp.ExternalPorts,
		Success:      resp.Error == nil,
		Error:        error_utils.GetErrorStringOrDefault(resp.Error, ""),
	}
	if !result.Success {
		fmt.Printf("Error creating vm: %v", result.Error)
	}

	c.JSON(http.StatusOK, result)
}

func StopServer(c *gin.Context) {
	var request StopServerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vmService := vm_service.Init()
	resp := vmService.StopVm(vm_service.VmStopRequest{
		VmId: request.VmId,
	})
	result := StopServerResponse{
		Success: resp.IsSuccess,
		Error:   error_utils.GetErrorStringOrDefault(resp.Error, ""),
	}
	c.JSON(http.StatusOK, result)
}

func DeleteServer(c *gin.Context) {
	var request RemoveServerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vmService := vm_service.Init()
	resp := vmService.DeleteVm(vm_service.VmDeleteRequest{
		VmId: request.VmId,
	})
	result := RemoveServerResponse{
		Success: resp.IsSuccess,
		Error:   error_utils.GetErrorStringOrDefault(resp.Error, ""),
	}
	c.JSON(http.StatusOK, result)
}
