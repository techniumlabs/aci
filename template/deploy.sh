#!/bin/bash -e

az group create --name "hiberapp" --location "australia east"
az group deployment create \
  --name hiberapp \
  --resource-group hiberapp \
  --template-file azuredeploy.json \
  --parameters azuredeploy.parameters.json
  