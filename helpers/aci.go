package helpers

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/containerinstance/mgmt/containerinstance"
)

// GetContainerFromSpec returns a container struct with provided config
func GetContainerFromSpec(containerSpec ContainerSpec) (container containerinstance.Container) {

	// Define Env Variables
	var envVars []containerinstance.EnvironmentVariable
	for key, value := range containerSpec.EnvironmentVariables {
		envVars = append(envVars, containerinstance.EnvironmentVariable{
			Name:  &key,
			Value: &value,
		})
	}
	// Define container's properties
	containerProperties := containerinstance.ContainerProperties{
		Image: &containerSpec.ContainerImage,
		Ports: setTCPPort(containerSpec.Ports),
		Resources: &containerinstance.ResourceRequirements{
			Requests: setResourceRequests(containerSpec.CPU, containerSpec.MemoryInGB),
		},
		EnvironmentVariables: &envVars,
	}

	// Define a container with given properties
	container = containerinstance.Container{
		ContainerProperties: &containerProperties,
		Name:                &containerSpec.ContainerName,
	}

	// return containers
	return
}

// GetContainersFromSpec returns an array of Container structs from the specs provided
func GetContainersFromSpec(containerSpecs []ContainerSpec) (containers *[]containerinstance.Container) {
	var myContainerSpecs []containerinstance.Container
	for _, containerSpec := range containerSpecs {
		myContainerSpecs = append(myContainerSpecs, GetContainerFromSpec(containerSpec))
	}

	return &myContainerSpecs
}

// GetContainerGroupFromSpec converts a ContainerGroupSpec and ContainerSpec struct to a containerinstance.ContainerGroupProperties struct
func GetContainerGroupFromSpec(containerGroupSpec ContainerGroupSpec, containerSpecs []ContainerSpec) (containerGroup *containerinstance.ContainerGroupProperties) {
	log.Println("Starting GetContainerGroupFromSpec...")

	cgroup := containerinstance.ContainerGroupProperties{
		Containers:    GetContainersFromSpec(containerSpecs),
		OsType:        containerGroupSpec.OsType,
		RestartPolicy: containerGroupSpec.RestartPolicy,
		IPAddress: &containerinstance.IPAddress{
			Type:         containerinstance.Public,
			DNSNameLabel: &containerGroupSpec.DNSNameLabel,
			Ports:        setContainerGroupTCPPort(containerGroupSpec.Ports),
		},
	}

	return &cgroup
}

func setTCPPort(ports []int32) (containerPorts *[]containerinstance.ContainerPort) {
	log.Println("Starting setTCPPort...")
	var portList []containerinstance.ContainerPort

	for i, port := range ports {
		log.Printf("%d, %d", i, port)
		portList = append(portList, containerinstance.ContainerPort{
			Port:     &port,
			Protocol: "tcp",
		})
	}

	return &portList
}

func setContainerGroupTCPPort(ports []int32) (containerPorts *[]containerinstance.Port) {

	var portList []containerinstance.Port

	for i, port := range ports {
		log.Printf("%d, %d", i, port)
		portList = append(portList, containerinstance.Port{
			Port:     &port,
			Protocol: "tcp",
		})
	}

	return &portList
}

func setResourceRequests(cpu float64, memoryInGB float64) (resourceRequirements *containerinstance.ResourceRequests) {
	requirements := containerinstance.ResourceRequests{
		CPU:        &cpu,
		MemoryInGB: &memoryInGB,
	}

	return &requirements
}

// ContainerSpec defines the details of the container to launch
type ContainerSpec struct {
	ContainerName        string
	Ports                []int32
	ContainerImage       string
	CPU                  float64
	MemoryInGB           float64
	EnvironmentVariables map[string]string
	VolumeMount          AzureFileMount
}

// ContainerGroupSpec defines the details of the container to launch
type ContainerGroupSpec struct {
	ResourceGroupName string
	Name              string
	Ports             []int32
	DNSNameLabel      string
	OsType            containerinstance.OperatingSystemTypes
	RestartPolicy     containerinstance.ContainerGroupRestartPolicy
	IPAddressType     containerinstance.ContainerGroupIPAddressType
}

// AzureFileMount describes the Azure File Mount for a container
type AzureFileMount struct {
	ShareName          string
	StorageAccountKey  string
	StorageAccountName string
}
