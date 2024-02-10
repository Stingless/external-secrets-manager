# EXTERNAL-SECRETS-MANAGER
### Tired of creating external secrets for vault? 
### Want a generic external secrets manager to automate everything? 
### You have come to the right place. 
![image](https://github.com/Stingless/external-secrets-manager/assets/83643646/9bacc4a7-94b4-4a8a-b182-cc57edf953cd)

To run , mention vault url and root key in plain text (Living in the edge) inside main.go and simply use `go run main.go` (I still am trying to make it more portable)\
external-secrets-manager basically creates all the external secret files for you based on vault secrets. All you need to do is to mention which path you want to deploy in which kubernetes namespace. This can be done through the vault UI itself in `es-manager` folder inside `namespace` file \
Here are the sample secrets I have created in vault and respective external secrets created in ./vault-es folder \

![image](https://github.com/Stingless/external-secrets-manager/assets/83643646/075868ec-a2e2-4c1e-8012-15c976bb91ca)
![image](https://github.com/Stingless/external-secrets-manager/assets/83643646/418b4efa-8460-46f7-ade2-7aa1b3ba4e5b)
![image](https://github.com/Stingless/external-secrets-manager/assets/83643646/7c5ee186-e921-4d03-b83a-221f1b210e42)
