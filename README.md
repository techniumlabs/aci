# Overview

Playing with Azure ACI. 

## Setup Creds

Run the following command to create an AZure Service Principle for this app:

```
az ad sp create-for-rbac -n \"<yourAppName>\"
```

And use the output to generate a credential file `~/.mycreds`:

```
export AZURE_TENANT_ID=<blah>
export AZURE_CLIENT_ID=<blah>
export AZURE_CLIENT_SECRET=<blah>
export AZURE_SUBSCRIPTION_ID=<blah>
```

## Testing the App

The app deploys an ACI app and then routes traffic to the ACI app once ready. App runs on `http://localhost:8080`

- dot source your creds:

```
source ~/.mycreds
```

- Run the app

```
go run main.go
```

