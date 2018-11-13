package helpers

import "github.com/Azure/azure-sdk-for-go/profiles/preview/containerinstance/mgmt/containerinstance"

// GetContainerStruct returns a container struct with provided config
func GetContainerStruct(containerSpec ContainerSpec) (container containerinstance.Container) {

	// Define container's properties
	containerProperties := containerinstance.ContainerProperties{
		Image: &containerSpec.ContainerImage,
		Ports: setTCPPort(containerSpec.ports),
		Resources: &containerinstance.ResourceRequirements{
			Requests: setResourceRequests(containerSpec.cpu, containerSpec.memoryInGB),
		},
	}

	// Define a container with given properties
	container = containerinstance.Container{
		ContainerProperties: &containerProperties,
		Name:                &containerSpec.ContainerName,
	}

	// return containers
	return
}

func setTCPPort(ports []int32) (containerPorts *[]containerinstance.ContainerPort) {

	var portList []containerinstance.ContainerPort

	for i, port := range ports {
		portList[i] = containerinstance.ContainerPort{
			Port:     &port,
			Protocol: "tcp",
		}
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
	ContainerName  string
	ports          []int32
	ContainerImage string
	cpu            float64
	memoryInGB     float64
}
