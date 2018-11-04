#!/bin/bash -e

#------------------------------------------------------------------------
# Globals
#------------------------------------------------------------------------
ACTION=""

#------------------------------------------------------------------------
# Functions
#------------------------------------------------------------------------

function create () {
  echo "Running create action:"
  az group create --name "$resourceGroup" --location "$location"
  az group deployment create \
    --name "$appName" \
    --resource-group "$resourceGroup" \
    --template-file "$templatefile" \
    --parameters "$parameterfile"
}

function delete() {
  echo "Running delete action:"
  az group deployment delete \
    --name "$appName" \
    --resource-group "$resourceGroup"
}

function checkargs(){
  echo "Checking args..."

  ACTION="$1"
  if [ -z "$ACTION" ]; 
  then
    echo "Must supply action: 'create' or 'delete'"
    exit 1
  fi

  if [ "$ACTION" != "create" ] && [ "$ACTION" != "delete" ];
  then
    echo "Action has to be 'create' or 'delete'"
    exit 1
  fi

  echo "args OK"
}
#------------------------------------------------------------------------
# MAIN
#------------------------------------------------------------------------

# Check arugments
checkargs $*

# Init vars
templatefile="example/azuredeploy.json"
parameterfile="example/azuredeploy.parameters.json"
resourceGroup="hiberapp"
appName="hiberapp"
location="australiaeast"

# Create deployment
if [ "$ACTION" == "create" ]; 
then
  create
elif [ "$ACTION" == "delete" ]
then
  delete
fi 
