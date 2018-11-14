package helpers

import (
	"log"
	"net/http"

	"github.com/writeameer/httphandlers/handlers"
)

func deployACIApp(request ArmDeploymentRequest) (siteFQDN string, err error) {
	// Get ARM template and params
	template, err := ReadJSON(request.Template)
	PrintError(err)
	templateParameters, _ := ReadJSON(request.Parameters)
	PrintError(err)

	// Deploy ARM Template
	log.Printf("Starting deployment...")

	groupName := request.GroupName
	location := request.Location
	deploymentName := request.DeploymentName

	// Get Deployment result
	result, err := DeployArmTemplate(groupName, location, deploymentName, template, templateParameters)
	if err != nil {
		log.Printf("Error %v", err)
	}

	// Parse Output
	properties := result.Properties.Outputs.(map[string]interface{})
	propInfo := properties["siteFQDN"].(map[string]interface{})
	siteFQDN = propInfo["value"].(string)

	return
}

func blah() {
	// Deploy ACI and get siteFQDN
	siteFQDN, err := deployACIApp(ArmDeploymentRequest{
		Template:       "./templates/example/wordpress/azuredeploy.json",
		Parameters:     "./templates/example/wordpress/azuredeploy.parameters.json",
		GroupName:      "hiberapp-we",
		Location:       "West Europe",
		DeploymentName: "hiberapp",
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
