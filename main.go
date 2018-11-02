package main

import (
	"log"
	"os"

	"github.com/Azure/go-autorest/autorest/azure/auth"
	helpers "github.com/writeameer/aci/helpers"
)

func main() {

	// Check env for creds and read env
	helpers.FatalError(helpers.CheckEnv())
	sid := os.Getenv("AZURE_SUBSCRIPTION_ID")

	// Authenticate with Azure
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	helpers.FatalError(err)

	// Get ARM template and params
	template, err := helpers.ReadJSON("./template/azuredeploy.json")
	helpers.FatalError(err)

	templateParameters, _ := helpers.ReadJSON("./template/azuredeploy.parameters.json")
	helpers.FatalError(err)

	// Deploy ARM Template
	log.Printf("Starting deployment...")
	helpers.DeployArmTemplate(sid, authorizer, "hiberapp", "Australia East", template, templateParameters)

}
