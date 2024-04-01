
# <p align="center"> EXTERNAL-SECRETS-MANAGER </p>
<p align="center"> <img src="https://github.com/Stingless/external-secrets-manager/assets/83643646/d307ec82-305d-460c-bd7a-cb3b772ee451" width="400" /> </p>

## <p align="center"> CRD generator for Vault External secrets. </p>
<p align="center">
Integrate vault secrets to your kubernetes all with just one pod! External secret manager works along side of external secrets to automatically deploy secrets which are added to vault. To help you understand what goes on behind the scenes, this is how the script results look like. You can observe the new changes which have been made to vault have been added to the external-secrets CRD file which can then be deployed as per your will. These files will be pushed to a repository of your will where you can review and maintain the vault secrets.
</p>
<p align="center">
 <img src="https://github.com/Stingless/external-secrets-manager/assets/83643646/9bacc4a7-94b4-4a8a-b182-cc57edf953cd" width="700" \>
</p>
<p align="center">
Provide the tool with necessary env variables and run it as a cronjob or a standalone script. The code basically interacts with Vault API to get the keys for generating the CRD files. We then push these changes to a remote repository which can later be deployed to kubernetes. All you need to do is to mention which path you want to deploy in which kubernetes namespace. This can be done through the vault UI itself in es-manager folder inside namespace file 
Here are the sample secrets I have created in vault and respective external secrets created in ./vault-es folder 
</p>
<p align="center">
 <img src="https://github.com/Stingless/external-secrets-manager/assets/83643646/075868ec-a2e2-4c1e-8012-15c976bb91ca" width="700" \>
</p>
<p align="center">
 <img src="https://github.com/Stingless/external-secrets-manager/assets/83643646/418b4efa-8460-46f7-ade2-7aa1b3ba4e5b" width="700" \>
</p>
<p align="center">
 <img src="https://github.com/Stingless/external-secrets-manager/assets/83643646/7c5ee186-e921-4d03-b83a-221f1b210e42" width="700" \>
</p>
<p align="center">
This tool is a consequence of a problem that I had faced working in CloudSEK. 
</p>
