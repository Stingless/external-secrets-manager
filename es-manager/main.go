package main
 
import (
	"fmt"
    "os"
    "generator"
    "kube-api"
)

func main() {
    _ = os.Mkdir("vault-es", os.ModePerm)
    _ = generator.EsGenerator()
    fmt.Println("Done !")

}
