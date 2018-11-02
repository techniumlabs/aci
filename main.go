package main

import (
	"context"
	"log"
	"os"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/writeameer/aci/helpers"
)

var (
	ctx = context.Background()
)

func main() {

	// Check env for creds and read env
	helpers.FatalError(helpers.CheckEnv())
	sid := os.Getenv("AZURE_SUBSCRIPTION_ID")

	// Authenticate with Azure
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	helpers.FatalError(err)

	log.Println(reflect.TypeOf(authorizer))

	// Deploy ARM Template
	deployArmTemplate(sid, authorizer, "hiberapp", "Australia East")

}

func deployArmTemplate(sid string, authorizer autorest.Authorizer, groupName string, location string) (deployment resources.DeploymentExtended, err error) {
	// Setup ARM Client
	armClient := resources.NewGroupsClient(sid)
	armClient.Authorizer = authorizer

	// Create ARM group
	params := resources.Group{
		Location: &location,
	}
	group, _ := armClient.CreateOrUpdate(ctx, groupName, params)
	log.Printf("%v arm group created", group)

	// Create deployment client
	dClient := resources.NewDeploymentsClient(sid)
	dClient.Authorizer = authorizer

	template, err := helpers.ReadJSON("./template/azuredeploy.json")
	helpers.FatalError(err)

	templateParameters, _ := helpers.ReadJSON("./template/azuredeploy.parameters.json")
	helpers.FatalError(err)

	// Deploy ARM template deployment
	deploymentFuture, err := dClient.CreateOrUpdate(
		ctx,
		groupName,
		"ACIDeployment",
		resources.Deployment{
			Properties: &resources.DeploymentProperties{
				Template:   template,
				Parameters: templateParameters,
				Mode:       resources.Incremental,
			},
		},
	)
	helpers.PrintError(err)

	// Wait for completion
	err = deploymentFuture.Future.WaitForCompletion(ctx, dClient.BaseClient.Client)
	if err != nil {
		return
	}
	return deploymentFuture.Result(dClient)
}
