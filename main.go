package main

import (
	"log"
	"net/http"

	helpers "github.com/writeameer/aci/helpers"
	handlers "github.com/writeameer/httphandlers/handlers"
)

func main() {

	// Deploy ACI and get siteFQDN
	siteFQDN, err := deployACIApp(helpers.ArmDeploymentRequest{
		Template:       "./templates/example/azuredeploy.json",
		Parameters:     "./templates/example/azuredeploy.parameters.json",
		GroupName:      "hiberapp",
		Location:       "Australia East",
		DeploymentName: "deploymentName",
	})

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

func deployACIApp(request helpers.ArmDeploymentRequest) (siteFQDN string, err error) {
	// Get ARM template and params
	template, err := helpers.ReadJSON(request.Template)
	helpers.PrintError(err)
	templateParameters, _ := helpers.ReadJSON(request.Parameters)
	helpers.PrintError(err)

	// Deploy ARM Template
	log.Printf("Starting deployment...")

	groupName := request.GroupName
	location := request.Location
	deploymentName := request.DeploymentName

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
