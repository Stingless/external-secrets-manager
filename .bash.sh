#!/bin/bash

# Read key-value pairs from file
secrets_file="secrets.txt"
secrets=$(cat "$secrets_file")

# Connect to Vault
export VAULT_ADDR="http://10.43.41.9:8200/"
export VAULT_TOKEN="hvs.QjdJ8dXOWhbB1poIjxu40jfu"

# Write secrets to Vault
secrets_dir="hello-world/one-secrets/" # specify the path to the secrets directory in Vault
out=""
for secret in $secrets; do
  key=$(echo "$secret" | cut -d= -f1)
  value=$(echo "$secret" | cut -d= -f2)
  out=$out' '$key=$value
done
vault kv put "$secrets_dir" $out
# vault kv put "$secrets_dir" $key=$value







