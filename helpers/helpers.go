package helpers

import (
	"errors"
	"log"
	"os"
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

func PrintError(err error) {
	if err != nil {
		log.Printf(err.Error())
	}
}
