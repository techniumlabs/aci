#!/bin/bash -e
ACI_PERS_RESOURCE_GROUP="hiberapp"
ACI_PERS_STORAGE_ACCOUNT_NAME="<blah>"
STORAGE_KEY="<blah>"
ACI_PERS_SHARE_NAME="mysql-share"


az container create \
    --resource-group $ACI_PERS_RESOURCE_GROUP \
    --name ameer \
    --image mysql:5.6.24 \
    --ports 3306 \
    --azure-file-volume-account-name $ACI_PERS_STORAGE_ACCOUNT_NAME \
    --azure-file-volume-account-key $STORAGE_KEY \
    --azure-file-volume-share-name $ACI_PERS_SHARE_NAME \
    --azure-file-volume-mount-path /var/lib/mysql \
    --environment-variables MYSQL_ROOT_PASSWORD=password