package controller_server

import (
	"net/http"
	vm_service "simple-hosting/compositor/app/services/vm-service"

	"github.com/gin-gonic/gin"
)

func GetServersList(c *gin.Context) {
	c.JSON(http.StatusOK, "{}") //TODO: implement
}

func CreateServer(c *gin.Context) {
	var request CreateServerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
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
	result := CreateServerResponse{
		VmId: response.VmId,
	}
	c.JSON(http.StatusCreated, result)
}

func StartServer(c *gin.Context) {
	var request StartServerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
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
		Success: resp.VmId != "",
	}
	c.JSON(http.StatusOK, result)
}
