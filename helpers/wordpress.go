package helpers

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/containerinstance/mgmt/containerinstance"
)

func RunWordPress() {
	// Define containers to run
	containerSpecs := []ContainerSpec{
		ContainerSpec{
			ContainerName:  "wordpress",
			ContainerImage: "wordpress",
			Ports:          []int32{80},
			CPU:            0.5,
			MemoryInGB:     0.5,
		},
		ContainerSpec{
			ContainerName:  "mysql",
			ContainerImage: "mysql",
			Ports:          []int32{3306},
			CPU:            0.5,
			MemoryInGB:     0.5,
			EnvironmentVariables: map[string]string{
				"MYSQL_ROOT_PASSWORD": "0rsmP@ssw0rd",
			},
		},
	}

	// Define the container group's specifications
	containerGroupSpecs := ContainerGroupSpec{
		ResourceGroupName: "aci-example",
		Name:              "wordpress",
		Ports:             []int32{80},
		DNSNameLabel:      "hiberapp",
		OsType:            containerinstance.Linux,
		IPAddressType:     containerinstance.Public,
	}

	deployedGroup, err := DeployContainer2("East US", "aci-example", "wordpress-app", containerSpecs, containerGroupSpecs)
	log.Printf(*deployedGroup.IPAddress.Fqdn)
	if err != nil {
		log.Printf(err.Error())
	}
}
