package main

import (
	"log"
	"net/http"

	helpers "github.com/writeameer/aci/helpers"
	handlers "github.com/writeameer/httphandlers/handlers"
)

func main() {

	// Deploy App
	siteFQDN, err := deployACIApp()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(siteFQDN)

	// Create new router
	mux := http.NewServeMux()
	originHost := siteFQDN

	mux.Handle("/", handlers.ReverseProxyHandler(originHost))

	// Listen and Serve
	port := ":8080"
	log.Println("Server started on port" + port)
	log.Fatal(http.ListenAndServe(port, mux))

}

func deployACIApp() (siteFQDN string, err error) {
	// Get ARM template and params
	template, err := helpers.ReadJSON("./template/azuredeploy.json")
	helpers.PrintError(err)
	templateParameters, _ := helpers.ReadJSON("./template/azuredeploy.parameters.json")
	helpers.PrintError(err)

	// Deploy ARM Template
	log.Printf("Starting deployment...")

	groupName := "hiberapp"
	location := "Australia East"
	deploymentName := "ACIDeployment"

	// Get Deployment result
	result, err := helpers.DeployArmTemplate(groupName, location, deploymentName, template, templateParameters)
	if err != nil {
		log.Printf("Error %v", err)
	}

	// Parse Output
	properties := result.Properties.Outputs.(map[string]interface{})
	propInfo := properties["siteFQDN"].(map[string]interface{})
	siteFQDN = propInfo["value"].(string)

	return
}
