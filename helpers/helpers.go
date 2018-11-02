package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/Azure/go-autorest/autorest"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

var (
	ctx = context.Background()
)

// CheckEnv Check the Azure creds are set in the environment variables
func CheckEnv() (err error) {

	azureCreds := []string{
		"AZURE_TENANT_ID",
		"AZURE_CLIENT_ID",
		"AZURE_CLIENT_SECRET",
		"AZURE_SUBSCRIPTION_ID",
	}

	for _, cred := range azureCreds {
		if os.Getenv(cred) == "" {
			log.Printf("credential variable " + cred + " has not be set")
			err = errors.New("error, missing envrionment variables. run `az ad sp create-for-rbac -n \"<yourAppName>\"' to create a service principal and generate the necessary credential variables")
		}
	}

	return err
}

// FatalError Quits if error is fatal
func FatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//PrintError Prints if error
func PrintError(err error) {
	if err != nil {
		log.Printf(err.Error())
	}
}

// ReadJSON Reads json and returns a map
func ReadJSON(path string) (*map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	contents := make(map[string]interface{})
	json.Unmarshal(data, &contents)
	return &contents, nil
}

// DeployArmTemplate Deploys and ARM template
func DeployArmTemplate(groupName string, location string, deploymentName string, template *map[string]interface{}, paramaters *map[string]interface{}) (deployment resources.DeploymentExtended, err error) {

	// Authenticate with Azure
	authorizer, sid := AzureAuth()

	// Setup ARM Client
	armClient := resources.NewGroupsClient(sid)
	armClient.Authorizer = authorizer

	// Create ARM group
	params := resources.Group{
		Location: &location,
	}
	group, _ := armClient.CreateOrUpdate(ctx, groupName, params)
	log.Printf("%v arm group created", *group.Name)

	// Create deployment client
	dClient := resources.NewDeploymentsClient(sid)
	dClient.Authorizer = authorizer

	// Deploy ARM template deployment
	deploymentFuture, err := dClient.CreateOrUpdate(
		ctx,
		groupName,
		deploymentName,
		resources.Deployment{
			Properties: &resources.DeploymentProperties{
				Template:   template,
				Parameters: paramaters,
				Mode:       resources.Incremental,
			},
		},
	)

	PrintError(err)

	// Wait for completion
	log.Printf("Wait for completion...")
	err = deploymentFuture.Future.WaitForCompletion(ctx, dClient.BaseClient.Client)
	if err != nil {
		return
	}
	log.Printf("Deployment completed...")
	return deploymentFuture.Result(dClient)
}

// AzureAuth Checks creds are provided in the ENV and returns an Azure token and Subscription ID
func AzureAuth() (authorizer autorest.Authorizer, sid string) {
	// Check env for creds and read env
	FatalError(CheckEnv())
	sid = os.Getenv("AZURE_SUBSCRIPTION_ID")

	// Authenticate with Azure
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	FatalError(err)

	return
}
