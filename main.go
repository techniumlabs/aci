package main

import (
	"github.com/writeameer/aci/apps"
)

func main() {

	resourceGroupName := "hiberapp"
	apps.RunMoodle(resourceGroupName, "moodle-app", "hiberappmoodle", "moodle-mysql", "moodle-web")
}
