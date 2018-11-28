package apps

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/containerinstance/mgmt/containerinstance"
	"github.com/writeameer/aci/azure"
	"github.com/writeameer/aci/helpers"
)

// RunMoodle Runs wordpress on ACI
func RunMoodle(resourceGroupName string, containerGroupName string) (err error) {

	//	wp_config_extra := "\n\tdefine('WP_HOME', 'http://localhost:8000/');\n\tdefine('WP_SITEURL', 'http://localhost:8000/');\n"

	// Define containers to run
	containerSpecs := []azure.ContainerSpec{
		azure.ContainerSpec{
			ContainerName:  "moodle",
			ContainerImage: "bitnami/moodle:latest",
			Ports:          []int32{80},
			CPU:            0.5,
			MemoryInGB:     0.5,
			EnvironmentVariables: map[string]string{
				"MARIADB_HOST":         "127.0.0.1",
				"MARIADB_PORT_NUMBER":  "3306",
				"MOODLE_DATABASE_USER": "bn_moodle",
				"MOODLE_DATABASE_NAME": "bitnami_moodle",
				"ALLOW_EMPTY_PASSWORD": "yes",
			},
		},
		azure.ContainerSpec{
			ContainerName:  "mariahdb",
			ContainerImage: "bitnami/mariadb:latest",
			Ports:          []int32{3306},
			CPU:            0.5,
			MemoryInGB:     0.5,
			EnvironmentVariables: map[string]string{
				"MARIADB_USER":         "bn_moodle",
				"MARIADB_DATABASE":     "bitnami_moodle",
				"ALLOW_EMPTY_PASSWORD": "yes",
			},
		},
	}

	// Define the container group's specifications
	containerGroupSpecs := azure.ContainerGroupSpec{
		ResourceGroupName: resourceGroupName,
		Name:              containerGroupName,
		Ports:             []int32{80},
		DNSNameLabel:      "hiberapp",
		OsType:            containerinstance.Linux,
		IPAddressType:     containerinstance.Public,
	}

	// Get ARM group to inspect location
	armGroup, err := azure.GetGroup(resourceGroupName)
	helpers.PrintError(err)

	deployedGroup, err := azure.DeployContainer(*armGroup.Location, resourceGroupName, containerGroupName, containerSpecs, containerGroupSpecs)
	log.Printf(*deployedGroup.IPAddress.Fqdn)
	helpers.PrintError(err)

	return
}
