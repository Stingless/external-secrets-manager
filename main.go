package main
 
import (
	"context"
	"fmt"
	"log"
    "os"
    "strings"

	vault "github.com/hashicorp/vault/api"
)
func EsGenerator() (string) {
	config := vault.DefaultConfig()

	config.Address = "http://10.43.41.9:8200"

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	client.SetToken("")
    
    var f = ""
    var remoteref = ""
    namespace, err := client.KVv2("es-manager").Get(context.Background(), "namespace")
	if err != nil {
		log.Fatalf("unable to read secret: %v", err)
	}

	for vaultpath, k8namespace := range namespace.Data {
        _ = os.Mkdir("vault-es/"+fmt.Sprintf("%v",vaultpath), os.ModePerm)
        
        vaultpathlist := strings.Split(vaultpath, "/")
        vaultkv := vaultpathlist[0]
        vaultpathtemp := vaultpathlist [1:]
        vaultpathnested := strings.Join([]string(vaultpathtemp), "/")
        secretpaths, err := client.Logical().List(vaultkv+"/metadata/"+vaultpathnested)
        if err != nil {
	        log.Fatalf("unable to read secret: %v", err)
	    }

        for _, value := range secretpaths.Data {
            for _,fname := range value.([]interface{}) {

            metadata := `
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: ` + fmt.Sprintf("%v",fname) + `-es
  namespace: `+ fmt.Sprintf("%v",k8namespace)+"\nspec:\n  data: "
            f = f + metadata

            secretkeymap, err := client.KVv2(vaultpath).Get(context.Background(),fmt.Sprintf("%v",fname))
	        if err != nil {
		        continue
	        }

	        for secretkey, _ := range secretkeymap.Data {
              remoteref = remoteref + `
  - remoteRef:
      conversionStrategy: Default
      decodingStrategy: None
      key: `+ fmt.Sprintf("%v",vaultpath)+"/"+fmt.Sprintf("%v",fname)+`
      property: `+fmt.Sprintf("%v",secretkey)+`
    secretKey: `+fmt.Sprintf("%v",secretkey)
            }
f = f + remoteref +`
  refreshInterval: 15s
  secretStoreRef:
    kind: ClusterSecretStore
    name: vault-backend 
  target:
    creationPolicy: Owner
    deletionPolicy: Retain
    name: `+fmt.Sprintf("%v",fname)+"-secret"
    
    esfile, err := os.Create("vault-es/"+fmt.Sprintf("%v",vaultpath)+"/"+fmt.Sprintf("%v",fname) + "-es.yaml")
    if err != nil {
        panic(err)
    }
    defer func() {
        if err := esfile.Close(); err != nil {
            panic(err)
        }
    }()
    if _, err := esfile.WriteString(f); err != nil {
        panic(err)
    }
    remoteref = ""
    f = ""
            }
        }
    }

    return f 

}

func main() {
    _ = os.Mkdir("vault-es", os.ModePerm)
    _ = EsGenerator()
    fmt.Println("Done !")

}
