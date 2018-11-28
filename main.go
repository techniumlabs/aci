package main

import (
	"log"
	"net/http"

	azure "github.com/writeameer/aci/azure"
	handlers "github.com/writeameer/httphandlers/handlers"
)

func main() {

	storageAccountName := "moodlesharebmggehql"

	resourceGroupName := "aci-example"

	key, _ := azure.CreateAzureFileShare(resourceGroupName, storageAccountName, "moodle-share")

	log.Println(key)
	// azure.CreateStorageAccount(resourceGroupName, "")

	// apps.RunMoodle(resourceGroupName, "moodle4-app")
	// apps.RunWordPress(resourceGroupName, "wordpress-app")
	//doReverseProxy("hiberapp.eastus.azurecontainer.io")
}

func doReverseProxy(originHost string) {
	mux := http.NewServeMux()

	mux.Handle("/", handlers.ReverseProxyHandler(originHost))

	// Listen and Server
	port := ":8080"
	log.Println("Server started on port" + port)
	log.Fatal(http.ListenAndServe(port, mux))
}

// HiberApps defines the interface HiberApps should implement
type HiberApps interface {
	init() bool
	run() bool
	sleep() bool
	destroy() bool
	poll() bool
}

// WordPressApp defines a wordpress hiber app
type WordPressApp struct{}

func (w WordPressApp) init() bool {
	log.Println("init")
	return true
}

func (w WordPressApp) run() bool {
	log.Println("run")
	return true
}

func (w WordPressApp) poll() bool {
	log.Println("run")
	return true
}

func (w WordPressApp) sleep() bool {
	log.Println("sleep")
	return true
}

func (w WordPressApp) destroy() bool {
	log.Println("destroy")
	return true
}
