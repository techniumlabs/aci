package azure

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-07-01/storage"
	"github.com/Azure/azure-storage-file-go/2017-07-29/azfile"
	"github.com/writeameer/aci/helpers"
)

// CreateStorageAccount Creates an Azure storage account
func CreateStorageAccount(resourceGroupName string, storageAccountName string) (account storage.Account, keys *[]storage.AccountKey, err error) {
	// Authenticate with Azure
	authorizer, sid := Auth()

	// Setup ARM Client

	client := storage.NewAccountsClient(sid)
	client.Authorizer = authorizer

	location := "Australia East"
	sku := storage.Sku{
		Name: storage.StandardLRS,
		Tier: storage.Standard,
	}
	kind := storage.StorageV2
	accountCreateFuture, err := client.Create(ctx, resourceGroupName, storageAccountName, storage.AccountCreateParameters{
		Location: &location,
		Kind:     kind,
		Sku:      &sku,
	})

	helpers.PrintError(err)

	log.Println("Start storage account creation")
	err = accountCreateFuture.WaitForCompletion(ctx, client.BaseClient.Client)
	log.Println("End storage account creation")
	helpers.PrintError(err)

	account, err = accountCreateFuture.Result(client)
	result, _ := client.ListKeys(ctx, resourceGroupName, storageAccountName)
	keys = result.Keys

	return
}

// CreateAzureFileShare Creates an Azure File Share
func CreateAzureFileShare(resourceGroupName string, storageAccountName string, shareName string) (key string, err error) {

	// Create storage account
	storageAccount, keys, err := CreateStorageAccount(resourceGroupName, storageAccountName)
	storageKey := (*keys)[0]
	key = *storageKey.Value

	// Generate the share url string and create a URL struct
	urlString := fmt.Sprintf("https://%s.file.core.windows.net/%v", *storageAccount.Name, shareName)
	u, _ := url.Parse(urlString)

	// Create new share
	credential, err := azfile.NewSharedKeyCredential(storageAccountName, key)
	shareURL := azfile.NewShareURL(*u, azfile.NewPipeline(credential, azfile.PipelineOptions{}))
	_, err = shareURL.Create(ctx, azfile.Metadata{"createdby": "HiberApp"}, 0)

	if strings.Contains(err.Error(), "The specified share already exists") {
		log.Printf("The share %v already exists", shareURL.String())
		err = nil
	}
	return
}
