#!/bin/bash -e

public_key=`cat ~/.ssh/id_rsa.pub`

az vm create \
  --resource-group hiberapp \
  --name myVM \
  --image UbuntuLTS \
  --admin-username writeameer \
  --ssh-key-value "${public_key}"