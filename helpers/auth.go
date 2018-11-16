package helpers

import (
	"log"
	"os"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// AzureAuth Checks creds are provided in the ENV and returns an Azure token and Subscription ID
func AzureAuth() (authorizer autorest.Authorizer, sid string) {
	// Check env for creds and read env
	FatalError(CheckEnv())
	sid = os.Getenv("AZURE_SUBSCRIPTION_ID")

	// Authenticate with Azure
	log.Println("Starting azure auth...")
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	FatalError(err)
	log.Println("After azure auth...")

	return
}
