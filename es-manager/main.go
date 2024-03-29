package main
 
import (
	"context"
	"external-secrets-manager/gitgo"
    "fmt"
	"log"
    "os"
    "strings"
    "io/ioutil"
    "sort" 

    diff "github.com/kylelemons/godebug/diff"
	vault "github.com/hashicorp/vault/api"
)
func EsGenerator() (string) {
    var filetoapply = ""
    var secretkeylist []string 
	config := vault.DefaultConfig()

	config.Address = os.Getenv("VAULT_ADDRESS")

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("Unable to initialize Vault client: %v", err)
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))
    
    var f = ""
    var remoteref = ""
    namespace, err := client.KVv2("es-manager").Get(context.Background(), "namespace")
	if err != nil {
		log.Fatalf("Unable to read secret from es-manager/namespace: %v", err)
	}

	for vaultpath, k8namespace := range namespace.Data {
        
        if _, err := os.Stat("vault-es/"+fmt.Sprintf("%v",vaultpath)); err == nil {
            fmt.Println("Directory already exists: vault-es/"+vaultpath)
        } else if os.IsNotExist(err) {
            _ = os.MkdirAll("vault-es/"+fmt.Sprintf("%v",vaultpath), os.ModePerm)
            fmt.Println("Created new directory: vault-es/"+vaultpath)
        } else {
            fmt.Println("Ambiguous error: vault-es/"+vaultpath)
 
        }
        
        vaultpathlist := strings.Split(vaultpath, "/")
        vaultkv := vaultpathlist[0]
        vaultpathtemp := vaultpathlist [1:]
        vaultpathnested := strings.Join([]string(vaultpathtemp), "/")
        secretpaths, err := client.Logical().List(vaultkv+"/metadata/"+vaultpathnested)
        if err != nil {
	        log.Fatalf("Unable to read secret file path: %v", err)
	    }

        for _, value := range secretpaths.Data {
            for _,fname := range value.([]interface{}) {

            metadata := `---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: ` + fmt.Sprintf("%v",fname) + `-es
  namespace: `+ fmt.Sprintf("%v",k8namespace)+"\nspec:\n  data: "
            f = f + metadata

            secretkeymap, err := client.KVv2(vaultkv).Get(context.Background(),vaultpathnested+"/"+fmt.Sprintf("%v",fname))
	        if err != nil {
                fmt.Println("Skipping: vault-es/"+vaultpath+"/"+fmt.Sprintf("%v",fname))
                f = ""
		        continue
	        }

            //sorting all keys got from vault to avoid random order get requests for keys
            secretkeylist = nil
	        for secretkey, _ := range secretkeymap.Data {
            secretkeylist = append(secretkeylist, fmt.Sprintf("%v",secretkey) )
            }
            sort.Sort(sort.StringSlice(secretkeylist))
	        for _, secretkey := range secretkeylist {
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
    name: `+os.Getenv("CLUSTER_SECRET_STORE")+` 
  target:
    creationPolicy: Owner
    deletionPolicy: Retain
    name: `+fmt.Sprintf("%v",fname)+"-secret\n---"
    
    filetoapply = filetoapply + f   
    if _, err := os.Stat("vault-es/"+fmt.Sprintf("%v",vaultpath)+"/"+fmt.Sprintf("%v",fname)+"-es.yaml"); err == nil {
        esread, err := ioutil.ReadFile("vault-es/"+fmt.Sprintf("%v",vaultpath)+"/"+fmt.Sprintf("%v",fname)+"-es.yaml")
        if err != nil {
            log.Fatal(err)
        }
        if string(esread) == f {
            fmt.Println("No changes in external secret file: vault-es/"+vaultpath+"/"+fmt.Sprintf("%v",fname)+"-es.yaml")
        } else {

            fmt.Println("Configured external secret file in path: vault-es/"+vaultpath+"/"+fmt.Sprintf("%v",fname)+"-es.yaml")
            fmt.Println(diff.Diff ( string(esread), f))
        } 
    } else if os.IsNotExist(err) {
       fmt.Println("Creating external secret file in path: vault-es/"+vaultpath+"/"+fmt.Sprintf("%v",fname)+"-es.yaml")
    } else {
        fmt.Println("Amogus path: vault-es/"+vaultpath+"/"+fmt.Sprintf("%v",fname))
    }
    
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
    gitgo.CloneRepository(os.Getenv("ES_GIT_REPOSITORY"),"vault-es")
    gitgo.CheckForChanges("vault-es")
    fmt.Println("Done !")

}
