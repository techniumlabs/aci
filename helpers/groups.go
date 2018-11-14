package helpers

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
)

// CreateARMGroup creates an Azure Resource Manager Group
func CreateARMGroup(groupName string, location string) (group resources.Group, err error) {

	// Get ARM group client
	client := initGroupsClient()

	// Create ARM group
	group, err = client.CreateOrUpdate(ctx, groupName,
		resources.Group{
			Location: &location,
		},
	)
	PrintError(err)

	return
}

// GetGroup Gets an Azure resource group by name
func GetGroup(resourceGroupName string) (group resources.Group, err error) {

	// Get ARM group client
	client := initGroupsClient()

	group, err = client.Get(ctx, resourceGroupName)
	return
}

func initGroupsClient() (client resources.GroupsClient) {
	// Authenticate with Azure
	authorizer, sid := AzureAuth()

	// Setup ARM Client
	client = resources.NewGroupsClient(sid)
	client.Authorizer = authorizer

	return
}
