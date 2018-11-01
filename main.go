package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/writeameer/aci/helpers"
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

func deployArmTemplate(sid string, authorizer autorest.Authorizer, groupName string, location string) (err error) {
	// Setup ARM Client
	armClient := resources.NewGroupsClient(sid)
	armClient.Authorizer = authorizer

	// Create ARM group
	params := resources.Group{
		Location: &location,
	}
	group, _ := armClient.CreateOrUpdate(context.Background(), groupName, params)
	log.Printf("%v arm group created", group)

	// Create deployment client
	dClient := resources.NewDeploymentsClient(sid)
	dClient.Authorizer = authorizer

	template, _ := ioutil.ReadFile("./template/azuredeploy.json")
	templateParameters, _ := ioutil.ReadFile("./azuredeploy.parameters.json")

	// Define ARM deployment template and params
	deployment := resources.Deployment{
		Properties: &resources.DeploymentProperties{
			Template:   string(template),
			Parameters: templateParameters,
			Mode:       resources.Incremental,
		},
	}

	// Create ARM template deployment
	_, err = dClient.CreateOrUpdate(context.Background(), groupName, "ACIDeployment", deployment)
	helpers.PrintError(err)

	return err
}
