#!/bin/bash -e

#------------------------------------------------------------------------
# Functions
#------------------------------------------------------------------------

function create () {
  az group create --name "hiberapp" --location "australia east"
}

function delete() {
  az group deployment create \
    --name hiberapp \
    --resource-group hiberapp \
    --template-file azuredeploy.json \
    --parameters azuredeploy.parameters.json
}

#------------------------------------------------------------------------
# MAIN
#------------------------------------------------------------------------

delete
