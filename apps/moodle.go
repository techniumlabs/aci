package apps

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/containerinstance/mgmt/containerinstance"
	"github.com/writeameer/aci/azure"
	"github.com/writeameer/aci/helpers"
)

// RunMoodle Runs wordpress on ACI
func RunMoodle(resourceGroupName string, containerGroupName string, storageAccountName string, mysqlShareName string, webShareName string) (err error) {

	// Create a storage account if not exists
	azure.CreateStorageAccount(resourceGroupName, storageAccountName)

	// Creat file share if not exists
	key, err := azure.CreateAzureFileShare(resourceGroupName, storageAccountName, mysqlShareName)
	if err != nil {
		return
	}

	key, err = azure.CreateAzureFileShare(resourceGroupName, storageAccountName, webShareName)
	if err != nil {
		return
	}
	log.Println("After fileshare creation")
	log.Println("The key was:" + key)

	//	wp_config_extra := "\n\tdefine('WP_HOME', 'http://localhost:8000/');\n\tdefine('WP_SITEURL', 'http://localhost:8000/');\n"

	// Define containers to run
	readonly := false
	mountPath1 := "/bitnami"
	mountPath2 := "/bitnami/mariadb"

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
			VolumeMounts: &[]containerinstance.VolumeMount{
				containerinstance.VolumeMount{
					Name:      &webShareName,
					MountPath: &mountPath1,
					ReadOnly:  &readonly,
				},
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
			VolumeMounts: &[]containerinstance.VolumeMount{
				containerinstance.VolumeMount{
					Name:      &mysqlShareName,
					MountPath: &mountPath2,
					ReadOnly:  &readonly,
				},
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
		Volumes: &[]containerinstance.Volume{
			containerinstance.Volume{
				Name: &webShareName,
				AzureFile: &containerinstance.AzureFileVolume{
					ShareName:          &webShareName,
					StorageAccountKey:  &key,
					StorageAccountName: &storageAccountName,
				},
			},
			containerinstance.Volume{
				Name: &mysqlShareName,
				AzureFile: &containerinstance.AzureFileVolume{
					ShareName:          &mysqlShareName,
					StorageAccountKey:  &key,
					StorageAccountName: &storageAccountName,
				},
			},
		},
	}

	// Get ARM group to inspect location
	armGroup, err := azure.GetGroup(resourceGroupName)
	helpers.PrintError(err)

	deployedGroup, err := azure.DeployContainer(*armGroup.Location, resourceGroupName, containerGroupName, containerSpecs, containerGroupSpecs)
	log.Printf(*deployedGroup.IPAddress.Fqdn)
	helpers.PrintError(err)

	return
}
